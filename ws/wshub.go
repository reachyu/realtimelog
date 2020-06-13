package ws

import "strings"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type WSHub struct {
	// Registered clients.
	clients map[*WSClient]bool
	// Inbound messages from the clients.
	broadcast chan []byte
	// Register requests from the clients.
	register chan *WSClient
	// Unregister requests from clients.
	unregister chan *WSClient
	// 唯一id key:client value:唯一id
	trainId map[*WSClient]string
}

// NewHub .
func NewHub() *WSHub {
	return &WSHub{
		broadcast:  make(chan []byte),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
		clients:    make(map[*WSClient]bool),
		trainId:    make(map[*WSClient]string),
	}
}

// Run .
func (h *WSHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true                 // 注册client端
			h.trainId[client] = string(client.trainId) // 给client端添加唯一id
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.trainId, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				// 使用“&”对message进行message切割 获取唯一id
				// 向信息所属的训练内的所有client 内添加send
				// msg[0]为唯一id msg[1]为打印内容
				msg := strings.Split(string(message), "&")
				if string(client.trainId) == msg[0] {
					select {
					case client.send <- []byte(msg[1]):
					default:
						close(client.send)
						delete(h.clients, client)
						delete(h.trainId, client)
					}
				}
			}
		}
	}
}
