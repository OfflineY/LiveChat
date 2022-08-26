package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/ini.v1"
)

// TODO:列入ini 为可变内容
// 2022.8：pingPeriod 省略
const (

	// 允许向对等方写入消息的时间。
	writeWait = 10 * time.Second

	// 允许从对等方读取下一条pong消息的时间。
	pongWait = 60 * time.Second

	// 将ping发送到此时段的对等节点。必须小于pongWait。
	// pingPeriod = (pongWait * 9) / 10

	// 允许来自对等方的最大消息大小。
	maxMessageSize = 10240
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// 给 http 升级为 websocket
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

// writePump 将消息从集线器抽到websocket连接上。
// 一个运行writePump的Goroutine为每个连接启动。
// 应用程序确保一个连接最多只有一个写入者。
// 从这个goroutine中执行所有的写操作。
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
				// 中枢关闭了通道。
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将排队的聊天信息添加到当前websocket信息中。
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

// 集线器维护活动客户的集合，并向这些客户广播消息到客户端。
type Hub struct {

	// 注册客户。
	clients map[*Client]bool

	// 来自客户的呼入信息。
	broadcast chan []byte

	// 登记来自客户的请求。
	register chan *Client

	// 取消对客户请求的登记。
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

// 启动 web
func (h *Hub) run() {

	for {

		// select 判断
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

// 客户端获取返回数据模块
func ChatUser(url string) {

	// url := "ws://localhost:8080/socket"
	// url 获取 ini 里面的内容
	// 连接 url 里的服务器
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("连接服务器失败:", err)
		Enter()
	}

	log.Printf("连接服务器成功")

	// 设置自己的用户名
	fmt.Print("请输入你的用户名: ")
	inputReader := bufio.NewReader(os.Stdin)
	clientName, _ := inputReader.ReadString('\n')

	// 新用户加入全群聊通知
	err = c.WriteMessage(websocket.TextMessage, []byte("【系统消息】 新的用户加入："+clientName))
	if err != nil {
		log.Printf("发生错误：%s", err)
	}

	go sendMsg(c, clientName)

	// 重复监听服务端发送过来的数据
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}

		// log形式打印出来
		log.Printf("%s", message)
	}
}

func sendMsg(c *websocket.Conn, clientName string) {
	for {
		// fmt.Print("Send> ")
		inputReader := bufio.NewReader(os.Stdin)
		sendMsg, _ := inputReader.ReadString('\n')

		// 简化服务端，在客户端进行数据处理
		// 拼贴用户名和发送的消息
		// msg := "【" + clientName + "】" + sendMsg
		msg := sendMsg
		// fmt.Print(msg)

		// 处理完成发送至服务端
		err := c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			fmt.Println(err)
		}
	}
}

// 开启时引导程序
// 包括：初始化、判断配置文件、启动服务
func bootstrap() {
	setPath := "set.ini"

	// 检测配置文件的存在性
	s, Seterr := ini.Load(setPath)

	// 无 初始化，创建ini
	// 有 跳过直接开始
	if Seterr != nil {

		// 初始化创建文件
		log.Print("未找到配置文件开始初始化\n")
		_, err := os.Create(setPath)
		if err != nil {
			log.Printf("初始化失败(%s)\n", err)
			Enter()
		}

		// 打开文件写入
		f, err := os.OpenFile(setPath, os.O_RDWR, 0600)
		if err != nil {
			log.Printf("初始化失败(%s)\n", err)
			Enter()
		}

		// 初始化写入ini
		f.WriteString(`# mode 为设置服务模式, 0 为服务端, 1 为客户端
mode = 0

# 开启服务端模式则必须设置此内容，客户端情况下删去保留没有影响
[Server]
# Url可以当作密码或者标签
Url = "/socket"
Port = "8080"
# 待补充配置
# 客户端模式必须设置否则闪退，举例ws://【ServerIp】:【Port】/【Url】
[User]
Url = "ws://localhost:8080/socket"
# 待补充配置`)

		// 写入完成后关闭文件
		defer f.Close()

		// 省略（感觉log写的太多了）
		// log.Print("创建配置文件完成...\n")

		log.Print("初始化结束，请更改配置文件，并重启应用\n")

		fmt.Print(`
		


【注意事项】
		
  + 请在 cmd 或 vscode 终端或一些常驻终端使用，否则可能会出现出现错误闪退的情况。

  + 目前仅处于开发中且项目未完全完工，如有错误请提出。

  + 虽已经实现效果但安全性有待测试，仅供学习，请勿用于其他用途。



`)
		Enter()

	} else {

		// 跳过初始化
		log.Printf("找到配置文件，跳过初始化过程\n")

		// 获取 ini 的 mode 数据，0为服务端，1为客户端
		mode := s.Section("").Key("mode").String()
		switch {
		case mode == "0":

			// mode 0 服务端
			log.Printf("已开启【服务端】模式，此模式是为服务器所准备\n")
			flag.Parse()
			hub := newHub()

			// 开启另一个 goroutine 启动 ws
			go hub.run()
			severUrl := s.Section("Server").Key("Url").String()
			severPort := s.Section("Server").Key("Port").String()
			log.Printf("程序已在运行，本地端口: http://localhost:%s%s", severPort, severUrl)
			// 服务端指令控制
			go serveCommand()
			http.HandleFunc(severUrl, func(w http.ResponseWriter, r *http.Request) {

				// serveWs(hub, w, r)

				conn, err := upgrader.Upgrade(w, r, nil)
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

		case mode == "1":

			// mode 1 客户端
			// 文字在客户端处理完成后发回服务端返回给其他用户
			// 可以尽可能减少服务器负担
			log.Print("已开启【客户端】模式，此模式是为使用者所准备。\n")

			// 检查版本的更新
			// 启动一个新的 goroutine 检查，防止github连接失败一直卡在连接步骤
			go CheckNew(version)
			ChatUser(s.Section("User").Key("Url").String())
		}
	}
}

// 服务端指令控制
func serveCommand() {
	fmt.Print("\n")
	for {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Print(">>> ")
		command, _ := inputReader.ReadString('\n')
		fmt.Printf("err command: %s", command)
	}
}

// 按任意键退出
func Enter() {
	fmt.Printf("\n\n\n\n按任意键退出...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}

// 打印 logo
func logo(version string) {
	fmt.Println(`
    _____     _ _         _____ _       _   
   |     |___| |_|___ ___|     | |_ ___| |_     OnlineChat
   |  |  |   | | |   | -_|   --|   | .'|  _|    Version ` + version + `
   |_____|_|_|_|_|_|_|___|_____|_|_|__,|_|      Powered by BillyYuan

========================================================================
   `)
}

// 设置全局 version 变量
var version string = "1.0"

func CheckNew(version string) {
	log.Print("正在连接 Github 检查更新...")

	// 用 gist 这个来记录更新
	url := "https://raw.githubusercontent.com/OfflineY/OfflineY/main/online-chat-version"

	// 将要读取的网站放入到 get 方法中
	resp, err := http.Get(url)
	if err != nil {
		log.Println("检查更新失败，无法连接 Github")
	} else {
		// 读取 body 里面的内容
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("检查更新失败，无法解析返回数据")
		}
		// 关闭链接
		defer resp.Body.Close()
		data := string(bytes)
		// 判断新的版本提示更新or不更新
		if data != version {
			log.Printf("检查更新完成，已经是最新版本：%s", data)
		} else {
			log.Printf("检查更新完成，有新的可更新版本：%s", data)
		}
	}
}

func main() {
	// 打印 logo
	logo(version)
	// 初始引导程序
	bootstrap()
}
