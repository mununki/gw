package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	for _, arg := range os.Args {
		log.Println(arg)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

MAIN:
	for {
		sig := <-sigs
		log.Println(sig)
		break MAIN
	}
	fmt.Println("runner terminated")
}
