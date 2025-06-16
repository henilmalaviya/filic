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

func TestReadFile(t *testing.T) {
	cleanup()

	tmpDir := filic.NewDirectory(getTempDirPath())
	tmpDir.Create()

	file, err := tmpDir.OpenFile("test.txt")
	if err != nil {
		t.Error(err)
	}

	// Test reading non-existent file
	_, err = file.Read()
	if err == nil {
		t.Error("Expected error when reading non-existent file")
	}

	file.Create()

	// Write data and test reading
	testData := []byte("test content for reading")
	err = file.Write(testData)
	if err != nil {
		t.Error(err)
	}

	data, err := file.Read()
	if err != nil {
		t.Error(err)
	}

	if len(data) != len(testData) {
		t.Errorf("Expected %d bytes, got %d", len(testData), len(data))
	}

	for i := range data {
		if data[i] != testData[i] {
			t.Errorf("Data mismatch at index %d: expected %v, got %v", i, testData, data)
		}
	}

	cleanup()
}

func TestReadStringFile(t *testing.T) {
	cleanup()

	tmpDir := filic.NewDirectory(getTempDirPath())
	tmpDir.Create()

	file, err := tmpDir.OpenFile("test.txt")
	if err != nil {
		t.Error(err)
	}

	// Test reading non-existent file
	_, err = file.ReadString()
	if err == nil {
		t.Error("Expected error when reading non-existent file as string")
	}

	// Write data and test reading as string
	testContent := "Hello, World!\nThis is a test file."
	err = file.Write([]byte(testContent))
	if err != nil {
		t.Error(err)
	}

	content, err := file.ReadString()
	if err != nil {
		t.Error(err)
	}

	if content != testContent {
		t.Errorf("Expected %q, got %q", testContent, content)
	}

	cleanup()
}

func TestAppendFile(t *testing.T) {
	cleanup()

	tmpDir := filic.NewDirectory(getTempDirPath())
	tmpDir.Create()

	file, err := tmpDir.OpenFile("test.txt")
	if err != nil {
		t.Error(err)
	}

	// Test appending to non-existent file
	err = file.Append([]byte("should fail"))
	if err == nil {
		t.Error("Expected error when appending to non-existent file")
	}

	// Create file with initial content
	initialContent := "Initial content"
	err = file.Write([]byte(initialContent))
	if err != nil {
		t.Error(err)
	}

	// Append additional content
	appendContent := "\nAppended content"
	err = file.Append([]byte(appendContent))
	if err != nil {
		t.Error(err)
	}

	// Read and verify combined content
	finalContent, err := file.ReadString()
	if err != nil {
		t.Error(err)
	}

	expected := initialContent + appendContent
	if finalContent != expected {
		t.Errorf("Expected %q, got %q", expected, finalContent)
	}

	cleanup()
}

func TestFileCreateWithParentDirectories(t *testing.T) {
	cleanup()

	// Create a file in a nested directory that doesn't exist
	nestedPath := path.Join(getTempDirPath(), "level1", "level2", "test.txt")
	file := filic.NewFile(nestedPath)

	if file.Exists() {
		t.Error("File should not exist yet")
	}

	err := file.Create()
	if err != nil {
		t.Error(err)
	}

	if !file.Exists() {
		t.Error("File should exist after creation")
	}

	// Verify parent directories were created
	parent := file.OpenParent()
	if !parent.Exists() {
		t.Error("Parent directory should have been created")
	}

	cleanup()
}

func TestFileOperationsChain(t *testing.T) {
	cleanup()

	tmpDir := filic.NewDirectory(getTempDirPath())
	tmpDir.Create()

	file, err := tmpDir.OpenFile("chain_test.txt")
	if err != nil {
		t.Error(err)
	}

	// Test complete workflow: Create -> Write -> Read -> Append -> Read
	err = file.Create()
	if err != nil {
		t.Error(err)
	}

	// Write initial content
	err = file.Write([]byte("Line 1"))
	if err != nil {
		t.Error(err)
	}

	// Read content
	content, err := file.ReadString()
	if err != nil {
		t.Error(err)
	}
	if content != "Line 1" {
		t.Errorf("Expected 'Line 1', got %q", content)
	}

	// Append more content
	err = file.Append([]byte("\nLine 2"))
	if err != nil {
		t.Error(err)
	}

	// Read final content
	finalContent, err := file.ReadString()
	if err != nil {
		t.Error(err)
	}
	expected := "Line 1\nLine 2"
	if finalContent != expected {
		t.Errorf("Expected %q, got %q", expected, finalContent)
	}

	cleanup()
}
