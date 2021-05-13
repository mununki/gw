package main

import (
	"fmt"
	"io/fs"
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
	cmd := Command{commands: os.Args}
	if err := cmd.Check(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}
	err = filepath.WalkDir(".", func(walkPath string, fi fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return nil
		}
		if fi.IsDir() {
			// check if dot directory
			if strings.HasPrefix(walkPath, ".") {
				return nil
			}
			if strings.HasPrefix(walkPath, "node_modules") {
				return nil
			}
			if err = watcher.Add(walkPath); err != nil {
				log.Println(err)
				return nil
				// return err
			}
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)

	p, err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		wg.Done()
		os.Exit(0)
	}

	tm.Clear()
	tm.MoveCursor(1, 1)
	tm.Println(tm.Color(tm.Bold("** Ctrl-C to exit **"), tm.RED))
	tm.Flush()

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

	wg.Wait()

	close(sigs)
	p.Kill()
}
