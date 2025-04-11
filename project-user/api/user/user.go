package user

import (
	"github.com/gin-gonic/gin"
)

type HandlerUser struct {
}

func (*HandlerUser) getCaptcha(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "User route",
	})
}
