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

// DatabaseConn 存储数据库信息
type DatabaseConn struct {
	Conn *mongo.Client
	Name string
}

// Rooms 储存房间信息
type Rooms struct {
	RoomHubs map[string]*ws.Hub
}

// RecoverRoom recover groups from db
func (r *Rooms) RecoverRoom(Room string, c *DatabaseConn) {
	hub := r.RoomHubs[Room]
	go hub.RunHub()

	http.HandleFunc("/"+Room, func(w http.ResponseWriter, r *http.Request) {
		ws.WebSocketServer(hub, w, r, c.Conn, c.Name, Room)
	})

	log.Println("Recover the Room:", Room)
}

// CreateRoom create new groups
func (r *Rooms) CreateRoom(Room string, c *DatabaseConn) (bool, string, string) {
	if _, ok := r.RoomHubs[Room]; ok {
		log.Println("Room exist and do not need to be created:", Room, util.MD5(Room))
		return !ok, "", ""
	} else {
		r.RoomHubs[Room] = ws.NewHub()

		roomId := util.MD5(Room)

		db.Add(&db.DatabaseConn{
			Conn: c.Conn,
			Name: c.Name,
		}, "groups",
			bson.D{
				{Key: "create_time", Value: time.Now()},
				{Key: "group_name", Value: Room},
				{Key: "group_id", Value: roomId},
			},
		)

		go r.RoomHubs[Room].RunHub()

		http.HandleFunc("/"+roomId, func(w http.ResponseWriter, req *http.Request) {
			ws.WebSocketServer(r.RoomHubs[Room], w, req, c.Conn, c.Name, Room)
		})
		log.Println("Create the Room:", Room, roomId)
		return !ok, Room, roomId
	}
}
