## Data model

```go
type Plan struct {
    Nodes    []Node  `json:"nodes"`
    Defaults Options `json:"defaults,omitzero"`
}

type Node struct {
    ID       NodeID   `json:"id"`        // unique within the plan
    Activity string   `json:"activity"`  // registered Temporal activity name
    Args     []string `json:"args,omitempty"`
    Requires []NodeID `json:"requires,omitempty"` // upstream nodes this one consumes
    Options  Options  `json:"options,omitzero"`
}
```

A plan is plain data, not Go closures. Temporal dispatches activities by registered name and records serializable arguments, and replay requires every activity to be registered in advance. So the graph has to be data. Anything that can produce JSON can produce a plan: a DSL translator, a build description, or an LLM (as long as the result is validated).

## Validation

`Submit` calls `Validate` and refuses to start an invalid plan. `Workflow` validates again, failing with the non-retryable error type `InvalidPlan`, so a plan started some other way is caught too. A plan is invalid if:

- it has no nodes
- a node has an empty `ID`, a duplicate `ID`, or an empty `Activity`
- a node `Requires` an unknown node, or itself
- the graph has a cycle
- any `Timeout`, `Heartbeat` or `MaxAttempts` is negative

`Validate` doesn't check that activities are registered. An unregistered name fails when that node is dispatched.

## Execution

[[/images/layers.svg|alt=Layer 0 runs side by side, a barrier, then layer 1]]


1. `Layers` sorts the nodes into topological layers. Layer 0 is every node with no `Requires`; layer *k* is every node whose dependencies are all in earlier layers. IDs within a layer are sorted, and nodes are always walked in plan order (never a map), so scheduling is deterministic. Temporal's `workflowcheck` analyzer verifies this on every push.
2. All nodes in a layer are dispatched in parallel. The next layer starts only after **every** node in the current layer finishes. `TestInv_LayerRunsInParallel` checks this: three nodes each wait until all three have started, so the plan fails if they run one at a time.
3. Each node receives a `plan.Input`:

   ```go
   type Input struct {
       NodeID   string            `json:"node_id"`
       Args     []string          `json:"args,omitempty"`
       Upstream map[string]string `json:"upstream,omitempty"` // upstream node ID -> its output
   }
   ```

   and returns a `string`. `Upstream` is a map, so it has no order. A node that joins several branches tells them apart by node ID, and any order it imposes is its own choice.
4. The first node failure, after its retries are used up, fails the workflow.
5. `Workflow` returns:

   ```go
   type Result struct {
       Outputs    map[string]string // every node's output
       Leaves     []string          // nodes nothing depends on, sorted
       DurationMs int64
   }
   ```

## Options

Options resolve in three levels, the later overriding the earlier field by field:
`plan.DefaultOptions`, then `Plan.Defaults`, then `Node.Options`. `OptionsFor(node)` gives the result.

| Field | Temporal setting | Default |
|---|---|---|
| `TaskQueue` | Activity task queue, to send a node to specific workers (a GPU host, a build box) | The workflow's own queue |
| `Timeout` | `StartToCloseTimeout`, per attempt | 3 minutes |
| `Heartbeat` | `HeartbeatTimeout` | none |
| `MaxAttempts` | `RetryPolicy.MaximumAttempts`, including the first | 3 |
| `NonRetryable` | `RetryPolicy.NonRetryableErrorTypes` | `InvalidInput`, `InvalidPlan`, `ConfigurationError` |

`NonRetryable` is the exception to overriding: each level **adds** to the list, so an error type marked never-worth-retrying stays that way. Retries back off from 1 second, doubling each time.

## Idempotent submission

Give the workflow ID a value derived from the work itself, and the same work runs only once:

```go
run, err := plan.Submit(ctx, c, p, client.StartWorkflowOptions{
    ID:                       requestHash, // e.g. a hash of the plan and its inputs
    TaskQueue:                "build",
    WorkflowIDConflictPolicy: enums.WORKFLOW_ID_CONFLICT_POLICY_USE_EXISTING,
})
```

Once warp exists, this hash will be the plan's content hash plus the hash of its input tree.

## Limits to design around

- **Outputs must stay small.** Every `Input` and output is stored in Temporal's workflow history, which has size limits. A node that produces something large should store it elsewhere (warp, once built) and return a reference.
- **Layers are barriers.** A slow node holds back every node in later layers, even ones that don't depend on it.
- **One output type.** Nodes exchange strings. Use JSON or a content hash for structured data.
- **Side effects aren't undone.** A failed plan doesn't roll back what earlier nodes did. Make activities idempotent, especially with `MaxAttempts > 1`.
