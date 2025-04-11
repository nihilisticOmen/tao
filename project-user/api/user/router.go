package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"tao.com/project-user/router"
)

type RouterUser struct {
}

func init() {
	log.Println("init user router")
	// Register the router
	router.Register(&RouterUser{})
}
func (*RouterUser) Router(r *gin.Engine) {
	// User routes
	h := New()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
}
