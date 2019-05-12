package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/fsnotify/fsnotify"
)

func TestExampleCode(t *testing.T) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Error(err)
	}
	defer watcher.Close()

	err = watcher.Add(".")
	if err != nil {
		t.Error(err)
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
MAIN:
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				// break MAIN
				continue
			}
			log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("modified file:", event.Name)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				// break MAIN
				continue
			}
			log.Println("error:", err)

		case signal := <-sigs:
			if signal == os.Interrupt {
				break MAIN
			}
		}
	}
}
