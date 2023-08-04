package service

import (
	"LiveChat/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// siteInfo 网站信息路由
func siteInfo(context *gin.Context, c *DatabaseConn) {
	siteName, err := db.Find(&db.DatabaseConn{Conn: c.Conn, Name: c.Name}, "settings",
		bson.M{"name": "site_name"},
	)

	var msg error = nil
	if err != nil {
		msg = err
	}

	context.JSON(http.StatusOK, gin.H{
		"data": gin.H{"site_name": siteName[0]["value"]}, "msg": msg,
	})
}

// groupsInfo 房间路由
func groupsInfo(context *gin.Context, c *DatabaseConn) {
	data, err := db.Find(&db.DatabaseConn{Conn: c.Conn, Name: c.Name}, "groups",
		bson.M{},
	)

	var msg error = nil
	if err != nil {
		msg = err
	}

	context.JSON(http.StatusOK, gin.H{"data": data, "msg": msg})
}

// findGroup 寻找群组
func findGroup(c *DatabaseConn, groupId string, groupName string) ([]bson.M, error) {
	if groupId == "" {
		data, err := db.Find(&db.DatabaseConn{Conn: c.Conn, Name: c.Name}, "groups",
			bson.M{"group_name": bson.M{"$regex": groupName, "$options": "im"}},
		)
		return data, err
	}
	if groupName == "" {
		data, err := db.Find(&db.DatabaseConn{Conn: c.Conn, Name: c.Name}, "groups",
			bson.M{"group_id": groupId},
		)
		return data, err
	}
	data, err := db.Find(&db.DatabaseConn{Conn: c.Conn, Name: c.Name}, "groups",
		bson.M{"group_id": groupId, "group_name": groupName},
	)
	return data, err
}

// createGroup 创造群组
func createGroup(context *gin.Context, c *DatabaseConn, r *Rooms, groupName string) {
	isOk, name, id := r.CreateRoom(groupName, c)
	if isOk {
		context.JSON(http.StatusOK, gin.H{
			"data": gin.H{"group_name": name, "group_status": isOk, "group_id": id},
			"msg":  nil,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{"data": nil, "group_status": isOk, "msg": "error"})
	}
}

// historicalMessages 历史消息获取
func historicalMessages(context *gin.Context, c *DatabaseConn, id string) {
	data, err := db.Find(&db.DatabaseConn{Conn: c.Conn, Name: c.Name}, "messages", bson.M{
		"group_id": id,
	})

	var msg error = nil
	if err != nil {
		msg = err
	}

	context.JSON(http.StatusOK, gin.H{"data": data, "msg": msg})
}
