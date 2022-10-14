package modules

import (
	// "embed"
	// "html/template"
	"io"
	"log"
	"net/http"
	"os"

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

	webTitle := "OnlineChat!"
	webSub := "一个开源的在线聊天应用程序"
	webPort := s.Section("Server").Key("WebPort").String()

	log.Printf("已开启【WEB】模式，指令已关闭，运行在端口：%s\n", webPort)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// 设置中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// t, _ := template.ParseFS(tmpl, "templates/*.html")
	// r.SetHTMLTemplate(t)
	r.LoadHTMLGlob("templates/*.html")

	// r.StaticFS("/assets", http.FS(static))
	r.Static("/assets", "./assets")

	r.GET("/", CheckLogin, func(c *gin.Context) {
		name, err := c.Cookie("cookie_username")
		if err != nil {
			log.Print("get_cookie:", err)
		}
		// data := c.Param("data")
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title":   webTitle,
			"sub":     webSub,
			"page":    "在线聊天",
			"routher": "index.html",
			"ws":      "ws://" + s.Section("Server").Key("Server").String() + s.Section("Server").Key("Url").String(),
			"name":    name,
		})
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": webTitle,
			"sub":   webSub,
			"page":  "登录",
		})
	})
	api := r.Group("/api")
	{
		// status
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": "RUNNING!",
			})
		})
		api.POST("/login", func(c *gin.Context) {
			name := c.PostForm("username")
			// *******
			// cookie = name 为临时举动
			// *******
			cookie := name
			c.SetCookie("cookie_username", cookie, 3600, "/", "", false, true)
			c.Redirect(http.StatusMovedPermanently, "/")
			// 向后台加入统计
			AddUser(name)
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
	err := r.Run(":" + webPort)
	if err != nil {
		log.Print("gin:", err)
	}
}

// 中间件检查 cookie 登录情况
func CheckLogin(c *gin.Context) {
	_, err := c.Cookie("cookie_username")
	if err != nil {
		log.Print("cookie_check:", err)
		c.Redirect(http.StatusMovedPermanently, "/login")
		// cookie = "NotSet"
		// c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
	}
	// log.Print(name, " join chat")
}
