package main

import (
	"github.com/gin-gonic/gin"
	"log"
	srv "project-common"
	"project-common/logs"
	_ "project-user/api"
	"project-user/router"
)

func main() {
	r := gin.Default()
	// init log
	lc := &logs.LogConfig{
		DebugFileName: "..\\logs\\debug\\project-debug.log",
		InfoFileName:  "..\\logs\\info\\project-info.log",
		WarnFileName:  "..\\logs\\error\\project-error.log",
		MaxSize:       5,
		MaxAge:        28,
		MaxBackups:    3,
	}
	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
	// Initialize the router
	router.InitRouter(r)
	srv.Run(r, "project-user", ":80")
}
