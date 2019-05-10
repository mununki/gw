package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

// Command struct for commands
type Command struct {
	commands []string
}

// Run func to run commands
func (c *Command) Run() *Process {
	if len(c.commands) > 1 {
		cmd := exec.Command(c.commands[1], c.commands[2:]...)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = "."
		cmd.Start()

		pgid, err := syscall.Getpgid(cmd.Process.Pid)
		if err == nil {
			return &Process{process: cmd, pgid: &pgid}
		}

		fmt.Println("can't get pgid")
		os.Exit(0)
		return nil
	}

	fmt.Println("need a command")
	os.Exit(0)
	return nil
}

// Process struct for process
type Process struct {
	process *exec.Cmd
	pgid    *int
}

// Kill func to kill process and children process
func (p *Process) Kill() {
	syscall.Kill(-*p.pgid, 15)
}

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
	defer p.Kill()

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

				p.Kill()
				err = p.process.Wait()

				log.Println("current process just terminated")

				p = cmd.Run()

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
