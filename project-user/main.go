package main

import (
	"github.com/gin-gonic/gin"
	srv "project-common"
	_ "project-user/api"
	"project-user/router"
)

func main() {
	r := gin.Default()
	// Initialize the router
	router.InitRouter(r)
	srv.Run(r, "project-user", ":80")
}
