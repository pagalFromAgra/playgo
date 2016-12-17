package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

const (
	PATH_BATTERY = "/Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/go/processors/abtest" + "/.battery"
)

func main() {
	cmd := "tail -1 " + PATH_BATTERY + " | cut -f2 -d\"(\" | cut -f1 -d\")\""
	log.Println("cmd = ", cmd)
	battery, _ := exec.Command("sh", "-c", cmd).Output()
	val, _ := strconv.ParseFloat(strings.TrimRight(string(battery), "\n"), 64)
	log.Println("battery = ", val)
}
