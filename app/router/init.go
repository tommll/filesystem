package router

import (
	"fs/app/usecase"
	"fs/module/fs/controller"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	authorizedAccounts := map[string]string{
		"tung": "123",
	} // TODO: add authorized accounts

	somePath := "fs"
	authorized := router.Group(somePath, gin.BasicAuth(authorizedAccounts))
	//authorized := router.Group(somePath)

	someModule := controller.NewFileSystemHandler(usecase.GetAppUseCase().TreeUC)
	someModule.SetupRouter(authorized)

	return router
}

func InitRouter() {
	router := SetupRouter()
	router.Run(":8080")
}
