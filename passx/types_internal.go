package passx

// A single atomic lecture/lab combination, where |lecture| must be taken with
// |lab|. |lecture| or |lab|, but not both, may be null.
type lecLab_ struct {
	lec *section
	lab *section
}

// Internal, stripped-down version of the eventual Schedule object
type schedule_ struct {
	lecLabs []*lecLab_

	// A helper-set of Times, for easily determining overlaps when adding a
	// Section.
	times map[time]bool
}

// Add |lectureLab| to |schedule|, and return whether or not the add succeeded.
func (s *schedule_) addLecLab(l *lecLab_) bool {
	if s.canAddSection(l.lec) && s.canAddSection(l.lab) {
		s.lecLabs = append(s.lecLabs, l)
		return true
	}

	return false
}

// Returns whether or not |sch| is valid once |sec| is added. Nil |sec| is OK.
func (sch *schedule_) canAddSection(sec *section) bool {
	if sec != nil {
		for _, day := range sec.Days {
			for t := sec.From; t < sec.To; t += 50 {
				newTime := time{day, t}
				if ok := sch.times[newTime]; ok {
					return false
				}

				sch.times[newTime] = true
			}
		}
	}

	return true
}

func newSchedule_(lls []lecLab_) (*schedule_, bool) {
	sch := new(schedule_)
	sch.times = make(map[time]bool)

	for i := range lls {
		if !sch.addLecLab(&lls[i]) {
			return nil, false
		}
	}

	return sch, true
}

// Helper to Schedule.
type time struct {
	// M, T, W, R, or F
	day rune

	// Military-time-like integer [0-2400)
	time int
}
