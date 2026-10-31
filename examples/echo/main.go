// Command echo runs a three-node plan on a local Temporal dev server.
//
// It starts a worker in-process, submits the plan, waits for the result and
// prints every node's output:
//
//	temporal server start-dev        # in another terminal
//	go run ./examples/echo
//
// The plan is a diamond: "hello" and "world" have no dependencies and run in
// parallel; "join" requires both and runs after them.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/heedle-ai/shuttle/plan"
)

const taskQueue = "shuttle-echo"

// demoPlan is the plan this example runs. It uses only the built-in Echo
// activity, which joins a node's args with its upstream outputs.
func demoPlan() plan.Plan {
	return plan.Plan{
		Nodes: []plan.Node{
			{ID: "hello", Activity: plan.EchoName, Args: []string{"hello"}},
			{ID: "world", Activity: plan.EchoName, Args: []string{"world"}},
			{ID: "join", Activity: plan.EchoName, Args: []string{"joined:"},
				Requires: []plan.NodeID{"hello", "world"}},
		},
	}
}

func main() {
	hostPort := flag.String("address", client.DefaultHostPort, "Temporal frontend address")
	flag.Parse()

	c, err := client.Dial(client.Options{HostPort: *hostPort})
	if err != nil {
		log.Fatalf("connect to Temporal at %s: %v (is `temporal server start-dev` running?)", *hostPort, err)
	}
	defer c.Close()

	w := worker.New(c, taskQueue, worker.Options{})
	plan.Register(w)
	if err := w.Start(); err != nil {
		log.Fatalf("start worker: %v", err)
	}
	defer w.Stop()

	ctx := context.Background()
	run, err := plan.Submit(ctx, c, demoPlan(), client.StartWorkflowOptions{TaskQueue: taskQueue})
	if err != nil {
		log.Fatalf("submit plan: %v", err)
	}
	fmt.Fprintf(os.Stderr, "workflow %s run %s\n", run.GetID(), run.GetRunID())

	var res plan.Result
	if err := run.Get(ctx, &res); err != nil {
		log.Fatalf("plan failed: %v", err)
	}
	ids := make([]string, 0, len(res.Outputs))
	for id := range res.Outputs {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		fmt.Printf("%-6s %s\n", id, res.Outputs[id])
	}
}
