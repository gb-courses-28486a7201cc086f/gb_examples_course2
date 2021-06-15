package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func testSetUp() (testDir string, err error) {
	testDir, err = os.MkdirTemp(os.TempDir(), "go-duplicates-tests-")
	if err != nil {
		return testDir, err
	}

	// we'd like tree nested levels of directories
	// with different files...
	for i := 0; i < 5; i++ {
		l1dir := path.Join(testDir, fmt.Sprintf("dir-1-%d", i))
		if err = os.Mkdir(l1dir, 0755); err != nil {
			return testDir, err
		}
		for j := 0; j < 5; j++ {
			l2dir := path.Join(l1dir, fmt.Sprintf("dir-2-%d", j))
			if err = os.Mkdir(l2dir, 0755); err != nil {
				return testDir, err
			}
			for k := 0; k < 5; k++ {
				l3dir := path.Join(l2dir, fmt.Sprintf("dir-3-%d", k))
				if err = os.Mkdir(l3dir, 0755); err != nil {
					return testDir, err
				}
				// make test file
				fileName := path.Join(l3dir, fmt.Sprintf("file-%d-%d-%d", i, j, k))
				if err = ioutil.WriteFile(fileName, []byte{1, 0, 0, 0, 0, 0, 0}, 0644); err != nil {
					return testDir, err
				}
			}
		}
	}
	// and two duplicates
	dupOne := path.Join(testDir, "dir-1-0", "dir-2-0", "dir-3-0", "duplicate")
	if err = ioutil.WriteFile(dupOne, []byte{1, 0, 0, 0, 0, 0, 0}, 0644); err != nil {
		return testDir, err
	}
	dupTwo := path.Join(testDir, "dir-1-4", "dir-2-0", "duplicate")
	if err = ioutil.WriteFile(dupTwo, []byte{1, 0, 0, 0, 0, 0, 0}, 0644); err != nil {
		return testDir, err
	}

	return testDir, nil
}

func testTeardown(testDir string) (err error) {
	return os.RemoveAll(testDir)
}

func TestDuplicates(t *testing.T) {
	testDir, err := testSetUp()
	defer testTeardown(testDir)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("ScanDir", func(t *testing.T) {
		expected := 127
		files, err := ScanDir(testDir)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(files) != expected {
			t.Errorf("invalid files count: got %d, expected %d", len(files), expected)
		}
	})

	t.Run("ScanDirParallel", func(t *testing.T) {
		expected := 127
		workers := 4
		files, err := ScanDirParallel(testDir, workers)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(files) != expected {
			t.Errorf("invalid files count: got %d, expected %d", len(files), expected)
		}
	})
}

func BenchmarkScanDir(b *testing.B) {
	testDir, err := testSetUp()
	defer testTeardown(testDir)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ScanDir(testDir)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkScanDirParallel(b *testing.B) {
	testDir, err := testSetUp()
	defer testTeardown(testDir)
	if err != nil {
		b.Fatal(err)
	}

	for _, workers := range []int{1, 2, 5, 10, 50, 100} {
		benchName := fmt.Sprintf("workers-%d", workers)
		b.Run(benchName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err := ScanDirParallel(testDir, workers)
				if err != nil {
					b.Error(err)
				}
			}
		})
	}
}
