package passx

import (
	"fmt"
)

type Times struct {
	// String of any of M, T, W, R, F
	Days string // From JSON

	// Military-time integers (e.g. 1200, 1330, 1800, etc), but converted to
	// to military-time-like integers during validation (half-hours are
	// represented by a 50 instead of a 30, for ease of iteration).
	From int // From JSON
	To   int // From JSON
}

func (times *Times) String() string {
	return fmt.Sprintf("%s %d - %d", times.Days, times.From, times.To)
}
