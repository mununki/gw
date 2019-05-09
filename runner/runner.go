package main

import (
	"log"
)

func main() {
	log.Println("I'm running...")

	done := make(chan bool)

	go func() {
		for {
		}
	}()

	<-done
}
