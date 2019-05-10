package main

import (
	"log"
	"os"
	"sync"
)

func main() {
	for _, arg := range os.Args {
		log.Println(arg)
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
		}
	}()

	wg.Wait()
}
