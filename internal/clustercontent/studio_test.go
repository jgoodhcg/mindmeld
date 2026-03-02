package clustercontent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStudioToLibrarySkipsNonReadyPrompts(t *testing.T) {
	src := StudioSource{
		Version:        "v1",
		CreatedByLabel: "cluster-studio",
		AxisSets: []AxisSet{
			{
				Slug:      "axis-a",
				XMinLabel: "Low",
				XMaxLabel: "High",
				YMinLabel: "Slow",
				YMaxLabel: "Fast",
				MinRating: 10,
			},
		},
		Prompts: []StudioPrompt{
			{Slug: "ready", Text: "Ready", MinRating: 20, AxisSlugs: []string{"axis-a"}, Status: "ready"},
			{Slug: "draft", Text: "Draft", MinRating: 20, AxisSlugs: []string{"axis-a"}, Status: "draft"},
			{Slug: "archived", Text: "Archived", MinRating: 20, AxisSlugs: []string{"axis-a"}, Status: "archived"},
			{Slug: "implicit-ready", Text: "Implicit", MinRating: 20, AxisSlugs: []string{"axis-a"}, Status: ""},
		},
	}

	lib := src.ToLibrary()
	if len(lib.Prompts) != 2 {
		t.Fatalf("expected 2 importable prompts, got %d", len(lib.Prompts))
	}
	if lib.Prompts[0].Slug != "ready" || lib.Prompts[1].Slug != "implicit-ready" {
		t.Fatalf("unexpected import prompt set: %+v", lib.Prompts)
	}
}

func TestAnalyzeStudioWarnings(t *testing.T) {
	src := StudioSource{
		Version:        "v1",
		CreatedByLabel: "cluster-studio",
		AxisSets: []AxisSet{
			{Slug: "axis-a", XMinLabel: "L", XMaxLabel: "H", YMinLabel: "S", YMaxLabel: "F", MinRating: 10},
			{Slug: "axis-unused", XMinLabel: "L", XMaxLabel: "H", YMinLabel: "S", YMaxLabel: "F", MinRating: 10},
		},
		Prompts: []StudioPrompt{
			{Slug: "p1", Text: "Same text", MinRating: 20, AxisSlugs: []string{"axis-a"}, Status: "ready"},
			{Slug: "p2", Text: " same   text ", MinRating: 20, AxisSlugs: []string{"missing-axis"}, Status: "draft"},
		},
	}

	diag := AnalyzeStudio(src)
	if diag.Summary.ExactDuplicateTextCount != 1 {
		t.Fatalf("expected one duplicate text group, got %+v", diag.Summary)
	}
	if diag.Summary.MissingAxisRefCount != 1 {
		t.Fatalf("expected one missing axis ref, got %+v", diag.Summary)
	}
	if diag.Summary.OrphanAxisCount != 1 {
		t.Fatalf("expected one orphan axis, got %+v", diag.Summary)
	}
	if len(diag.Warnings) == 0 {
		t.Fatal("expected warnings")
	}
}

func TestLoadStudioOrLibraryFallsBackToLibrary(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "library.json")
	content := `{
  "version":"v1",
  "created_by_label":"cluster-library-v1",
  "axis_sets":[{"slug":"axis-a","x_min_label":"L","x_max_label":"H","y_min_label":"S","y_max_label":"F","min_rating":10}],
  "prompts":[{"slug":"p1","text":"Prompt","min_rating":20,"axis_slugs":["axis-a"]}]
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	src, err := LoadStudioOrLibrary(path)
	if err != nil {
		t.Fatalf("load studio or library: %v", err)
	}
	if len(src.Prompts) != 1 {
		t.Fatalf("unexpected source prompts: %+v", src.Prompts)
	}
	if len(src.ToLibrary().Prompts) != 1 {
		t.Fatalf("expected library conversion to treat empty status as ready: %+v", src.Prompts)
	}
}
