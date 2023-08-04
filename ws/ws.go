package ws

import (
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

// upGrader 升级为 websocket
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Group 群组信息
type Group struct {
	DatabaseConn *mongo.Client
	DatabaseName string
	RoomName     string
	RoomId       string
}

// WebSocketServer 初始化 websocket 连接
func WebSocketServer(hub *Hub, w http.ResponseWriter, r *http.Request, Conn *mongo.Client, DatabaseName string, RoomName string) {
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client

	go client.WritePump()
	go client.ReadPump(&Group{DatabaseConn: Conn, DatabaseName: DatabaseName, RoomName: RoomName, RoomId: RoomName})
}
