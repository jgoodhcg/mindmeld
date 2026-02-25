package clustercontent

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadPromptSourceRowsCSVAndTSV(t *testing.T) {
	dir := t.TempDir()

	csvContent := "slug,text,min_rating,axis_slugs,theme,status,notes\n" +
		"alpha,Prompt alpha,20,axis-a|axis-b,team,ready,\n" +
		"beta,Prompt beta,mild,axis-c,fun,draft,hold\n"
	if err := os.WriteFile(filepath.Join(dir, "a.csv"), []byte(csvContent), 0o644); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	tsvContent := "slug\ttext\tmin_rating\taxis_slugs\tstatus\n" +
		"gamma\tPrompt gamma\tadults\taxis-x|axis-y\tready\n"
	if err := os.WriteFile(filepath.Join(dir, "b.tsv"), []byte(tsvContent), 0o644); err != nil {
		t.Fatalf("write tsv: %v", err)
	}

	rows, report, err := LoadPromptSourceRows(dir)
	if err != nil {
		t.Fatalf("load prompt rows: %v", err)
	}

	if report.SourceFiles != 2 || report.RowsRead != 3 || report.RowsReady != 2 || report.RowsDraft != 1 {
		t.Fatalf("unexpected report: %+v", report)
	}

	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}

	if rows[0].Slug != "alpha" || rows[0].Status != "ready" || rows[0].MinRating != 20 {
		t.Fatalf("unexpected first row: %+v", rows[0])
	}
	if !reflect.DeepEqual(rows[0].AxisSlugs, []string{"axis-a", "axis-b"}) {
		t.Fatalf("unexpected axis slugs: %+v", rows[0].AxisSlugs)
	}

	if rows[1].Slug != "beta" || rows[1].Status != "draft" || rows[1].MinRating != 10 {
		t.Fatalf("unexpected second row: %+v", rows[1])
	}

	if rows[2].Slug != "gamma" || rows[2].MinRating != 30 {
		t.Fatalf("unexpected third row: %+v", rows[2])
	}
}

func TestLoadPromptSourceRowsRejectsUnknownColumns(t *testing.T) {
	dir := t.TempDir()
	content := "slug,text,min_rating,axis_slugs,oops\nalpha,Prompt,20,axis-a,x\n"
	if err := os.WriteFile(filepath.Join(dir, "prompts.csv"), []byte(content), 0o644); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	_, _, err := LoadPromptSourceRows(dir)
	if err == nil {
		t.Fatal("expected error for unknown column")
	}
}

func TestBuildLibraryFromPromptRowsSkipsDrafts(t *testing.T) {
	template := Library{
		Version:        "v1",
		CreatedByLabel: "cluster-library-v1",
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
		Prompts: []Prompt{{Slug: "old", Text: "old", MinRating: 10, AxisSlugs: []string{"axis-a"}}},
	}
	rows := []PromptSourceRow{
		{Slug: "ready-one", Text: "Ready one", MinRating: 20, AxisSlugs: []string{"axis-a"}, Status: "ready"},
		{Slug: "draft-one", Text: "Draft one", MinRating: 20, AxisSlugs: []string{"axis-a"}, Status: "draft"},
		{Slug: "ready-two", Text: "Ready two", MinRating: 10, AxisSlugs: []string{"axis-a"}, Status: "ready"},
	}

	lib := BuildLibraryFromPromptRows(template, rows)

	if lib.Version != "v1" || lib.CreatedByLabel != "cluster-library-v1" {
		t.Fatalf("unexpected metadata: %+v", lib)
	}
	if len(lib.AxisSets) != 1 || lib.AxisSets[0].Slug != "axis-a" {
		t.Fatalf("unexpected axis sets: %+v", lib.AxisSets)
	}
	if len(lib.Prompts) != 2 {
		t.Fatalf("expected 2 prompts, got %d", len(lib.Prompts))
	}
	if lib.Prompts[0].Slug != "ready-one" || lib.Prompts[1].Slug != "ready-two" {
		t.Fatalf("unexpected prompt order: %+v", lib.Prompts)
	}
}
