package bridges

import (
	"slices"
	"testing"
)

func TestExpandZshCompdefPattern(t *testing.T) {
	actual := expandZshCompdef("(ruby|[ei]rb)[0-9.]#")
	for _, expected := range []string{"ruby", "erb", "irb"} {
		if !slices.Contains(actual, expected) {
			t.Fatalf("expected %q in %v", expected, actual)
		}
	}
}
