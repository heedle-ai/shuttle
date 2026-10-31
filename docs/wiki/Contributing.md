## Ground rules

- **`plan` imports the Temporal SDK and nothing else.** `options_test.go` enforces this.
- **Workflow code is deterministic.** No I/O, wall-clock time, randomness or unordered map iteration inside `plan.Workflow`. Use `workflow.Now`, and sort keys before iterating.
- **New behaviour comes with a test** using `testsuite.WorkflowTestSuite`. No test should need a running server.
- **Documented output is tested.** `examples/echo/main_test.go` pins the output shown in [Getting Started](Getting-Started); change both together.
- **Update this wiki in the same change.** The pages live in `docs/wiki/` in this repo, and every push to `main` publishes them here. Don't edit pages on GitHub directly; the next sync overwrites them.

- Supported platforms are Linux (amd64, arm64) and macOS on Apple Silicon.

## Before opening a pull request

```bash
gofmt -l .        # prints nothing
go vet ./...
go test -race ./...
```
