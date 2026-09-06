//go:build !race

package rational_test

import (
	"errors"
	"strings"
	"testing"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/rational"
)

func TestParserRejectsHostileInputWithoutScalingAllocations(t *testing.T) {
	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 8
	measure := func(input string) int64 {
		result := testing.Benchmark(func(benchmark *testing.B) {
			for benchmark.Loop() {
				if _, err := rational.Parse(input, limits); !errors.Is(err, gomath.ErrLimitExceeded) {
					panic("hostile rational did not retain its limit identity")
				}
			}
		})

		return result.AllocedBytesPerOp()
	}
	boundary := measure(strings.Repeat("9", limits.MaxInputDigits+1) + "/1")
	hostile := measure(strings.Repeat("9", 1<<20) + "/1")
	if hostile > boundary+1024 {
		t.Fatalf("attacker-sized input allocated %d bytes; near-boundary rejection allocated %d", hostile, boundary)
	}
}
