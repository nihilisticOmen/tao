package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"project-api/api/midd"
	"project-api/api/rpc"
	"project-api/router"
)

type RouterUser struct {
}

func init() {
	log.Println("init user router")
	ru := &RouterUser{}
	router.Register(ru)
}

func (*RouterUser) Route(r *gin.Engine) {
	//初始化grpc的客户端连接
	rpc.InitRpcUserClient()
	h := New()
	r.POST("/project/login/getCaptcha", h.getCaptcha)
	r.POST("/project/login/register", h.register)
	r.POST("/project/login", h.login)
	org := r.Group("/project/organization")
	org.Use(midd.TokenVerify())
	org.POST("/_getOrgList", h.myOrgList)
}
