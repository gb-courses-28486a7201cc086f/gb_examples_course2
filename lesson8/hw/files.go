package main

import (
	"io/fs"
	"io/ioutil"
	"path"
	"path/filepath"

	"geekbrains/examples/lesson8/hw/workerpool"
)

type FileDescr struct {
	Name string
	Size int64
}

func ScanDir(startDir string) (files map[string]FileDescr, err error) {
	files = make(map[string]FileDescr)
	err = filepath.Walk(startDir, func(currDir string, currFile fs.FileInfo, err error) error {
		if currFile.Mode().IsRegular() {
			key := path.Join(currDir, currFile.Name())
			files[key] = FileDescr{currFile.Name(), currFile.Size()}
		}
		return nil
	})
	return
}

type IndexDirJob struct {
	dirname   string
	childDirs []string
	files     map[string]FileDescr
	err       error
}

func (idj *IndexDirJob) Run() {
	files, err := ioutil.ReadDir(idj.dirname)
	if err != nil {
		idj.err = err
		return
	}
	for _, file := range files {
		key := path.Join(idj.dirname, file.Name())
		if file.Mode().IsRegular() {
			idj.files[key] = FileDescr{file.Name(), file.Size()}
		} else if file.IsDir() {
			idj.childDirs = append(idj.childDirs, key)
		}
	}
}

func ScanDirParallel(startDir string, workers int) (files map[string]FileDescr, err error) {
	files = make(map[string]FileDescr)
	pool, err := workerpool.NewPool(workers)
	if err != nil {
		return files, err
	}

	// start from one dir and go 'deeper' to files tree
	// until exists child directories.
	// each iteration is one level of files hierarhy
	childDirs := []string{startDir}
	for len(childDirs) > 0 {
		// create index jobs
		jobs := make([]workerpool.Job, 0, len(childDirs))
		for _, dir := range childDirs {
			jobs = append(jobs, &IndexDirJob{
				dirname:   dir,
				childDirs: make([]string, 0),
				files:     make(map[string]FileDescr),
			})
		}
		// run jobs in concurrent mode
		err := pool.RunBatch(jobs...)
		if err != nil {
			return files, err
		}
		// parse results
		childDirs = []string{}
		for _, j := range jobs {
			job := j.(*IndexDirJob)
			if job.err != nil {
				return files, job.err
			}
			childDirs = append(childDirs, job.childDirs...)
			for key, val := range job.files {
				files[key] = val
			}
		}
	}
	return
}
