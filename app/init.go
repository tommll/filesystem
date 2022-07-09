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
}
