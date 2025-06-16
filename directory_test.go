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

func TestListAsEntities(t *testing.T) {
	cleanup()

	dir := filic.NewDirectory(getTempDirPath())
	err := dir.Create()
	if err != nil {
		t.Error(err)
	}

	// Create test files and directories
	subDir := path.Join(dir.Path, "subdir")
	os.MkdirAll(subDir, 0755)

	file1, err := os.Create(path.Join(dir.Path, "file1.txt"))
	if err != nil {
		t.Error(err)
	}
	file1.Close()

	file2, err := os.Create(path.Join(dir.Path, "file2.txt"))
	if err != nil {
		t.Error(err)
	}
	file2.Close()

	// Test ListAsEntities
	entities, err := dir.ListAsEntities()
	if err != nil {
		t.Error(err)
	}

	if len(entities) != 3 {
		t.Errorf("Expected 3 entities, got %d", len(entities))
	}

	// Verify entities have correct paths
	expectedNames := map[string]bool{"subdir": false, "file1.txt": false, "file2.txt": false}
	for _, entity := range entities {
		name := path.Base(entity.Path)
		if _, exists := expectedNames[name]; exists {
			expectedNames[name] = true
		} else {
			t.Errorf("Unexpected entity: %s", name)
		}
	}

	// Check all expected entities were found
	for name, found := range expectedNames {
		if !found {
			t.Errorf("Expected entity not found: %s", name)
		}
	}

	cleanup()
}

func TestListDirectories(t *testing.T) {
	cleanup()

	dir := filic.NewDirectory(getTempDirPath())
	err := dir.Create()
	if err != nil {
		t.Error(err)
	}

	// Create test directories and files
	subDir1 := path.Join(dir.Path, "subdir1")
	subDir2 := path.Join(dir.Path, "subdir2")
	os.MkdirAll(subDir1, 0755)
	os.MkdirAll(subDir2, 0755)

	file1, err := os.Create(path.Join(dir.Path, "file1.txt"))
	if err != nil {
		t.Error(err)
	}
	file1.Close()

	// Test ListDirectories
	directories, err := dir.ListDirectories()
	if err != nil {
		t.Error(err)
	}

	if len(directories) != 2 {
		t.Errorf("Expected 2 directories, got %d", len(directories))
	}

	// Verify directory names
	expectedDirs := map[string]bool{"subdir1": false, "subdir2": false}
	for _, directory := range directories {
		name := path.Base(directory.Path)
		if _, exists := expectedDirs[name]; exists {
			expectedDirs[name] = true
		} else {
			t.Errorf("Unexpected directory: %s", name)
		}
	}

	// Check all expected directories were found
	for name, found := range expectedDirs {
		if !found {
			t.Errorf("Expected directory not found: %s", name)
		}
	}

	cleanup()
}

func TestListFiles(t *testing.T) {
	cleanup()

	dir := filic.NewDirectory(getTempDirPath())
	err := dir.Create()
	if err != nil {
		t.Error(err)
	}

	// Create test files and directories
	subDir := path.Join(dir.Path, "subdir")
	os.MkdirAll(subDir, 0755)

	file1, err := os.Create(path.Join(dir.Path, "file1.txt"))
	if err != nil {
		t.Error(err)
	}
	file1.Close()

	file2, err := os.Create(path.Join(dir.Path, "file2.log"))
	if err != nil {
		t.Error(err)
	}
	file2.Close()

	// Test ListFiles
	files, err := dir.ListFiles()
	if err != nil {
		t.Error(err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}

	// Verify file names
	expectedFiles := map[string]bool{"file1.txt": false, "file2.log": false}
	for _, file := range files {
		name := path.Base(file.Path)
		if _, exists := expectedFiles[name]; exists {
			expectedFiles[name] = true
		} else {
			t.Errorf("Unexpected file: %s", name)
		}
	}

	// Check all expected files were found
	for name, found := range expectedFiles {
		if !found {
			t.Errorf("Expected file not found: %s", name)
		}
	}

	cleanup()
}

func TestListingNonExistentDirectory(t *testing.T) {
	cleanup()

	dir := filic.NewDirectory(getTempDirPath())

	// Test ListAsEntities on non-existent directory
	_, err := dir.ListAsEntities()
	if err == nil {
		t.Error("Expected error when listing non-existent directory")
	}

	// Test ListDirectories on non-existent directory
	_, err = dir.ListDirectories()
	if err == nil {
		t.Error("Expected error when listing directories in non-existent directory")
	}

	// Test ListFiles on non-existent directory
	_, err = dir.ListFiles()
	if err == nil {
		t.Error("Expected error when listing files in non-existent directory")
	}

	cleanup()
}
