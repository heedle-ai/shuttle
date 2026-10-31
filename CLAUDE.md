# shuttle

Durable execution of heddle plans on Temporal. Part of heddle:
weft (algebra) → warp (object store) → shuttle (execution) → drawdown (UI).

## Rules for every change

1. **Keep the wiki current.** Any change that affects how a developer
   installs, builds, tests, uses or understands this repo updates
   `docs/wiki/` in the same commit or PR. That folder is the source of truth;
   `.github/workflows/wiki.yml` mirrors it to
   https://github.com/heedle-ai/shuttle/wiki on every push to `main`. Never edit
   the wiki on GitHub directly: the next sync overwrites it. Code, README and
   wiki must never disagree.
2. **Wiki quality bar.** Every command is copy-pasteable and has been run.
   Every code snippet compiles and is checked by a test; prefer quoting a
   Go `Example` test with an `// Output:` comment. State required versions.
   No placeholders, TODOs or "coming soon" on published pages. Mark design
   that isn't built yet as **Planned**, plainly.
3. **Platforms: Linux (amd64, arm64) and macOS on Apple Silicon only.** Don't
   write Windows instructions, scripts, CI jobs or path handling.
4. **Bare minimum.** Add only what the task needs. No speculative packages,
   dependencies or abstractions.
5. **Tests pass before commit:** `go vet ./...` and `go test -race ./...`.

## Repo-specific

- **`plan` imports the Temporal SDK and nothing else.** `options_test.go`
  enforces this. Don't add LLM clients, metrics or other integrations here.
- **Workflow code is deterministic.** No I/O, clocks, randomness or map
  iteration order inside `plan.Workflow`; all of that goes in activities.
- **Outputs stay small.** Node outputs live in Temporal history; large data
  goes to warp and the node returns its content hash.
- `go.mod` replaces `github.com/heedle-ai/weft` with `../weft` until weft is
  published; clone the repos as siblings.
- `examples/echo` is the wiki's quickstart; its test pins the documented output.
