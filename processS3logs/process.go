package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wearkinetic/awss3"
	"github.com/wearkinetic/logging"
)

const (
	LOGBUCKET = "kinetic-logging"
)

func main() {

	if len(os.Args)-1 != 4 {
		log.Fatalf("Arguments: <deviceKey> <start 2016-12-14> <end 2017-02-31> <out.csv>")
	}
	// Get these as arguments
	deviceKey := os.Args[1]
	startDate := os.Args[2]
	endDate := os.Args[3]
	currentDir, _ := os.Getwd()
	outfile := currentDir + "/" + os.Args[4]

	dates, dateErr := getDateRange(startDate, endDate)
	if dateErr != nil {
		log.Fatalln(dateErr)
	}

	session := awss3.NewSession(awss3.REGION_US_EAST_1)

	f, err := os.Create(outfile)
	if err != nil {
		log.Fatalf("Could not open file", outfile)
	}
	defer f.Close()

	msg := &logging.Message{}
	for _, date := range dates {
		// fmt.Println(date)
		list, err := session.List(LOGBUCKET, deviceKey+"/"+date)
		if err != nil {
			log.Println("Couldn't read file list")
		}

		for _, key := range list {
			fmt.Printf("Processing: %s\n", key)
			got, getError := session.Get(LOGBUCKET, key)
			if getError != nil {
				log.Fatalf("Get error", getError)
			}

			// now read the file and prepare a CSV
			scanner := bufio.NewScanner(got.Body)
			for scanner.Scan() {
				err = json.Unmarshal(scanner.Bytes(), msg)
				if err != nil {
					continue // skip this line
				}
				fmt.Fprintf(f, "%s,%s,%s\n", msg.Time, msg.Type, msg.Log)
			}
		}
	}
}

// func dateRange(startDate, endDate string) ([]string, error) {
func getDateRange(s, e string) ([]string, error) {
	dates := []string{}
	start, startErr := time.Parse("2006-1-2", s)
	if startErr != nil {
		return dates, startErr
	}
	end, endErr := time.Parse("2006-1-2", e)
	if endErr != nil {
		return dates, endErr
	}
	if start.Unix() > end.Unix() {
		return dates, errors.New("start date is after end date")
	}

	for d := start; d.Unix() <= end.Unix(); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}

	return dates, nil
}
