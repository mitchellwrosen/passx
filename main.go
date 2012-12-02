package main

import (
	"fmt"
	"os"

	"github.com/mitchellwrosen/passx-go/passx"
)

func main() {
	if len(os.Args) != 2 {
		panic("Usage: main <filename>")
	}

	classes, err := passx.ParseClassesFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Printf("Classes: %v\n", classes)

	schedules := passx.GenerateSchedules(classes)
	fmt.Printf("Schedules: %v\n", schedules)
}
