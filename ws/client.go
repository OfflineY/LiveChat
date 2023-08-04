package ws

import (
	"LiveChat/db"
	"LiveChat/util"
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// writeWait 最大等待写入时限
	writeWait = 10 * time.Second
	// pongWait 心跳包最大等待时限
	pongWait = 60 * time.Second
	// pingPeriod 发送心跳包间隔时间
	pingPeriod = (pongWait * 9) / 10
	// maxMessageSize 最大消息长度
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	// hub 创建管道
	hub *Hub
	// conn 创建websocket连接
	conn *websocket.Conn
	// send 消息管道
	send chan []byte
}

// UserMessage 用户消息结构体
type UserMessage struct {
	UserName string `json:"user_name"`
	MsgType  string `json:"msg_type"`
	Url      string `json:"url"`
	Msg      string `json:"msg"`
}

// ReadPump 接收器
func (client *Client) ReadPump(g *Group) {
	// 结束长连接处理
	defer func() {
		client.hub.unregister <- client
		err := client.conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	client.conn.SetReadLimit(maxMessageSize)
	err := client.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Println(err)
	}
	client.conn.SetPongHandler(
		func(string) error {
			err := client.conn.SetReadDeadline(time.Now().Add(pongWait))
			if err != nil {
				log.Println(err)
			}
			return nil
		})
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return
				//log.Println("ConnClose: UnexpectedClose:", err)
			}
			break
		}
		message = bytes.TrimSpace(
			bytes.Replace(message, newline, space, -1),
		)
		var msg UserMessage
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println(err)
		}
		db.Add(&db.DatabaseConn{Conn: g.DatabaseConn, Name: g.DatabaseName}, "messages", bson.D{
			{Key: "group_name", Value: g.RoomName}, {Key: "group_id", Value: util.MD5(g.RoomName)},
			{Key: "user_name", Value: msg.UserName}, {Key: "send_time", Value: time.Now()},
			{Key: "msg_type", Value: msg.MsgType}, {Key: "url", Value: msg.Url}, {Key: "msg", Value: msg.Msg},
		})
		client.hub.broadcast <- message
	}
}

// WritePump 写入器
func (client *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	// 结束长连接处理
	defer func() {
		ticker.Stop()
		err := client.conn.Close()
		if err != nil {
			return
		}
	}()
	for {
		select {
		case message, ok := <-client.send:
			err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Println(err)
			}
			if !ok {
				err := client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					return
				}
				return
			}
			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err)
			}
			_, err = w.Write(message)
			if err != nil {
				log.Println(err)
			}
			n := len(client.send)
			for i := 0; i < n; i++ {
				_, err := w.Write(newline)
				if err != nil {
					log.Println(err)
				}
				_, err = w.Write(<-client.send)
				if err != nil {
					log.Println(err)
				}
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Println(err)
			}
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
			}
		}
	}
}
