package main

import (
	"os/exec"
	"syscall"
)

// Process struct for process
type Process struct {
	process *exec.Cmd
	pgid    *int
}

// Kill func to kill process and children process
func (p *Process) Kill() {
	syscall.Kill(-*p.pgid, 15)
}
