package main

import (
	"fs/app"
	"fs/app/database"
	"fs/app/router"
)

func main() {
	database.Init()
	app.InitAppRepositories()
	app.InitAppUseCase()
	router.InitRouter()
}
