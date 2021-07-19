package main

import (
	"log"
	"sync"
	"tcsrdr/internal/tcsrdrserver"
)

func main() {

	wg := &sync.WaitGroup{}

	err := tcsrdrserver.Init(false)
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go tcsrdrserver.HttpTickerServer(wg)
	wg.Wait()

	log.Println("Utility has been stopped")
}
