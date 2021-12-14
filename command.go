package main

import (
	"errors"
	"io"
	"os/exec"
)

type Shell struct {
	Proc   *exec.Cmd
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

// from "github.com/kylefeng28/go-shell"
func NewShell(command string) (Shell, error) {
	var err error

	shell := Shell{}
	shell.Proc = exec.Command(command)
	if shell.Stdin, err = shell.Proc.StdinPipe(); err != nil {
		return shell, errors.New("could not get a pipe to stdin")
	}
	if shell.Stdout, err = shell.Proc.StdoutPipe(); err != nil {
		return shell, errors.New("could not get a pipe to stdout")
	}
	if shell.Stderr, err = shell.Proc.StderrPipe(); err != nil {
		return shell, errors.New("could not get a pipe to stderr")
	}

	if err = shell.Proc.Start(); err != nil {
		return shell, errors.New("could not start process")
	}

	return shell, nil
}

func (shell Shell) Close() error {
	shell.Stdout.Close()
	shell.Stderr.Close()
	return shell.Proc.Process.Kill()
}
