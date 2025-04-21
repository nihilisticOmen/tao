package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"project-api/router"
)

type RouterUser struct {
}

func init() {
	log.Println("init user router")
	// Register the router
	router.Register(&RouterUser{})
}
func (*RouterUser) Router(r *gin.Engine) {
	InitRpcUserClient()
	// User routes
	h := New()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
	r.POST("/project/login/register", h.register)
}
