# shuttle

Durable execution for heddle on [Temporal](https://temporal.io).

- **`plan`**: a `Plan` is a serializable DAG whose nodes name registered
  Temporal activities. `plan.Workflow` is the one workflow that runs any plan,
  layer by layer, threading each node's output to its dependents by node ID.
  It imports the Temporal SDK and nothing else.
- **`shuttle.ArrowAsActivity`**: registers a composed `weft.Arrow` as a
  Temporal activity.

Categorically, a `Plan` is syntax (a term built from activity names) and
`plan.Workflow` is one interpretation of it as durable executions.

```go
p := plan.Plan{
    Defaults: plan.Options{TaskQueue: "build"},
    Nodes: []plan.Node{
        {ID: "fetch", Activity: "Fetch", Args: []string{"src"}},
        {ID: "make", Activity: "Assemble", Requires: []plan.NodeID{"fetch"}},
    },
}
plan.Register(w)
run, err := plan.Submit(ctx, c, p, client.StartWorkflowOptions{TaskQueue: "build"})
```

Every node input and output is kept in Temporal's workflow history, so outputs
must stay small. A node that produces something large stores it in
[warp](https://github.com/heddle/warp) and returns its content hash.

```bash
go test -race ./...   # uses Temporal's in-process test suite; no server needed
```

During development, `go.mod` points `github.com/heddle/weft` at `../weft`.
Remove the `replace` line once weft is published.

## Provenance

`plan/` is copied from [`vinodhalaharvi/sibyl`](https://github.com/vinodhalaharvi/sibyl)
at `8b0ccfa`, with import paths rewritten. `ArrowAsActivity` comes from Sibyl's
`agent/lift.go`. Sibyl's agents, convergence and supervisor workflows, LLM
clients, OAuth, Slack channels, metrics, tracing and commands were left behind.
