package QuickRes

import (
	"github.com/gin-gonic/gin"
	"main/StatusCode"
	"net/http"
)

func SetOrigin(c *gin.Context) {
	or := c.GetHeader("Origin")
	c.Header("Access-Control-Allow-Origin", or)
	c.Header("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
	c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	c.Header("Access-Control-Allow-Credentials", "true")
	if c.Request.Method == "OPTIONS" {
		c.Status(http.StatusOK)
	}
}

func BadRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, gin.H{
		"status":  StatusCode.BadRequest,
		"message": "参数错误或不完整",
	})
}

func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  StatusCode.InternalError,
		"message": "服务器处理错误",
	})
}

func ProcessOK(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  StatusCode.Success,
		"message": "OK",
	})
}

func NotPermitted(c *gin.Context) {
	c.JSON(http.StatusForbidden, gin.H{
		"status":  StatusCode.NotPermitted,
		"message": "没有权限",
	})
}
