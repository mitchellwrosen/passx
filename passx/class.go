package passx

import "fmt"

type Class struct {
	Subject  string    // From JSON
	Number   string    // From JSON
	Lectures []Section // From JSON
	Labs     []Section // From JSON
}

func (class *Class) String() string {
	return fmt.Sprintf("Foo") // TODO(mwrosen) why isn't this working?
	//return fmt.Sprintf("%s %s", class.Subject, class.Number)
}

func (class *Class) LectureLabs() (lectureLabs []LectureLab) {
	// Lecture-only class.
	if len(class.Labs) == 0 {
		for i := range class.Lectures {
			lectureLabs = append(lectureLabs, LectureLab{&class.Lectures[i],
				nil})
		}

		return
	}

	// Lab-only class.
	if len(class.Lectures) == 0 {
		for i := range class.Labs {
			lectureLabs = append(lectureLabs, LectureLab{nil, &class.Labs[i]})
		}

		return
	}

	// Lecture/lab class.
	for i, lecture := range class.Lectures {
		for _, labPtr := range lecture.labPtrs {
			lectureLabs = append(lectureLabs, LectureLab{&class.Lectures[i],
				labPtr})
		}
	}

	return
}
