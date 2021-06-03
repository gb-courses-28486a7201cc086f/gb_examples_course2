// Package main simulates WaitGroup usage
package main

import (
	"log"
	"sync"
)

const workers = 10

func main() {
	wg := sync.WaitGroup{}

	log.Println("Started")

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			log.Printf("Goroutine %d working...", workerId)
		}(i)
	}
	wg.Wait()

	log.Println("Finished")
}
