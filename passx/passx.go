package passx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// Class ///////////////////////////////////////////////////////////////////////
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

// Section /////////////////////////////////////////////////////////////////////
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

// LectureLab //////////////////////////////////////////////////////////////////
// A single atomic lecture/lab combination, where |lecture| must be taken with
// |lab|. |lecture| or |lab|, but not both, may be null.
type LectureLab struct {
	lecture *Section
	lab     *Section
}

// Times ///////////////////////////////////////////////////////////////////////
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

// Time ////////////////////////////////////////////////////////////////////////
type Time struct {
	// M, T, W, R, or F
	day rune

	// Military-time-like integer (see comment in Times).
	time int
}

// Schedule ////////////////////////////////////////////////////////////////////
type Schedule struct {
	sections []*Section

	times map[Time]bool
}

// Adds |section| to |schedule|, and returns whether or not the add succeeded.
// Currently, this function may modify |schedule| and still return false, so
// even upon failing |schedule| must be considered garbage. Perhaps it may be
// useful to first determine if |section| may be added to |schedule|, and THEN
// add it, instead of on-the-go.
func (schedule *Schedule) addSection(section *Section) bool {
	for _, times := range section.Times {
		for _, day := range times.Days {
			for time := times.From; time != times.To; time += 50 {
				newTime := Time{day, time}
				if _, ok := schedule.times[newTime]; ok {
					return false
				}

				schedule.times[newTime] = true
			}
		}
	}

	schedule.sections = append(schedule.sections, section)
	return true
}

func NewSchedule(lectureLabs []LectureLab) (*Schedule, bool) {
	schedule := new(Schedule)
	schedule.times = make(map[Time]bool)

	for _, lectureLab := range lectureLabs {
		if lectureLab.lecture != nil && !schedule.addSection(lectureLab.lecture) ||
			lectureLab.lab != nil && !schedule.addSection(lectureLab.lab) {
			return schedule, false
		}
	}

	return schedule, true
}

////////////////////////////////////////////////////////////////////////////////
func ParseClassesFile(filename string) ([]Class, error) {
	var classes []Class
	var jsonBlob []byte
	var err error

	if jsonBlob, err = ioutil.ReadFile(filename); err != nil {
		return classes, err
	}

	if err = json.Unmarshal(jsonBlob, &classes); err != nil {
		return classes, err
	}

	if err = validateClasses(classes); err != nil {
		return classes, err
	}

	return classes, nil
}

func GenerateSchedules(classes []Class) (schedules []*Schedule) {
	for numClasses := 1; numClasses <= len(classes); numClasses++ {
		schedules = append(schedules,
			generateSchedules(
				classes, make([]*Class, 0), 0, numClasses, numClasses)...)
	}

	return
}

// Fills |subClasses| with |classesLeft| classes (all possible iterations), and
// then generates all possible schedules from the classes in |subClasses|.
func generateSchedules(classes []Class, subClasses []*Class, curClass,
	totalClasses, classesLeft int) (schedules []*Schedule) {

	if curClass < len(classes) && len(classes)-curClass >= classesLeft {
		// Add the current class to |subClasses| and fill the rest of it out.
		// We can safely ignore the return value of generateSchedules here,
		// since it only adds to |schedules| if totalClasses == classesLeft,
		// i.e. we're at the top level.
		subClasses = append(subClasses, &classes[curClass])
		generateSchedules(classes, subClasses, curClass+1, totalClasses,
			classesLeft-1)

		// |subClasses| is full of |classesLeft| classes. It's only "full" if
		// we're at the top level, where |classesLeft| == |totalClasses|.
		if totalClasses == classesLeft {
			schedules = append(schedules, generateSchedules_(subClasses, 0,
				make([]LectureLab, len(subClasses)))...)

			// Also generate schedules of |classesLeft| length that don't
			// include the current class.
			schedules = append(schedules, generateSchedules(classes,
				make([]*Class, 0), curClass+1, totalClasses, classesLeft)...)
		}
	}

	return
}

// Given a slice of *Class, fill parallel slice |lectureLabs|, generating every
// combination of each of |subClasses|'s LectureLabs.
//
// i.e. if |subClasses| == [A, B] where A has LectureLabs [1, 2, 3] and B has
// LectureLabs [foo, bar], then |lectureLabs| will consist of [1, foo], then
// [1, bar], then [2, foo], then [2, bar], then [3, foo], then [3, bar]. For
// each of these, a schedule will be made. If it's valid, it's added to the list
// and returned.
func generateSchedules_(subClasses []*Class, cur int,
	lectureLabs []LectureLab) (schedules []*Schedule) {

	if cur == len(subClasses) {
		if schedule, ok := NewSchedule(lectureLabs); ok {
			schedules = append(schedules, schedule)
		}

		return
	}

	for _, lectureLab := range subClasses[cur].LectureLabs() {
		lectureLabs[cur] = lectureLab
		schedules = append(schedules,
			generateSchedules_(subClasses, cur+1, lectureLabs)...)
	}

	return
}

func validateClasses(classes []Class) error {
	for i := range classes {
		if err := validateClass(&classes[i]); err != nil {
			return err
		}
	}

	return nil
}

// Validates the Class's JSON. Associates each lecture Section with its Class,
// a slice of *Section representing labs (if it's a lecture), and modifies each
// Section's Times from military to military-like (half-hours represented by a
// 50 instead of a 30).
func validateClass(class *Class) error {
	if class.Subject == "" {
		return errors.New("Error validating class: missing \"subject\"")
	}

	if class.Number == "" {
		return errors.New("Error validating class: missing \"number\"")
	}

	// (Probably) modifying each lecture -- get a pointer to each one.
	for i := range class.Lectures {
		lecture := &class.Lectures[i]

		if lecture.Section == "" {
			return fmt.Errorf("Error validating %s %s lecture: "+
				"missing \"section\"", class.Subject, class.Number)
		}

		if err := validateTimes(class.Lectures[i].Times); err != nil {
			return fmt.Errorf("Error validating %s %s lecture %s: %s",
				class.Subject, class.Number, lecture.Section, err)
		}

		for _, referencedLab := range lecture.Labs {
			var found bool

			for j := range class.Labs {
				lab := &class.Labs[j]

				if referencedLab == lab.Section {
					lecture.labPtrs = append(lecture.labPtrs, lab)

					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("Error validating %s %s lecture %s: lab %s "+
					"not found", class.Subject, class.Number, lecture.Section,
					referencedLab)
			}
		}
	}

	// Not modifying lab -- just refer to copies.
	for _, lab := range class.Labs {
		if lab.Section == "" {
			return fmt.Errorf("Error validating %s %s lab: "+
				"missing \"section\"", class.Subject, class.Number)
		}

		if err := validateTimes(lab.Times); err != nil {
			return fmt.Errorf("Error validating %s %s lab %s: %s",
				class.Subject, class.Number, lab.Section, err)
		}
	}

	return nil
}

func validateTimes(times []Times) error {
	if len(times) == 0 {
		return errors.New("missing \"times\"")
	}

	for i := range times {
		// Pass the address because |time| may be modified (to military-like).
		if err := validateTime(&times[i]); err != nil {
			return err
		}
	}

	return nil
}

func validateTime(times *Times) error {
	if times.Days == "" {
		return fmt.Errorf("missing \"days\"")
	}

	daysSet := make(map[rune]bool)
	for _, day := range times.Days {
		if strings.IndexRune("MTWRF", day) == -1 {
			return fmt.Errorf("invalid day \"%c\"", day)
		}

		if ok := daysSet[day]; ok {
			return fmt.Errorf("duplicate day \"%c\"", day)
		}

		daysSet[day] = true
	}

	if times.From < 0 || times.From >= 2400 ||
		(times.From%100 != 0 && times.From%100 != 30) {
		return fmt.Errorf("invalid from-time %d", times.From)
	}

	// Fill |from|, adjust to military-like time.
	if times.From%100 == 30 {
		times.From += 20
	}

	if times.To < 0 || times.To >= 2400 ||
		(times.To%100 != 0 && times.To%100 != 30) || times.To <= times.From {
		return fmt.Errorf("invalid to-time %d", times.To)
	}

	// Fill |to|, adjust to military-like time.
	if times.To%100 == 30 {
		times.To += 20
	}

	return nil
}
