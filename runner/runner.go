package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	for _, arg := range os.Args {
		log.Println(arg)
	}

	var wg sync.WaitGroup

	wg.Add(1)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Println(sig)
		wg.Done()
	}()

	wg.Wait()
}
