package main

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

// FileDescr contains some file attributes
// used to check if it is unique
type FileDescr struct {
	Name string
	Size int64
}

// ScanDir recursively scan directories and returns list of files
func ScanDir(startDir string) (files map[string]FileDescr, err error) {
	if _, err := ioutil.ReadDir(startDir); err != nil {
		return files, err
	}

	files = make(map[string]FileDescr)
	err = filepath.Walk(startDir, func(currPath string, currFile fs.FileInfo, err error) error {
		if currFile.Mode().IsRegular() {
			files[currPath] = FileDescr{currFile.Name(), currFile.Size()}
		}
		return nil
	})
	return
}

// DuplicatesDescr contains description
// of origin file and list of it's copies
type DuplicatesDescr struct {
	Origin     string
	Size       int64
	Duplicates []string
}

// duplicateKey is unique key of file
type duplicateKey struct {
	name string
	size int64
}

// groupFiles groups files by unique key into map,
// where key is unique key of file and value is list of the file paths
func groupFiles(files map[string]FileDescr) map[duplicateKey][]string {
	groupFiles := make(map[duplicateKey][]string)
	for filePath, fileDescr := range files {
		key := duplicateKey{fileDescr.Name, fileDescr.Size}
		_, exists := groupFiles[key]
		if !exists {
			groupFiles[key] = []string{filePath}
		} else {
			// if already exists - add new name to list of copies
			groupFiles[key] = append(groupFiles[key], filePath)
		}
	}
	return groupFiles
}

// CheckDuplicates search for files which has more than one copy
func CheckDuplicates(files map[string]FileDescr) (duplicates []DuplicatesDescr) {
	for key, group := range groupFiles(files) {
		if len(group) > 1 {
			// files has more than 1 copy
			descr := DuplicatesDescr{
				Origin:     group[0],
				Size:       key.size,
				Duplicates: group[1:],
			}
			duplicates = append(duplicates, descr)
		}
	}
	return duplicates
}

// DeleteDuplicates removes all extra copies of files
func DeleteDuplicates(duplicates []DuplicatesDescr) error {
	for _, group := range duplicates {
		for _, filename := range group.Duplicates {
			err := os.Remove(filename)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
