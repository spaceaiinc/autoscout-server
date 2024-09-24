package responses

import "github.com/spaceaiinc/autoscout-server/domain/entity"

type FilePath struct {
	FilePath *entity.FilePath `json:"file_path"`
}

func NewFilePath(filePath *entity.FilePath) FilePath {
	return FilePath{
		FilePath: filePath,
	}
}
