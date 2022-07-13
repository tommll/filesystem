package router

import (
	"fs/app/usecase"
	"fs/module/fs/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// authorizedAccounts := map[string]string{
	// 	"tung": "123",
	// } // TODO: add authorized accounts

	somePath := "fs"
	// authorized := router.Group(somePath, gin.BasicAuth(authorizedAccounts))
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))

	router.Use(static.Serve("/", static.LocalFile("./static", false)))

	authorized := router.Group(somePath)

	someModule := controller.NewFileSystemHandler(usecase.GetAppUseCase().TreeUC)
	someModule.SetupRouter(authorized)

	return router
}

func InitRouter() {
	router := SetupRouter()
	router.Run(":8080")
}
