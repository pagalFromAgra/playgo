package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	drs "github.com/ironbay/drs/drs-go"
	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/drs/drs-go/transports/ws"
	"github.com/ironbay/dynamic"
	"github.com/wearkinetic/go/core/domains/raw"
)

const (
	API       = "http://api.wearkinetic.com/"
	HOME      = "/Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/go/processors/abtest/renamefile"
	DATAFILES = HOME + "/datafiles/"
)

// TODO: Use streams to prevent memory use
func main() {

	files, err := ioutil.ReadDir(DATAFILES)
	if err != nil {
		log.Println("could not read the directory")
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), "out") || file.IsDir() {
			continue
		}
		abs := DATAFILES + file.Name()

		f, err := os.Open(abs)
		if err != nil {
			log.Println("could not read")
		}

		startTimeMS, endTimeMS, _ := raw.ExtractStartEndTimes(f)
		filename := strconv.FormatUint(startTimeMS, 10) + "-" + strconv.FormatUint(endTimeMS, 10) + ".out"

		log.Printf("filename: %s", filename)

		// data, err := ioutil.ReadFile(abs)
		// if err != nil {
		// 	log.Println("Could not read the file")
		// }
		// info, err := os.Stat(abs)
		// if err != nil {
		// 	log.Println("Could not read the stat")
		// }

		// _, err = conn.Call(&drs.Command{
		// 	Action: "data.raw",
		// 	Body: dynamic.Build(
		// 		"created", info.ModTime(),
		// 		"data", base64.StdEncoding.EncodeToString(data),

		// 	),
		// })
		// if err != nil {
		// 	log.Println("could not connect")
		// } else {
		// 	log.Println("Sent: ", abs)
		// }
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
