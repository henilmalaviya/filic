package filic

type File struct {
	Entity
}

func NewFile(path string) *File {
	return &File{
		Entity: Entity{
			Path: path,
		},
	}
}
