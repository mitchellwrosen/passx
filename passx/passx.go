package passx

import (
	"encoding/json"
	"io/ioutil"
)

func generateSchedulesJSON(classesJSON []byte, minUnits,
	maxUnits int) ([]byte, error) {

	var classes []class

	err := json.Unmarshal(classesJSON, &classes)
	if err != nil {
		return nil, err
	}

	setSectionClassPointers(classes)

	schedules := generateSchedules(classes, minUnits, maxUnits)
	schedulesJSON, err := json.Marshal(schedules)
	if err != nil {
		return nil, err
	}

	return schedulesJSON, nil
}

// For testing.
func generateSchedulesJSONFromFile(filename string, minUnits,
	maxUnits int) ([]byte, error) {

	jsonBlob, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return generateSchedulesJSON(jsonBlob, minUnits, maxUnits)
}

func setSectionClassPointers(classes []class) {
	for classIndex := range classes {
		for _, group := range classes[classIndex].Groups {
			if group.Lecs != nil {
				for lecIndex := range *group.Lecs {
					(*group.Lecs)[lecIndex].class = &classes[classIndex]
				}
			}

			if group.Labs != nil {
				for labIndex := range *group.Labs {
					(*group.Labs)[labIndex].class = &classes[classIndex]
				}
			}
		}
	}
}

func generateSchedules(classes []class, minUnits, maxUnits int) (schedules []*schedule) {
	for numClasses := 1; numClasses <= len(classes); numClasses++ {
		schedules = append(schedules,
			generateSchedules_(
				classes, make([]*class, numClasses), 0, 0, minUnits, maxUnits)...)
	}

	return
}

// Fills |subClasses| with Classes from |classes| (all possible iterations), and
// then generates all possible schedules from the classes in |subClasses|.
func generateSchedules_(classes []class, subClasses []*class, curClass,
	curSubClass, minUnits, maxUnits int) (schedules []*schedule) {

	if len(subClasses)-curSubClass <= len(classes)-curClass {
		// Add the current class to |subClasses| and fill the rest of it out.
		// We can safely ignore the return value of generateSchedules here,
		// since it only adds to |schedules| if |topLevel|.
		subClasses[curSubClass] = &classes[curClass]
		if curSubClass == len(subClasses)-1 {
			schedules = append(schedules, generateSchedules__(subClasses, 0,
				make([]lecLab_, len(subClasses)), minUnits, maxUnits)...)
		} else {
			schedules = append(schedules, generateSchedules_(classes,
				subClasses, curClass+1, curSubClass+1, minUnits, maxUnits)...)
		}

		schedules = append(schedules, generateSchedules_(classes, subClasses,
			curClass+1, curSubClass, minUnits, maxUnits)...)
	}

	return
}

// Given a slice of *class, fill parallel slice |lecLabs|, generating every
// combination of each of |subClasses|'s lecLabs_.
//
// i.e. if |subClasses| == [A, B] where A has lecLabs_ [1, 2, 3] and B has
// lecLabs_ [foo, bar], then |lecLabs| will consist of [1, foo], then [1, bar], 
// then [2, foo], then [2, bar], then [3, foo], then [3, bar]. For each of 
// these, a schedule_ will be made. If it's valid, an equivalent schedule is
// made and added to the returned slice.
func generateSchedules__(subClasses []*class, cur int, lecLabs []lecLab_,
	minUnits, maxUnits int) (schedules []*schedule) {

	if cur == len(subClasses) {
		var units int
		for _, class := range subClasses {
			units += class.Units
		}

		if units >= minUnits && units <= maxUnits {
			if sch_, ok := newSchedule_(lecLabs); ok {
				schedules = append(schedules, newSchedule(sch_))
			}
		}

		return
	}

	for _, lecLab := range subClasses[cur].lecLabs() {
		lecLabs[cur] = lecLab
		schedules = append(schedules,
			generateSchedules__(subClasses, cur+1, lecLabs, minUnits,
				maxUnits)...)
	}

	return
}
