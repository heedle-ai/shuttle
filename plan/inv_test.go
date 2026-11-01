package plan_test

import (
	"math/rand"
	"reflect"
	"testing"

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
