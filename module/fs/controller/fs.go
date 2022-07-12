package controller

import (
	"fs/app/entity"
	"fs/app/usecase"
	"fs/lib"
	"fs/module/fs/composer"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FileSystemHandler interface {
	SetupRouter(r *gin.RouterGroup)
	Navigate(c *gin.Context)
	Create(c *gin.Context)
	Show(c *gin.Context)
	List(c *gin.Context)
	FindMany(c *gin.Context)
	Update(c *gin.Context)
	Move(c *gin.Context)
	Remove(c *gin.Context)
}

func NewFileSystemHandler(ucase usecase.TreeUseCase) FileSystemHandler {
	return &handlerImpl{TreeUC: ucase}
}

type handlerImpl struct {
	TreeUC usecase.TreeUseCase
}

type NavigateRequest struct {
	Path string `json:"path"`
}

func (g *handlerImpl) Navigate(c *gin.Context) {
	var req NavigateRequest
	err := c.ShouldBind(&req)
	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, "invalid parameter")
	}

	_, err = g.TreeUC.Navigate(req.Path)

	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.BaseDataResponse{
		ReturnCode: 1,
		Message:    "okay",
	})
}

type CreateRequest struct {
	Path string `json:"path"`
	Data string `json:"data"`
}

func (g *handlerImpl) Create(c *gin.Context) {
	var req CreateRequest
	err := c.ShouldBind(&req)
	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, "invalid parameter")
	}

	err = g.TreeUC.Create(req.Path, req.Data)

	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.BaseDataResponse{
		ReturnCode: 1,
		Message:    "okay",
	})
}

type ShowRequest struct {
	Path string `json:"path"`
}

func (g *handlerImpl) Show(c *gin.Context) {
	var req ShowRequest
	err := c.ShouldBind(&req)
	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, "invalid parameter")
	}

	data, err := g.TreeUC.ShowData(req.Path)

	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.BaseDataResponse{
		ReturnCode: 1,
		Message:    "okay",
		Data:       data,
	})
}

type ListRequest struct {
	Path string `json:"path"`
}

func (g *handlerImpl) List(c *gin.Context) {
	var req ListRequest
	err := c.ShouldBind(&req)
	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, "invalid parameter")
	}

	data, err := g.TreeUC.List(req.Path)

	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, err.Error())
		return
	}

	resp := composer.NodeListToFileList(data)

	c.JSON(http.StatusOK, entity.BaseDataResponse{
		ReturnCode: 1,
		Message:    "okay",
		Data:       resp,
	})
}

type FindManyRequest struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func (g *handlerImpl) FindMany(c *gin.Context) {
	var req FindManyRequest
	err := c.ShouldBind(&req)
	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, "invalid parameter")
	}

	data, err := g.TreeUC.FindByName(req.Name, req.Path)

	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, err.Error())
		return
	}

	resp := composer.NodeListToFileList(data)

	c.JSON(http.StatusOK, entity.BaseDataResponse{
		ReturnCode: 1,
		Message:    "okay",
		Data:       resp,
	})
}

type UpdateRequest struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Data string `json:"data"`
}

func (g *handlerImpl) Update(c *gin.Context) {
	var req UpdateRequest
	err := c.ShouldBind(&req)
	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, "invalid parameter")
	}

	err = g.TreeUC.Update(req.Path, req.Name, req.Data)

	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.BaseDataResponse{
		ReturnCode: 1,
		Message:    "okay",
	})
}

type MoveRequest struct {
	SrcPath  string `json:"src_path"`
	DestPath string `json:"dest_path"`
}

func (g *handlerImpl) Move(c *gin.Context) {
	var req MoveRequest
	err := c.ShouldBind(&req)
	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, "invalid parameter")
	}

	err = g.TreeUC.Move(req.SrcPath, req.DestPath)

	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.BaseDataResponse{
		ReturnCode: 1,
		Message:    "okay",
	})
}

type RemoveRequest struct {
	Paths []string `json:"paths"`
}

func (g *handlerImpl) Remove(c *gin.Context) {
	var req RemoveRequest
	err := c.ShouldBind(&req)
	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, "invalid parameter")
	}

	err = g.TreeUC.Remove(req.Paths)

	if err != nil {
		lib.HandleError(c, http.StatusInternalServerError, -1, err.Error())
		return
	}

	c.JSON(http.StatusOK, entity.BaseDataResponse{
		ReturnCode: 1,
		Message:    "okay",
	})
}

func (g *handlerImpl) SetupRouter(r *gin.RouterGroup) {
	r.POST("/cd", g.Navigate)
	r.POST("/cr", g.Create)
	r.POST("/cat", g.Show)
	r.POST("/ls", g.List)
	r.POST("/find", g.FindMany)
	r.POST("/up", g.Update) // TODO: fix update function
	r.POST("/mv", g.Move)
	r.POST("/rm", g.Remove)
}
