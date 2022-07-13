package app

import (
	"fs/app/database"
	"fs/app/entity"
	"fs/app/repositories"
	"fs/app/usecase"
	item "fs/module/fs/repositories"
	ucase "fs/module/fs/usecase"
	"log"
)

func InitAppUseCase() {
	root := entity.Node{}
	treeUC := ucase.NewTreeUseCase(&root, repositories.GetAppRepo().ItemRepo)
	err := treeUC.Import()
	if err != nil {
		log.Fatal(err)
	}
	usecase.InitAppUseCase(treeUC)
}

func InitAppRepositories() {
	db := database.GetDB()
	itemRepo := item.NewItemRepo(db)

	repositories.InitAppRepositories(itemRepo)
	insertInitialData(itemRepo)
}

func insertInitialData(repo repositories.ItemRepository) {
	data := make([]entity.Node, 0)
	repo.DeleteAll()

	root := entity.Node{
		Id:       "0",
		Type:     0,
		Name:     "~",
		Data:     "",
		Path:     "/",
		Children: []*entity.Node{},
	}
	data = append(data, root)

	x1 := entity.Node{
		Id:       "1",
		Type:     0,
		Parent:   &root,
		Name:     "tung",
		Data:     "",
		Path:     "/tung",
		Children: []*entity.Node{},
		Root:     &root,
	}
	root.Children = append(root.Children, &x1)
	data = append(data, x1)

	x2a := entity.Node{
		Id:       "2",
		Type:     1,
		Parent:   &x1,
		Name:     "file1",
		Path:     "/tung/file1",
		Data:     "hello world",
		Children: []*entity.Node{},
		Root:     &root,
	}
	x1.Children = append(x1.Children, &x2a)
	data = append(data, x2a)

	x2b := entity.Node{
		Id:       "3",
		Type:     1,
		Parent:   &x1,
		Name:     "file2",
		Path:     "/tung/file2",
		Data:     "as98d9asd",
		Children: []*entity.Node{},
		Root:     &root,
	}
	x1.Children = append(x1.Children, &x2b)
	data = append(data, x2b)

	x2c := entity.Node{
		Id:       "4",
		Type:     0,
		Parent:   &x1,
		Name:     "folder1",
		Path:     "/tung/folder1",
		Data:     "",
		Children: []*entity.Node{},
		Root:     &root,
	}
	x1.Children = append(x1.Children, &x2c)
	data = append(data, x2c)

	x2d := entity.Node{
		Id:       "5",
		Type:     1,
		Parent:   &x2c,
		Name:     "file3",
		Path:     "/tung/folder1/file3",
		Data:     "as98d9asd",
		Children: []*entity.Node{},
		Root:     &root,
	}
	x2c.Children = append(x2c.Children, &x2d)
	data = append(data, x2d)

	x3 := entity.Node{
		Id:       "6",
		Type:     0,
		Parent:   &root,
		Name:     "usr",
		Path:     "/usr",
		Data:     "",
		Children: []*entity.Node{},
		Root:     &root,
	}
	root.Children = append(root.Children, &x3)
	data = append(data, x3)

	x4 := entity.Node{
		Id:       "7",
		Type:     1,
		Parent:   &x3,
		Name:     "file4",
		Path:     "/usr/file4",
		Data:     "hello tung",
		Children: []*entity.Node{},
		Root:     &root,
	}
	x3.Children = append(x3.Children, &x4)
	data = append(data, x4)

	for _, x := range data[1:] {
		item, _ := x.ToItem()
		repo.InsertNewItem(item)
	}
}
