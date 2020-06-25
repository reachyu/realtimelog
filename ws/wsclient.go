package ws

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	_const "reallog/common/const"
	mq "reallog/msgmq"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var clientMaps map[string]*WSClient
var once sync.Once

func GetClientMaps() map[string]*WSClient {
	once.Do(func() {
		clientMaps = make(map[string]*WSClient)
	})
	return clientMaps
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type WSClient struct {
	hub *WSHub
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
	trainId []byte
}

func SendLogsToWeb(trainId string,msg string) {
	clientMaps := GetClientMaps()
	c := clientMaps[trainId]

	if c == nil{
		return
	}
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	message := bytes.TrimSpace(bytes.Replace([]byte(msg), newline, space, -1))
	message = []byte(trainId + "&" + msg)
	fmt.Println("websocket读取到的消息====="+string(message))
	c.hub.broadcast <- []byte(message)

}

func (c *WSClient) readPump() {
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
                // todo:断开连接后，会到这里来
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		message = []byte(string(c.trainId) + "&" + string(message))
		fmt.Println("websocket读取到的消息====="+string(message))
		c.hub.broadcast <- []byte(message)
	}
}

func (c *WSClient) writePump() {
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
			if err := w.Close(); err != nil {
				log.Printf("error: %v", err)
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

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *WSHub, c *gin.Context) {
	GetClientMaps()
	trainId := c.Param("trainId")
	// 将网络请求变为websocket
	var upgrader = websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("websocket接收到前端的trainId======"+trainId)
	client := &WSClient{hub: hub, conn: conn, send: make(chan []byte, 256), trainId: []byte(trainId)}
	client.hub.register <- client

	clientMaps[trainId] = client

	// 监听日志消息
	go mq.ConsumeMsgRout(_const.RABBITMQ_ROUT_EXCHANGE_NAME,_const.RABBITMQ_ROUT_QUEUE + trainId,_const.RABBITMQ_ROUT_ROUTING_KEY + trainId,SendLogsToWeb)
	//go mq.ConsumeMsgFanout(_const.RABBITMQ_FANOUT_EXCHANGE_NAME,_const.RABBITMQ_FANOUT_QUEUE + trainId,SendLogsToWeb)

	// kafka：相同的group.id的消费者将视为同一个消费者组, 每个消费者都需要设置一个组id,
	// 每条消息只能被 consumer group 中的一个 Consumer 消费,但可以被多个 consumer group 消费
    //go mq.ConsumeMsgKafka(_const.KAFKA_LOGS_GROUP +  trainId,_const.KAFKA_LOGS_TOPIC + trainId,SendLogsToWeb)

	//go mq.TestKafkaComs(_const.KAFKA_LOGS_GROUP + strconv.FormatInt(rand.Int63n(10000000),10),_const.KAFKA_LOGS_TOPIC + trainId,SendLogsToWeb)


	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	//go client.readPump()
}
