package filic

import (
	"os"
	"path"
)

// FileSystemEntity interface that both Directory and File implement
type FileSystemEntity interface {
	IsDirectory() (bool, error)
	Join(name string) string
	Exists() bool
	OpenParent() Directory
}

type Entity struct {
	FileSystemEntity
	Path string
}

func (e *Entity) IsDirectory() (bool, error) {
	info, err := os.Stat(e.Path)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func (e *Entity) Join(name string) string {
	return path.Join(e.Path, name)
}

func (e *Entity) OpenParent() Directory {
	return *NewDirectory(path.Dir(e.Path))
}

func (e *Entity) Exists() bool {
	_, err := os.Stat(e.Path)
	return err == nil
}

func NewEntity(path string) *Entity {
	return &Entity{
		Path: path,
	}
}
