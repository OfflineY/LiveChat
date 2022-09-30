package modules

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 服务端 web 模块
func WebServer() {

	b, err := ioutil.ReadFile("db.txt")
	if err != nil {
		log.Print(err)
	}

	webPort := 1234

	log.Printf("已开启【WEB】模式，指令已关闭，运行在%s\n", strconv.Itoa(webPort))
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	r.GET("/", func(c *gin.Context) {
		// data := c.Param("data")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "OnlineChat!",
			"sub":   "一个开源的在线聊天应用程序",
			"page":  "在线聊天",
		})
	})
	api := r.Group("/api")
	{
		api.GET("/past_msg", func(c *gin.Context) {
			c.JSON(http.StatusOK, string(b))
		})
		// v1 := r.Group("/v1")
		// {
		// 	v1.POST("/test")
		// }
	}

	r.Run(":" + strconv.Itoa(webPort))
}
