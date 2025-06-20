# filic

A simple Go library for file system operations with a clean API for working with files and directories.

## Installation

```bash
go get github.com/henilmalaviya/filic
```

## Usage

### Working with Directories

```go
package main

import "github.com/henilmalaviya/filic"

func main() {
    // Create a new directory instance
    dir := filic.NewDirectory("/path/to/directory")

    // Create the directory if it doesn't exist
    err := dir.Create()
    if err != nil {
        // handle error
    }

    // Check if directory exists
    if dir.Exists() {
        // Open a subdirectory
        subdir, err := dir.OpenDir("subdir")
        if err != nil {
            // handle error
        }

        // Open a file in the directory
        file, err := dir.OpenFile("file.txt")
        if err != nil {
            // handle error
        }
    }

    // Get parent directory
    parent := dir.OpenParent()

    // Join paths
    fullPath := dir.Join("filename.txt")
}
```

### Working with Files

```go
package main

import "github.com/henilmalaviya/filic"

func main() {
    // Create a new file instance
    file := filic.NewFile("/path/to/file.txt")

    // Create the file if it doesn't exist
    err := file.Create()
    if err != nil {
        // handle error
    }

    // Write data to file
    err = file.Write([]byte("Hello, World!"))
    if err != nil {
        // handle error
    }

    // Check if file exists
    if file.Exists() {
        // File operations...
    }

    // Get parent directory
    parent := file.OpenParent()
}
```

## License

MIT
