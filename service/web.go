package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// RunWeb 运行Web应用
func RunWeb(port string, c *DatabaseConn, r *Rooms) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// 中间件 允许跨域
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	api := router.Group("/api")
	{
		api.GET("/", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"data": nil,
				"msg":  "LiveChat api v3.0",
			})
		})

		api.GET("/settings", func(context *gin.Context) {
			siteInfo(context, c)
		})

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
					context.JSON(http.StatusBadRequest, gin.H{
						"data": nil,
						"msg":  err.Error(),
					})
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
				context.JSON(http.StatusOK, gin.H{
					"data": data,
					"msg":  msg,
				})
			})

			groups.GET("/:id/info", func(context *gin.Context) {
				groupId := context.Param("id")
				context.JSON(http.StatusOK, gin.H{
					"data": gin.H{
						"group_id": groupId,
					},
					"msg": http.StatusOK,
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
