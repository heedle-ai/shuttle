Every invariant below is guarded by a test or a CI step, and CI checks this page against the code on every push with `bash .github/invariants.sh`. Only the **Guarded by** column counts. Each test named there must exist in exactly one package and is run, and the gate fails if it fails or skips; each `ci:NAME` must be a step in `.github/workflows/ci.yml`. Every `TestLaw_*`, `TestInv_*` and `Fuzz*` test must be listed, and no test file may have build tags, which would hide tests from CI. A mistake isn't fixed until it has a row on this page.

[[/images/invariant-loop.svg|alt=From a mistake to a codified invariant: mistake, named invariant, guard, CI gate]]

| ID | Invariant | Guarded by | Came from |
|---|---|---|---|
| SHUT-1 | Workflow code is deterministic and does no I/O: no clocks, randomness or map iteration (Temporal's analyzer), and `plan` can't import anything that does I/O (SHUT-3) | `ci:workflowcheck`, `TestInv_ImportsOnlyTemporal` | Its first run flagged `Layers` iterating a map |
| SHUT-2 | Layers depend only on the graph, never on node order | `TestInv_LayersIgnoreNodeOrder` | Same finding |
| SHUT-3 | `plan` imports only an allowlist: `context`, `fmt`, `sort`, `strings`, `time` and the Temporal SDK | `TestInv_ImportsOnlyTemporal` | Keep the execution core free of integrations; review 3 found the old ban-list test skipped itself when `go list` failed |
| SHUT-4 | Options layer per field: node over plan over defaults; non-retryable types only accumulate | `TestOptions_LayeredOverDefaults`, `TestOptions_DefaultsUntouched`, `TestWorkflow_PerNodeOptions` | Retry policy must be predictable |
| SHUT-5 | Plans survive JSON unchanged, and invalid plans are refused before they run | `TestPlan_JSONRoundTrip`, `TestValidate_NegativeOptions`, `TestSubmit_RequiresTaskQueue` | Plans cross the Temporal boundary as data |
| SHUT-6 | A weft arrow registered with `ArrowAsActivity` runs as a plan node unchanged | `TestArrowAsPlanActivity`, `TestArrowAsActivity_PreservesBehaviour` | The bridge from weft to Temporal |
| SHUT-7 | The quickstart's documented output is exact | `TestDemoPlan` | Docs must never claim untested behaviour |
| SHUT-8 | Nodes in one layer run at the same time | `TestInv_LayerRunsInParallel` | Review 3: nothing checked the parallelism the docs promise; the test fails if nodes are awaited one by one |
| SHUT-9 | This page is checked strictly against the code: every listed guard exists once and runs without skipping, every `TestLaw_`, `TestInv_` and `Fuzz` test is listed, and no test file has build tags | `ci:invariants` | Review 3: a reviewer satisfied the old gate with a guard named only in prose, a skipped test and a build-tagged test |
