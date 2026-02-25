package clustercontent

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jgoodhcg/mindmeld/internal/contentrating"
)

type PromptSourceRow struct {
	SourceFile string
	SourceLine int

	Slug      string
	Text      string
	MinRating int16
	AxisSlugs []string

	Theme  string
	Status string
	Notes  string
}

type PromptSourceLoadReport struct {
	SourceFiles int
	RowsRead    int
	RowsReady   int
	RowsDraft   int
}

func LoadPromptSourceRows(dir string) ([]PromptSourceRow, PromptSourceLoadReport, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, PromptSourceLoadReport{}, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		switch strings.ToLower(filepath.Ext(name)) {
		case ".csv", ".tsv":
			files = append(files, filepath.Join(dir, name))
		}
	}
	slices.SortFunc(files, comparePromptSourceFiles)

	if len(files) == 0 {
		return nil, PromptSourceLoadReport{}, fmt.Errorf("no CSV/TSV prompt source files found in %s", dir)
	}

	var (
		rows   []PromptSourceRow
		errs   []error
		report PromptSourceLoadReport
	)
	report.SourceFiles = len(files)

	for _, path := range files {
		fileRows, fileReport, err := loadPromptSourceFile(path)
		report.RowsRead += fileReport.RowsRead
		report.RowsReady += fileReport.RowsReady
		report.RowsDraft += fileReport.RowsDraft
		if err != nil {
			errs = append(errs, err)
		}
		rows = append(rows, fileRows...)
	}

	if len(errs) > 0 {
		return rows, report, errors.Join(errs...)
	}
	return rows, report, nil
}

func BuildLibraryFromPromptRows(axisTemplate Library, rows []PromptSourceRow) Library {
	lib := Library{
		Version:        axisTemplate.Version,
		CreatedByLabel: axisTemplate.CreatedByLabel,
		AxisSets:       slices.Clone(axisTemplate.AxisSets),
		Prompts:        make([]Prompt, 0, len(rows)),
	}

	for _, row := range rows {
		if row.Status != "ready" {
			continue
		}
		lib.Prompts = append(lib.Prompts, Prompt{
			Slug:      row.Slug,
			Text:      row.Text,
			MinRating: row.MinRating,
			AxisSlugs: slices.Clone(row.AxisSlugs),
		})
	}

	return lib
}

func SaveLibrary(path string, lib Library) error {
	out, err := json.MarshalIndent(lib, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(path, out, 0o644)
}

func loadPromptSourceFile(path string) ([]PromptSourceRow, PromptSourceLoadReport, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, PromptSourceLoadReport{}, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	switch strings.ToLower(filepath.Ext(path)) {
	case ".tsv":
		reader.Comma = '\t'
	}

	header, err := reader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, PromptSourceLoadReport{}, fmt.Errorf("%s: missing header row", path)
		}
		return nil, PromptSourceLoadReport{}, fmt.Errorf("%s: read header: %w", path, err)
	}

	colIndex, err := parsePromptSourceHeader(path, header)
	if err != nil {
		return nil, PromptSourceLoadReport{}, err
	}

	var (
		rows   []PromptSourceRow
		errs   []error
		report PromptSourceLoadReport
	)
	for rowNum := 2; ; rowNum++ {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("%s:%d: read row: %w", path, rowNum, err))
			continue
		}

		if isBlankRecord(record) {
			continue
		}

		row, rowErr := parsePromptSourceRow(path, rowNum, record, colIndex)
		report.RowsRead++
		if rowErr != nil {
			errs = append(errs, rowErr)
			continue
		}

		switch row.Status {
		case "draft":
			report.RowsDraft++
		case "ready":
			report.RowsReady++
		}
		rows = append(rows, row)
	}

	if len(errs) > 0 {
		return rows, report, errors.Join(errs...)
	}
	return rows, report, nil
}

func parsePromptSourceHeader(path string, rawHeader []string) (map[string]int, error) {
	allowed := map[string]bool{
		"slug":       true,
		"text":       true,
		"min_rating": true,
		"axis_slugs": true,
		"theme":      true,
		"status":     true,
		"notes":      true,
	}
	required := []string{"slug", "text", "min_rating", "axis_slugs"}

	colIndex := make(map[string]int, len(rawHeader))
	var errs []error
	for i, col := range rawHeader {
		name := normalizePromptSourceHeader(col)
		if name == "" {
			continue
		}
		if !allowed[name] {
			errs = append(errs, fmt.Errorf("%s: unknown column %q", path, strings.TrimSpace(col)))
			continue
		}
		if _, exists := colIndex[name]; exists {
			errs = append(errs, fmt.Errorf("%s: duplicate column %q", path, name))
			continue
		}
		colIndex[name] = i
	}

	for _, key := range required {
		if _, ok := colIndex[key]; !ok {
			errs = append(errs, fmt.Errorf("%s: missing required column %q", path, key))
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return colIndex, nil
}

func normalizePromptSourceHeader(value string) string {
	value = strings.TrimSpace(strings.TrimPrefix(value, "\ufeff"))
	value = strings.ToLower(value)
	value = strings.ReplaceAll(value, " ", "_")
	return value
}

func parsePromptSourceRow(path string, rowNum int, record []string, colIndex map[string]int) (PromptSourceRow, error) {
	var errs []error
	row := PromptSourceRow{
		SourceFile: filepath.Base(path),
		SourceLine: rowNum,
		Theme:      fieldValue(record, colIndex, "theme"),
		Notes:      fieldValue(record, colIndex, "notes"),
	}

	row.Slug = strings.TrimSpace(fieldValue(record, colIndex, "slug"))
	if row.Slug == "" {
		errs = append(errs, fmt.Errorf("%s:%d: slug is required", path, rowNum))
	}

	row.Text = strings.TrimSpace(fieldValue(record, colIndex, "text"))
	if row.Text == "" {
		errs = append(errs, fmt.Errorf("%s:%d: text is required", path, rowNum))
	}

	minRating, err := parsePromptSourceMinRating(fieldValue(record, colIndex, "min_rating"))
	if err != nil {
		errs = append(errs, fmt.Errorf("%s:%d: %w", path, rowNum, err))
	} else {
		row.MinRating = minRating
	}

	row.AxisSlugs = splitPipeList(fieldValue(record, colIndex, "axis_slugs"))
	if len(row.AxisSlugs) == 0 {
		errs = append(errs, fmt.Errorf("%s:%d: axis_slugs is required", path, rowNum))
	}

	status := strings.ToLower(strings.TrimSpace(fieldValue(record, colIndex, "status")))
	if status == "" {
		status = "ready"
	}
	switch status {
	case "draft", "ready":
		row.Status = status
	default:
		errs = append(errs, fmt.Errorf("%s:%d: unsupported status %q (want draft|ready)", path, rowNum, status))
	}

	if len(errs) > 0 {
		return row, errors.Join(errs...)
	}
	return row, nil
}

func parsePromptSourceMinRating(raw string) (int16, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	switch value {
	case "mild", "kids", "kid":
		return contentrating.Kids, nil
	case "polite", "work":
		return contentrating.Work, nil
	case "adult", "adults":
		return contentrating.Adults, nil
	}

	id, err := contentrating.ParseID(value)
	if err != nil {
		return 0, fmt.Errorf("invalid min_rating %q", raw)
	}
	return id, nil
}

func splitPipeList(raw string) []string {
	parts := strings.Split(raw, "|")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		values = append(values, value)
	}
	return values
}

func fieldValue(record []string, colIndex map[string]int, name string) string {
	i, ok := colIndex[name]
	if !ok || i >= len(record) {
		return ""
	}
	return record[i]
}

func isBlankRecord(record []string) bool {
	for _, field := range record {
		if strings.TrimSpace(field) != "" {
			return false
		}
	}
	return true
}

func comparePromptSourceFiles(a, b string) int {
	aBase := strings.ToLower(strings.TrimSuffix(filepath.Base(a), filepath.Ext(a)))
	bBase := strings.ToLower(strings.TrimSuffix(filepath.Base(b), filepath.Ext(b)))

	aRank := promptSourceFileRank(aBase)
	bRank := promptSourceFileRank(bBase)
	if aRank != bRank {
		return aRank - bRank
	}
	if aBase != bBase {
		return strings.Compare(aBase, bBase)
	}
	return strings.Compare(a, b)
}

func promptSourceFileRank(base string) int {
	switch base {
	case "mild":
		return 0
	case "polite":
		return 1
	case "adults":
		return 2
	default:
		return 100
	}
}
