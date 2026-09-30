//go:build !race

package integer_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/integer"
)

func TestHighDegreeRootDoesNotMaterializeOversizedPowers(t *testing.T) {
	limits := gomath.DefaultLimits()
	limits.MaxIntermediateBits = 64
	limits.MaxRootDegree = 1_000_000

	result := testing.Benchmark(func(benchmark *testing.B) {
		for benchmark.Loop() {
			root, err := integer.New(2).Root(context.Background(), limits.MaxRootDegree, limits)
			if err != nil || root.String() != "1" {
				panic("high-degree root did not return one")
			}
		}
	})
	if allocated := result.AllocedBytesPerOp(); allocated > 4096 {
		t.Fatalf("high-degree root allocated %d bytes, want at most 4096", allocated)
	}
}

func TestParserRejectsHostileInputWithoutScalingAllocations(t *testing.T) {
	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 8
	options := integer.ParseOptions{Base: 10, Limits: limits}
	measure := func(input string) int64 {
		result := testing.Benchmark(func(benchmark *testing.B) {
			for benchmark.Loop() {
				if _, err := integer.Parse(input, options); !errors.Is(err, gomath.ErrLimitExceeded) {
					panic("hostile integer did not retain its limit identity")
				}
			}
		})

		return result.AllocedBytesPerOp()
	}
	boundary := measure(strings.Repeat("9", limits.MaxInputDigits+1))
	hostile := measure(strings.Repeat("9", 1<<20))
	if hostile > boundary+1024 {
		t.Fatalf("attacker-sized input allocated %d bytes; near-boundary rejection allocated %d", hostile, boundary)
	}
}
