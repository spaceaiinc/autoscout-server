package entity

type FilePath struct {
	PathName string `json:"path_name"`
}

func NewFilePath(
	pathName string,
) *FilePath {
	return &FilePath{
		PathName: pathName,
	}
}
