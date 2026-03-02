package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jgoodhcg/mindmeld/internal/clustercontent"
	"github.com/jgoodhcg/mindmeld/internal/contentrating"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "bootstrap-studio":
		if err := runBootstrapStudio(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
	case "build":
		if err := runBuild(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
	case "review":
		if err := runReview(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
	case "validate":
		if err := runValidate(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
	case "import":
		if err := runImport(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println("Cluster Content Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run ./cmd/cluster-content bootstrap-studio [-source-dir content/cluster/source] [-file content/cluster/studio.v1.json]")
	fmt.Println("  go run ./cmd/cluster-content build [-studio-file content/cluster/studio.v1.json | -source-dir content/cluster/source] [-file content/cluster/library.v1.json]")
	fmt.Println("  go run ./cmd/cluster-content review [-file content/cluster/studio.v1.json] [-listen 127.0.0.1:8097]")
	fmt.Println("  go run ./cmd/cluster-content validate [-studio-file content/cluster/studio.v1.json | -source-dir content/cluster/source | -file content/cluster/library.v1.json]")
	fmt.Println("  go run ./cmd/cluster-content import [-studio-file content/cluster/studio.v1.json | -source-dir content/cluster/source | -file content/cluster/library.v1.json] [flags]")
	fmt.Println()
	fmt.Println("Studio Flags:")
	fmt.Println("  bootstrap-studio:")
	fmt.Println("    -file string         Output studio JSON path (default content/cluster/studio.v1.json)")
	fmt.Println("    -source-dir string   Existing TSV source directory (preferred bootstrap source)")
	fmt.Println("    -library-file string Fallback import library JSON path")
	fmt.Println("  review:")
	fmt.Println("    -file string         Studio JSON path (default content/cluster/studio.v1.json)")
	fmt.Println("    -listen string       HTTP listen address (default 127.0.0.1:8097)")
	fmt.Println("    -allow-non-local     Allow binding review server to non-local interfaces")
	fmt.Println()
	fmt.Println("Build Flags:")
	fmt.Println("  -file string           Output library JSON path (default content/cluster/library.v1.json)")
	fmt.Println("  -studio-file string    Studio JSON source path (preferred for review-first workflow)")
	fmt.Println("  -source-dir string     Canonical source directory (meta.json, axes.tsv, prompts.tsv)")
	fmt.Println()
	fmt.Println("Import Flags:")
	fmt.Println("  -database-url string   Explicit DB URL (fallback: DATABASE_URL env)")
	fmt.Println("  -env string            Target environment: dev|prod (default dev)")
	fmt.Println("  -dry-run               Preview DB changes without writes")
	fmt.Println("  -allow-production      Required for production-like DB URLs")
}

func runBuild(args []string) error {
	fs := flag.NewFlagSet("build", flag.ContinueOnError)
	file := fs.String("file", "content/cluster/library.v1.json", "Output cluster library JSON path")
	studioFile := fs.String("studio-file", "", "Studio JSON source path")
	sourceDir := fs.String("source-dir", "content/cluster/source", "Source directory containing meta.json, axes.tsv, and prompts.tsv")
	if err := fs.Parse(args); err != nil {
		return err
	}

	var (
		lib              clustercontent.Library
		promptLoadReport clustercontent.PromptSourceLoadReport
		axisRowCount     int
		err              error
	)
	if strings.TrimSpace(*studioFile) != "" {
		studio, loadErr := clustercontent.LoadStudio(strings.TrimSpace(*studioFile))
		if loadErr != nil {
			return loadErr
		}
		lib = studio.ToLibrary()
		axisRowCount = len(studio.AxisSets)
		for _, p := range studio.Prompts {
			promptLoadReport.RowsRead++
			switch normalizeStatusForReview(p.Status) {
			case "draft":
				promptLoadReport.RowsDraft++
			case "ready":
				promptLoadReport.RowsReady++
			}
		}
		promptLoadReport.SourceFiles = 1
	} else {
		var sourceReport clustercontent.SourceDirLoadReport
		lib, sourceReport, err = clustercontent.LoadSourceDir(strings.TrimSpace(*sourceDir))
		if err != nil {
			return err
		}
		promptLoadReport = sourceReport.PromptRows
		axisRowCount = sourceReport.AxisRows
	}

	report, _, err := clustercontent.Validate(lib)
	if err != nil {
		return err
	}

	if err := clustercontent.SaveLibrary(strings.TrimSpace(*file), lib); err != nil {
		return err
	}

	if strings.TrimSpace(*studioFile) != "" {
		fmt.Printf("Studio source: %s\n", strings.TrimSpace(*studioFile))
	} else {
		fmt.Printf("Source dir: %s\n", strings.TrimSpace(*sourceDir))
	}
	fmt.Printf("Axis rows read: %d\n", axisRowCount)
	fmt.Printf("Prompt rows read: %d (ready=%d, draft=%d)\n", promptLoadReport.RowsRead, promptLoadReport.RowsReady, promptLoadReport.RowsDraft)
	fmt.Printf("Generated prompts: %d\n", len(lib.Prompts))
	fmt.Printf("Wrote: %s\n", strings.TrimSpace(*file))
	printReport(report)
	printTargetGap(report)
	fmt.Println("Build: OK")
	return nil
}

func runValidate(args []string) error {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	file := fs.String("file", "content/cluster/library.v1.json", "Path to cluster library JSON")
	studioFile := fs.String("studio-file", "", "Studio JSON source path")
	sourceDir := fs.String("source-dir", "", "Source directory containing meta.json, axes.tsv, and prompts.tsv")
	if err := fs.Parse(args); err != nil {
		return err
	}

	lib, report, _, err := loadAndValidate(*file, *sourceDir, *studioFile)
	if err != nil {
		return err
	}

	fmt.Printf("Library version: %s\n", lib.Version)
	fmt.Printf("Created by label: %s\n", lib.CreatedByLabel)
	printReport(report)
	printTargetGap(report)
	fmt.Println("Validation: OK")
	return nil
}

func runImport(args []string) error {
	fs := flag.NewFlagSet("import", flag.ContinueOnError)
	file := fs.String("file", "content/cluster/library.v1.json", "Path to cluster library JSON")
	studioFile := fs.String("studio-file", "", "Studio JSON source path")
	sourceDir := fs.String("source-dir", "", "Source directory containing meta.json, axes.tsv, and prompts.tsv")
	databaseURLFlag := fs.String("database-url", "", "Explicit database URL (fallback: DATABASE_URL)")
	targetEnv := fs.String("env", "dev", "Target environment: dev|prod")
	dryRun := fs.Bool("dry-run", false, "Preview change plan without DB writes")
	allowProduction := fs.Bool("allow-production", false, "Required for production-like DB URLs")
	if err := fs.Parse(args); err != nil {
		return err
	}

	databaseURL, err := resolveDatabaseURL(*databaseURLFlag)
	if err != nil {
		return err
	}
	if err := validateImportSafety(strings.TrimSpace(*targetEnv), databaseURL, *allowProduction); err != nil {
		return err
	}

	lib, report, pairs, err := loadAndValidate(*file, *sourceDir, *studioFile)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	plan, err := clustercontent.Analyze(ctx, pool, lib, pairs)
	if err != nil {
		return err
	}

	fmt.Printf("Target env: %s\n", normalizeEnv(*targetEnv))
	printReport(report)
	printTargetGap(report)
	printPlan(plan)

	if *dryRun {
		fmt.Println("Dry-run: no writes applied.")
		return nil
	}

	if err := clustercontent.Import(ctx, pool, lib, pairs); err != nil {
		return err
	}
	fmt.Println("Import: OK")
	return nil
}

func loadAndValidate(path string, sourceDir string, studioFile string) (clustercontent.Library, clustercontent.Report, []clustercontent.Pair, error) {
	var (
		lib clustercontent.Library
		err error
	)
	switch {
	case strings.TrimSpace(studioFile) != "":
		var src clustercontent.StudioSource
		src, err = clustercontent.LoadStudio(strings.TrimSpace(studioFile))
		if err == nil {
			lib = src.ToLibrary()
		}
	case strings.TrimSpace(sourceDir) != "":
		lib, _, err = clustercontent.LoadSourceDir(strings.TrimSpace(sourceDir))
	default:
		lib, err = clustercontent.Load(path)
	}
	if err != nil {
		return clustercontent.Library{}, clustercontent.Report{}, nil, err
	}

	report, pairs, err := clustercontent.Validate(lib)
	if err != nil {
		return clustercontent.Library{}, clustercontent.Report{}, nil, err
	}
	return lib, report, pairs, nil
}

func resolveDatabaseURL(flagValue string) (string, error) {
	if strings.TrimSpace(flagValue) != "" {
		return strings.TrimSpace(flagValue), nil
	}
	if strings.TrimSpace(os.Getenv("DATABASE_URL")) != "" {
		return strings.TrimSpace(os.Getenv("DATABASE_URL")), nil
	}
	return "", errors.New("database URL is required (set -database-url or DATABASE_URL)")
}

func validateImportSafety(targetEnv string, databaseURL string, allowProduction bool) error {
	env := normalizeEnv(targetEnv)
	isProdLikeURL := isProductionLikeURL(databaseURL)

	if isProdLikeURL && !allowProduction {
		return errors.New("refusing to run against production-like DB without -allow-production")
	}
	if env == "prod" && !allowProduction {
		return errors.New("refusing prod import without -allow-production")
	}
	if isProdLikeURL && env != "prod" {
		return errors.New("production-like DB requires -env=prod")
	}
	if env == "prod" && !isProdLikeURL {
		return errors.New("-env=prod provided but database URL looks local; use -env=dev instead")
	}

	return nil
}

func normalizeEnv(env string) string {
	value := strings.ToLower(strings.TrimSpace(env))
	switch value {
	case "prod", "production":
		return "prod"
	default:
		return "dev"
	}
}

func isProductionLikeURL(databaseURL string) bool {
	parsed, err := url.Parse(strings.TrimSpace(databaseURL))
	if err != nil {
		return true
	}

	host := strings.ToLower(strings.TrimSpace(parsed.Hostname()))
	if host == "" {
		return true
	}

	localHosts := map[string]bool{
		"localhost": true,
		"127.0.0.1": true,
		"0.0.0.0":   true,
		"::1":       true,
	}
	if localHosts[host] {
		return false
	}
	if strings.HasSuffix(host, ".local") {
		return false
	}

	return true
}

func printReport(report clustercontent.Report) {
	fmt.Printf("Axis sets: %d\n", report.AxisSetCount)
	fmt.Printf("Prompts: %d\n", report.PromptCount)
	fmt.Printf("Pairs: %d\n", report.PairCount)
	fmt.Printf("Pairs available for Mild (%d): %d\n", contentrating.Kids, report.PairCountByRating[contentrating.Kids])
	fmt.Printf("Pairs available for Polite (%d): %d\n", contentrating.Work, report.PairCountByRating[contentrating.Work])
	fmt.Printf("Pairs available for Adults (%d): %d\n", contentrating.Adults, report.PairCountByRating[contentrating.Adults])
}

func printTargetGap(report clustercontent.Report) {
	target := 500
	if report.PairCount >= target {
		fmt.Printf("Target reached: %d/%d pairs\n", report.PairCount, target)
		return
	}
	fmt.Printf("Target gap: need %d more pairs to reach %d\n", target-report.PairCount, target)
}

func printPlan(plan clustercontent.ImportPlan) {
	fmt.Println("Planned DB changes:")
	printEntityPlan("Prompts", plan.Prompts)
	printEntityPlan("Axis sets", plan.AxisSets)
	printEntityPlan("Pairs", plan.Pairs)
}

func printEntityPlan(label string, plan clustercontent.EntityPlan) {
	fmt.Printf("- %s: desired=%d, managed_existing=%d, create=%d, upsert=%d, reactivate=%d, deactivate=%d\n",
		label,
		plan.DesiredCount,
		plan.ManagedExistingCount,
		plan.CreateCount,
		plan.UpsertCount,
		plan.ReactivateCount,
		plan.DeactivateCount,
	)
}
