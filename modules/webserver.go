package modules

// web 服务的启动基于 gin

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 服务端 web 模块
func WebServer() {
	log.Print("已开启【WEB】模式，指令已关闭\n")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("assets/*")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		// data := c.Param("data")
		c.HTML(http.StatusOK, "index.html", nil)
	})
	// v1 := r.Group("/v1")
	// {
	// 	v1.POST("/test")
	// }
	r.Run(":1234")
}
