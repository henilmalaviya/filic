package filic

import (
	"os"
	"path"
	"testing"
)

// cross-platform temp dir
func getTempDirPath() string {
	return path.Join(os.TempDir(), "filic")
}

func cleanup() {
	os.RemoveAll(getTempDirPath())
}

func TestNew(t *testing.T) {
	NewDirectory(getTempDirPath())
}

func TestExists(t *testing.T) {

	dir := NewDirectory(getTempDirPath())

	// remove dir first
	cleanup()

	if dir.Exists() {
		t.Error("Directory should not exist")
	}

	// create dir to test the true case
	os.MkdirAll(dir.Path, 0755)

	if !dir.Exists() {
		t.Error("Directory should exist")
	}

	cleanup()
}

func TestOpen(t *testing.T) {

	cleanup()

	dir := NewDirectory(getTempDirPath())

	file, err := dir.Open("somefile_that_doesnt_exist")

	if err != nil {
		t.Errorf("Error opening file. error: %v", err)
	}

	if f, ok := file.(*File); !ok {
		t.Error("File should be returned")
	} else {
		if f.Exists() {
			t.Error("File should not exist")
		}
	}

	cleanup()
}

func TestOpenParent(t *testing.T) {

	cleanup()

	dir := NewDirectory(getTempDirPath())

	parent := dir.OpenParent()

	if parent.Path != path.Dir(getTempDirPath()) {
		t.Errorf("Parent path should be %v, got %v", getTempDirPath(), parent.Path)
	}

	if !parent.Exists() {
		t.Error("Parent should exist")
	}

	if isDir, ok := parent.IsDirectory(); ok != nil {
		t.Error("Error checking if parent is a directory")
	} else {
		if !isDir {
			t.Error("Parent should be a directory")
		}
	}

	cleanup()
}

func TestJoin(t *testing.T) {

	dir := NewDirectory(getTempDirPath())

	joined := dir.Join("somefile_that_doesnt_exist")

	expected := path.Join(getTempDirPath(), "somefile_that_doesnt_exist")

	if joined != expected {
		t.Errorf("Joined path should be %v, got %v", expected, joined)
	}
}
