package modules

// 负责 websocket 的后端处理

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
	maxMessageSize = 10240
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
	hub *Hub
	// websocket连接
	conn *websocket.Conn
	// 消息的缓冲通道
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
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		// 设置连接循环
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WS错误: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

// writePump 将消息从集线器抽到websocket连接上，一个运行writePump的Goroutine为每个连接启动
// 应用程序确保一个连接最多只有一个写入者，从这个goroutine中执行所有的写操作
func (c *Client) writePump() {

	ticker := time.NewTicker((pongWait * 9) / 10)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// 中枢关闭了通道
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将排队的聊天信息添加到当前websocket信息中
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 集线器维护活动客户的集合，并向这些客户广播消息到客户端
type Hub struct {

	// 注册客户
	clients map[*Client]bool

	// 来自客户的呼入信息
	broadcast chan []byte

	// 登记来自客户的请求
	register chan *Client

	// 取消对客户请求的登记
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
	// mode 0 服务端
	log.Printf("已开启【服务端】模式，此模式是为服务器所准备")
	flag.Parse()
	hub := NewHub()

	// 开启另一个 goroutine 启动 ws
	go hub.run()
	severUrl := s.Section("Server").Key("Url").String()
	severPort := s.Section("Server").Key("Port").String()
	log.Printf("程序已在运行，本地端口: http://localhost:%s%s", severPort, severUrl)
	http.HandleFunc(severUrl, func(w http.ResponseWriter, r *http.Request) {

		// serveWs(hub, w, r)

		conn, err := Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
		client.hub.register <- client

		// 给多个 goroutine 启动
		go client.writePump()
		go client.readPump()
	})

	// 在主要端口启动 severPort -> ini
	ListenAndServe := "localhost:" + severPort
	log.Fatal(http.ListenAndServe(ListenAndServe, nil))
}