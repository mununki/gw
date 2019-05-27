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

	tm "github.com/buger/goterm"
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

	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Println(tm.Color(tm.Bold("** Ctrl-C to exit **"), tm.RED))
	tm.Flush()

	cmd := Command{commands: os.Args}
	p, err := cmd.Run()
	if err != nil {
		wg.Done()
		os.Exit(0)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					wg.Done()
					close(sigs)
					os.Exit(0)
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

				p, err = cmd.Run()
				if err != nil {
					fmt.Println(err)
					wg.Done()
					close(sigs)
					os.Exit(0)
				}
				tm.Clear()
				tm.MoveCursor(1, 1)
				tm.Println(tm.Color(tm.Bold("Trying to run the command..."), tm.GREEN))
				tm.Flush()

			case err, ok := <-watcher.Errors:
				if !ok {
					wg.Done()
					close(sigs)
					os.Exit(0)
				}
				log.Println("error: ", err)

			}
		}
	}()

	go func() {
		<-sigs
		wg.Done()
	}()

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

	wg.Wait()

	close(sigs)
	p.Kill()
}
