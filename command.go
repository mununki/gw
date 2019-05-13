package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Command struct for commands
type Command struct {
	commands []string
}

// Run func to run commands
func (c *Command) Run() (*Process, error) {
	if len(c.commands) > 1 {
		if strings.HasPrefix(c.commands[1], "-") {
			if c.commands[1] == "-v" {
				fmt.Println(OptionVersion)
				return nil, fmt.Errorf("version check")
			} else if c.commands[1] == "-h" {
				fmt.Println(OptionHelp)
				return nil, fmt.Errorf("print help")
			}
		} else {
			cmd := exec.Command(c.commands[1], c.commands[2:]...)
			cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Dir = "."
			err := cmd.Start()
			if err != nil {
				log.Fatal(err)
				return nil, fmt.Errorf("can't run the command")
			}

			pgid, err := syscall.Getpgid(cmd.Process.Pid)
			if err == nil {
				return &Process{process: cmd, pgid: &pgid}, nil
			}

			return nil, fmt.Errorf("error to get a excuted process id")
		}
	} else {
		fmt.Println(OptionHelp)
		return nil, fmt.Errorf("no command to run")
	}
	return nil, fmt.Errorf("unexpected error")
}
