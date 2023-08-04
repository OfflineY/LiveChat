package ws

type Hub struct {
	// clients 注册表
	clients map[*Client]bool
	// broadcast 广播消息管道
	broadcast chan []byte
	// register 注册管道
	register chan *Client
	// unregister 注销管道
	unregister chan *Client
}

// NewHub 创建监听管道
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// RunHub 创建管道监听，处理管道操作
func (hub *Hub) RunHub() {
	for {
		select {
		// 注册管道
		case client := <-hub.register:
			hub.clients[client] = true
		// 注销管道
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}
		// 消息广播管道
		case message := <-hub.broadcast:
			for client := range hub.clients {
				select {
				case client.send <- message:
				// 如果管道不能立即写入数据，就认为该 client 出故障
				default:
					close(client.send)
					delete(hub.clients, client)
				}
			}
		}
	}
}
