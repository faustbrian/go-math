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

func TestRootAcceptsExactPerfectPowerAtBitBudget(t *testing.T) {
	limits := gomath.DefaultLimits()
	limits.MaxIntermediateBits = 5

	for _, test := range []struct {
		value int64
		want  string
	}{
		{15, "1"},
		{16, "2"},
		{17, "2"},
	} {
		root, err := integer.New(test.value).Root(context.Background(), 4, limits)
		if err != nil || root.String() != test.want {
			t.Errorf("Root(%d,4) = %s, %v; want %s", test.value, root, err, test.want)
		}
	}
}

func TestRootHugeDegreeRemainsOneAcrossHostIntWidths(t *testing.T) {
	limits := gomath.DefaultLimits()
	limits.MaxIntermediateBits = 5
	limits.MaxRootDegree = ^uint32(0)

	root, err := integer.New(16).Root(context.Background(), limits.MaxRootDegree, limits)
	if err != nil || root.String() != "1" {
		t.Fatalf("Root(16, MaxUint32) = %s, %v; want 1", root, err)
	}
}
