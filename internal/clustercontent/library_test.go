package clustercontent

import "testing"

func TestValidateSuccess(t *testing.T) {
	lib := Library{
		Version:        "v1",
		CreatedByLabel: "cluster-library-test",
		AxisSets: []AxisSet{
			{
				Slug:      "decision-style",
				XMinLabel: "Consensus",
				XMaxLabel: "Contrarian",
				YMinLabel: "Low stakes",
				YMaxLabel: "High stakes",
				MinRating: 20,
			},
		},
		Prompts: []Prompt{
			{
				Slug:      "kickoff-sync",
				Text:      "Best way to start a weekly sync",
				MinRating: 20,
				AxisSlugs: []string{"decision-style"},
			},
		},
	}

	report, pairs, err := Validate(lib)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	if len(pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(pairs))
	}
	if report.PairCount != 1 {
		t.Fatalf("expected pair count 1, got %d", report.PairCount)
	}
	if report.PairCountByRating[20] != 1 || report.PairCountByRating[30] != 1 {
		t.Fatalf("unexpected rating counts: %+v", report.PairCountByRating)
	}
}

func TestValidateFailsOnUnknownAxisSlug(t *testing.T) {
	lib := Library{
		Version:        "v1",
		CreatedByLabel: "cluster-library-test",
		AxisSets: []AxisSet{
			{
				Slug:      "decision-style",
				XMinLabel: "Consensus",
				XMaxLabel: "Contrarian",
				YMinLabel: "Low stakes",
				YMaxLabel: "High stakes",
				MinRating: 20,
			},
		},
		Prompts: []Prompt{
			{
				Slug:      "kickoff-sync",
				Text:      "Best way to start a weekly sync",
				MinRating: 20,
				AxisSlugs: []string{"missing-axis"},
			},
		},
	}

	_, _, err := Validate(lib)
	if err == nil {
		t.Fatal("expected validation error for missing axis slug")
	}
}

func TestPairUUIDDeterministic(t *testing.T) {
	a := PairUUID("prompt-a", "axis-a")
	b := PairUUID("prompt-a", "axis-a")
	if a != b {
		t.Fatalf("expected deterministic UUIDs to match: %s vs %s", a, b)
	}
}
