package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// simulates ping tool
func ping(hostname string) bool {
	response := rand.Intn(100)
	time.Sleep(time.Duration(response * 1000))
	if response == 100 {
		return false
	}
	return true
}

func main() {
	// to simulate race condition,
	// try to store ping results in global map

	pingResults := make(map[string]bool, 1000)

	wg := sync.WaitGroup{}
	go func() {
		for i := 0; i < 500; i++ {
			hostname := fmt.Sprintf("host%4.d", i)
			pingResults[hostname] = ping(hostname)
		}
	}()
	go func() {
		for i := 500; i < 1000; i++ {
			hostname := fmt.Sprintf("host%4.d", i)
			pingResults[hostname] = ping(hostname)
		}
	}()

	wg.Wait()
}
