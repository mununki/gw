package main

// OptionHelp of -h
const OptionHelp string = `Go Watcher is a command wrapper to run a command every time the filesystem changes.

Usage:

	gw [COMMAND arg1 arg2, ...]

OPTIONS:

	-v	show version
	-h	help

e.g.

	1) gw node server.js

	2) gw go run server.go
`

// OptionVersion of -v
const OptionVersion string = `
Go Watcher v0.2.0
`
