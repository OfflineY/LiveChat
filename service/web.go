package service

import (
	"LiveChat/util"
	"log"

	"github.com/gin-gonic/gin"
)

// RunWeb 运行Web应用
func RunWeb(port string, c *DatabaseConn, r *Rooms) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(util.AllowCORS())
	createRouter(router, c, r)

	err := router.Run(port)
	if err != nil {
		log.Panicln(err)
	}
}
