package passx

import "fmt"

// Input types /////////////////////////////////////////////////////////////////
type class struct {
	Subject string
	Number  string
	Units   int
	Groups  []group
}

func (c *class) String() string {
	return fmt.Sprintf("%s %s", c.Subject, c.Number)
}

func (c *class) lecLabs() (lls []lecLab_) {
	for _, group := range c.Groups {
		if group.Lecs == nil {
			labs := *group.Labs
			for i := range labs {
				lls = append(lls, lecLab_{nil, &labs[i]})
			}
		} else if group.Labs == nil {
			lecs := *group.Lecs
			for i := range lecs {
				lls = append(lls, lecLab_{&lecs[i], nil})
			}
		} else {
			lecs := *group.Lecs
			labs := *group.Labs
			for i := range lecs {
				for j := range labs {
					lls = append(lls, lecLab_{&lecs[i], &labs[j]})
				}
			}
		}
	}

	return
}

// A group of associated lectures and labs, where any lecture combined with any
// lab is a valid LectureLab for this class.
type group struct {
	Lecs *[]section
	Labs *[]section
}

//func (schedule *Schedule) String() string {
//var buffer bytes.Buffer

//for _, section := range schedule.Sections {
//buffer.WriteString(section.String())
//buffer.WriteByte('\n')
//}

//return buffer.String()
//}

type section struct {
	Section string

	// String of any of M, T, W, R, F
	Days string

	// Military-time integers (e.g. 1200, 1330, 1800, etc), but converted to
	// to military-time-like integers during validation (half-hours are
	// represented by a 50 instead of a 30, for ease of iteration).
	From int
	To   int

	// The Class this Section belongs to (filled in during validation).
	// This is used to fill a schedule, distinct from a schedule_, which is a
	// lighter-weight version.
	class *class
}

func (sec *section) String() string {
	from, to := sec.From, sec.To
	if from%100 == 50 {
		from -= 20
	}
	if to%100 == 50 {
		to -= 20
	}

	return fmt.Sprintf("%s-%s %s %d-%d", sec.class, sec.Section, sec.Days, from,
		to)
}

// Output types ////////////////////////////////////////////////////////////////
type schedule struct {
	Classes []*scheduleClass
}

// Creates a schedule, given a schedule_
func newSchedule(sch_ *schedule_) *schedule {
	sch := new(schedule)
	sch.Classes = make([]*scheduleClass, len(sch_.lecLabs))

	for i, ll := range sch_.lecLabs {
		var classPtr *class
		if ll.lec == nil {
			classPtr = ll.lab.class
		} else {
			classPtr = ll.lec.class
		}

		sch.Classes[i] = &scheduleClass{
			classPtr.Subject,
			classPtr.Number,
			classPtr.Units,
			ll.lec,
			ll.lab}
	}

	return sch
}

type scheduleClass struct {
	Subject string
	Number  string
	Units   int
	Lec     *section
	Lab     *section
}
