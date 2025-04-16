package main

import (
	"github.com/gin-gonic/gin"
	srv "project-common"
	_ "project-user/api"
	"project-user/config"
	"project-user/router"
)

func main() {
	r := gin.Default()
	gc := router.RegisterGrpc()
	stop := func() {
		gc.Stop()
	}
	// Initialize the router
	router.InitRouter(r)
	srv.Run(r, config.AppConf.SC.Name, config.AppConf.SC.Addr, stop)
}
