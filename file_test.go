package filic_test

import (
	"os"
	"path"
	"testing"

	"github.com/henilmalaviya/filic"
)

func TestNewFile(t *testing.T) {
	filic.NewFile(path.Join(getTempDirPath(), "abc.txt"))
}

func TestCreateFile(t *testing.T) {

	cleanup()

	tmpDir := filic.NewDirectory(getTempDirPath())

	tmpDir.Create()

	file, err := tmpDir.OpenFile("abc.txt")

	if err != nil {
		t.Error(err)
	}

	if file.Exists() {
		t.Error("File should not exist yet")
	}

	err = file.Create()

	if err != nil {
		t.Error(err)
	}

	if !file.Exists() {
		t.Error("File should exist")
	}

	cleanup()
}

func TestWriteFile(t *testing.T) {

	cleanup()

	tmpDir := filic.NewDirectory(getTempDirPath())

	tmpDir.Create()

	file, err := tmpDir.OpenFile("abc.txt")

	if err != nil {
		t.Error(err)
	}

	err = file.Write([]byte("hello world"))

	if err != nil {
		t.Error(err)
	}

	data, err := os.ReadFile(file.Path)

	if err != nil {
		t.Error(err)
	}

	expected := []byte("hello world")

	for i := range data {
		if data[i] != expected[i] {
			t.Errorf("Data should be %v, got %v", expected, data)
		}
	}
}
