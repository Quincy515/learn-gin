https://www.jtthink.com/course/146 Go WebSocket实战速学篇

[toc]



### 01. 最基本的 WebSocket 代码

WebSocket 概念

> WebSocket  是 HTML5 提供的一个浏览器与服务器间进行全双工通讯的网络技术。

所谓的全工通信是指：

在同一时刻信息可以进行双向传输，和打电话一样，说的同时也能听，边听边说。

基于 http 协议

```bash
GET /HTTP/1.1
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: xxxxxxxxx(是一个Base64加密的秘钥)
Sec-WebSocket-Version: 13 (告诉服务器 ws 的版本)
Origin: http://xxxx.com (来源)

其中 upgrade websocket 用于告诉服务器此链接需要升级到 websocket。
这意味着你的服务端需要支持 websocket 协议。
```

服务端响应

```bash
HTTP/1.1 101 Switching Protocols
Content-Length: 0
Upgrade: websocket
Sec-Websocket-Accept: ZEs+c+VBk8Aj01+wJGN7Y15796g=
Connection: Upgrade

101 代表协议切换
Sec-WebSocket-Accept 表示服务器同意握手建立链接
接下来就没 http 什么事情了。
```

第三方库

Https://github.com/gorilla/websocket

```bash
go get github.com/gorilla/websocket
```

创建 upgrade 对象

```go 
var upgrade = websocket.Upgrade {
  CheckOrigin: func(r *http.Request) bool {
    return true
  }
}
```

示例代码

```go
package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		client, _ := upgrader.Upgrade(w, r, nil) // 升级
		for {
			err := client.WriteMessage(websocket.TextMessage, []byte("hello"))
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second * 2)
		}
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

使用 websocket 在线测试工具 http://coolaf.com/tool/chattest

连接 ws://localhost:8080/echo ，可以发现连接成功，并接受到服务端发送的消息 "hello"

代码变动 [git commit]()