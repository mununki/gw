package main

import (
	// "fmt"
	"log"
	"os"
	"os/exec"
	"sync"

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

	var wg sync.WaitGroup

	p := runTest()

	defer p.Process.Kill()

	wg.Add(1)

	go func() {
		defer wg.Done()

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

	wg.Wait()
}
