package main

import (
	// "fmt"
	"log"
	"os"
	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func runTest() *exec.Cmd {
	p := exec.Command("./test-run")
	p.Stdout = os.Stdout
	p.Stderr = os.Stderr
	p.Dir = "."
	p.Start()

	return p
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	p := runTest()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file: ", event.Name)
				}
				p.Process.Kill()

				p = runTest()
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error: ", err)
			}
		}
	}()

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}

	<-done
}
