package clustercontent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/jgoodhcg/mindmeld/internal/contentrating"
)

type StudioSource struct {
	Version        string         `json:"version"`
	CreatedByLabel string         `json:"created_by_label"`
	AxisSets       []AxisSet      `json:"axis_sets"`
	Prompts        []StudioPrompt `json:"prompts"`
}

type StudioPrompt struct {
	Slug      string   `json:"slug"`
	Text      string   `json:"text"`
	MinRating int16    `json:"min_rating"`
	AxisSlugs []string `json:"axis_slugs"`
	Theme     string   `json:"theme,omitempty"`
	Status    string   `json:"status,omitempty"`
	Notes     string   `json:"notes,omitempty"`
}

type StudioReviewSummary struct {
	TotalPrompts            int
	ReadyPrompts            int
	DraftPrompts            int
	ArchivedPrompts         int
	UnknownStatusPrompts    int
	AxisCount               int
	OrphanAxisCount         int
	MissingAxisRefCount     int
	ExactDuplicateTextCount int
}

type StudioWarning struct {
	Code    string
	Message string
	Items   []string
}

type StudioDiagnostics struct {
	Summary        StudioReviewSummary
	Warnings       []StudioWarning
	AxisUsageCount map[string]int
}

func LoadStudio(path string) (StudioSource, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return StudioSource{}, err
	}

	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.DisallowUnknownFields()

	var src StudioSource
	if err := dec.Decode(&src); err != nil {
		return StudioSource{}, err
	}
	return src, nil
}

func SaveStudio(path string, src StudioSource) error {
	out, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')
	return os.WriteFile(path, out, 0o644)
}

func StudioFromLibrary(lib Library) StudioSource {
	src := StudioSource{
		Version:        lib.Version,
		CreatedByLabel: lib.CreatedByLabel,
		AxisSets:       slices.Clone(lib.AxisSets),
		Prompts:        make([]StudioPrompt, 0, len(lib.Prompts)),
	}
	for _, prompt := range lib.Prompts {
		src.Prompts = append(src.Prompts, StudioPrompt{
			Slug:      prompt.Slug,
			Text:      prompt.Text,
			MinRating: prompt.MinRating,
			AxisSlugs: slices.Clone(prompt.AxisSlugs),
			Status:    "ready",
		})
	}
	return src
}

func (src StudioSource) ToLibrary() Library {
	lib := Library{
		Version:        src.Version,
		CreatedByLabel: src.CreatedByLabel,
		AxisSets:       slices.Clone(src.AxisSets),
		Prompts:        make([]Prompt, 0, len(src.Prompts)),
	}
	for _, prompt := range src.Prompts {
		if normalizeStudioStatus(prompt.Status) != "ready" {
			continue
		}
		lib.Prompts = append(lib.Prompts, Prompt{
			Slug:      prompt.Slug,
			Text:      prompt.Text,
			MinRating: prompt.MinRating,
			AxisSlugs: slices.Clone(prompt.AxisSlugs),
		})
	}
	return lib
}

func ValidateStudio(src StudioSource) (StudioDiagnostics, error) {
	diag := AnalyzeStudio(src)

	// Reuse existing runtime validation for the importable subset (ready prompts).
	if _, _, err := Validate(src.ToLibrary()); err != nil {
		return diag, err
	}
	return diag, nil
}

func AnalyzeStudio(src StudioSource) StudioDiagnostics {
	axisBySlug := make(map[string]AxisSet, len(src.AxisSets))
	axisUsage := make(map[string]int, len(src.AxisSets))
	for _, axis := range src.AxisSets {
		axisBySlug[axis.Slug] = axis
		axisUsage[axis.Slug] = 0
	}

	var (
		missingRefs []string
	)
	dupTextMap := map[string][]string{}
	summary := StudioReviewSummary{
		AxisCount: len(src.AxisSets),
	}

	for _, p := range src.Prompts {
		summary.TotalPrompts++
		switch normalizeStudioStatus(p.Status) {
		case "ready":
			summary.ReadyPrompts++
		case "draft":
			summary.DraftPrompts++
		case "archived":
			summary.ArchivedPrompts++
		default:
			summary.UnknownStatusPrompts++
		}

		if textKey := normalizeDuplicateTextKey(p.Text); textKey != "" {
			dupTextMap[textKey] = append(dupTextMap[textKey], p.Slug)
		}

		for _, axisSlug := range p.AxisSlugs {
			if _, ok := axisBySlug[axisSlug]; !ok {
				missingRefs = append(missingRefs, fmt.Sprintf("%s -> %s", p.Slug, axisSlug))
				continue
			}
			axisUsage[axisSlug]++
		}

		// Flag obviously invalid ratings early for review even before Validate reports.
		if !contentrating.IsValid(p.MinRating) {
			missingRefs = append(missingRefs, fmt.Sprintf("%s -> invalid min_rating %d", p.Slug, p.MinRating))
		}
	}

	var warnings []StudioWarning
	if len(missingRefs) > 0 {
		summary.MissingAxisRefCount = len(missingRefs)
		warnings = append(warnings, StudioWarning{
			Code:    "missing-axis-ref",
			Message: "Prompts reference axis slugs that are not defined in axis_sets",
			Items:   limitStrings(missingRefs, 12),
		})
	}

	var duplicateGroups []string
	for _, slugs := range dupTextMap {
		if len(slugs) < 2 {
			continue
		}
		slices.Sort(slugs)
		summary.ExactDuplicateTextCount++
		duplicateGroups = append(duplicateGroups, strings.Join(slugs, ", "))
	}
	if len(duplicateGroups) > 0 {
		slices.Sort(duplicateGroups)
		warnings = append(warnings, StudioWarning{
			Code:    "duplicate-text",
			Message: "Exact duplicate prompt text detected (normalized)",
			Items:   limitStrings(duplicateGroups, 12),
		})
	}

	var orphanAxes []string
	for slug, count := range axisUsage {
		if count == 0 {
			orphanAxes = append(orphanAxes, slug)
		}
	}
	if len(orphanAxes) > 0 {
		slices.Sort(orphanAxes)
		summary.OrphanAxisCount = len(orphanAxes)
		warnings = append(warnings, StudioWarning{
			Code:    "orphan-axis",
			Message: "Axes defined but unused by any prompt",
			Items:   limitStrings(orphanAxes, 12),
		})
	}

	return StudioDiagnostics{
		Summary:        summary,
		Warnings:       warnings,
		AxisUsageCount: axisUsage,
	}
}

func normalizeStudioStatus(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "ready":
		return "ready"
	case "draft":
		return "draft"
	case "archived":
		return "archived"
	default:
		return "unknown"
	}
}

func normalizeDuplicateTextKey(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	return strings.Join(strings.Fields(value), " ")
}

func limitStrings(values []string, n int) []string {
	if len(values) <= n {
		return values
	}
	out := slices.Clone(values[:n])
	out = append(out, fmt.Sprintf("...and %d more", len(values)-n))
	return out
}

func LoadStudioFromSourceDir(dir string) (StudioSource, error) {
	meta, err := loadSourceMeta(dir + "/meta.json")
	if err != nil {
		return StudioSource{}, err
	}
	axes, _, err := loadAxisSourceFile(dir + "/axes.tsv")
	if err != nil {
		return StudioSource{}, err
	}
	promptsPath, err := findSourcePromptsFile(dir)
	if err != nil {
		return StudioSource{}, err
	}
	rows, _, err := loadPromptSourceFile(promptsPath)
	if err != nil {
		return StudioSource{}, err
	}

	src := StudioSource{
		Version:        meta.Version,
		CreatedByLabel: meta.CreatedByLabel,
		AxisSets:       axes,
		Prompts:        make([]StudioPrompt, 0, len(rows)),
	}
	for _, row := range rows {
		src.Prompts = append(src.Prompts, StudioPrompt{
			Slug:      row.Slug,
			Text:      row.Text,
			MinRating: row.MinRating,
			AxisSlugs: slices.Clone(row.AxisSlugs),
			Theme:     strings.TrimSpace(row.Theme),
			Status:    normalizeStudioStatus(row.Status),
			Notes:     strings.TrimSpace(row.Notes),
		})
	}
	return src, nil
}

func LoadStudioOrLibrary(path string) (StudioSource, error) {
	src, err := LoadStudio(path)
	if err == nil {
		return src, nil
	}

	// Fallback for one-time migration convenience: existing import library JSON.
	var parseErr *json.SyntaxError
	if errors.As(err, &parseErr) || strings.Contains(err.Error(), "unknown field") || strings.Contains(err.Error(), "cannot unmarshal") {
		lib, libErr := Load(path)
		if libErr != nil {
			return StudioSource{}, err
		}
		return StudioFromLibrary(lib), nil
	}

	lib, libErr := Load(path)
	if libErr == nil {
		return StudioFromLibrary(lib), nil
	}
	return StudioSource{}, err
}
