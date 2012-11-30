package main

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"log"

	"github.com/mitchellwrosen/passx-go/passx"
)

func main() {
	jsonBlob, err := ioutil.ReadFile("classes.txt")
	if err != nil {
		panic(err)
	}

	//log.Println(jsonBlob)

	var classes []passx.Class
	err = json.Unmarshal(jsonBlob, &classes)
	if err != nil {
		log.Fatal("Error parsing classes.txt: ", err)
	}
	log.Println(classes)

	err = passx.ValidateClasses(classes)
	if err != nil {
		log.Fatal(err)
	}
}
