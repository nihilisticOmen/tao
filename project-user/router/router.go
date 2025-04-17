package router

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"net"
	"project-user/config"
	loginServiceV1 "project-user/pkg/service/login.service.v1"
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

type gRPCConfig struct {
	Addr         string
	RegisterFunc func(*grpc.Server)
}

func RegisterGrpc() *grpc.Server {
	c := gRPCConfig{
		Addr: config.AppConf.GC.Addr,
		RegisterFunc: func(g *grpc.Server) {
			loginServiceV1.RegisterLoginServiceServer(g, loginServiceV1.New())
		}}
	// 创建grpc服务
	s := grpc.NewServer()
	// 注册服务
	c.RegisterFunc(s)
	lis, err := net.Listen("tcp", config.AppConf.GC.Addr)
	if err != nil {
		log.Println("cannot listen")
	}
	go func() {
		err = s.Serve(lis)
		if err != nil {
			log.Println("server started error", err)
			return
		}
	}()
	return s
}
