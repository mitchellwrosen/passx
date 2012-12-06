package passx

import (
	"fmt"
	"io/ioutil"
	"mytesting"
	"testing"
)

func TestSimple(te *testing.T) {
	t := mytesting.T{te}

	for i := 0; i < 1; i++ {
		schedules, err := generateSchedulesJSONFromFile(
			fmt.Sprintf("tests/simple%d.in", i), 4, 20)
		t.ExpectEq(err, nil)

		expected, err := ioutil.ReadFile(fmt.Sprintf("tests/simple%d.out", i))
		t.ExpectEq(err, nil)

		t.ExpectEq(string(schedules), string(expected))
	}
}
