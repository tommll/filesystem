package repositories

import (
	"errors"
	"fs/app/entity"
	"fs/app/repositories"
	"github.com/jinzhu/gorm"
	"strconv"
)

type repoImpl struct {
	db *gorm.DB
}

func NewItemRepo(db *gorm.DB) repositories.ItemRepository {
	return &repoImpl{db: db}
}

func (g *repoImpl) InsertNewItem(item entity.Item) error {
	if len(item.ParentId) == 0 {
		return ErrorNewItemMustHaveParentId
	}
	var tmp entity.Item
	query := g.db.First(&tmp, "id = ?", item.ParentId)
	if query.Error != nil {
		return ErrorNewItemMustHaveParentId
	}

	err := item.CheckValid()
	if err != nil {
		return err
	}

	if len(item.Id) == 0 {
		res := []entity.Item{}
		query = g.db.Find(&res)
		n := int(query.RowsAffected)
		item.Id = strconv.Itoa(n + 1)
	}

	return g.db.Save(&item).Error
}

func (g *repoImpl) GetById(id string) (entity.Item, error) {
	var res entity.Item
	query := g.db.First(&res, "id = ?", id)
	return res, query.Error
}

func (g *repoImpl) GetByOrderNum(num int) (entity.Item, error) {
	var res entity.Item
	query := g.db.First(&res, "order_num = ?", num)
	return res, query.Error
}

func (g *repoImpl) Update(item entity.Item) error {
	err := item.CheckValid()
	if err != nil {
		return err
	}

	return g.db.Save(&item).Error
}

func (g *repoImpl) Delete(id string) error {
	return g.db.Delete(&entity.Item{}, id).Error
}

func (g *repoImpl) DeleteAll() error {
	return g.db.Delete(&entity.Item{}, "name <> ?", "~").Error
}

func (g *repoImpl) GetAll() ([]entity.Item, error) {
	res := []entity.Item{}
	query := g.db.Find(&res)
	return res, query.Error
}

var (
	ErrorNewItemMustHaveParentId = errors.New("newly inserted item must have a valid parent id")
)
