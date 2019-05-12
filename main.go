package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(".")
	if err != nil {
		fmt.Println("[Error!] Can't watch the root directory.")
		os.Exit(0)
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
		os.Exit(0)
	}

	cmd := Command{commands: os.Args}
	p := cmd.Run()
	fmt.Println("** Ctrl-C to exit **")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

MAIN:
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				break MAIN
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
				break MAIN
			}
			log.Println("error: ", err)

		case signal := <-sigs:
			if signal == os.Interrupt {
				break MAIN
			}
		}
	}

	p.Kill()
}
