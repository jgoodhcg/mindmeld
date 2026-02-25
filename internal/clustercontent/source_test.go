package clustercontent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSourceDir(t *testing.T) {
	dir := t.TempDir()

	meta := `{"version":"v1","created_by_label":"cluster-library-v1"}`
	if err := os.WriteFile(filepath.Join(dir, "meta.json"), []byte(meta), 0o644); err != nil {
		t.Fatalf("write meta: %v", err)
	}

	axes := "slug\tx_min_label\tx_max_label\ty_min_label\ty_max_label\tmin_rating\n" +
		"axis-a\tLow\tHigh\tSlow\tFast\t10\n"
	if err := os.WriteFile(filepath.Join(dir, "axes.tsv"), []byte(axes), 0o644); err != nil {
		t.Fatalf("write axes: %v", err)
	}

	prompts := "slug\ttext\tmin_rating\taxis_slugs\ttheme\tstatus\tnotes\n" +
		"prompt-ready\tPrompt ready\tpolite\taxis-a\tteam\tready\t\n" +
		"prompt-draft\tPrompt draft\t20\taxis-a\tteam\tdraft\t\n"
	if err := os.WriteFile(filepath.Join(dir, "prompts.tsv"), []byte(prompts), 0o644); err != nil {
		t.Fatalf("write prompts: %v", err)
	}

	lib, loadReport, err := LoadSourceDir(dir)
	if err != nil {
		t.Fatalf("load source dir: %v", err)
	}

	if loadReport.AxisRows != 1 {
		t.Fatalf("expected 1 axis row, got %d", loadReport.AxisRows)
	}
	if loadReport.PromptRows.RowsRead != 2 || loadReport.PromptRows.RowsReady != 1 || loadReport.PromptRows.RowsDraft != 1 {
		t.Fatalf("unexpected prompt load report: %+v", loadReport.PromptRows)
	}
	if lib.Version != "v1" || lib.CreatedByLabel != "cluster-library-v1" {
		t.Fatalf("unexpected meta in library: %+v", lib)
	}
	if len(lib.AxisSets) != 1 || lib.AxisSets[0].Slug != "axis-a" {
		t.Fatalf("unexpected axis sets: %+v", lib.AxisSets)
	}
	if len(lib.Prompts) != 1 || lib.Prompts[0].Slug != "prompt-ready" {
		t.Fatalf("unexpected prompts: %+v", lib.Prompts)
	}
}

func TestLoadSourceDirFailsWhenAxisFileMissing(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "meta.json"), []byte(`{"version":"v1","created_by_label":"x"}`), 0o644); err != nil {
		t.Fatalf("write meta: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "prompts.tsv"), []byte("slug\ttext\tmin_rating\taxis_slugs\n"), 0o644); err != nil {
		t.Fatalf("write prompts: %v", err)
	}

	if _, _, err := LoadSourceDir(dir); err == nil {
		t.Fatal("expected missing axes.tsv to fail")
	}
}
