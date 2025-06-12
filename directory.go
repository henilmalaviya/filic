package filic

type Directory struct {
	Entity
}

func (d *Directory) Open(name string) (FileSystemEntity, error) {

	path := d.Join(name)

	entity := NewEntity(path)

	exists := entity.Exists()
	// if the entity doesn't exists, we assume it's a file
	if !exists {
		return NewFile(path), nil
	}

	isDir, err := entity.IsDirectory()

	if err != nil {
		return nil, err
	}

	if isDir {
		return NewDirectory(path), nil
	}

	return NewFile(path), nil
}

func NewDirectory(path string) *Directory {
	return &Directory{
		Entity: Entity{
			Path: path,
		},
	}
}
