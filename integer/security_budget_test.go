package integer_test

import (
	"context"
	"testing"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/integer"
)

func TestRootBudgetedComparisonsPreserveTruncation(t *testing.T) {
	limits := gomath.DefaultLimits()
	limits.MaxIntermediateBits = 5
	for _, test := range []struct {
		value int64
		want  string
	}{
		{-2, "-1"},
		{0, "0"},
		{15, "2"},
		{20, "2"},
		{27, "3"},
	} {
		root, err := integer.New(test.value).Root(context.Background(), 3, limits)
		if err != nil || root.String() != test.want {
			t.Fatalf("Root(%d,3) = %s, %v; want %s", test.value, root, err, test.want)
		}
	}
}
