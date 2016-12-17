// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package main

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

func main() {
	go watch("/tmp/foo")
	// go watch("/tmp/foo2")
	for {
		errC1 := make(chan error, 1)
		go func() {
			// errC1 <- Do()
			errC1 <- testtimeout()
		}()

		select {
		case err := <-errC1:
			if err != nil {
				log.Println("not success")
			} else {
				log.Println("success")
			}
		case <-time.After(time.Second * 5):
			log.Println("timeout 1")
		}
	}
}

func testtimeout() error {
	time.Sleep(3 * time.Second)
	return nil
}

func watch(filename string) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("removed file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

WATCHAGAIN:
	err = watcher.Add(filename)
	if err != nil {
		log.Println("file doesn't exist yet")
		time.Sleep(2 * time.Second)
		goto WATCHAGAIN
	}

	<-done
}
