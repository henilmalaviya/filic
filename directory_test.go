package filic_test

import (
	"os"
	"path"
	"testing"

	"github.com/henilmalaviya/filic"
)

// cross-platform temp dir
func getTempDirPath() string {
	return path.Join(os.TempDir(), "filic")
}

func cleanup() {
	os.RemoveAll(getTempDirPath())
}

func TestNew(t *testing.T) {
	filic.NewDirectory(getTempDirPath())
}

func TestExists(t *testing.T) {

	dir := filic.NewDirectory(getTempDirPath())

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

func TestOpenParent(t *testing.T) {

	cleanup()

	dir := filic.NewDirectory(getTempDirPath())

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

	dir := filic.NewDirectory(getTempDirPath())

	joined := dir.Join("abc")

	expected := path.Join(getTempDirPath(), "abc")

	if joined != expected {
		t.Errorf("Joined path should be %v, got %v", expected, joined)
	}
}

func TestOpenDir(t *testing.T) {

	cleanup()

	dir := filic.NewDirectory(getTempDirPath())

	os.MkdirAll(dir.Path, 0755)

	f, err := os.Create(dir.Join("subfile"))

	if err != nil {
		t.Error(err)
	}

	_, err = dir.OpenDir("subfile")

	// we expect error
	if err == nil {
		t.Error("Opening a directory should return an error")
	}

	f.Close()

	cleanup()
}

func TestOpenFile(t *testing.T) {

	cleanup()

	dir := filic.NewDirectory(getTempDirPath())

	os.MkdirAll(dir.Join("subdir"), 0755)

	_, err := dir.OpenFile("subdir")

	// we expect error
	if err == nil {
		t.Error("Opening a file should return an error")
	}

	cleanup()

}

func TestCreate(t *testing.T) {

	cleanup()

	dir := filic.NewDirectory(path.Join(getTempDirPath(), "abc"))

	if dir.Exists() {
		t.Error("Directory should not exist yet")
	}

	err := dir.Create()

	if err != nil {
		t.Error(err)
	}

	if !dir.Exists() {
		t.Error("Directory should exist")
	}

	cleanup()
}
