package passx

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
