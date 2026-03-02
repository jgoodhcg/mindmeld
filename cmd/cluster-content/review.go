package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/jgoodhcg/mindmeld/internal/clustercontent"
	"github.com/jgoodhcg/mindmeld/internal/contentrating"
)

func runBootstrapStudio(args []string) error {
	fs := flag.NewFlagSet("bootstrap-studio", flag.ContinueOnError)
	outFile := fs.String("file", "content/cluster/studio.v1.json", "Output studio JSON path")
	sourceDir := fs.String("source-dir", "content/cluster/source", "TSV source directory to bootstrap from")
	libraryFile := fs.String("library-file", "content/cluster/library.v1.json", "Fallback library JSON path")
	if err := fs.Parse(args); err != nil {
		return err
	}

	var (
		src clustercontent.StudioSource
		err error
	)
	if _, statErr := os.Stat(strings.TrimSpace(*sourceDir)); statErr == nil {
		src, err = clustercontent.LoadStudioFromSourceDir(strings.TrimSpace(*sourceDir))
		if err != nil {
			return err
		}
	} else {
		src, err = clustercontent.LoadStudioOrLibrary(strings.TrimSpace(*libraryFile))
		if err != nil {
			return err
		}
	}

	if _, err := clustercontent.ValidateStudio(src); err != nil {
		return err
	}
	if err := clustercontent.SaveStudio(strings.TrimSpace(*outFile), src); err != nil {
		return err
	}

	fmt.Printf("Wrote studio source: %s\n", strings.TrimSpace(*outFile))
	fmt.Printf("Axis sets: %d\n", len(src.AxisSets))
	fmt.Printf("Prompts: %d\n", len(src.Prompts))
	fmt.Println("Bootstrap: OK")
	return nil
}

func runReview(args []string) error {
	fs := flag.NewFlagSet("review", flag.ContinueOnError)
	file := fs.String("file", "content/cluster/studio.v1.json", "Studio JSON path")
	listen := fs.String("listen", "127.0.0.1:8097", "HTTP listen address")
	allowNonLocal := fs.Bool("allow-non-local", false, "Allow binding to non-local interfaces")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := validateReviewListenAddr(strings.TrimSpace(*listen), *allowNonLocal); err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		renderStudioReview(w, r, strings.TrimSpace(*file))
	})
	mux.HandleFunc("/__meta", func(w http.ResponseWriter, r *http.Request) {
		info, err := os.Stat(strings.TrimSpace(*file))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"file":          strings.TrimSpace(*file),
			"size":          info.Size(),
			"modified_unix": info.ModTime().Unix(),
		})
	})

	fmt.Printf("Cluster review UI: http://%s\n", strings.TrimSpace(*listen))
	fmt.Printf("Studio file: %s\n", strings.TrimSpace(*file))
	fmt.Println("Review server binds localhost by default; use -allow-non-local to override.")
	return http.ListenAndServe(strings.TrimSpace(*listen), mux)
}

func validateReviewListenAddr(addr string, allowNonLocal bool) error {
	host, _, err := net.SplitHostPort(strings.TrimSpace(addr))
	if err != nil {
		return fmt.Errorf("invalid -listen address %q: %w", addr, err)
	}
	if allowNonLocal {
		return nil
	}

	switch strings.ToLower(strings.TrimSpace(host)) {
	case "127.0.0.1", "localhost", "::1":
		return nil
	default:
		return errors.New("review server must bind localhost by default; use -allow-non-local to override")
	}
}

type reviewFilters struct {
	Query  string
	Status string
	Theme  string
	Rating string
	Sort   string
}

type reviewPromptRow struct {
	Slug        string
	Text        string
	Rating      int16
	RatingLabel string
	Status      string
	Theme       string
	AxisSlugs   []string
	AxisCount   int
	Notes       string
}

type reviewPageData struct {
	SourcePath     string
	SourceVersion  string
	CreatedByLabel string
	Filters        reviewFilters
	StatusOptions  []string
	ThemeOptions   []string
	PromptRows     []reviewPromptRow
	PromptRowCount int
	Diagnostics    clustercontent.StudioDiagnostics
	LoadError      string
	ValidateError  string
	FileModUnix    int64
}

func renderStudioReview(w http.ResponseWriter, r *http.Request, path string) {
	page := reviewPageData{
		SourcePath: path,
		Filters: reviewFilters{
			Query:  strings.TrimSpace(r.URL.Query().Get("q")),
			Status: strings.ToLower(strings.TrimSpace(r.URL.Query().Get("status"))),
			Theme:  strings.TrimSpace(r.URL.Query().Get("theme")),
			Rating: strings.TrimSpace(r.URL.Query().Get("rating")),
			Sort:   strings.TrimSpace(r.URL.Query().Get("sort")),
		},
	}

	if page.Filters.Sort == "" {
		page.Filters.Sort = "slug"
	}

	src, err := clustercontent.LoadStudio(path)
	if err != nil {
		page.LoadError = err.Error()
		writeReviewPage(w, page)
		return
	}
	page.SourceVersion = src.Version
	page.CreatedByLabel = src.CreatedByLabel
	if info, err := os.Stat(path); err == nil {
		page.FileModUnix = info.ModTime().Unix()
	}

	diag, validateErr := clustercontent.ValidateStudio(src)
	page.Diagnostics = diag
	if validateErr != nil {
		page.ValidateError = validateErr.Error()
	}

	statusSet := map[string]bool{}
	themeSet := map[string]bool{}
	rows := make([]reviewPromptRow, 0, len(src.Prompts))
	for _, p := range src.Prompts {
		status := normalizeStatusForReview(p.Status)
		theme := strings.TrimSpace(p.Theme)
		statusSet[status] = true
		if theme != "" {
			themeSet[theme] = true
		}

		if !matchesReviewFilters(p, page.Filters) {
			continue
		}

		rows = append(rows, reviewPromptRow{
			Slug:        p.Slug,
			Text:        p.Text,
			Rating:      p.MinRating,
			RatingLabel: contentrating.Label(p.MinRating),
			Status:      status,
			Theme:       theme,
			AxisSlugs:   slices.Clone(p.AxisSlugs),
			AxisCount:   len(p.AxisSlugs),
			Notes:       strings.TrimSpace(p.Notes),
		})
	}

	for status := range statusSet {
		page.StatusOptions = append(page.StatusOptions, status)
	}
	sort.Strings(page.StatusOptions)
	for theme := range themeSet {
		page.ThemeOptions = append(page.ThemeOptions, theme)
	}
	sort.Strings(page.ThemeOptions)

	sortReviewPromptRows(rows, page.Filters.Sort)
	page.PromptRows = rows
	page.PromptRowCount = len(rows)

	writeReviewPage(w, page)
}

func normalizeStatusForReview(value string) string {
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

func matchesReviewFilters(p clustercontent.StudioPrompt, filters reviewFilters) bool {
	query := strings.ToLower(strings.TrimSpace(filters.Query))
	if query != "" {
		slug := strings.ToLower(strings.TrimSpace(p.Slug))
		text := strings.ToLower(strings.TrimSpace(p.Text))
		if !strings.Contains(slug, query) && !strings.Contains(text, query) {
			return false
		}
	}

	if status := strings.TrimSpace(filters.Status); status != "" && status != "all" {
		if normalizeStatusForReview(p.Status) != status {
			return false
		}
	}

	if theme := strings.TrimSpace(filters.Theme); theme != "" {
		if strings.TrimSpace(p.Theme) != theme {
			return false
		}
	}

	if rating := strings.TrimSpace(filters.Rating); rating != "" {
		want, err := strconv.Atoi(rating)
		if err != nil || int16(want) != p.MinRating {
			return false
		}
	}

	return true
}

func sortReviewPromptRows(rows []reviewPromptRow, key string) {
	switch key {
	case "rating":
		sort.SliceStable(rows, func(i, j int) bool {
			if rows[i].Rating != rows[j].Rating {
				return rows[i].Rating < rows[j].Rating
			}
			return rows[i].Slug < rows[j].Slug
		})
	case "status":
		sort.SliceStable(rows, func(i, j int) bool {
			if rows[i].Status != rows[j].Status {
				return rows[i].Status < rows[j].Status
			}
			return rows[i].Slug < rows[j].Slug
		})
	case "axes":
		sort.SliceStable(rows, func(i, j int) bool {
			if rows[i].AxisCount != rows[j].AxisCount {
				return rows[i].AxisCount > rows[j].AxisCount
			}
			return rows[i].Slug < rows[j].Slug
		})
	default:
		sort.SliceStable(rows, func(i, j int) bool { return rows[i].Slug < rows[j].Slug })
	}
}

func writeReviewPage(w http.ResponseWriter, page reviewPageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := reviewPageTemplate.Execute(w, page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var reviewPageTemplate = template.Must(template.New("review").Funcs(template.FuncMap{
	"join": strings.Join,
}).Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Cluster Content Review</title>
  <style>
    :root { --bg:#f3efe7; --ink:#18211f; --muted:#5b655f; --card:#fffdf9; --line:#d8d1c6; --accent:#1f6f5e; --warn:#9a5a00; --err:#9b1c1c; }
    * { box-sizing: border-box; }
    body { margin:0; font-family: ui-sans-serif, system-ui, sans-serif; background: radial-gradient(circle at 0% 0%, #fffdf9, #f3efe7 60%); color:var(--ink); }
    .wrap { max-width: 1400px; margin: 0 auto; padding: 20px; }
    .hero { display:grid; gap:10px; margin-bottom:16px; }
    .hero h1 { margin:0; font-size: 28px; letter-spacing: -0.02em; }
    .muted { color: var(--muted); }
    .banner { border:1px solid var(--line); background: var(--card); border-radius: 12px; padding:12px 14px; margin-bottom: 12px; }
    .banner.err { border-color: #efc6c6; background: #fff4f4; color: var(--err); }
    .banner.warn { border-color: #efd7b3; background: #fff8ee; color: #744100; }
    .grid { display:grid; gap:12px; grid-template-columns: repeat(4, minmax(0,1fr)); margin-bottom: 16px; }
    .card { background: rgba(255,253,249,.9); border:1px solid var(--line); border-radius: 14px; padding: 12px; }
    .metric { font-size: 24px; font-weight: 700; }
    .label { font-size: 12px; text-transform: uppercase; letter-spacing: .08em; color:var(--muted); }
    form.filters { display:grid; gap:10px; grid-template-columns: 2fr repeat(4, minmax(0,1fr)); margin-bottom: 12px; }
    .filters input, .filters select { width:100%; border:1px solid var(--line); border-radius: 10px; padding:10px 12px; background:#fff; }
    table { width:100%; border-collapse: separate; border-spacing: 0; background:#fff; border:1px solid var(--line); border-radius: 14px; overflow:hidden; }
    thead th { position: sticky; top: 0; background:#f7f2ea; color:var(--muted); font-size: 12px; text-transform: uppercase; letter-spacing: .06em; text-align:left; padding:10px 12px; border-bottom:1px solid var(--line); }
    tbody td { vertical-align: top; padding:10px 12px; border-bottom:1px solid #eee6da; }
    tbody tr:last-child td { border-bottom:none; }
    .pill { display:inline-block; padding:2px 8px; border-radius:999px; border:1px solid var(--line); background:#faf7f0; font-size: 12px; }
    .status-ready { border-color:#b7ddd5; background:#ecfaf6; color:#0e5e4f; }
    .status-draft { border-color:#e9d6a2; background:#fff6df; color:#825100; }
    .status-archived { border-color:#d8d8d8; background:#f3f3f3; color:#5a5a5a; }
    .axis-list { color: var(--muted); font-size: 12px; line-height:1.4; }
    .warn-list { margin: 8px 0 0 18px; padding:0; }
    .warn-list li { margin: 2px 0; }
    @media (max-width: 960px) {
      .grid { grid-template-columns: 1fr 1fr; }
      form.filters { grid-template-columns: 1fr 1fr; }
      table { display:block; overflow:auto; }
    }
  </style>
</head>
<body>
  <div class="wrap">
    <div class="hero">
      <h1>Cluster Content Review</h1>
      <div class="muted">{{.SourcePath}} · version {{.SourceVersion}} · {{.CreatedByLabel}}</div>
    </div>

    {{if .LoadError}}
      <div class="banner err"><strong>Load error:</strong> {{.LoadError}}</div>
    {{else}}
      {{if .ValidateError}}<div class="banner err"><strong>Validation error:</strong> {{.ValidateError}}</div>{{end}}
      {{range .Diagnostics.Warnings}}
        <div class="banner warn">
          <div><strong>{{.Message}}</strong></div>
          {{if .Items}}
          <ul class="warn-list">{{range .Items}}<li>{{.}}</li>{{end}}</ul>
          {{end}}
        </div>
      {{end}}

      <div class="grid">
        <div class="card"><div class="label">Visible Prompts</div><div class="metric">{{.PromptRowCount}}</div></div>
        <div class="card"><div class="label">Ready / Draft</div><div class="metric">{{.Diagnostics.Summary.ReadyPrompts}} / {{.Diagnostics.Summary.DraftPrompts}}</div></div>
        <div class="card"><div class="label">Axes / Orphans</div><div class="metric">{{.Diagnostics.Summary.AxisCount}} / {{.Diagnostics.Summary.OrphanAxisCount}}</div></div>
        <div class="card"><div class="label">Missing Refs / Duplicates</div><div class="metric">{{.Diagnostics.Summary.MissingAxisRefCount}} / {{.Diagnostics.Summary.ExactDuplicateTextCount}}</div></div>
      </div>

      <form class="filters" method="get">
        <input type="search" name="q" placeholder="Search slug or phrase" value="{{.Filters.Query}}">
        <select name="status">
          <option value="">All statuses</option>
          {{range .StatusOptions}}<option value="{{.}}" {{if eq $.Filters.Status .}}selected{{end}}>{{.}}</option>{{end}}
        </select>
        <select name="theme">
          <option value="">All themes</option>
          {{range .ThemeOptions}}<option value="{{.}}" {{if eq $.Filters.Theme .}}selected{{end}}>{{.}}</option>{{end}}
        </select>
        <select name="rating">
          <option value="">All ratings</option>
          <option value="10" {{if eq .Filters.Rating "10"}}selected{{end}}>Mild (10)</option>
          <option value="20" {{if eq .Filters.Rating "20"}}selected{{end}}>Polite (20)</option>
          <option value="30" {{if eq .Filters.Rating "30"}}selected{{end}}>Adults (30)</option>
        </select>
        <select name="sort">
          <option value="slug" {{if eq .Filters.Sort "slug"}}selected{{end}}>Sort: Slug</option>
          <option value="rating" {{if eq .Filters.Sort "rating"}}selected{{end}}>Sort: Rating</option>
          <option value="status" {{if eq .Filters.Sort "status"}}selected{{end}}>Sort: Status</option>
          <option value="axes" {{if eq .Filters.Sort "axes"}}selected{{end}}>Sort: Axis Count</option>
        </select>
      </form>

      <table>
        <thead>
          <tr>
            <th>Status</th>
            <th>Rating</th>
            <th>Slug</th>
            <th>Prompt</th>
            <th>Theme</th>
            <th>Axes</th>
            <th>Notes</th>
          </tr>
        </thead>
        <tbody>
          {{range .PromptRows}}
          <tr>
            <td><span class="pill status-{{.Status}}">{{.Status}}</span></td>
            <td>{{.RatingLabel}} ({{.Rating}})</td>
            <td><code>{{.Slug}}</code></td>
            <td>{{.Text}}</td>
            <td>{{if .Theme}}{{.Theme}}{{else}}<span class="muted">-</span>{{end}}</td>
            <td><div>{{.AxisCount}}</div><div class="axis-list">{{join .AxisSlugs ", "}}</div></td>
            <td>{{if .Notes}}{{.Notes}}{{else}}<span class="muted">-</span>{{end}}</td>
          </tr>
          {{else}}
          <tr><td colspan="7" class="muted">No prompts matched the current filters.</td></tr>
          {{end}}
        </tbody>
      </table>
    {{end}}
  </div>
</body>
</html>`))
