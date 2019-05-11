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
func (c *Command) Run() *Process {
	if len(c.commands) > 1 {
		if strings.HasPrefix(c.commands[1], "-") {
			if c.commands[1] == "-v" {
				fmt.Println(OptionVersion)
			} else if c.commands[1] == "-h" {
				fmt.Println(OptionHelp)
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
			}

			pgid, err := syscall.Getpgid(cmd.Process.Pid)
			if err == nil {
				return &Process{process: cmd, pgid: &pgid}
			}

			fmt.Println(`
[ERROR!] Can't find a executed process id.
		`)
		}
	} else {
		fmt.Println(OptionHelp)
	}
	os.Exit(0)
	return nil
}
