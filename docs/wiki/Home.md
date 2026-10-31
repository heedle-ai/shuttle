**shuttle** runs heddle computations durably on [Temporal](https://temporal.io). A computation is a **plan**: a DAG of named activities, stored as plain data. One generic workflow, `plan.Workflow`, can run any plan. Temporal gives you retries, crash recovery and a full event history for free.

```go
p := plan.Plan{Nodes: []plan.Node{
    {ID: "hello", Activity: plan.EchoName, Args: []string{"hello"}},
    {ID: "world", Activity: plan.EchoName, Args: []string{"world"}},
    {ID: "join",  Activity: plan.EchoName, Args: []string{"joined:"},
        Requires: []plan.NodeID{"hello", "world"}},
}}
```

## What's in the module

| Package | Contents | Imports |
|---|---|---|
| `github.com/heedle-ai/shuttle/plan` | `Plan`, `Node`, `Options`, `Validate`, `Layers`, `Workflow`, `Register`, `Submit`, and the built-in `Echo` activity | Temporal SDK only |
| `github.com/heedle-ai/shuttle` | `ArrowAsActivity`, which registers a [weft](https://github.com/heedle-ai/weft/wiki) arrow as an activity | Temporal SDK, weft |
| `examples/echo` | Runnable quickstart | |

## Where shuttle fits in heddle

| Repo | Role |
|---|---|
| [weft](https://github.com/heedle-ai/weft/wiki) | Composition algebra: what a step *is* |
| [warp](https://github.com/heedle-ai/warp/wiki) | Content-addressed store: what the system *knows* |
| **shuttle** | Durable execution: how the work *ran* (this repo) |
| [drawdown](https://github.com/heedle-ai/drawdown/wiki) | UI for commits, diffs and provenance |

Categorically, a `Plan` is **syntax**: a term built from activity names, with `Requires` edges for composition and independent nodes for parallel products. `plan.Workflow` is one **interpretation** of that syntax as durable Temporal executions.

## Pages

- [Getting Started](Getting-Started): install Go and Temporal, run the example
- [Plans](Plans): the data model, the execution rules, options and limits
- [Writing Activities](Writing-Activities): add your own nodes, including weft arrows
- [Contributing](Contributing)

## Requirements

- Go **1.24** or later
- Temporal CLI, for the local dev server
- Linux (amd64 or arm64), or macOS on Apple Silicon
