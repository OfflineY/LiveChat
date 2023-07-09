package service

import (
	"LiveChat/db"
	"LiveChat/util"
	"LiveChat/ws"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// RoomDetail 房间细节信息包括存储数据库位置
type RoomDetail struct {
	RoomHubs     map[string]*ws.Hub
	Conn         *mongo.Client
	DatabaseName string
}

// RecoverRoom 恢复聊天房间
func RecoverRoom(Room string, dbConn *RoomDetail) {
	hub := dbConn.RoomHubs[Room]
	go hub.RunHub()

	http.HandleFunc("/"+Room, func(w http.ResponseWriter, r *http.Request) {
		ws.WebSocketServer(hub, w, r, dbConn.Conn, dbConn.DatabaseName, Room)
	})

	log.Println("Recover the Room:", Room)
}

// CreateRoom 创建新房间
func CreateRoom(Room string, dbConn *RoomDetail) (bool, string, string) {
	if _, ok := dbConn.RoomHubs[Room]; ok {
		log.Println("Room exist and do not need to be created:", Room, util.MD5(Room))
		return !ok, "", ""
	} else {
		dbConn.RoomHubs[Room] = ws.NewHub()

		roomId := util.MD5(Room)

		db.Add(dbConn.Conn, dbConn.DatabaseName, "groups",
			bson.D{
				{"create_time", time.Now()},
				{"group_name", Room},
				// TODO：改名不改变
				{"group_id", roomId},
			},
		)

		go dbConn.RoomHubs[Room].RunHub()

		http.HandleFunc("/"+roomId, func(w http.ResponseWriter, r *http.Request) {
			ws.WebSocketServer(dbConn.RoomHubs[Room], w, r, dbConn.Conn, dbConn.DatabaseName, Room)
		})
		log.Println("Create the Room:", Room, roomId)
		return !ok, Room, roomId
	}
}
