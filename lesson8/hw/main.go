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
)

// Config contains tool parameters, provided by user
type Config struct {
	searchPath                string
	deleteDuplicates          bool
	disableDeleteConfirmation bool
	workers                   int
}

// RunningConfig is common config object
// populated during initial setup
var RunningConfig Config

// setUp makes input flags parsing and validation.
// parametes stored in global config object
func setUp() error {
	flag.StringVar(&RunningConfig.searchPath, "p", ".", "path to start directory, if not exits tool return error")
	flag.BoolVar(&RunningConfig.deleteDuplicates, "remove", false, "optionally: remove second (and other) copy of file")
	flag.BoolVar(&RunningConfig.disableDeleteConfirmation, "y", false, "if 'd' flag provided, deletes files without confirmation")
	flag.IntVar(&RunningConfig.workers, "w", 1, "parallel workers to run")
	flag.Parse()

	if RunningConfig.workers < 1 {
		return errors.New("-w: count of workers should be greater or equal 1")
	}

	return nil
}

func exitOnErr(err error) {
	msg := fmt.Sprintf("ERROR: %s\n", err.Error())
	os.Stderr.WriteString(msg)
	os.Exit(1)
}

func report(duplicates []DuplicatesDescr) {
	for _, group := range duplicates {
		fmt.Printf("File %s, size %d has duplications:\n", group.Origin, group.Size)
		for _, item := range group.Duplicates {
			fmt.Printf("\t%s\n", item)
		}
	}
	fmt.Println("")
}

func confirmDeleteDuplicates() bool {
	if !RunningConfig.deleteDuplicates {
		return false
	}
	if RunningConfig.disableDeleteConfirmation {
		fmt.Println("Remove duplicates confirmed")
		return true
	}

	var deleteConfirmation string
	fmt.Println("Remove all duplicates (cannot be undone)? type 'yes' to confirm")
	fmt.Scanln(&deleteConfirmation)
	if deleteConfirmation != "yes" {
		fmt.Println("Remove did not confirmed")
		return false
	}
	return true
}

func main() {
	if err := setUp(); err != nil {
		exitOnErr(err)
	}

	files, err := ScanDir(RunningConfig.searchPath)
	if err != nil {
		exitOnErr(err)
	}

	duplicates := CheckDuplicates(files)
	if len(duplicates) > 0 {
		report(duplicates)

		if confirmDeleteDuplicates() {
			err = DeleteDuplicates(duplicates)
			if err != nil {
				exitOnErr(err)
			}
		}
	}

	fmt.Println("Done")
}
