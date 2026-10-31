package shuttle

import (
	"context"
	"strings"
	"testing"

	"github.com/heedle-ai/weft"
)

func TestArrowAsActivity_PreservesBehaviour(t *testing.T) {
	a := weft.Compose(weft.Pure(strings.TrimSpace), weft.Pure(strings.ToUpper))
	act := ArrowAsActivity(a)
	got, err := act(context.Background(), "  hi ")
	if err != nil || got != "HI" {
		t.Fatalf("got %q, %v; want %q, nil", got, err, "HI")
	}
}
