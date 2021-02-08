package core

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type WsClient struct {
	conn      *websocket.Conn
	readChan  chan *WsMessage //读队列 (chan)
	closeChan chan byte       // 失败队列
}

func NewWsClient(conn *websocket.Conn) *WsClient {
	return &WsClient{conn: conn, readChan: make(chan *WsMessage), closeChan: make(chan byte)}
}

func (w *WsClient) Ping(waittime time.Duration) {
	for {
		time.Sleep(waittime)
		err := w.conn.WriteMessage(websocket.TextMessage, []byte("ping"))
		if err != nil {
			ClientMap.Remove(w.conn)
			return
		}
	}
}

func (w *WsClient) ReadLoop() {
	for {
		t, data, err := w.conn.ReadMessage()
		if err != nil {
			w.conn.Close()
			ClientMap.Remove(w.conn)
			w.closeChan <- 1
			break
		}
		w.readChan <- NewWsMessage(t, data)
	}
}

func (w *WsClient) HandlerLoop() {
loop:
	for {
		select {
		case msg := <-w.readChan:
			fmt.Println(string(msg.MessageData))
		case <-w.closeChan:
			log.Println("已经关闭")
			break loop
		}
	}
}
