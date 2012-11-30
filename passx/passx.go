package passx

import (
	"errors"
	"fmt"
	"strings"
)

type Class struct {
	Subject  string
	Number   string
	Lectures []Lecture
	Labs     []Lab
}

type Lecture struct {
	Section string
	Times   []Time
	Labs    []string
}

type Lab struct {
	Section string
	Times   []Time
}

type Time struct {
	Day  string
	Time string
}

func ValidateClasses(classes []Class) error {
	for _, class := range classes {
		err := validateClass(class)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateClass(class Class) error {
	if class.Subject == "" {
		return errors.New("Error validating class: missing \"subject\"")
	}

	if class.Number == "" {
		return errors.New("Error validating class: missing \"number\"")
	}

	// Validate labs first, to keep a list of all labs (referenced by lectures)
	labSections := make([]string, len(class.Labs))
	for index, lab := range class.Labs {
		if lab.Section == "" {
			return fmt.Errorf("Error validating %s %s lab: "+
				"missing \"section\"", class.Subject, class.Number)
		}

		err := validateTimes(lab.Times)
		if err != nil {
			return fmt.Errorf("Error validating %s %s lab %s: %s",
				class.Subject, class.Number, lab.Section, err)
		}

		labSections[index] = lab.Section
	}

	for _, lecture := range class.Lectures {
		if lecture.Section == "" {
			return fmt.Errorf("Error validating %s %s lecture: "+
				"missing \"section\"", class.Subject, class.Number)
		}

		err := validateTimes(lecture.Times)
		if err != nil {
			return fmt.Errorf("Error validating %s %s lecture %s: %s",
				class.Subject, class.Number, lecture.Section, err)
		}

		for _, lab := range lecture.Labs {
			var found bool
			for _, labSection := range labSections {
				if lab == labSection {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("Error validating %s %s lecture %s: lab %s "+
					"not found", class.Subject, class.Number, lecture.Section,
					lab)
			}
		}
	}

	return nil
}

func validateTimes(times []Time) error {
	for _, time := range times {
		err := validateTime(time)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateTime(time Time) error {
	if time.Day == "" {
		return fmt.Errorf("missing \"day\"")
	}
	for _, day := range time.Day {
		if strings.IndexRune("MTWRF", day) == -1 {
			return fmt.Errorf("invalid day \"%c\"", day)
		}
	}

	return nil
}
