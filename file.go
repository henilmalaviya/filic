package filic

import "os"

// File represents a file in the file system. It embeds Entity to inherit
// common file system operations and adds file-specific functionality like
// creating files and writing data to them.
type File struct {
	Entity
}

// Create creates an empty file at the specified path. If the file already
// exists, this method does nothing and returns nil. If the file doesn't exist,
// it creates an empty file with 0644 permissions (rw-r--r--).
func (f *File) Create() error {
	if f.Exists() {
		return nil
	}
	return f.Write([]byte{})
}

// Write writes the provided data to the file, replacing any existing content.
// The file is created if it doesn't exist, and parent directories are not
// automatically created. The file is written with 0644 permissions (rw-r--r--).
func (f *File) Write(data []byte) error {
	return os.WriteFile(f.Path, data, 0644)
}

// NewFile creates a new File instance with the specified path.
// The file doesn't need to exist at the time of creation - it can be
// created later using the Create method or written to using the Write method.
func NewFile(path string) *File {
	return &File{
		Entity: Entity{
			Path: path,
		},
	}
}
