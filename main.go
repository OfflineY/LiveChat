package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/ini.v1"
)

const (
	// 允许向对等方写入消息的时间。
	writeWait = 10 * time.Second
	// 允许从对等方读取下一条pong消息的时间。
	pongWait = 60 * time.Second
	// 将ping发送到此时段的对等节点。必须小于pongWait。
	pingPeriod = (pongWait * 9) / 10
	// 允许来自对等方的最大消息大小。
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub *Hub
	// websocket连接。
	conn *websocket.Conn
	// 消息的缓冲通道。
	send chan []byte
}

//	Gorilla Websocket Example START
//
// ***********************************************
// readPump 将消息从 websocket 连接泵送到集线器。
// 应用程序在每个连接 goroutine 中运行 readPump 应用程序
// 通过执行 all 确保连接上最多有一个读取器
// 从这里读到 goroutine。
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
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

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

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

//       Gorilla Websocket Example END
// ***********************************************

func ChatUser(url string) {
	// url := "ws://localhost:8080/socket"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("连接服务器失败:", err)
		Enter()
	}
	log.Printf("连接服务器成功")
	// fmt.Print(res)
	for {
		// var msg string
		//
		// err = c.WriteMessage(websocket.TextMessage, []byte(msg))
		// if err != nil {
		// 	fmt.Println(err)
		// }
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[User] %s", message)
	}
}

func bootstrap() {
	setPath := "set.ini"
	s, Seterr := ini.Load(setPath)
	if Seterr != nil {
		log.Print("未找到配置文件开始初始化。\n")
		_, err := os.Create(setPath)
		if err != nil {
			log.Printf("初始化失败(%s)。\n", err)
			Enter()
		}
		f, err := os.OpenFile(setPath, os.O_RDWR, 0600)
		if err != nil {
			log.Printf("初始化失败(%s)。\n", err)
			Enter()
		}
		f.WriteString(`# mode 为设置服务模式, 0 为服务端, 1 为客户端
mode = 0

[Server]
Url = "/socket"
Port = "8080"
# 待补充配置
[User]
Url = "ws://localhost:8080/socket"
# 待补充配置`)
		defer f.Close()
		log.Print("创建配置文件完成...\n")
		log.Print("初始化结束，请更改配置文件，并重启应用。\n")
		fmt.Print(`
		
【注意事项】
		
  + 请在 cmd 或 vscode 终端或一些常驻终端使用，否则可能会出现出现错误闪退的情况。
  + 目前仅处于 BETA 版本，测试中且项目未未完工。
  + 虽已经实现效果但安全性有待测试，仅供学习，请勿用于其他用途。




`)
		Enter()
	} else {
		log.Printf("找到配置文件，跳过初始化过程.\n")
		// 获取 ini 的 mode 数据，0为服务端，1为客户端
		mode := s.Section("").Key("mode").String()
		switch {
		case mode == "0":
			log.Printf("已开启【服务端】模式，此模式是为服务器所准备。\n")
			flag.Parse()
			hub := newHub()
			// 开启另一个 goroutine
			go hub.run()
			severUrl := s.Section("Server").Key("Url").String()
			severPort := s.Section("Server").Key("Port").String()
			log.Printf("程序已在运行，本地端口: http://localhost:%s%s", severPort, severUrl)
			http.HandleFunc(severUrl, func(w http.ResponseWriter, r *http.Request) {
				// serveWs(hub, w, r)
				conn, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					log.Println(err)
					return
				}
				client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
				client.hub.register <- client
				// Allow collection of memory referenced by the caller by doing all work in
				// new goroutines.
				go client.writePump()
				go client.readPump()
			})
			// 在端口启动
			ListenAndServe := "localhost:" + severPort
			log.Fatal(http.ListenAndServe(ListenAndServe, nil))
		case mode == "1":
			log.Print("已开启【客户端】模式，此模式是为使用者所准备。\n")
			ChatUser(s.Section("User").Key("Url").String())
		}
	}
}

func Enter() {
	fmt.Printf("\n\n\n\n按任意键退出...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}

func logo(version string) {
	fmt.Println(`
    _____     _ _         _____ _       _   
   |     |___| |_|___ ___|     | |_ ___| |_     OlineChat
   |  |  |   | | |   | -_|   --|   | .'|  _|    Version ` + version + `
   |_____|_|_|_|_|_|_|___|_____|_|_|__,|_|      Powered by BillyYuan

========================================================================
   `)
}

func main() {
	// 打印 logo
	logo("0.1 BETA")
	// 初始引导程序
	bootstrap()
}
