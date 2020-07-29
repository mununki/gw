package main

import (
	"fmt"
	// "log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Command struct for commands
type Command struct {
	commands []string
}

func (c *Command) Check() error {
	if len(c.commands) > 1 {
		if strings.HasPrefix(c.commands[1], "-") {
			if c.commands[1] == "-v" {
				// fmt.Println(OptionVersion)
				return fmt.Errorf(OptionVersion)
			} else if c.commands[1] == "-h" {
				// fmt.Println(OptionHelp)
				return fmt.Errorf(OptionHelp)
			}
		} else {
			return nil
		}
	}
	return fmt.Errorf(OptionHelp)
}

// Run func to run commands
func (c *Command) Run() (*Process, error) {
	cmd := exec.Command(c.commands[1], c.commands[2:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = "."
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return &Process{process: cmd, pgid: &pgid}, nil
}
