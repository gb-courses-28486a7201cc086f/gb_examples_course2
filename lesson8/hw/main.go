// Package implements cli tool for
// searching duplicate files in files tree.
//
// Usage:
// go run * -p <path to start directory>
//
// Parameters:
// -h: show help message
// -p: path to start directory, if not exits tool return error
// -remove: optionally: remove second (and other) copy of file
// -y: if 'd' flag provided, deletes files without confirmation
// -w: parallel workers to run, default is 1
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"geekbrains/examples/lesson8/hw/workerpool"
)

// Config contains tool parameters, provided by user
type Config struct {
	searchPath         string
	deleteDuplicates   bool
	notConfirmDeletion bool
	workers            int
}

// RunningConfig is common config object
// populated during initial setup
var RunningConfig Config

func setUp() error {
	flag.StringVar(&RunningConfig.searchPath, "p", ".", "path to start directory, if not exits tool return error")
	flag.BoolVar(&RunningConfig.deleteDuplicates, "remove", false, "optionally: remove second (and other) copy of file")
	flag.BoolVar(&RunningConfig.notConfirmDeletion, "y", false, "if 'd' flag provided, deletes files without confirmation")
	flag.IntVar(&RunningConfig.workers, "w", 1, "parallel workers to run")
	flag.Parse()

	if RunningConfig.workers < 1 {
		return errors.New("count of workers should be greater or equal 1")
	}

	return nil
}

func main() {
	if err := setUp(); err != nil {
		fmt.Printf("invalid parameters: %s\nexit", err.Error())
		os.Exit(1)
	}

	pool, err := workerpool.NewPool(RunningConfig.workers)
	if err != nil {
		msg := fmt.Sprintf("unexpected error: %e", err)
		os.Stderr.WriteString(msg)
		os.Exit(1)
	}

	// cleanup
	pool.Join()
}
