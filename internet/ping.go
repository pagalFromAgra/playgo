package main

import (
	"log"
	"net"
	"time"
)

func main() {
	_, err := net.DialTimeout("tcp", "google.com:80", 5*time.Second)
	if err != nil {
		log.Println("Not connected")
		// handle error
	} else {
		log.Println("connected")
	}
}
