package src

import (
	"log"
	"os"
	"os/signal"
)

var ServerSigChan chan os.Signal

func init() {
	ServerSigChan = make(chan os.Signal) // 创建信号 chan
}

func ShutDownServer(err error) {
	log.Println(err)
	ServerSigChan <- os.Interrupt // 发送 SIGINT 信号
}

func ServerNotify() {
	signal.Notify(ServerSigChan, os.Interrupt)  // 监听所有信号
	<-ServerSigChan
}
