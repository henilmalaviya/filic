package filic

import (
	"fmt"
	"os"
)

// Directory represents a directory in the file system. It embeds Entity
// to inherit common file system operations and adds directory-specific
// functionality like creating directories and opening child files/directories.
type Directory struct {
	Entity
}

// Create creates the directory at the specified path, including any necessary
// parent directories. If the directory already exists, this method does nothing
// and returns nil. It uses permissions 0755 (rwxr-xr-x) for created directories.
func (d *Directory) Create() error {
	if d.Exists() {
		return nil
	}
	return os.MkdirAll(d.Path, 0755)
}

// OpenDir opens or prepares to open a subdirectory with the given name.
// It returns a Directory instance for the child directory. If a file system
// entity already exists at the target path but is not a directory, it returns
// an error. If the path doesn't exist, it returns a Directory that can be
// created later using the Create method.
func (d *Directory) OpenDir(name string) (*Directory, error) {

	path := d.Join(name)

	entity := NewEntity(path)

	// if the path exists, and is not a directory
	// return error
	if entity.Exists() {

		isDir, err := entity.IsDirectory()
		if err != nil {
			return nil, err
		}

		if !isDir {
			return nil, fmt.Errorf("path %v exists but its not a directory", path)
		}

	}

	return NewDirectory(path), nil
}

// OpenFile opens or prepares to open a file with the given name within this directory.
// It returns a File instance for the target file. If a file system entity already
// exists at the target path but is a directory, it returns an error. If the path
// doesn't exist, it returns a File that can be created later using file operations.
func (d *Directory) OpenFile(name string) (*File, error) {

	path := d.Join(name)

	entity := NewEntity(path)

	// if the path exists, and is not a file
	// return error
	if entity.Exists() {

		isDir, err := entity.IsDirectory()
		if err != nil {
			return nil, err
		}

		if isDir {
			return nil, fmt.Errorf("path %v exists but its a directory", path)
		}

	}

	return NewFile(path), nil
}

// List returns a list of all the files and directories in the directory.
// It returns an error if the directory doesn't exist or cannot be read.
func (d *Directory) List() ([]string, error) {
	files, err := os.ReadDir(d.Path)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, file := range files {
		names = append(names, file.Name())
	}

	return names, nil
}

// ListAsEntities returns a list of all items in the directory as Entity instances.
// This provides a unified way to work with both files and directories, allowing
// you to check their type using the IsDirectory method. Returns an error if the
// directory doesn't exist or cannot be read.
func (d *Directory) ListAsEntities() ([]Entity, error) {
	names, err := d.List()
	if err != nil {
		return nil, err
	}

	var entities []Entity
	for _, name := range names {
		entity := NewEntity(d.Join(name))
		entities = append(entities, *entity)
	}

	return entities, nil
}

// ListDirectories returns a list of all subdirectories within this directory
// as Directory instances. Only items that are confirmed to be directories are
// included in the result. Returns an error if the directory doesn't exist,
// cannot be read, or if there's an issue checking the type of any item.
func (d *Directory) ListDirectories() ([]*Directory, error) {
	entities, err := d.ListAsEntities()
	if err != nil {
		return nil, err
	}

	var directories []*Directory
	for _, entity := range entities {
		isDir, err := entity.IsDirectory()

		if err != nil {
			return nil, err
		}

		if isDir {
			directories = append(directories, NewDirectory(entity.Path))
		}
	}
	return directories, nil
}

// ListFiles returns a list of all files within this directory as File instances.
// Only items that are confirmed to be files (not directories) are included in
// the result. Returns an error if the directory doesn't exist, cannot be read,
// or if there's an issue checking the type of any item.
func (d *Directory) ListFiles() ([]*File, error) {
	entities, err := d.ListAsEntities()
	if err != nil {
		return nil, err
	}

	var files []*File
	for _, entity := range entities {
		isDir, err := entity.IsDirectory()

		if err != nil {
			return nil, err
		}

		if !isDir {
			files = append(files, NewFile(entity.Path))
		}
	}
	return files, nil
}

// NewDirectory creates a new Directory instance with the specified path.
// The directory doesn't need to exist at the time of creation - it can be
// created later using the Create method.
func NewDirectory(path string) *Directory {
	return &Directory{
		Entity: Entity{
			Path: path,
		},
	}
}
