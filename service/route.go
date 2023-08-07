package service

import (
	"LiveChat/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

// createRouter 创建总路由组
func createRouter(router *gin.Engine, conn *DatabaseConn, roomsGroup *Rooms) {
	api := router.Group("/api")
	{
		api.GET("/", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{"data": nil, "msg": "__BETA_3.0__"})
		})
		api.GET("/settings", func(context *gin.Context) {
			getSiteSettings(context, conn)
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
					status, err := loginAuth(conn, userData.UserName, userData.Password)
					context.JSON(http.StatusBadRequest, gin.H{
						"data": gin.H{
							"login_status": status,
							"user":         gin.H{"name": userData.UserName, "password_md5": util.MD5(userData.Password)},
						}, "msg": err,
					})
				}
			})
			users.POST("/register", func(context *gin.Context) {
				type userData struct {
					UserName string `json:"user_name"`
					Password string `json:"password"`
				}
				var ud userData
				err := context.BindJSON(&ud)
				if err != nil {
					context.JSON(http.StatusBadRequest, gin.H{"data": nil, "msg": err.Error()})
				} else {
					status, err := registerAuth(conn, ud.UserName, ud.Password)
					context.JSON(http.StatusBadRequest, gin.H{
						"data": gin.H{
							"register_status": status,
							"user":            gin.H{"name": ud.UserName, "password_md5": util.MD5(ud.Password)},
						}, "msg": err,
					})
				}
			})
		}
		groups := api.Group("/groups")
		{
			groups.GET("/", func(context *gin.Context) {
				getGroupsList(context, conn)
			})
			groups.POST("/create", func(context *gin.Context) {
				type createGroupStruct struct {
					GroupName string `json:"name"`
				}
				var cgs createGroupStruct
				err := context.BindJSON(&cgs)
				if err != nil {
					context.JSON(http.StatusBadRequest, gin.H{"data": nil, "msg": err.Error()})
				} else {
					createGroup(context, conn, roomsGroup, cgs.GroupName)
				}
			})
			groups.GET("/search", func(context *gin.Context) {
				groupId := context.Query("id")
				groupName := context.Query("name")

				data, err := findGroup(conn, groupId, groupName)

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
				getHistoricalMessages(context, conn, id)
			})
		}
	}
}
