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

type Library struct {
	Version        string    `json:"version"`
	CreatedByLabel string    `json:"created_by_label"`
	AxisSets       []AxisSet `json:"axis_sets"`
	Prompts        []Prompt  `json:"prompts"`
}

type AxisSet struct {
	Slug      string `json:"slug"`
	XMinLabel string `json:"x_min_label"`
	XMaxLabel string `json:"x_max_label"`
	YMinLabel string `json:"y_min_label"`
	YMaxLabel string `json:"y_max_label"`
	MinRating int16  `json:"min_rating"`
}

type Prompt struct {
	Slug      string   `json:"slug"`
	Text      string   `json:"text"`
	MinRating int16    `json:"min_rating"`
	AxisSlugs []string `json:"axis_slugs"`
}

type Pair struct {
	PromptSlug string
	AxisSlug   string
}

type Report struct {
	PromptCount       int
	AxisSetCount      int
	PairCount         int
	PairCountByRating map[int16]int
}

func Load(path string) (Library, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Library{}, err
	}

	dec := json.NewDecoder(strings.NewReader(string(raw)))
	dec.DisallowUnknownFields()

	var lib Library
	if err := dec.Decode(&lib); err != nil {
		return Library{}, err
	}

	return lib, nil
}

func Validate(lib Library) (Report, []Pair, error) {
	var errs []error
	report := Report{
		PairCountByRating: map[int16]int{
			contentrating.Kids:   0,
			contentrating.Work:   0,
			contentrating.Adults: 0,
		},
	}

	if strings.TrimSpace(lib.Version) == "" {
		errs = append(errs, errors.New("version is required"))
	}
	if strings.TrimSpace(lib.CreatedByLabel) == "" {
		errs = append(errs, errors.New("created_by_label is required"))
	}
	if len(lib.AxisSets) == 0 {
		errs = append(errs, errors.New("at least one axis set is required"))
	}
	if len(lib.Prompts) == 0 {
		errs = append(errs, errors.New("at least one prompt is required"))
	}

	axisBySlug := make(map[string]AxisSet, len(lib.AxisSets))
	for i, axis := range lib.AxisSets {
		prefix := fmt.Sprintf("axis_sets[%d]", i)
		if !isSlug(axis.Slug) {
			errs = append(errs, fmt.Errorf("%s has invalid slug: %q", prefix, axis.Slug))
		}
		if _, exists := axisBySlug[axis.Slug]; exists {
			errs = append(errs, fmt.Errorf("duplicate axis slug %q", axis.Slug))
		}
		if strings.TrimSpace(axis.XMinLabel) == "" || strings.TrimSpace(axis.XMaxLabel) == "" || strings.TrimSpace(axis.YMinLabel) == "" || strings.TrimSpace(axis.YMaxLabel) == "" {
			errs = append(errs, fmt.Errorf("%s has empty axis labels", prefix))
		}
		if !contentrating.IsValid(axis.MinRating) {
			errs = append(errs, fmt.Errorf("%s has invalid min_rating %d", prefix, axis.MinRating))
		}
		axisBySlug[axis.Slug] = axis
	}

	promptBySlug := make(map[string]Prompt, len(lib.Prompts))
	pairSeen := make(map[string]bool)
	pairs := make([]Pair, 0)
	for i, prompt := range lib.Prompts {
		prefix := fmt.Sprintf("prompts[%d]", i)
		if !isSlug(prompt.Slug) {
			errs = append(errs, fmt.Errorf("%s has invalid slug: %q", prefix, prompt.Slug))
		}
		if _, exists := promptBySlug[prompt.Slug]; exists {
			errs = append(errs, fmt.Errorf("duplicate prompt slug %q", prompt.Slug))
		}
		if strings.TrimSpace(prompt.Text) == "" {
			errs = append(errs, fmt.Errorf("%s has empty text", prefix))
		}
		if !contentrating.IsValid(prompt.MinRating) {
			errs = append(errs, fmt.Errorf("%s has invalid min_rating %d", prefix, prompt.MinRating))
		}
		if len(prompt.AxisSlugs) == 0 {
			errs = append(errs, fmt.Errorf("%s must include at least one axis slug", prefix))
		}
		if hasDuplicateStrings(prompt.AxisSlugs) {
			errs = append(errs, fmt.Errorf("%s contains duplicate axis slugs", prefix))
		}

		for _, axisSlug := range prompt.AxisSlugs {
			axis, ok := axisBySlug[axisSlug]
			if !ok {
				errs = append(errs, fmt.Errorf("%s references missing axis slug %q", prefix, axisSlug))
				continue
			}

			pairKey := prompt.Slug + "|" + axisSlug
			if pairSeen[pairKey] {
				errs = append(errs, fmt.Errorf("duplicate prompt/axis pair %q", pairKey))
				continue
			}
			pairSeen[pairKey] = true
			pairs = append(pairs, Pair{PromptSlug: prompt.Slug, AxisSlug: axisSlug})

			pairMinRating := max(prompt.MinRating, axis.MinRating)
			for _, rating := range []int16{contentrating.Kids, contentrating.Work, contentrating.Adults} {
				if pairMinRating <= rating {
					report.PairCountByRating[rating]++
				}
			}
		}

		promptBySlug[prompt.Slug] = prompt
	}

	report.PromptCount = len(promptBySlug)
	report.AxisSetCount = len(axisBySlug)
	report.PairCount = len(pairs)

	if len(errs) > 0 {
		return report, pairs, errors.Join(errs...)
	}
	return report, pairs, nil
}

func hasDuplicateStrings(values []string) bool {
	seen := make(map[string]bool, len(values))
	for _, value := range values {
		if seen[value] {
			return true
		}
		seen[value] = true
	}
	return false
}

func isSlug(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			continue
		}
		return false
	}
	return true
}

func max(a, b int16) int16 {
	return maxInt16([]int16{a, b})
}

func maxInt16(values []int16) int16 {
	if len(values) == 0 {
		return 0
	}
	sorted := slices.Clone(values)
	slices.Sort(sorted)
	return sorted[len(sorted)-1]
}
