package repositories

import (
	"fs/app/entity"
)

type ItemRepository interface {
	//InsertNewItem auto create id for item if not given
	//
	//return error if parent_id not valid
	InsertNewItem(item entity.Item) error
	GetById(id string) (entity.Item, error)
	GetByOrderNum(num int) (entity.Item, error)
	Update(item entity.Item) error
	Delete(id string) error
	DeleteAll() error
	GetAll() ([]entity.Item, error)
}
