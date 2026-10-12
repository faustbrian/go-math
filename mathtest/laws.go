// Package mathtest provides synchronous reusable assertions for numeric laws
// and codecs. Assertions compare values through Equal, not representation
// identity. Callers own the input domain, bounded set size, reporter and
// callback behavior; helpers start no goroutines and add no cancellation or
// panic recovery.
package mathtest

// T is the reporting subset of testing.T used by the helpers. Calls are
// synchronous; the reporter must support Helper and nonfatal Errorf calls.
type T interface {
	// Helper marks the caller as a test helper for failure locations.
	Helper()
	// Errorf records a nonfatal failure; helpers may continue with later inputs.
	Errorf(format string, arguments ...any)
}

// EqualValue supplies the numeric equivalence used by the law assertions. The
// helpers do not separately compare representation, scale or precision.
type EqualValue[V any] interface {
	// Equal reports numeric equivalence between the receiver and another V.
	Equal(V) bool
}

// Operation is a binary operation supplied by the caller. Helpers report each
// returned error as a failure; they do not recover panics, inject cancellation,
// or skip out-of-domain operands. Use a bounded set within the operation's
// domain.
type Operation[V any] func(V, V) (V, error)

// Commutative checks operation(a, b) == operation(b, a) for every ordered pair
// using Equal. The synchronous quadratic traversal includes self-pairs.
func Commutative[V EqualValue[V]](t T, values []V, operation Operation[V]) {
	t.Helper()
	for _, left := range values {
		for _, right := range values {
			forward, err := operation(left, right)
			if err != nil {
				t.Errorf("forward operation failed: %v", err)
				continue
			}
			reverse, err := operation(right, left)
			if err != nil {
				t.Errorf("reverse operation failed: %v", err)
				continue
			}
			if !forward.Equal(reverse) {
				t.Errorf("operation is not commutative for %v and %v", left, right)
			}
		}
	}
}

// Associative checks operation(operation(a, b), c) equals operation(a,
// operation(b, c)) for every ordered triple, using Equal. The synchronous cubic
// traversal is appropriate only for bounded input sets; rounded
// finite-precision operations need not satisfy this exact-arithmetic law.
func Associative[V EqualValue[V]](t T, values []V, operation Operation[V]) {
	t.Helper()
	for _, first := range values {
		for _, second := range values {
			for _, third := range values {
				leftPair, err := operation(first, second)
				if err != nil {
					t.Errorf("left pair failed: %v", err)
					continue
				}
				left, err := operation(leftPair, third)
				if err != nil {
					t.Errorf("left association failed: %v", err)
					continue
				}
				rightPair, err := operation(second, third)
				if err != nil {
					t.Errorf("right pair failed: %v", err)
					continue
				}
				right, err := operation(first, rightPair)
				if err != nil {
					t.Errorf("right association failed: %v", err)
					continue
				}
				if !left.Equal(right) {
					t.Errorf("operation is not associative for %v, %v, and %v", first, second, third)
				}
			}
		}
	}
}

// Identity checks both left and right identity behavior for every supplied
// value, using Equal and reporting callback errors as failures.
func Identity[V EqualValue[V]](t T, values []V, identity V, operation Operation[V]) {
	t.Helper()
	for _, value := range values {
		left, err := operation(identity, value)
		if err != nil {
			t.Errorf("left identity failed: %v", err)
			continue
		}
		right, err := operation(value, identity)
		if err != nil {
			t.Errorf("right identity failed: %v", err)
			continue
		}
		if !left.Equal(value) || !right.Equal(value) {
			t.Errorf("identity does not preserve %v", value)
		}
	}
}

// RoundTrip checks decoded.Equal(original) for each supplied value. It does not
// prove representation preservation: add a separate SameRepresentation,
// precision or policy assertion when the codec contract requires those
// distinctions.
func RoundTrip[V EqualValue[V]](t T, values []V, encode func(V) ([]byte, error), decode func([]byte) (V, error)) {
	t.Helper()
	for _, value := range values {
		encoded, err := encode(value)
		if err != nil {
			t.Errorf("encode failed: %v", err)
			continue
		}
		decoded, err := decode(encoded)
		if err != nil {
			t.Errorf("decode failed: %v", err)
			continue
		}
		if !decoded.Equal(value) {
			t.Errorf("round trip changed %v", value)
		}
	}
}
