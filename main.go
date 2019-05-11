package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	cmd := Command{commands: os.Args}
	p := cmd.Run()
	fmt.Println("** Ctrl-C to exit **")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println(`Go Watcher terminated.`)
		wg.Done()
	}()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if strings.HasSuffix(event.Name, "~") {
					continue
				}
				// log.Println("event: ", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					// log.Println("modified file: ", event.Name)
				}

				p.Kill()
				err = p.process.Wait()

				p = cmd.Run()
				fmt.Println("Trying to run the command...")

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
		fmt.Println("[Error!] Can't watch the root directory.")
	}

	err = filepath.Walk(".", func(walkPath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			// check if dot directory
			if strings.HasPrefix(walkPath, ".") {
				return nil
			}
			if err = watcher.Add(walkPath); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	wg.Wait()

	p.Kill()
}
