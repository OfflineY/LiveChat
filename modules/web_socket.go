package modules

import (
	"bytes"
	"flag"
	"net/http"

	"log"
	"time"

	"gopkg.in/ini.v1"

	"github.com/gorilla/websocket"
)

const (
	// 允许向对等方写入消息的时间
	writeWait = 10 * time.Second
	// 允许从对等方读取下一条pong消息的时间
	pongWait = 60 * time.Second

	// 将ping发送到此时段的对等节点。必须小于pongWait
	// 目前已经取消此功能但会随着pongWait改变
	// pingPeriod = ----

	// 允许来自对等方的最大消息大小
	maxMessageSize = 1024000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// 给 http 升级为 websocket
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// 其中借鉴了websocket的官方文档中的示例
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	// 设置服务端基本配置
	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Print("ws_conn:", err)
	}
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Print("ws_conn:", err)
		}
		return nil
	})
	for {
		// 设置连接循环
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Websocket内部错误: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		go DbSave(string(message))
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {

	ticker := time.NewTicker((pongWait * 9) / 10)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Print("ws_conn:", err)
			}
			if !ok {
				// 中枢关闭了通道
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Print("ws_conn:", err)
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, err = w.Write(message)
			if err != nil {
				log.Print("ws_write:", err)
			}

			// 将排队的聊天信息添加到当前websocket信息中
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, err = w.Write(newline)
				if err != nil {
					log.Print("ws_write:", err)
				}
				_, err = w.Write(<-c.send)
				if err != nil {
					log.Print("ws_write:", err)
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Print("ws_conn:", err)
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 集线器维护活动客户的集合，并向这些客户广播消息到客户端
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// 启动 web
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// websocket 基础功能启动
func WorkPrimaryServer(s *ini.File) {

	log.Printf("已开启【服务端】模式，此模式是为服务器所准备")
	flag.Parse()
	hub := NewHub()

	go hub.run()
	severUrl := s.Section("Server").Key("Url").String()
	sever := s.Section("Server").Key("Server").String()
	log.Printf("程序已在运行，本地端口: ws://%s%s", sever, severUrl)
	http.HandleFunc(severUrl, func(w http.ResponseWriter, r *http.Request) {

		// serveWs(hub, w, r)

		conn, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
		client.hub.register <- client

		go client.writePump()
		go client.readPump()
	})

	ListenAndServe := sever
	log.Fatal(http.ListenAndServe(ListenAndServe, nil))
}
