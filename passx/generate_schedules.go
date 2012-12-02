package passx

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
