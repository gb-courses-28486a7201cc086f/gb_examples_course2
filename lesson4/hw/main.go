package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const stopTimeout = 1 * time.Second

func main() {
	log.Printf("Started pid=%d\n", os.Getpid())
	// setup signal handler
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM)

	workUntilTerm(ctx, superWork10Sec)

	log.Println("Finished")
}

func workUntilTerm(ctx context.Context, task func()) {
	workDone := make(chan struct{})
	for {
		// do work
		go func() {
			task()
			workDone <- struct{}{}
		}()
		// wait done or cancel
		select {
		case <-ctx.Done():
			log.Println("break on signal!")
			// after cancel signal
			// wait {stopTimeout} to work done or exit
			timeout, _ := context.WithTimeout(context.Background(), stopTimeout)
			select {
			case <-timeout.Done():
				log.Println("timeout expires")
			case <-workDone:
				log.Println("work done")
			}
			return
		case <-workDone:
		}
	}
}

// simulates long work
func superWork10Sec() {
	log.Println("I'm working!")
	time.Sleep(10 * time.Second)
}
