Every invariant below is guarded by a test or a CI step, and CI checks this page against the code on every push: each test named here must exist, and every `TestLaw_*`, `TestInv_*` and `Fuzz*` test must be listed here. A mistake isn't fixed until it has a row on this page.

[[/images/invariant-loop.svg|alt=From a mistake to a codified invariant: mistake, named invariant, guard, CI gate]]

| ID | Invariant | Guarded by | Came from |
|---|---|---|---|
| SHUT-1 | Workflow code is deterministic: no clocks, randomness, map iteration or I/O in `plan.Workflow` | `ci:workflowcheck` (Temporal's own analyzer) | Its first run flagged `Layers` iterating a map |
| SHUT-2 | Layers depend only on the graph, never on node order | `TestInv_LayersIgnoreNodeOrder` | Same finding |
| SHUT-3 | `plan` imports the Temporal SDK and nothing else | `TestImportsOnlyTemporal` | Keep the execution core free of integrations |
| SHUT-4 | Options layer per field: node over plan over defaults; non-retryable types only accumulate | `TestOptions_LayeredOverDefaults`, `TestOptions_DefaultsUntouched`, `TestWorkflow_PerNodeOptions` | Retry policy must be predictable |
| SHUT-5 | Plans survive JSON unchanged, and invalid plans are refused before they run | `TestPlan_JSONRoundTrip`, `TestValidate_NegativeOptions`, `TestSubmit_RequiresTaskQueue` | Plans cross the Temporal boundary as data |
| SHUT-6 | A weft arrow registered with `ArrowAsActivity` runs as a plan node unchanged | `TestArrowAsPlanActivity`, `TestArrowAsActivity_PreservesBehaviour` | The bridge from weft to Temporal |
| SHUT-7 | The quickstart's documented output is exact | `TestDemoPlan` | Docs must never claim untested behaviour |
