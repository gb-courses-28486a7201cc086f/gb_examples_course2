package main

import (
	"errors"
	"fmt"
	"io/fs"
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

func TestScanDir(t *testing.T) {
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
}

func TestScanDirNotExisting(t *testing.T) {
	notExistingDir := "somebadpath"
	_, err := ScanDir(notExistingDir)
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("unexpected error: %v", err)
	}

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

func TestGroupFiles(t *testing.T) {
	t.Run("with duplicates", func(t *testing.T) {
		testData := map[string]FileDescr{
			"/tmp/d1/file1": {"file1", 5},
			"/tmp/d2/file1": {"file1", 5},
			"/tmp/d3/file1": {"file1", 6},
			"/tmp/d4/file2": {"file2", 5},
		}
		expected := map[duplicateKey]map[string]struct{}{
			{"file1", 5}: {
				"/tmp/d1/file1": {},
				"/tmp/d2/file1": {},
			},
			{"file1", 6}: {
				"/tmp/d3/file1": {},
			},
			{"file2", 5}: {
				"/tmp/d4/file2": {},
			},
		}

		groups := groupFiles(testData)
		for key, values := range groups {
			// value (path) should be in expected set
			for _, val := range values {
				_, ok := expected[key][val]
				if !ok {
					t.Errorf("path %s for file %s does not exists in result", val, key.name)
				}
			}
		}
	})
	t.Run("empty map", func(t *testing.T) {
		testData := map[string]FileDescr{}
		expectedLen := 0

		groups := groupFiles(testData)
		if len(groups) != expectedLen {
			t.Errorf("got %d items, expected empty result", len(groups))
		}
	})
}

func TestCheckDuplicates(t *testing.T) {
	t.Run("with duplicates", func(t *testing.T) {
		testData := map[string]FileDescr{
			"/tmp/d1/file1": {"file1", 5},
			"/tmp/d2/file1": {"file1", 5},
			"/tmp/d3/file1": {"file1", 6},
			"/tmp/d4/file2": {"file2", 5},
		}
		expectedCount := 1
		expectedFiles := map[string]struct{}{
			"/tmp/d1/file1": {},
			"/tmp/d2/file1": {},
			"/tmp/d5/file1": {},
		}

		duplicates := CheckDuplicates(testData)
		if len(duplicates) != expectedCount {
			t.Errorf("invalid duplicates count: got %d, expected %d", len(duplicates), expectedCount)
			return
		}
		if _, ok := expectedFiles[duplicates[0].Origin]; !ok {
			t.Errorf("missed copy of file: %s", duplicates[0].Origin)
		}
		for _, copy := range duplicates[0].Duplicates {
			if _, ok := expectedFiles[copy]; !ok {
				t.Errorf("missed copy of file: %s", copy)
			}
		}
	})
	t.Run("with no duplicates", func(t *testing.T) {
		testData := map[string]FileDescr{
			"/tmp/d1/file1": {"file1", 5},
			"/tmp/d2/file2": {"file2", 5},
			"/tmp/d3/file3": {"file3", 5},
			"/tmp/d4/file4": {"file4", 5},
		}
		expectedCount := 0

		duplicates := CheckDuplicates(testData)
		if len(duplicates) != expectedCount {
			t.Errorf("invalid duplicates count: got %d, expected %d", len(duplicates), expectedCount)
			return
		}
	})
	t.Run("empty", func(t *testing.T) {
		testData := map[string]FileDescr{}
		expectedCount := 0

		duplicates := CheckDuplicates(testData)
		if len(duplicates) != expectedCount {
			t.Errorf("invalid duplicates count: got %d, expected %d", len(duplicates), expectedCount)
			return
		}
	})
}

func TestDeleteDuplicates(t *testing.T) {
	testDir, err := testSetUp()
	defer testTeardown(testDir)
	if err != nil {
		t.Fatal(err)
	}

	files, err := ScanDir(testDir)
	if err != nil {
		t.Fatal(err)
	}
	duplicates := CheckDuplicates(files)
	err = DeleteDuplicates(duplicates)
	if err != nil {
		t.Errorf("unexpected delete error: %v", err)
	}

	// try delete twice
	err = DeleteDuplicates(duplicates)
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("unexpected error: got %v, expected %v", err, fs.ErrNotExist)
	}

	// second scan to check files not exists
	files, err = ScanDir(testDir)
	if err != nil {
		t.Fatal(err)
	}
	duplicates = CheckDuplicates(files)
	expected := 0
	if len(duplicates) != expected {
		t.Errorf("duplicates still exists: got %d, expected %d", len(duplicates), expected)
	}
}
