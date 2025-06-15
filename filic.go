package filic

import (
	"os"
	"path"
)

// FileSystemEntity defines the common interface for file system operations
// that both Directory and File types implement. It provides methods for
// checking if an entity is a directory, joining paths, checking existence,
// and accessing the parent directory.
type FileSystemEntity interface {
	IsDirectory() (bool, error)
	Join(name string) string
	Exists() bool
	OpenParent() Directory
}

// Entity represents a file system entity (file or directory) with a specific path.
// It embeds the FileSystemEntity interface and provides concrete implementations
// for common file system operations.
type Entity struct {
	FileSystemEntity
	Path string
}

// IsDirectory checks whether the entity at the current path is a directory.
// It returns true if the path points to a directory, false if it's a file,
// and an error if the path cannot be accessed or doesn't exist.
func (e *Entity) IsDirectory() (bool, error) {
	info, err := os.Stat(e.Path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

// Join creates a new path by joining the current entity's path with the given name.
// This is useful for creating paths to child files or directories.
// The resulting path uses the appropriate path separator for the operating system.
func (e *Entity) Join(name string) string {
	return path.Join(e.Path, name)
}

// OpenParent returns a Directory instance representing the parent directory
// of the current entity. This allows navigation up the directory tree.
func (e *Entity) OpenParent() Directory {
	return *NewDirectory(path.Dir(e.Path))
}

// Exists checks whether the entity exists at the specified path.
// It returns true if the file or directory exists, false otherwise.
// This method does not distinguish between files and directories.
func (e *Entity) Exists() bool {
	_, err := os.Stat(e.Path)
	return err == nil
}

// NewEntity creates a new Entity instance with the specified path.
// The path can point to either a file or directory - the actual type
// can be determined later using the IsDirectory method.
func NewEntity(path string) *Entity {
	return &Entity{
		Path: path,
	}
}
