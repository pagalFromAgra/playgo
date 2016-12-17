package main

import (
	"log"

	"github.com/wearkinetic/awss3"
)

const (
	LOGBUCKET = "kinetic-logging"
)

func main() {
	session := awss3.NewSession(awss3.REGION_US_EAST_1)

	// Get these as arguments
	deviceKey := "oU04IJOaetfDLfco"
	date := "2016-12-14"

	list, err := session.List(LOGBUCKET, deviceKey+"/"+date)
	if err != nil {
		log.Println("Couldn't read")
	}

	for file, _ := range list {
		log.Println(file)
	}

}
