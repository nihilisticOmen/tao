package main

import (
	"github.com/gin-gonic/gin"
	srv "tao.com/project-common"
	_ "tao.com/project-user/api"
	"tao.com/project-user/router"
)

func main() {
	r := gin.Default()
	// Initialize the router
	router.InitRouter(r)
	srv.Run(r, "project-user", ":80")
}
