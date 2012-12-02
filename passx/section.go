package passx

import (
	"fmt"
)

type Section struct {
	Section string   // From JSON
	Times   []Times  // From JSON
	Labs    []string // From JSON

	// The Class this Lecture belongs to.
	class *Class

	// Pointers to actual lab Section (filled in inside of validateClass)
	// This field will only be filled in for lectures.
	labPtrs []*Section
}

func (section *Section) String() string {
	var timesStr string
	for _, time := range section.Times {
		timesStr += time.String()
	}

	return fmt.Sprintf("%s %s %s", section.class, section.Section, timesStr)
}
