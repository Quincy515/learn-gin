package core

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

// ClientMap 外部公共使用
var ClientMap *ClientMapStruct

func init() {
	ClientMap = &ClientMapStruct{}
}

type ClientMapStruct struct {
	data sync.Map // key 是客户端 IP  Value 就是 WebSocket 连接对象
}

func (c *ClientMapStruct) Store(conn *websocket.Conn) {
	wsClient := NewWsClient(conn)
	c.data.Store(conn.RemoteAddr().String(), wsClient)
	go wsClient.Ping(time.Second * 1)
}

// SendAll 向所有客户端 发送消息
func (c *ClientMapStruct) SendAll(msg string) {
	c.data.Range(func(key, value interface{}) bool {
		err := value.(*WsClient).conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println(err)
		}
		return true
	})
}

// Remove
func (c *ClientMapStruct) Remove(conn *websocket.Conn) {
	c.data.Delete(conn.RemoteAddr().String())
}
