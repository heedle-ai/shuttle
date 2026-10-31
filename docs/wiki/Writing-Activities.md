A plan node can be any registered Temporal activity with this shape:

```go
func(ctx context.Context, in plan.Input) (string, error)
```

## From a weft arrow

`weft.Arrow[plan.Input, string]` already has that shape. `shuttle.ArrowAsActivity` makes it registrable:

```go
shout := weft.Pure(func(in plan.Input) string {
    return strings.ToUpper(in.Upstream["greet"])
})

p := plan.Plan{Nodes: []plan.Node{
    {ID: "greet", Activity: plan.EchoName, Args: []string{"hello", "heddle"}},
    {ID: "shout", Activity: "Shout", Requires: []plan.NodeID{"greet"}},
}}

plan.Register(w) // PlanWorkflow and Echo
w.RegisterActivityWithOptions(shuttle.ArrowAsActivity(shout), activity.RegisterOptions{Name: "Shout"})
```

`shout` produces `HELLO HEDDLE`. This code is `TestArrowAsPlanActivity` in `plan_activity_test.go`.

Because the activity is an arrow, you can build it with weft combinators: `Pipe3` for parse → call → render, or `Apply` for timeouts and taps. Prefer Temporal's retry policy (the node's `MaxAttempts`) over `weft.WithRetry` for activity-level retries, so each attempt shows up in the history.

## Rules for activities

- **Register by name.** The name in `Node.Activity` must match `activity.RegisterOptions{Name: …}` on some worker polling that node's task queue.
- **Return small strings.** See [Plans → Limits](Plans#limits-to-design-around).
- **Fail on purpose with typed errors.** If retrying can't help (bad input, missing configuration), return `temporal.NewNonRetryableApplicationError(msg, "InvalidInput", err)`, or use a type listed in the node's `NonRetryable`.
- **Be idempotent.** An activity can run more than once: retries, worker crashes, timeouts.
- **Heartbeat long work.** Call `activity.RecordHeartbeat(ctx)` and set `Heartbeat` on the node, so a hung attempt is detected before `Timeout`.

## The built-in `Echo`

`plan.Echo` (registered as `"Echo"`) needs no credentials, which makes it useful for trying out plan shapes:

| Args | Upstream | Output |
|---|---|---|
| set | none | args joined with spaces |
| none | set | upstream outputs, sorted by node ID, joined with spaces |
| set | set | `args \| upstream` |
