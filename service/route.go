package service

import (
	"LiveChat/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// siteInfo 网站信息路由
func siteInfo(context *gin.Context, dbConn *RoomDetail) {
	siteName, err := db.Find(dbConn.Conn, dbConn.DatabaseName, "settings", bson.M{
		"name": "site_name",
	})
	var msg error = nil
	if err != nil {
		msg = err
	}
	context.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"site_name": siteName[0]["value"],
		},
		"msg": msg,
	})
}

// groupsInfo 房间路由
func groupsInfo(context *gin.Context, dbConn *RoomDetail) {
	data, err := db.Find(
		dbConn.Conn,
		dbConn.DatabaseName,
		"groups",
		bson.M{},
	)
	var msg error = nil
	if err != nil {
		msg = err
	}
	context.JSON(http.StatusOK, gin.H{
		"data": data,
		"msg":  msg,
	})
}

func findGroup(dbConn *RoomDetail, groupId string, groupName string) ([]bson.M, error) {
	if groupId == "" {
		data, err := db.Find(dbConn.Conn, dbConn.DatabaseName, "groups",
			bson.M{
				"group_name": bson.M{"$regex": groupName, "$options": "im"},
			},
		)
		return data, err
	}
	if groupName == "" {
		data, err := db.Find(dbConn.Conn, dbConn.DatabaseName, "groups",
			bson.M{
				"group_id": groupId,
			},
		)
		return data, err
	}
	data, err := db.Find(dbConn.Conn, dbConn.DatabaseName, "groups",
		bson.M{
			"group_id":   groupId,
			"group_name": groupName,
		},
	)
	return data, err
}

func createGroup(context *gin.Context, dbConn *RoomDetail, groupName string) {
	isOk, name, id := CreateRoom(groupName, dbConn)
	if isOk {
		context.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"group_name": name,
				"group_id":   id,
			},
			"msg": nil,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"data": nil,
			"msg":  "error",
		})
	}
}

// historicalMessages 历史消息获取
func historicalMessages(context *gin.Context, dbConn *RoomDetail, groupId string) {
	data, err := db.Find(dbConn.Conn, dbConn.DatabaseName, "messages", bson.M{
		"group_id": groupId,
	})
	var msg error = nil
	if err != nil {
		msg = err
	}
	context.JSON(http.StatusOK, gin.H{
		"data": data,
		"msg":  msg,
	})
}
