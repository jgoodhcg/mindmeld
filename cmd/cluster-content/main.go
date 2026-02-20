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
	fmt.Println("  go run ./cmd/cluster-content validate -file content/cluster/library.v1.json")
	fmt.Println("  go run ./cmd/cluster-content import -file content/cluster/library.v1.json [flags]")
	fmt.Println()
	fmt.Println("Import Flags:")
	fmt.Println("  -database-url string   Explicit DB URL (fallback: DATABASE_URL env)")
	fmt.Println("  -env string            Target environment: dev|prod (default dev)")
	fmt.Println("  -dry-run               Preview DB changes without writes")
	fmt.Println("  -allow-production      Required for production-like DB URLs")
}

func runValidate(args []string) error {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	file := fs.String("file", "content/cluster/library.v1.json", "Path to cluster library JSON")
	if err := fs.Parse(args); err != nil {
		return err
	}

	lib, report, _, err := loadAndValidate(*file)
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

	lib, report, pairs, err := loadAndValidate(*file)
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

func loadAndValidate(path string) (clustercontent.Library, clustercontent.Report, []clustercontent.Pair, error) {
	lib, err := clustercontent.Load(path)
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
	fmt.Printf("Pairs available for Kids (%d): %d\n", contentrating.Kids, report.PairCountByRating[contentrating.Kids])
	fmt.Printf("Pairs available for Work (%d): %d\n", contentrating.Work, report.PairCountByRating[contentrating.Work])
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
