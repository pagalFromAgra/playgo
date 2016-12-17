package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"strings"

	drs "github.com/ironbay/drs/drs-go"
	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/drs/drs-go/transports/ws"
	"github.com/ironbay/dynamic"
	"github.com/wearkinetic/go/core/domains/device"
)

const (
	API       = "http://api.wearkinetic.com/"
	HOME      = "/home/kinetic/"
	DATAFILES = HOME + "datafiles/"
)

// TODO: Use streams to prevent memory use
func main() {

	key := device.Key()
	conn, err := connection(key)
	if err != nil {
		log.Println("could not create connection")
	}

	files, err := ioutil.ReadDir(DATAFILES)
	if err != nil {
		log.Println("could not read the directory")
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), "out") || file.IsDir() {
			continue
		}
		abs := DATAFILES + file.Name()

		// f, err := os.Open(abs)
		// if err != nil {
		// 	log.Println("could not read")
		// }
		//
		// startTimeMS, endTimeMS, _ := ExtractStartEndTimes(f)
		// filename := strconv.FormatUint(startTimeMS, 10) + "-" + strconv.FormatUint(endTimeMS, 10)
		//
		// log.Printf("filename: %s", filename)

		data, err := ioutil.ReadFile(abs)
		if err != nil {
			log.Println("Could not read the file")
		}
		info, err := os.Stat(abs)
		if err != nil {
			log.Println("Could not read the stat")
		}

		_, err = conn.Call(&drs.Command{
			Action: "data.raw",
			Body: dynamic.Build(
				"created", info.ModTime(),
				"data", base64.StdEncoding.EncodeToString(data),
			),
		})
		if err != nil {
			log.Println("could not connect")
		} else {
			log.Println("Sent: ", abs)
		}
		// os.Rename(abs, abs+".processed")
	}
}

func connection(key string) (*drs.Connection, error) {
	ws := ws.New(dynamic.Build(
		"token", key,
	))
	conn, err := drs.Dial(protocol.JSON, ws, "delta.wearkinetic.com")
	if err != nil {
		return nil, err
	}
	go func() {
		conn.Read()
	}()
	conn.Call(&drs.Command{
		Action:  "delta.subscribe",
		Version: 1,
	})
	return conn, nil
}

/*
// min of 2 uint64
func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

// max of 2 uint64
func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

// Extracts start and end time from an IO Reader
// Typically, feed it os.Open("...filename...")
// and close the the file once this has been used
func ExtractStartEndTimes(body io.Reader) (uint64, uint64, error) {

	// Skim trough the data
	d := new(data.Data)
	stream, err := d.Read(body)
	if err != nil {
		return 0, 0, err
	}
	firstRow := true
	var startTimeMS, endTimeMS uint64

	for rc := range stream {

		// Save file times
		if firstRow {
			startTimeMS = rc.Time
			endTimeMS = rc.Time
			firstRow = false
			continue
		}

		startTimeMS = min(startTimeMS, rc.Time)
		endTimeMS = max(endTimeMS, rc.Time)
	}

	return startTimeMS, endTimeMS, nil

}
*/
