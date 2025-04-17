package main

import (
	"github.com/gin-gonic/gin"
	_ "project-api/api"
	"project-api/config"
	"project-api/router"
	srv "project-common"
)

func main() {
	r := gin.Default()
	// Initialize the router
	router.InitRouter(r)
	srv.Run(r, config.AppConf.SC.Name, config.AppConf.SC.Addr, nil)
}
