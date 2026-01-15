# filic

`filic` is a small, idiomatic Go library that provides a higher-level, object-style API for working with the filesystem. It wraps common file and directory operations in simple `Directory` and `File` types, so you can manipulate paths and perform I/O in a more expressive way than using the raw `os` and path utilities directly.

Typical use cases include:

- Creating and working with nested directory structures
- Reading, writing, and appending to files
- Listing directory contents, filtered as files or subdirectories
- Navigating between parent and child paths in a clean, composable manner

---

## Installation

Install the module using `go get`:

```bash
go get github.com/henilmalaviya/filic
```

Then import it in your code:

```go
import "github.com/henilmalaviya/filic"
```

---

## Overview

The library exposes three main abstractions:

- `Entity`  
  A basic filesystem entity that holds a `Path` and provides shared methods such as `Exists`, `IsDirectory`, `Join`, and `OpenParent`.

- `Directory`  
  Represents a directory on disk. Allows you to:

  - Ensure a directory exists (`Create`)
  - Open child directories and files (`OpenDir`, `OpenFile`)
  - List contents (`List`, `ListDirectories`, `ListFiles`, `ListAsEntities`)

- `File`  
  Represents a file on disk. Allows you to:
  - Ensure a file and its parent directories exist (`Create`)
  - Write content (`Write`)
  - Read content as bytes or string (`Read`, `ReadString`)
  - Append to an existing file (`Append`)

All types live in the root `filic` package.

---

## Getting Started

### Working with Directories

The `Directory` type helps you work with directory trees in a fluent, safe manner.

```go
package main

import (
    "fmt"
    "log"

    "github.com/henilmalaviya/filic"
)

func main() {
    // Create a new directory instance (it may or may not exist yet).
    dir := filic.NewDirectory("/path/to/directory")

    // Ensure the directory (and any missing parents) exists.
    if err := dir.Create(); err != nil {
        log.Fatalf("failed to create directory: %v", err)
    }

    // Check if the directory now exists.
    if dir.Exists() {
        fmt.Println("Directory exists:", dir.Path)
    }

    // Open a subdirectory relative to the current directory.
    subdir, err := dir.OpenDir("subdir")
    if err != nil {
        log.Fatalf("failed to open subdir: %v", err)
    }

    // You can also ensure the subdirectory exists.
    if err := subdir.Create(); err != nil {
        log.Fatalf("failed to create subdir: %v", err)
    }

    // Open a file inside the directory (does not create it yet).
    file, err := dir.OpenFile("file.txt")
    if err != nil {
        log.Fatalf("failed to open file within directory: %v", err)
    }

    // Work with the file (see examples below).
    _ = file

    // Get the parent directory of the current directory.
    parent := dir.OpenParent()
    fmt.Println("Parent directory:", parent.Path)

    // Join a child path to this directory.
    fullPath := dir.Join("another-file.txt")
    fmt.Println("Joined path:", fullPath)
}
```

#### Listing Directory Contents

`Directory` provides several methods for listing its contents:

```go
// List returns the names of all entries (files and directories) within dir.
names, err := dir.List()
if err != nil {
    log.Fatalf("failed to list directory: %v", err)
}
fmt.Println("Entries:", names)

// ListDirectories returns only the subdirectories as *filic.Directory values.
dirs, err := dir.ListDirectories()
if err != nil {
    log.Fatalf("failed to list subdirectories: %v", err)
}
for _, d := range dirs {
    fmt.Println("Subdirectory:", d.Path)
}

// ListFiles returns only the files as *filic.File values.
files, err := dir.ListFiles()
if err != nil {
    log.Fatalf("failed to list files: %v", err)
}
for _, f := range files {
    fmt.Println("File:", f.Path)
}
```

---

### Working with Files

The `File` type wraps common file operations with a straightforward API.

```go
package main

import (
    "fmt"
    "log"

    "github.com/henilmalaviya/filic"
)

func main() {
    // Create a new file instance (it may or may not exist yet).
    file := filic.NewFile("/path/to/file.txt")

    // Ensure the file exists. This will:
    // - create any missing parent directories
    // - create an empty file if it doesn't exist
    if err := file.Create(); err != nil {
        log.Fatalf("failed to create file: %v", err)
    }

    // Write data to the file (overwrites existing contents).
    if err := file.Write([]byte("Hello, World!\n")); err != nil {
        log.Fatalf("failed to write file: %v", err)
    }

    // Check if the file exists.
    if file.Exists() {
        fmt.Println("File exists:", file.Path)
    }

    // Read file contents as bytes.
    data, err := file.Read()
    if err != nil {
        log.Fatalf("failed to read file: %v", err)
    }
    fmt.Printf("File content (bytes): %v\n", data)

    // Read file contents as a string.
    text, err := file.ReadString()
    if err != nil {
        log.Fatalf("failed to read file as string: %v", err)
    }
    fmt.Printf("File content (string): %s", text)

    // Append additional data to the existing file.
    // Note: Append expects the file to already exist.
    if err := file.Append([]byte("Appended line.\n")); err != nil {
        log.Fatalf("failed to append to file: %v", err)
    }

    // Get the parent directory as a filic.Directory.
    parent := file.OpenParent()
    fmt.Println("Parent directory:", parent.Path)
}
```

---

### Combining Directories and Files

A common pattern is to start from a base directory and then navigate to specific files:

```go
base := filic.NewDirectory("/var/log/myapp")

// Ensure base directory and nested structure exist.
if err := base.Create(); err != nil {
    log.Fatalf("failed to create base directory: %v", err)
}

logsDir, err := base.OpenDir("logs")
if err != nil {
    log.Fatalf("failed to open logs dir: %v", err)
}
if err := logsDir.Create(); err != nil {
    log.Fatalf("failed to create logs dir: %v", err)
}

// Open a log file inside the logs directory.
logFile, err := logsDir.OpenFile("app.log")
if err != nil {
    log.Fatalf("failed to open log file: %v", err)
}

// Create it if necessary, then append a line.
if !logFile.Exists() {
    if err := logFile.Create(); err != nil {
        log.Fatalf("failed to create log file: %v", err)
    }
}

if err := logFile.Append([]byte("application started\n")); err != nil {
    log.Fatalf("failed to append to log file: %v", err)
}
```

---

## API Summary

Below is a brief summary of the primary public methods. For full details, refer to the GoDoc comments in the source.

### `Entity`

Created indirectly via `NewDirectory` / `NewFile`, or directly via:

```go
e := filic.NewEntity("/path/to/something")
```

Methods:

- `func (e *Entity) Exists() bool`  
  Returns `true` if the path exists, `false` otherwise.

- `func (e *Entity) IsDirectory() (bool, error)`  
  Returns `true` if the path exists and is a directory, `false` if it exists and is a file.  
  Returns an error if the path cannot be inspected (e.g., permission issues).

- `func (e *Entity) Join(name string) string`  
  Joins `name` to `e.Path` and returns the resulting path string.

- `func (e *Entity) OpenParent() filic.Directory`  
  Returns a `Directory` representing the parent directory of `e.Path`.

### `Directory`

Constructed via:

```go
dir := filic.NewDirectory("/path/to/dir")
```

Key methods:

- `func (d *Directory) Create() error`  
  Creates the directory and all missing parents if they do not exist.  
  Safe to call multiple times (idempotent).

- `func (d *Directory) OpenDir(name string) (*Directory, error)`  
  Returns a `Directory` for the child path `d.Join(name)`.  
  If the path already exists and is not a directory, returns an error.

- `func (d *Directory) OpenFile(name string) (*File, error)`  
  Returns a `File` for the child path `d.Join(name)`.  
  If the path already exists and is a directory, returns an error.

- `func (d *Directory) List() ([]string, error)`  
  Returns the names of all entries (files and directories) inside `d`.

- `func (d *Directory) ListAsEntities() ([]filic.Entity, error)`  
  Returns all entries as `Entity` values.

- `func (d *Directory) ListDirectories() ([]*Directory, error)`  
  Returns only subdirectories.

- `func (d *Directory) ListFiles() ([]*File, error)`  
  Returns only files.

### `File`

Constructed via:

```go
file := filic.NewFile("/path/to/file.txt")
```

Key methods:

- `func (f *File) Create() error`  
  Ensures the parent directory exists and creates an empty file if it does not already exist.

- `func (f *File) Write(data []byte) error`  
  Writes `data` to the file, creating it if necessary and overwriting any existing contents.

- `func (f *File) Read() ([]byte, error)`  
  Reads the file contents into memory and returns them as a byte slice.

- `func (f *File) ReadString() (string, error)`  
  Convenience wrapper around `Read` that returns a string.

- `func (f *File) Append(data []byte) error`  
  Appends `data` to the file. The file must already exist; otherwise an error is returned.

---

## Error Handling

All operations that touch the filesystem (`Create`, `OpenDir`, `OpenFile`, `List*`, `Read`, `Write`, `Append`, etc.) can fail and return errors. In production code, you should always check and handle these errors appropriately.

For example:

```go
file := filic.NewFile("/path/to/file.txt")

if err := file.Create(); err != nil {
    // handle error: log, wrap, or propagate
}
```

---

## Testing

If you clone this repository and want to run the library’s own test suite:

```bash
go test ./...
```

This will execute tests for both `Directory` and `File` behaviors.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
