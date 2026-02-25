# Cluster Source Files

Canonical editable source for Cluster content.

Workflow:

1. Edit `prompts.tsv` (main review/edit file)
2. Edit `axes.tsv` when adding/changing axis definitions
3. Run `validate -source-dir content/cluster/source`
4. Run `import -source-dir content/cluster/source -dry-run`
5. Run `import -source-dir content/cluster/source`
6. Optionally regenerate `library.v1.json` with `build -source-dir ...`

Notes:

- `prompts.tsv` is prompt-centric (one row per prompt)
- `axis_slugs` is a pipe-separated list of references into `axes.tsv`
- Missing axis references fail validation/import (guards against prompt TSV typos)
- `status=draft` rows are kept in source but skipped from generated/imported prompt sets
