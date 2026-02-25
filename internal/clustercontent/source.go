package clustercontent

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type SourceMeta struct {
	Version        string `json:"version"`
	CreatedByLabel string `json:"created_by_label"`
}

type SourceDirLoadReport struct {
	AxisRows   int
	PromptRows PromptSourceLoadReport
}

func LoadSourceDir(dir string) (Library, SourceDirLoadReport, error) {
	meta, err := loadSourceMeta(filepath.Join(dir, "meta.json"))
	if err != nil {
		return Library{}, SourceDirLoadReport{}, err
	}

	axes, axisCount, err := loadAxisSourceFile(filepath.Join(dir, "axes.tsv"))
	if err != nil {
		return Library{}, SourceDirLoadReport{}, err
	}

	promptsPath, err := findSourcePromptsFile(dir)
	if err != nil {
		return Library{}, SourceDirLoadReport{}, err
	}
	promptRows, promptReport, err := loadPromptSourceFile(promptsPath)
	if err != nil {
		return Library{}, SourceDirLoadReport{}, err
	}
	if promptReport.SourceFiles == 0 {
		promptReport.SourceFiles = 1
	}

	lib := BuildLibraryFromPromptRows(Library{
		Version:        meta.Version,
		CreatedByLabel: meta.CreatedByLabel,
		AxisSets:       axes,
	}, promptRows)

	return lib, SourceDirLoadReport{
		AxisRows:   axisCount,
		PromptRows: promptReport,
	}, nil
}

func loadSourceMeta(path string) (SourceMeta, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return SourceMeta{}, err
	}

	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.DisallowUnknownFields()

	var meta SourceMeta
	if err := dec.Decode(&meta); err != nil {
		return SourceMeta{}, fmt.Errorf("%s: %w", path, err)
	}
	if strings.TrimSpace(meta.Version) == "" {
		return SourceMeta{}, fmt.Errorf("%s: version is required", path)
	}
	if strings.TrimSpace(meta.CreatedByLabel) == "" {
		return SourceMeta{}, fmt.Errorf("%s: created_by_label is required", path)
	}
	return meta, nil
}

func findSourcePromptsFile(dir string) (string, error) {
	candidates := []string{
		filepath.Join(dir, "prompts.tsv"),
		filepath.Join(dir, "prompts.csv"),
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("missing prompts source file in %s (expected prompts.tsv or prompts.csv)", dir)
}

func loadAxisSourceFile(path string) ([]AxisSet, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, fmt.Errorf("%s: missing header row", path)
		}
		return nil, 0, fmt.Errorf("%s: read header: %w", path, err)
	}

	colIndex, err := parseAxisSourceHeader(path, header)
	if err != nil {
		return nil, 0, err
	}

	var (
		axes  []AxisSet
		errs  []error
		count int
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

		axis, rowErr := parseAxisSourceRow(path, rowNum, record, colIndex)
		count++
		if rowErr != nil {
			errs = append(errs, rowErr)
			continue
		}
		axes = append(axes, axis)
	}

	if len(errs) > 0 {
		return axes, count, errors.Join(errs...)
	}
	return axes, count, nil
}

func parseAxisSourceHeader(path string, rawHeader []string) (map[string]int, error) {
	allowed := map[string]bool{
		"slug":        true,
		"x_min_label": true,
		"x_max_label": true,
		"y_min_label": true,
		"y_max_label": true,
		"min_rating":  true,
	}
	required := []string{"slug", "x_min_label", "x_max_label", "y_min_label", "y_max_label", "min_rating"}

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

func parseAxisSourceRow(path string, rowNum int, record []string, colIndex map[string]int) (AxisSet, error) {
	var errs []error
	axis := AxisSet{
		Slug:      strings.TrimSpace(fieldValue(record, colIndex, "slug")),
		XMinLabel: strings.TrimSpace(fieldValue(record, colIndex, "x_min_label")),
		XMaxLabel: strings.TrimSpace(fieldValue(record, colIndex, "x_max_label")),
		YMinLabel: strings.TrimSpace(fieldValue(record, colIndex, "y_min_label")),
		YMaxLabel: strings.TrimSpace(fieldValue(record, colIndex, "y_max_label")),
	}
	if axis.Slug == "" {
		errs = append(errs, fmt.Errorf("%s:%d: slug is required", path, rowNum))
	}
	if axis.XMinLabel == "" || axis.XMaxLabel == "" || axis.YMinLabel == "" || axis.YMaxLabel == "" {
		errs = append(errs, fmt.Errorf("%s:%d: all axis labels are required", path, rowNum))
	}

	minRating, err := parsePromptSourceMinRating(fieldValue(record, colIndex, "min_rating"))
	if err != nil {
		errs = append(errs, fmt.Errorf("%s:%d: %w", path, rowNum, err))
	} else {
		axis.MinRating = minRating
	}

	if len(errs) > 0 {
		return axis, errors.Join(errs...)
	}
	return axis, nil
}
