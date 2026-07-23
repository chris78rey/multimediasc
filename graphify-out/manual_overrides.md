# Manual Graphify Overrides

Source of truth for hand-maintained graph relations in this checkout:

- `graphify-out/manual_overrides.json`
- `graphify-out/apply_manual_overrides.py`

Reapply after any `graphify update` that regenerates `graphify-out/graph.json`:

```bash
/home/crrb/.local/share/uv/tools/graphifyy/bin/python graphify-out/apply_manual_overrides.py
graphify cluster-only . --graph graphify-out/graph.json
```

The overlay is idempotent:

- nodes are added only if missing by `id` or `(label, source_file, source_location)`
- links are deduplicated by `(source, target, relation)`
- removal specs are applied before additions
