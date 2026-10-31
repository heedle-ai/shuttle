package main

import (
	"testing"

	"go.temporal.io/sdk/testsuite"

	"github.com/heddle/shuttle/plan"
)

// TestDemoPlan runs the example's plan in Temporal's in-process test
// environment, so the output documented in the wiki is checked by CI.
func TestDemoPlan(t *testing.T) {
	var s testsuite.WorkflowTestSuite
	env := s.NewTestWorkflowEnvironment()
	plan.Register(env)

	env.ExecuteWorkflow(plan.WorkflowName, demoPlan())
	if !env.IsWorkflowCompleted() {
		t.Fatal("workflow did not complete")
	}
	if err := env.GetWorkflowError(); err != nil {
		t.Fatal(err)
	}
	var res plan.Result
	if err := env.GetWorkflowResult(&res); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"hello": "hello", "world": "world", "join": "joined: | hello world"}
	for id, w := range want {
		if got := res.Outputs[id]; got != w {
			t.Errorf("%s = %q, want %q", id, got, w)
		}
	}
	if len(res.Leaves) != 1 || res.Leaves[0] != "join" {
		t.Errorf("leaves = %v, want [join]", res.Leaves)
	}
}
