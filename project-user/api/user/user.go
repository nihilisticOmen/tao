package user

import (
	"github.com/gin-gonic/gin"
	common "test.com/project-common"
)

type HandlerUser struct {
}

func (*HandlerUser) getCaptcha(c *gin.Context) {
	rsp := &common.Result{}
	c.JSON(200, rsp.Success("123456"))
}
