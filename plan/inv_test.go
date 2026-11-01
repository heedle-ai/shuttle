package plan_test

import (
	"context"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"

	"github.com/heedle-ai/shuttle/plan"
)

// TestInv_LayersIgnoreNodeOrder: a plan's layers depend only on its graph,
// never on the order nodes are listed in, so execution order is the same
// for every equal graph. (Layers walks nodes in plan order, never a map;
// Temporal's workflowcheck enforces determinism in CI.)
func TestInv_LayersIgnoreNodeOrder(t *testing.T) {
	nodes := []plan.Node{
		{ID: "a", Activity: plan.EchoName}, {ID: "b", Activity: plan.EchoName},
		{ID: "c", Activity: plan.EchoName, Requires: []plan.NodeID{"a"}},
		{ID: "d", Activity: plan.EchoName, Requires: []plan.NodeID{"a", "b"}},
		{ID: "e", Activity: plan.EchoName, Requires: []plan.NodeID{"c", "d"}},
		{ID: "f", Activity: plan.EchoName},
	}
	want, err := plan.Plan{Nodes: nodes}.Layers()
	if err != nil {
		t.Fatal(err)
	}
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 200; i++ {
		shuffled := append([]plan.Node(nil), nodes...)
		r.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
		got, err := plan.Plan{Nodes: shuffled}.Layers()
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("order %v gave layers %v, want %v", shuffled, got, want)
		}
	}
}

// TestInv_LayerRunsInParallel: nodes in the same layer run at the same
// time. Each node waits until all three have started; if the workflow ran
// them one after another, the first would time out and the plan would fail.
func TestInv_LayerRunsInParallel(t *testing.T) {
	var started sync.WaitGroup
	started.Add(3)
	all := make(chan struct{})
	go func() { started.Wait(); close(all) }()
	barrier := func(ctx context.Context, in plan.Input) (string, error) {
		started.Done()
		select {
		case <-all:
			return in.NodeID, nil
		case <-time.After(5 * time.Second):
			return "", temporal.NewNonRetryableApplicationError("node ran alone: layer not parallel", "Sequential", nil)
		}
	}
	p := plan.Plan{Nodes: []plan.Node{
		{ID: "a", Activity: "Barrier"}, {ID: "b", Activity: "Barrier"}, {ID: "c", Activity: "Barrier"},
	}, Defaults: plan.Options{MaxAttempts: 1}}

	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestWorkflowEnvironment()
	plan.Register(env)
	env.RegisterActivityWithOptions(barrier, activity.RegisterOptions{Name: "Barrier"})
	env.ExecuteWorkflow(plan.WorkflowName, p)
	if err := env.GetWorkflowError(); err != nil {
		t.Fatalf("same-layer nodes didn't run in parallel: %v", err)
	}
}
