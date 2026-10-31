package shuttle_test

import (
	"strings"
	"testing"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"

	"github.com/heedle-ai/shuttle"
	"github.com/heedle-ai/shuttle/plan"
	"github.com/heedle-ai/weft"
)

// TestArrowAsPlanActivity is the wiki's "Write your own activity" example:
// a weft arrow from plan.Input to string, registered by name and used as a
// plan node downstream of the built-in Echo.
func TestArrowAsPlanActivity(t *testing.T) {
	// Any Arrow[plan.Input, string] can be a plan node.
	shout := weft.Pure(func(in plan.Input) string {
		return strings.ToUpper(in.Upstream["greet"])
	})

	p := plan.Plan{Nodes: []plan.Node{
		{ID: "greet", Activity: plan.EchoName, Args: []string{"hello", "heddle"}},
		{ID: "shout", Activity: "Shout", Requires: []plan.NodeID{"greet"}},
	}}

	var s testsuite.WorkflowTestSuite
	env := s.NewTestWorkflowEnvironment()
	plan.Register(env)
	env.RegisterActivityWithOptions(shuttle.ArrowAsActivity(shout), activity.RegisterOptions{Name: "Shout"})

	env.ExecuteWorkflow(plan.WorkflowName, p)
	if err := env.GetWorkflowError(); err != nil {
		t.Fatal(err)
	}
	var res plan.Result
	if err := env.GetWorkflowResult(&res); err != nil {
		t.Fatal(err)
	}
	if got := res.Outputs["shout"]; got != "HELLO HEDDLE" {
		t.Fatalf("shout = %q, want %q", got, "HELLO HEDDLE")
	}
}
