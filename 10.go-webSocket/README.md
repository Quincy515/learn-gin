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

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/05de031f80ae38126adcfb5f12bbd6e15f0e3d4c#diff-f05c20371760668299b744a5d25fd5e4c2e3c7fb1621bec98c95a759672fa6e7R1)

### 02. JS脚本简学、保存客户端对象、代码封装

复制下面的 html 到文件 `src/htmls/a.html` 中

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>客户端A</title>

</head>
<body>
 <div>
    <div id="message" style="width: 500px;height:300px;border:solid 1px gray;overflow:auto">

    </div>

 </div>
 <script>
     var ws = new WebSocket("ws://localhost:8080/echo");
     ws.onopen = function(){
         //当WebSocket创建成功时，触发onopen事件
         console.log("open");
         ws.send("i am user-a"); //将消息发送到服务端
     }
     ws.onmessage = function(e){
         //当客户端收到服务端发来的消息时，触发onmessage事件，参数e.data包含server传递过来的数据
        let html=document.getElementById("message").innerHTML;

         html+='<p>服务端消息:' + e.data + '</p>'
         document.getElementById("message").innerHTML=html
     }
     ws.onclose = function(e){
         //当客户端收到服务端发送的关闭连接请求时，触发onclose事件
         console.log("close");
     }
     ws.onerror = function(e){
         //如果出现连接、处理、接收、发送数据失败的时候触发onerror事件
         console.log(error);
     }

 </script>
</body>
</html>
```

运行 `go run main.go`，在 goland 中右键运行 `a.html` 可以在浏览器中看到连接到服务端。

新建文件 `src/core/Common.go`

```go
package core

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var Upgrader websocket.Upgrader

func init() {
	Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}
```

新建文件 `src/core/ClientMap.go`

```go
package core

import (
	"github.com/gorilla/websocket"
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
```

新建文件 `src/handlers/Echo.go`

```go
package handlers

import (
	"learn-websocket/src/core"
	"log"
	"net/http"
)

func Echo(w http.ResponseWriter, r *http.Request) {
	client, err := core.Upgrader.Upgrade(w, r, nil) // 升级
	if err != nil {
		log.Println()
	} else {
		core.ClientMap.Store(client.RemoteAddr().String(), client)
	}
}
```

修改 `main.go`

```go
package main

import (
	"github.com/gorilla/websocket"
	"learn-websocket/src/handlers"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/echo", handlers.Echo)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

对所有连接进来的客户端发送消息，扩展 `main.go`

```go
package main

import (
	"github.com/gorilla/websocket"
	"learn-websocket/src/core"
	"learn-websocket/src/handlers"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/echo", handlers.Echo)

	http.HandleFunc("/sendall", func(w http.ResponseWriter, req *http.Request) {
		msg:=req.URL.Query().Get("msg")
		core.ClientMap.SendAll(msg)
		w.Write([]byte("OK"))
	})

	
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```

修改 `src/core/ClientMap.go`

```go
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
```

在 `Goland` 中右键点击运行 `a.html`，打开 rest 工具，请求 http://localhost:8080/sendall?msg=abc，可以看到所有 `a.html` 客户端接收到 `服务端消息:abc`

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/5318e25ec46b8d1b7f419c7537f5693407e2e122?branch=5318e25ec46b8d1b7f419c7537f5693407e2e122&diff=split#diff-f05c20371760668299b744a5d25fd5e4c2e3c7fb1621bec98c95a759672fa6e7L2)

### 03. 简易封装Client对象、Ping客户端

```go
type ClientMapStruct struct {
	data sync.Map // key 是客户端 IP  Value 就是 WebSocket 连接对象
}
```

这样在后面做一些扩展，可能就不是很方便，首先新建文件 `src/core/WsClient.go` ，做封装为了后面扩展方法方便。

```go
package core

import (
	"github.com/gorilla/websocket"
	"time"
)

type WsClient struct {
	conn *websocket.Conn
}

func NewWsClient(conn *websocket.Conn) *WsClient {
	return &WsClient{conn: conn}
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
```

增加 `ClientMap.Remove` 方法

```go
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
```

修改 `src/handlers/Echo.go`

```go
package handlers

import (
	"learn-websocket/src/core"
	"log"
	"net/http"
)

func Echo(w http.ResponseWriter, r *http.Request) {
	client, err := core.Upgrader.Upgrade(w, r, nil) // 升级
	if err != nil {
		log.Println()
	} else {
		core.ClientMap.Store(client)
	}
}
```

修改前端不显示 `ping` 消息

```html
    ws.onmessage = function (e) {
        if(e.data==='ping'){
            return
        }
        //当客户端收到服务端发来的消息时，触发onmessage事件，参数e.data包含server传递过来的数据
        let html = document.getElementById("message").innerHTML;

        html += '<p>服务端消息:' + e.data + '</p>'
        document.getElementById("message").innerHTML = html
    }
```

这样就可以尽可能的保持 `ClientMap` 里面 `data` 里面存储的客户端尽可能有效，毕竟有 `ping` 的等待时间。

代表变动 [git commit](https://github.com/custer-go/learn-gin/commit/b763cc228057e049d054df6c615d0716fed8d04e?branch=b763cc228057e049d054df6c615d0716fed8d04e&diff=split#diff-29ed2b750da16e3e8bead96ec1397ab459c8ebe790c335fc8efee93fbfe5491dL4)

### 04. 简单整合vue、发送JSON数据

新增文件 `src/models/PodModel.go`

```go
package models

type PodModel struct {
	PodName  string
	PodImage string
	PodNode  string
}

func MockPodList() []*PodModel {
	return []*PodModel{
		{PodName: "pod-101", PodImage: "nginx:1.18", PodNode: "node1"},
		{PodName: "pod-76xs", PodImage: "alpine:3.12", PodNode: "node2"},
		{PodName: "pod-F#ff3", PodImage: "tomcat:8", PodNode: "node3"},
	}
}
```

修改 `src/core/ClientMap.go`

```go
// SendAll 向所有客户端 发送消息
func (c *ClientMapStruct) SendAllPods() {
	c.data.Range(func(key, value interface{}) bool {
		conn := value.(*WsClient).conn
		//err := c.WriteMessage(websocket.TextMessage, []byte(msg))
		err := conn.WriteJSON(models.MockPodList())
		if err != nil {
			c.Remove(conn)
			log.Println(err)
		}
		return true
	})
}
```

修改 `main.go`

```go
	http.HandleFunc("/sendall", func(w http.ResponseWriter, req *http.Request) {
		core.ClientMap.SendAllPods()
		w.Write([]byte("OK"))
	})
```



