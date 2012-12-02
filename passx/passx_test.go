package passx

import "testing"

type TWrapper struct {
	*testing.T
}

func (t *TWrapper) expectEq(a, b interface{}) {
	if a != b {
		t.Errorf("expectEq(%v, %v)", a, b)
	}
}

func TestParseClassesFile(te *testing.T) {
	var err error
	t := TWrapper{te}

	classes, err := ParseClassesFile("files/good2.txt")

	t.expectEq(err, nil)
	t.expectEq(1, len(classes))

	class := classes[0]
	t.expectEq("CSC", class.Subject)
	t.expectEq("101", class.Number)
	t.expectEq(1, len(class.Lectures))

	t.expectEq("01", class.Lectures[0].Section)
	t.expectEq(1, len(class.Lectures[0].Times))
	t.expectEq("M", class.Lectures[0].Times[0].Days)
	t.expectEq(1200, class.Lectures[0].Times[0].From)
	t.expectEq(1350, class.Lectures[0].Times[0].To)

	t.expectEq(0, len(class.Labs))
}
