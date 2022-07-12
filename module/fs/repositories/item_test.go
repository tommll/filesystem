package repositories

import (
	"fmt"
	"fs/app/database"
	"fs/app/entity"
	"fs/app/repositories"
	"github.com/stretchr/testify/assert"
	"testing"
)

var repo repositories.ItemRepository

type TestData struct {
	X1 entity.Item
	X2 entity.Item
	X3 entity.Item
	X4 entity.Item
	X5 entity.Item
}

var testData TestData

func init() {
	database.Init()
	db := database.GetDB()
	repo = NewItemRepo(db)

	testData = TestData{}

	testData.X1 = entity.Item{
		ItemType: 0,
		OrderNum: 1,
		ParentId: "0",
		Name:     "qqq",
		Data:     "",
	}

	testData.X2 = entity.Item{
		ItemType: 1,
		OrderNum: 2,
		ParentId: "0",
		Name:     "aaa",
		Data:     "asd",
	}
}

func TestInsertItemWithoutID(t *testing.T) {
	database.Init()
	db := database.GetDB()
	repo = NewItemRepo(db)

	x := entity.Item{ParentId: "0", Name: "2222", Id: "13"}
	gotErr := repo.InsertNewItem(x)
	assert.NoError(t, gotErr)
}

func TestBasic(t *testing.T) {
	t.Run("insert", func(t *testing.T) {
		gotErr := repo.InsertNewItem(testData.X1)
		assert.NoError(t, gotErr)
	})

	t.Run("get", func(t *testing.T) {
		gotItem, gotErr := repo.GetById("0")
		assert.NoError(t, gotErr)
		assert.Equal(t, "/", gotItem.Name)
	})

	t.Run("update", func(t *testing.T) {
		item := entity.Item{Id: "2", ItemType: 1, ParentId: "0", Name: "this is a new name", Data: "new data"}
		gotErr := repo.Update(item)
		assert.NoError(t, gotErr)
	})

	t.Run("delete", func(t *testing.T) {
		gotErr := repo.Delete("2")
		assert.NoError(t, gotErr)

		_, gotErr = repo.GetById("2")
		assert.NotEqualf(t, nil, gotErr, "")
	})

	t.Run("get total", func(t *testing.T) {
		gotNum, gotErr := repo.GetAll()
		assert.NoError(t, gotErr)
		fmt.Printf("total of %d\n", len(gotNum))
	})

	t.Run("get by order num", func(t *testing.T) {
		gotItem, gotErr := repo.GetByOrderNum(1)
		assert.NoError(t, gotErr)
		assert.Equal(t, "asd", gotItem.Name)
	})

	t.Run("delete all", func(t *testing.T) {
		gotErr := repo.DeleteAll()
		assert.NoError(t, gotErr)
	})

}
