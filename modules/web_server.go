package modules

import (
	// "embed"
	// "html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

// //go:embed assets
// var static embed.FS

// 服务端 web 模块
func WebServer(s *ini.File) {

	// d, err := os.Open("db.txt") // 打开文件 OpenFile(name, O_RDONLY, 0)
	// if err != nil {
	// 	log.Println("os:", err)
	// }

	// // 关闭文件
	// defer d.Close()

	// b, err := ioutil.ReadFile("db.txt")
	// if err != nil {
	// 	log.Print(err)
	// }

	webPort := 1234

	log.Printf("已开启【WEB】模式，指令已关闭，运行在端口：%s\n", strconv.Itoa(webPort))
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// t, _ := template.ParseFS(tmpl, "templates/*.html")
	// r.SetHTMLTemplate(t)
	r.LoadHTMLGlob("templates/*.html")
	// r.StaticFS("/assets", http.FS(static))
	r.Static("/assets", "./assets")
	r.GET("/", func(c *gin.Context) {
		// data := c.Param("data")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "OnlineChat!",
			"sub":   "一个开源的在线聊天应用程序",
			"page":  "在线聊天",
			"ws":    "ws://" + s.Section("Server").Key("Server").String() + s.Section("Server").Key("Url").String(),
		})
	})
	api := r.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "RUNNING!",
			})
		})

		api.GET("/msg", func(c *gin.Context) {
			// 循环读取
			f, _ := os.Open("db.txt")
			var dbData string
			for {
				buf := make([]byte, 6)
				_, err := f.Read(buf)
				dbData += string(buf)
				if err == io.EOF {
					c.JSON(http.StatusOK, string(dbData))
					break
				}
			}
			f.Close()
		})
	}

	err := r.Run(":" + strconv.Itoa(webPort))
	if err != nil {
		log.Print("gin:", err)
	}
}
