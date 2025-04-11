package main

import (
	"github.com/gin-gonic/gin"
	_ "tao.com/project-user/api"
	"tao.com/project-user/router"
	srv "test.com/project-common"
)

func main() {
	r := gin.Default()
	// Initialize the router
	router.InitRouter(r)
	srv.Run(r, "project-user", ":80")
}
