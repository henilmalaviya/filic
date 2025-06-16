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
// it creates an empty file with 0644 permissions (rw-r--r--). Parent directories
// are automatically created if they don't exist.
func (f *File) Create() error {
	if f.Exists() {
		return nil
	}

	parent := f.OpenParent()

	if !parent.Exists() {
		err := parent.Create()

		if err != nil {
			return err
		}
	}

	return f.Write([]byte{})
}

// Write writes the provided data to the file, replacing any existing content.
// The file is created if it doesn't exist, and parent directories are not
// automatically created. The file is written with 0644 permissions (rw-r--r--).
func (f *File) Write(data []byte) error {
	return os.WriteFile(f.Path, data, 0644)
}

// Read reads the entire contents of the file and returns it as a byte slice.
// It returns an error if the file doesn't exist or cannot be read.
func (f *File) Read() ([]byte, error) {
	return os.ReadFile(f.Path)
}

// ReadString reads the entire contents of the file and returns it as a string.
// It's a convenience method that calls Read() and converts the result to a string.
// It returns an error if the file doesn't exist or cannot be read.
func (f *File) ReadString() (string, error) {
	data, err := f.Read()
	return string(data), err
}

// Append appends the provided data to the end of the file. If the file doesn't
// exist, this method will return an error. The file must already exist before
// calling this method. Use Create() or Write() to create the file first if needed.
// The data is appended with write-only permissions.
func (f *File) Append(data []byte) error {
	file, err := os.OpenFile(f.Path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
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
