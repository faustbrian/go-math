package integer_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	gomath "github.com/faustbrian/go-math"
	"github.com/faustbrian/go-math/integer"
)

func TestParseStopsAtFirstConclusiveRejection(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 3
	options := integer.ParseOptions{Base: 10, Limits: limits}
	_, err := integer.Parse("1234x", options)
	assertIntegerInputError(
		t, err, gomath.ErrLimitExceeded,
		"math: resource limit exceeded: input digits",
	)
	_, err = integer.Parse("12x345", options)
	assertIntegerInputError(t, err, gomath.ErrInvalidSyntax, "math: invalid syntax")
	value, err := integer.Parse("123", options)
	if err != nil || value.String() != "123" {
		t.Fatalf("Parse(exact boundary) = %s, %v", value, err)
	}
}

func TestParseLeadingZeroSyntaxPrecedesDigitBudget(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 1
	_, err := integer.Parse("00", integer.ParseOptions{Base: 10, Limits: limits})
	assertIntegerInputError(t, err, gomath.ErrInvalidSyntax, "math: invalid syntax")
}

func TestParseAppliesRawInputBoundBeforeTrimming(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = 3
	options := integer.ParseOptions{
		Base: 10, AllowWhitespace: true, Limits: limits,
	}
	boundary := strings.Repeat(" ", 67) + "123"
	value, err := integer.Parse(boundary, options)
	if err != nil || value.String() != "123" {
		t.Fatalf("Parse(raw boundary) = %s, %v", value, err)
	}
	_, err = integer.Parse(" "+boundary, options)
	assertIntegerInputError(
		t, err, gomath.ErrLimitExceeded,
		"math: resource limit exceeded: integer input bytes",
	)
}

func TestParseRawBoundCalculationSaturates(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	limits.MaxInputDigits = int(^uint(0) >> 1)
	value, err := integer.Parse("1", integer.ParseOptions{Base: 10, Limits: limits})
	if err != nil || value.String() != "1" {
		t.Fatalf("Parse(maximum limit) = %s, %v", value, err)
	}
}

func TestContextPrecedesRandomArgumentsAndLimits(t *testing.T) {
	t.Parallel()

	invalid := gomath.DefaultLimits()
	invalid.MaxInputDigits = 0
	var nilContext context.Context
	_, err := integer.Random(nilContext, nil, integer.New(1), integer.Zero(), invalid)
	if err == nil || err.Error() != "math: invalid argument: nil context" ||
		!errors.Is(err, gomath.ErrInvalidArgument) {
		t.Fatalf("Random(nil context) error = %v", err)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = integer.Random(canceled, nil, integer.New(1), integer.Zero(), invalid)
	if err != context.Canceled {
		t.Fatalf("Random(canceled) error = %v", err)
	}
	deadline, deadlineCancel := context.WithDeadline(context.Background(), time.Unix(1, 0))
	defer deadlineCancel()
	_, err = integer.Random(deadline, nil, integer.New(1), integer.Zero(), invalid)
	if err != context.DeadlineExceeded {
		t.Fatalf("Random(expired deadline) error = %v", err)
	}
}

type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }

type matchingCause struct{ target error }

func (matchingCause) Error() string           { return "matching cause" }
func (cause matchingCause) Is(err error) bool { return err == cause.target }

type structuredCause struct{ details []byte }

func (structuredCause) Error() string { return "structured cause" }

type hostileCause struct{}

func (hostileCause) Error() string { panic("hostile Error called") }

type largeCause struct{}

func (largeCause) Error() string { return strings.Repeat("secret", 1<<18) }

func TestRandomSourceFailurePreservesSafeCauseTraversal(t *testing.T) {
	t.Parallel()

	limits := gomath.DefaultLimits()
	sentinel := errors.New("reader sentinel")
	customTarget := errors.New("custom target")
	tests := []struct {
		name   string
		cause  error
		target error
	}{
		{name: "comparable sentinel", cause: sentinel, target: sentinel},
		{name: "custom match", cause: matchingCause{target: customTarget}, target: customTarget},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := integer.Random(
				context.Background(), failingReader{err: test.cause},
				integer.Zero(), integer.New(2), limits,
			)
			assertRandomSourceError(t, err)
			if !errors.Is(err, test.target) {
				t.Fatalf("errors.Is(%v) = false", test.target)
			}
		})
	}

	cause := structuredCause{details: []byte("caller-owned")}
	_, err := integer.Random(
		context.Background(), failingReader{err: cause},
		integer.Zero(), integer.New(2), limits,
	)
	assertRandomSourceError(t, err)
	var got structuredCause
	if !errors.As(err, &got) || string(got.details) != "caller-owned" {
		t.Fatalf("errors.As() = %#v, want exact structured cause", got)
	}

	multi, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("error type %T does not expose multi-cause traversal", err)
	}
	first := multi.Unwrap()
	if len(first) != 2 || first[0] != gomath.ErrRandomSource || !reflect.DeepEqual(first[1], cause) {
		t.Fatalf("Unwrap() = %#v", first)
	}
	first[0] = errors.New("mutated")
	second := multi.Unwrap()
	if len(second) != 2 || second[0] != gomath.ErrRandomSource {
		t.Fatalf("Unwrap() retained caller mutation: %#v", second)
	}
}

func TestRandomSourceFailureRetainsCauseOwnedInputCategory(t *testing.T) {
	t.Parallel()

	_, err := integer.Random(
		context.Background(), failingReader{err: gomath.ErrInvalidSyntax},
		integer.Zero(), integer.New(2), gomath.DefaultLimits(),
	)
	if err == nil || err.Error() != "math: random source failed" ||
		!errors.Is(err, gomath.ErrRandomSource) || !errors.Is(err, gomath.ErrInvalidSyntax) {
		t.Fatalf("Random(input-category cause) error = %v", err)
	}
	multi, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("error type %T does not expose multi-cause traversal", err)
	}
	causes := multi.Unwrap()
	if len(causes) != 2 || causes[0] != gomath.ErrRandomSource || causes[1] != gomath.ErrInvalidSyntax {
		t.Fatalf("Unwrap() = %#v", causes)
	}
}

func TestRandomSourceFailureNeverFormatsHostileCause(t *testing.T) {
	t.Parallel()

	_, err := integer.Random(
		context.Background(), failingReader{err: hostileCause{}},
		integer.Zero(), integer.New(2), gomath.DefaultLimits(),
	)
	assertRandomSourceError(t, err)
	var got hostileCause
	if !errors.As(err, &got) {
		t.Fatal("errors.As() did not reach hostile cause")
	}
}

func TestRandomSourceFailureNeverIncludesLargeCauseText(t *testing.T) {
	t.Parallel()

	_, err := integer.Random(
		context.Background(), failingReader{err: largeCause{}},
		integer.Zero(), integer.New(2), gomath.DefaultLimits(),
	)
	assertRandomSourceError(t, err)
	var got largeCause
	if !errors.As(err, &got) {
		t.Fatal("errors.As() did not reach large cause")
	}
}

func assertRandomSourceError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("Random() succeeded")
	}
	if got := err.Error(); got != "math: random source failed" {
		t.Fatalf("Error() = %q", got)
	}
	if !errors.Is(err, gomath.ErrRandomSource) {
		t.Fatal("error does not match ErrRandomSource")
	}
	if errors.Is(err, gomath.ErrInvalidSyntax) ||
		errors.Is(err, gomath.ErrLimitExceeded) ||
		errors.Is(err, gomath.ErrInvalidArgument) {
		t.Fatalf("random source error overlaps an input category: %v", err)
	}
}

func assertIntegerInputError(t *testing.T, err, want error, text string) {
	t.Helper()
	if err == nil || !errors.Is(err, want) || err.Error() != text {
		t.Fatalf("error = %v, want %q matching %v", err, text, want)
	}
	for _, category := range []error{
		gomath.ErrInvalidArgument, gomath.ErrInvalidSyntax, gomath.ErrLimitExceeded,
	} {
		if category != want && errors.Is(err, category) {
			t.Fatalf("error %v unexpectedly matches %v", err, category)
		}
	}
}
