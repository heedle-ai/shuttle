// Package shuttle runs heddle computations durably on Temporal.
//
// The execution graph lives in package plan: a serializable DAG whose
// nodes name registered activities, interpreted by plan.Workflow. This
// root package holds the bridge from weft's arrow algebra to Temporal.
package shuttle

import (
	"context"

	"github.com/heddle/weft"
)

// ArrowAsActivity adapts a weft.Arrow into the function shape Temporal
// expects for an activity registration.
//
// A weft.Arrow[A, B] is structurally `func(ctx, A) (B, error)`, already
// the right shape, but Go's type system treats the two as nominally
// distinct. This thin wrapper makes the conversion explicit and lets
// Temporal's reflection-based RegisterActivity see a plain function.
func ArrowAsActivity[A, B any](a weft.Arrow[A, B]) func(context.Context, A) (B, error) {
	return func(ctx context.Context, in A) (B, error) {
		return a(ctx, in)
	}
}
