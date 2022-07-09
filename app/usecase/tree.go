package usecase

import (
	"fs/app/entity"
)

type TreeUseCase interface {
	Navigate(folderPath string) (*entity.Node, error)
	Create(folderPath string, data string) error
	ShowData(filePath string) (string, error)
	List(folderPath string) ([]*entity.Node, error)
	FindByName(name, folderPath string) ([]*entity.Node, error)
	Update(path, name, data string) error
	Move(sourcePath, destPath string) error
	Remove(paths []string) error

	Import() error
	Export() error
}
