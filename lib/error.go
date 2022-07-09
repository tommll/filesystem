package lib

import (
	"github.com/gin-gonic/gin"
)

//DisableErrorLogger ... Manual disable log for unit testing
var DisableErrorLogger *bool

func SendInternalError(c *gin.Context, err error) {
	HandleError(c, 500, 0, err.Error())
}

// HandleError .
func HandleError(c *gin.Context, statusCode int, errorCode int, message string) {
	c.JSON(statusCode, gin.H{
		"return_code": errorCode,
		"message":     message,
	})
}
