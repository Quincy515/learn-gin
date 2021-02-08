package core

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// ClientMap 外部公共使用
var ClientMap *ClientMapStruct

func init() {
	ClientMap = &ClientMapStruct{}
}

type ClientMapStruct struct {
	data sync.Map // key 是客户端 IP  Value 就是 WebSocket 连接对象
}

func (c *ClientMapStruct) Store(key string, conn *websocket.Conn) {
	c.data.Store(key, conn)
}

//向所有客户端 发送消息
func (c *ClientMapStruct) SendAll(msg string) {
	c.data.Range(func(key, value interface{}) bool {
		err := value.(*websocket.Conn).WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println(err)
		}
		return true
	})
}
