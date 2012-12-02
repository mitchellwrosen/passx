package passx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

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
