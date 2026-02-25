# Cluster Prompt Sources

Row-based authoring files for Cluster prompts.

- Supported formats: `.csv`, `.tsv`
- `cluster-content build` loads `mild`, `polite`, `adults` first (in that order), then other files lexically
- `status=draft` rows are kept in source but skipped from generated `content/cluster/library.v1.json`

Columns:

- `slug`
- `text`
- `min_rating`
- `axis_slugs` (pipe-separated)
- `theme` (optional)
- `status` (`draft|ready`)
- `notes` (optional)
