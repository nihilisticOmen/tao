package router

import (
	"github.com/gin-gonic/gin"
)

type Router interface {
	Router(r *gin.Engine)
}
type RegisterRouter struct {
}

func (*RegisterRouter) Router(ro Router, r *gin.Engine) {
	ro.Router(r)
}

//func New() *RegisterRouter {
//	return &RegisterRouter{}
//}

var routers []Router

func InitRouter(r *gin.Engine) {
	//rg := New()
	//rg.Router(&user.RouterUser{}, r)
	for _, ro := range routers {
		ro.Router(r)
	}
}
func Register(ro ...Router) {
	routers = append(routers, ro...)
}
