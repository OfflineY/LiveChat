package service

import (
	"LiveChat/util"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RunWeb 运行Web应用
func RunWeb(port string, c *DatabaseConn, r *Rooms) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(util.AllowCORS())
	api := router.Group("/api")
	{
		api.GET("/", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{"data": nil, "msg": "__BETA_3.0__"})
		})
		api.GET("/settings", func(context *gin.Context) {
			siteInfo(context, c)
		})
		users := api.Group("users")
		{
			users.POST("/login", func(context *gin.Context) {
				type user struct {
					UserName string `json:"user_name"`
					Password string `json:"password"`
				}
				var userData user
				err := context.BindJSON(&userData)
				if err != nil {
					context.JSON(http.StatusBadRequest, gin.H{
						"data": nil,
						"msg":  err.Error(),
					})
				} else {
					status, err := loginAuth(c, userData.UserName, userData.Password)
					context.JSON(http.StatusBadRequest, gin.H{
						"data": gin.H{
							"login_status": status,
							"user":         gin.H{"name": userData.UserName, "password_md5": util.MD5(userData.Password)},
						}, "msg": err,
					})
				}
			})
			users.POST("/register", func(context *gin.Context) {
				type user struct {
					UserName string `json:"user_name"`
					Password string `json:"password"`
				}
				var userData user
				err := context.BindJSON(&userData)
				if err != nil {
					context.JSON(http.StatusBadRequest, gin.H{"data": nil, "msg": err.Error()})
				} else {
					status, err := registerAuth(c, userData.UserName, userData.Password)
					context.JSON(http.StatusBadRequest, gin.H{
						"data": gin.H{
							"register_status": status,
							"user":            gin.H{"name": userData.UserName, "password_md5": util.MD5(userData.Password)},
						}, "msg": err,
					})
				}
			})
		}
		groups := api.Group("/groups")
		{
			groups.GET("/", func(context *gin.Context) {
				groupsInfo(context, c)
			})
			groups.POST("/create", func(context *gin.Context) {
				type group struct {
					GroupName string `json:"name"`
				}
				var cg group
				err := context.BindJSON(&cg)
				if err != nil {
					context.JSON(http.StatusBadRequest, gin.H{"data": nil, "msg": err.Error()})
				} else {
					createGroup(context, c, r, cg.GroupName)
				}
			})
			groups.GET("/search", func(context *gin.Context) {
				groupId := context.Query("id")
				groupName := context.Query("name")

				data, err := findGroup(c, groupId, groupName)

				var msg error = nil
				if err != nil {
					msg = err
				}

				context.JSON(http.StatusOK, gin.H{"data": data, "msg": msg})
			})
			groups.GET("/:id/info", func(context *gin.Context) {
				groupId := context.Param("id")

				context.JSON(http.StatusOK, gin.H{
					"data": gin.H{"group_id": groupId},
					"msg":  http.StatusOK,
				})
			})
			groups.GET("/:id/messages", func(context *gin.Context) {
				id := context.Param("id")
				historicalMessages(context, c, id)
			})
		}
	}

	err := router.Run(port)
	if err != nil {
		log.Panicln(err)
	}
}
