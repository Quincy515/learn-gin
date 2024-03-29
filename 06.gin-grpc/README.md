grpc 实战入门

[toc]

### 01. RPC入门、创建中间文件

#### 基本原理

RPC（Remote Procedure Call）远程过程调用，RPC 就是要像调用本地的函数一样去调远程函数。

<img src="../imgs/24_grpc.jpg" style="zoom:90%;" />



#### 整体来讲

- 整个过程

 1、客户端 发送 数据（以字节流的方式）

 2、服务端接收，并解析。 根据 约定 知道要知道执行什么。然后把结果返回客户端

-  RPC就是 把  

1、上述过程封装下，使其操作更加优化

2、使用一些大家都认可的协议 使其规范化

3、做成一些框架。直接或间接产生利益

#### grpc

grpc就是一款语言中立、平台中立、开源的远程过程调用(RPC)框架gRpc 。

https://github.com/grpc/grpc-go

安装：`go get -u google.golang.org/grpc`

#### Protobuf

Google Protocol Buffer( 简称 Protobuf)

轻便高效的序列化数据结构的协议，可以用于网络通信和数据存储。

 特点：性能高、传输快、维护方便，反正就是**各种好，各种棒**

一些第三方rpc库都会支持protobuf  

- github地址：https://github.com/protocolbuffers/protobuf

- golang库所属地址：https://github.com/golang/protobuf

#### 安装

- 第一步来到这：

查看 https://github.com/protocolbuffers/protobuf/blob/master/src/README.md#c-installation---windows

继而安装https://github.com/protocolbuffers/protobuf/releases/latest

- protobuf相关文档

https://developers.google.com/protocol-buffers/docs/gotutorial

这是protobuf编译器，将.proto文件，转译成protobuf的原生数据结构

- go 插件

`go get github.com/golang/protobuf/protoc-gen-go`

此时会在你的 `GOPATH` 的 `bin` 目录下生成可执行文件。

`protobuf` 的编译器插件 `protoc-gen-go`

等下我们执行 `protoc` 命令时 就会自动调用这个插件。

-  `Goland` 插件

#### 创建中间文件

以 `.proto` 结尾的中间文件，`gin-grpc/pbfile/Prod.proto`

```
syntax = "proto3";
package services;
message  ProdRequest {
  int32 prod_id = 1;   // 传入的商品ID
}
message ProdResponse{
  int32 prod_stock = 1;// 商品库存
}
```

新建文件 `gin-grpc/services`，在目录 `gin-grpc/pbfile` 下

执行 `protoc --go_out=../services/ Prod.proto`。

`--go_out` 表示调用 `go` 插件， 指定生成的目录，最后是对应的原文件 `Prod.proto`。

这样在 `services` 目录下就生成了 `Prod.pb.go` 文件。

或者在根目录下执行 `protoc --proto_path=pbfile --go_out=services pbfile/Prod.proto` 指定 `proto` 文件的路径。

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/a744fd338089e1093db7b8d6bd7369284b2fc89f#diff-91be4e775403e5048a6fb19c931abd2fe591204782b7aab11e3ec104bec5825eR1)

### 02. 创建gRPC服务端并运行

上面做了一个 “中间文件”并生成对应的go文件，现在创建真正的服务。

#### 第1步：修改 `.proto` 文件

```protobuf
syntax = "proto3";
package services;
message  ProdRequest {
  int32 prod_id = 1;   // 传入的商品ID
}
message ProdResponse{
  int32 prod_stock = 1;// 商品库存
}
service ProdService {
  rpc GetProdStock (ProdRequest) returns (ProdResponse);
}
```

#### 第2步：重新生成 `.pb.go` 文件

之前执行的是 `protoc --go_out=../services/ Prod.proto`

现在执行的是 `protoc --go_out=plugins=grpc:../services Prod.proto`

或者在根目录下执行

 `protoc --proto_path=pbfile --go_out=plugins=grpc:services pbfile/Prod.proto`。

会自动覆盖之前生成的 `Prod.pb.go` 文件。

此时生成的文件 `Prod.pb.go` 主要变动是有两个接口

```go
// ProdServiceClient is the client API for ProdService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProdServiceClient interface {
	GetProdStock(ctx context.Context, in *ProdRequest, opts ...grpc.CallOption) (*ProdResponse, error)
}
...
func NewProdServiceClient(cc *grpc.ClientConn) ProdServiceClient {
	return &prodServiceClient{cc}
}
...
// ProdServiceServer is the server API for ProdService service.
type ProdServiceServer interface {
	GetProdStock(context.Context, *ProdRequest) (*ProdResponse, error)
}
...
func RegisterProdServiceServer(s *grpc.Server, srv ProdServiceServer) {
	s.RegisterService(&_ProdService_serviceDesc, srv)
}
```

发布服务的时候，就需要新建个 `struct` 来继承这个 `interface{}` 接口，即实现 `GetProdStock(context.Context, *ProdRequest) (*ProdResponse, error)` 方法。

#### 第3步：新建具体的实现类

新建文件 `gin-grpc/services/ProdService.go`

```go
package services

import "context"

type ProdService struct{}

func (this *ProdService) GetProdStock(ctx context.Context, request *ProdRequest) (*ProdResponse, error) {
	return &ProdResponse{ProdStock: 20}, nil
}
```

#### 第4步：创建 grpc 服务端并运行

新建文件 `gin-grpc/server.go`

```go
package main

import (
	"gin-grpc/services"
	"google.golang.org/grpc"
	"net"
)

func main() {
	rpcServer := grpc.NewServer()
	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))
	listen, _ := net.Listen("tcp", ":8081")
	rpcServer.Serve(listen)
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/0fe343ad58a6696c1cb6607f7e0c6f3f767cf084#diff-91be4e775403e5048a6fb19c931abd2fe591204782b7aab11e3ec104bec5825eL5)

### 03. 创建客户端调用

客户端可以新建一个工程，或者在当前工程下新建 `client` 文件夹表示客户端代码。

在客户端的代码中不需要使用中间 `.proto` 文件，只引用生成的 `.pb.go` 文件。

新建 `main.go` 来完成客户端调用代码

```go
package main

import (
	"context"
	"fmt"
	"gin-grpc/services"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	prodClient := services.NewProdServiceClient(conn)
	prodRes, err := prodClient.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 12})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(prodRes.ProdStock)
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/9646a03fa175f50299c830c8294159f9968cb061#diff-dc576b33b5093f4c968f2943df65b7a64afda74e81f771e62d310a3c77e525a5R1)

### 04. 自签证书、服务加入证书验证

免费证书申请 https://freessl.cn/

在生成环境中不能使用自签证书，在云服务器中，单域名可以免费申请 ssl，或者购买。

Windows 下载 `openssl` 工具： http://slproweb.com/products/Win32OpenSSL.html

#### 第1步：生成 `.key` 私钥文件

`openssl genrsa -des3 -out custer.key 2048`

- `genrsa` : 生成 `rsa` 私钥
- `-des3`: `des3` 算法
- `2048`: 表示 2048 位强度
- `custer.key`: 私钥文件名

输入密码，这里输入两次。填写一样即可。随意填写一个。后续就会删除这个密码。

此时会生成 `custer.key` 这个文件。

#### 第2步： 删除密码

`openssl rsa -in custer.key -out custer.key`

注意这里目录和生成私钥的目录一致，会输入一遍密码。

#### 第3步：创建证书签名请求，生成 `.csr ` 文件

`openssl req -new -key custer.key -out custer.csr`

根据刚刚生成的 `key` 文件来生成证书请求文件。

执行以上命令后，需要依次输入国家、地区、城市、组织、组织单位、Common Name、Email和密码。其中Common Name应该与域名保持一致。密码我们已经删掉了,直接回车即可。

**温馨提示**Common Name就是证书对应的域名地址。

#### 第4步：生成自签名证书

根据以上2个文件生成crt证书文件，终端执行下面命令：

`openssl x509 -req -days 3650 -in custer.csr -signkey custer.key -out ssl.crt`

这里3650是证书有效期(单位：天)。这个随意。最后使用到的文件是key和crt文件。

到这里我们的证书就已经创建成功了(custer.key 和 custer.crt) 可以直接用到https的server中了。

> 需要注意的是，在使用自签名的证书时，浏览器会提示证书的颁发机构是未知的

#### 服务加入证书验证

创建新文件夹 `keys`，把没有密码的 `.key` 文件和 `.crt` 文件放入 `keys` 目录下。

#### 加入证书代码：服务端

```go
package main

import (
	"gin-grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

func main() {
	creds, err := credentials.NewServerTLSFromFile("keys/grpc.crt", "keys/grpc.key")
	if err != nil {
		log.Fatal(err)
	}

	rpcServer := grpc.NewServer(grpc.Creds(creds))
	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))
	listen, _ := net.Listen("tcp", ":8081")
	rpcServer.Serve(listen)
}
```

#### 加入证书代码：客户端

```go
package main

import (
	"gin-grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

func main() {
	creds, err := credentials.NewServerTLSFromFile("keys/grpc.crt", "keys/grpc.key")
	if err != nil {
		log.Fatal(err)
	}
	rpcServer:=grpc.NewServer(grpc.Creds(creds))

	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))
	listen, _ := net.Listen("tcp", ":8081")
	rpcServer.Serve(listen)
}
```

运行会报错

```bash
time="2020-11-18T12:54:48+08:00" level=fatal msg="rpc error: code = Unavailable desc = connection error: desc = \"transport: authentication handshake failed: x509: certificate relies on legacy Common Name field, use SANs or temporarily enable Common Name matching with GODEBUG=x509ignoreCN=0\""
```

如果出现上述报错，是因为 go 1.15 版本开始[废弃 CommonName](https://golang.org/doc/go1.15#commonname)，因此推荐使用 SAN 证书。 如果想兼容之前的方式，需要设置环境变量 GODEBUG 为 `x509ignoreCN=0`。

下面简单示例如何用 openssl 生成 ca 和双方 SAN 证书。

准备默认 OpenSSL 配置文件于当前目录

linux系统在 : `/etc/pki/tls/openssl.cnf`

Mac系统在: `/System/Library/OpenSSL/openssl.cnf`

第1步：cp 目录到项目目录进行修改设置

`cp /System/Library/OpenSSL/openssl.cnf /learn-gin/06.gin-grpc/keys`

第2步：找到 [ CA_default ],打开 copy_extensions = copy

第3步：找到[ req ],打开 req_extensions = v3_req # The extensions to add to a certificate request

第4步：找到[ v3_req ],添加 subjectAltName = @alt_names

第5步：添加新的标签 [ alt_names ] , 和标签字段 

```
[ alt_names ]
DNS.1 = *.org.custer.fun
DNS.2 = *.custer.fun
```

这里填入需要加入到 Subject Alternative Names 段落中的域名名称，可以写入多个。

第6步：生成证书私钥 test.key：

`openssl genpkey -algorithm RSA -out test.key`

第7步：通过私钥test.key生成证书请求文件test.csr：

`openssl req -new -nodes -key test.key -out test.csr -days 3650 -subj "/C=cn/OU=custer/O=custer/CN=custer.fun" -config ./openssl.cnf -extensions v3_req`

第8步：：test.csr是上面生成的证书请求文件。custer.crt/custer.key是CA证书文件和key，用来对test.csr进行签名认证。这两个文件在之前生成的。

第9步：生成SAN证书：

`openssl x509 -req -days 3650 -in test.csr -out test.pem -CA custer.crt -CAkey custer.key -CAcreateserial -extfile ./openssl.cnf -extensions v3_req`

现在 Go 1.15 以上版本的 GRPC 通信，就可以使用了

第10步：服务端 tls 加载

`creds, err := credentials.NewServerTLSFromFile("test.pem", "test.key")`

第11步：客户端加载

`creds,err := credentials.NewClientTLSFromFile("test.pem","*.custer.fun")`

学习参考链接

1. https://www.cnblogs.com/jackluo/p/13841286.html
2. https://blog.csdn.net/cuichenghd/article/details/109230584

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/44939287307e2434dff0ea176688447398fac992#diff-dc576b33b5093f4c968f2943df65b7a64afda74e81f771e62d310a3c77e525a5L5)

### 05. 让gRPC提供Http服务(初步)

之前在 创建gRPC服务端并运行时的代码是 `rpcServer.Serve(listen)`，

现在替换成提供 HTTP 服务。

```go
package main

import (
	"gin-grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net/http"
)

func main() {
	creds, err := credentials.NewServerTLSFromFile("keys/test.pem", "keys/test.key")
	if err != nil {
		log.Fatal(err)
	}
	rpcServer := grpc.NewServer(grpc.Creds(creds))
	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))

	//listen, _ := net.Listen("tcp", ":8081")
	//rpcServer.Serve(listen)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		rpcServer.ServeHTTP(writer, request)
	})
	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	httpServer.ListenAndServe()
}
```

直接使用浏览器访问 http://localhost:8081/ 可以看到 `gRPC requires HTTP/2`

直接使用上面的客户端代码也是访问不了的，会报一个错误

```bash
time="2020-11-18T13:32:21+08:00" level=fatal msg="rpc error: code = Unavailable desc = connection error: desc = \"transport: authentication handshake failed: tls: first record does not look like a TLS handshake\""
```

在服务端使用另外的方法启动 `http server`

```go
package main

import (
	"fmt"
	"gin-grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net/http"
)

func main() {
	creds, err := credentials.NewServerTLSFromFile("keys/test.pem", "keys/test.key")
	if err != nil {
		log.Fatal(err)
	}
	rpcServer := grpc.NewServer(grpc.Creds(creds))
	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))

	//listen, _ := net.Listen("tcp", ":8081")
	//rpcServer.Serve(listen)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request)
		rpcServer.ServeHTTP(writer, request)
	})
	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	httpServer.ListenAndServeTLS("keys/test.pem", "keys/test.key")
}
```

这时再访问浏览器 https://localhost:8081/ 会出现 `invalid gRPC request method`

使用客户端请求可以正常访问，查看控制台。

```go
&{POST 
  /services.ProdService/GetProdStock 
  HTTP/2.0 2 0 
  map[Content-Type:[application/grpc] 
      Te:[trailers] 
      User-Agent:[grpc-go/1.33.2]] 0xc00008e240 <nil> -1 [] false *.custer.fun map[] map[] <nil> map[] 127.0.0.1:56154 /services.ProdService/GetProdStock 0xc000134bb0 <nil> <nil> 0xc00009a080}
```

以客户端访问，思考：

1. 这个地址 `/services.ProdService/GetProdStock ` 是否可以改变？
2. 使用普通的 http client 是佛可以调用？
3. 在 linux 中怎么使用工具进行测试？

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/040cc1785cb88c8a046ee20bf9cafafc2d2f8fe0#diff-dc576b33b5093f4c968f2943df65b7a64afda74e81f771e62d310a3c77e525a5L15)

### 06. 使用自签CA、server、Client证书和双向认证

免费证书申请 https://freessl.cn/

之前在客户端代码中也是使用的是服务端 `.crt` 证书或 `.pem`。

在实际开发中，内置服务的调用，需要双向验证，客户端和服务端都必须要有各个的证书。

新建目录来生成证书 `key`

#### 使用CA证书

- 根证书（root certificate）是属于根证书颁发机构（CA）的公钥证书。 用以验证它所签发的证书（客户端、服务端）
- 1、`openssl genrsa -out ca.key 2048`
- 2、`openssl req -new -x509 -days 3650 -key ca.key -out ca.pem`

```bash
👍 openssl req -new -x509 -days 3650 -key ca.key -out ca.pem
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:cn
State or Province Name (full name) [Some-State]:shanghai
Locality Name (eg, city) []:shanghai
Organization Name (eg, company) [Internet Widgits Pty Ltd]:custer
Organizational Unit Name (eg, section) []:custer
Common Name (e.g. server FQDN or YOUR name) []:localhost
Email Address []:
```

生成 `ca.pem` 文件。

#### 重新生成服务端证书

- 1、`openssl genrsa -out server.key 2048`
- 2、`openssl req -new -key server.key -out server.csr`
-  注意 `common name` 请填写 `localhost`
- 3、`openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in server.csr -out server.pem`

 #### 生成客户端

- 1、`openssl ecparam -genkey -name secp384r1 -out client.key`
- 2、`openssl req -new -key client.key -out client.csr`
- 3、`openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.pem`

程序中重新覆盖 `server.crt` 和 `server.key`

#### 服务端拷贝证书文件

新建文件夹 `cert`，用来存放自签 ca 双向认证证书

在服务端，需要拷贝 `server.key` `server.pem` `ca.pem` 这三个证书。

#### 服务端代码改造

```go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"gin-grpc/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"net/http"
)

func main() {
	//creds, err := credentials.NewServerTLSFromFile("keys/client.pem", "keys/client.key")
	//if err != nil {
	//	log.Fatal(err)
	//}

	cert, _ := tls.LoadX509KeyPair("cert/server.pem", "cert/server.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	rpcServer := grpc.NewServer(grpc.Creds(creds))

	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))

	//listen, _ := net.Listen("tcp", ":8081")
	//rpcServer.Serve(listen)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(request)
		rpcServer.ServeHTTP(writer, request)
	})
	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	httpServer.ListenAndServeTLS("keys/client.pem", "keys/client.key")
}
```

#### 客户端证书拷贝

拷贝 `client.key` `client.pem` `ca.pem` 证书到客户端 `cert` 目录下

#### 客户端代码改造

```go
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"gin-grpc/services"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

func main() {
	//creds, err := credentials.NewClientTLSFromFile("keys/client.pem", "*.custer.fun")
	//if err != nil {
	//	log.Fatal(err)
	//}

	cert, _ := tls.LoadX509KeyPair("cert/client.pem", "cert/client.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
		RootCAs:      certPool,
	})

	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(creds))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	prodClient := services.NewProdServiceClient(conn)
	prodRes, err := prodClient.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 12})
	if err != nil {
		log.Fatal(err)
	}
	log.Info(prodRes.ProdStock)
}
```

按照 Go 1.15 生成 SAN 证书

```bash
第1步：生成 CA 根证书
👍 openssl genrsa -out ca.key 2048
Generating RSA private key, 2048 bit long modulus (2 primes)
.............+++++
..................................................................................................................+++++
e is 65537 (0x010001)
👍 openssl req -new -x509 -days 3650 -key ca.key -out ca.pem
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [AU]:cn
State or Province Name (full name) [Some-State]:shanghai
Locality Name (eg, city) []:shanghai
Organization Name (eg, company) [Internet Widgits Pty Ltd]:custer
Organizational Unit Name (eg, section) []:custer
Common Name (e.g. server FQDN or YOUR name) []:localhost
Email Address []:
      
第2步：生成服务端证书      
👍 openssl genpkey -algorithm RSA -out server.key
........................................................................................+++++
.......................................+++++
👍 openssl req -new -nodes -key server.key -out server.csr -days 3650 -subj "/C=cn/OU=custer/O=custer/CN=localhost" -config ./openssl.cnf -extensions v3_req
Ignoring -days; not generating a certificate
👍 openssl x509 -req -days 3650 -in server.csr -out server.pem -CA ca.pem -CAkey ca.key -CAcreateserial -extfile ./openssl.cnf -extensions v3_req
Signature ok
subject=C = cn, OU = custer, O = custer, CN = localhost
Getting CA Private Key

第3步：生成客户端证书
👍 openssl genpkey -algorithm RSA -out client.key
........+++++
...........+++++
👍 openssl req -new -nodes -key client.key -out client.csr -days 3650 -subj "/C=cn/OU=custer/O=custer/CN=localhost" -config ./openssl.cnf -extensions v3_req
Ignoring -days; not generating a certificate
👍 openssl x509 -req -days 3650 -in client.csr -out client.pem -CA ca.pem -CAkey ca.key -CAcreateserial -extfile ./openssl.cnf -extensions v3_req
Signature ok
subject=C = cn, OU = custer, O = custer, CN = localhost
Getting CA Private Key
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/b26d6870d8704d49829756b5b1f281d61eb9f802#diff-dc576b33b5093f4c968f2943df65b7a64afda74e81f771e62d310a3c77e525a5L2)

### 07. 双向认证下rpc-gateway使用（同时提供rpc和http接口)

第三方库 https://github.com/grpc-ecosystem/grpc-gateway

![architecture introduction diagram](https://camo.githubusercontent.com/5fc816f4575582674ed5f7216b7169e1a8496b531007faf2aab07a3b01484d7e/68747470733a2f2f646f63732e676f6f676c652e636f6d2f64726177696e67732f642f3132687034435071724e5046686174744c5f63496f4a707446766c41716d35774c513067677149356d6b43672f7075623f773d37343926683d333730)

在 `grpc` 之上加一层代理并转发，转变成 `protobuf` 格式来访问 `grpc` 服务。

#### 安装

```go
// +build tools

package tools

import (
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
    _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
    _ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
```

Run `go mod tidy` to resolve the versions. Install by running

```
$ go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

This will place four binaries in your `$GOBIN`;

- `protoc-gen-grpc-gateway`
- `protoc-gen-openapiv2`
- `protoc-gen-go`
- `protoc-gen-go-grpc`

Make sure that your `$GOBIN` is in your `$PATH`.

#### 修改 `proto` 文件

为了 `import "google/api/annotations.proto";` 路径

把 `go mod` 中的文件 `/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway/v2@v2.0.1/third_party/googleapis/google` 放到目录 `pbfile` 中

```bash
(base)  👍  ~/Work/2020/study/learn-gin/06.gin-grpc   main ±✚  tree .                               
.
├── README.md
├── cert
├── client
│   └── main.go
├── go.mod
├── go.sum
├── keys
├── pbfile
│   ├── Prod.proto
│   └── google
│       ├── api
│       │   ├── annotations.proto
│       │   ├── http.proto
│       │   └── httpbody.proto
│       └── rpc
│           ├── code.proto
│           ├── error_details.proto
│           └── status.proto
├── server.go
└── services
    ├── Prod.pb.go
    ├── Prod.pb.gw.go
    ├── ProdService.go
    └── pbfile
        └── Prod
            └── Prod.pb.gw.go
```

修改 `.proto` 以实现，比如访问的 url 是 `GET /prod/stock/{}`

```protobuf
syntax = "proto3";
package services;
option go_package = ".;services"; // .代表当前文件夹，分号后面是生成go文件引入的包名
import "google/api/annotations.proto";

message  ProdRequest {
  int32 prod_id = 1;   //传入的商品ID
}
message ProdResponse{
  int32 prod_stock = 1;//商品库存
}

service ProdService {
  rpc GetProdStock (ProdRequest) returns (ProdResponse){
    option (google.api.http) = {
      get: "/v1/prod/{prod_id}"
    };
  }
}
```

#### 生成两个文件

首先 `cd` 进入 `pbfiles` 目录

- 生成 `Prod.pb.go`

`protoc --go_out=plugins=grpc:../services Prod.proto`

- 生成 `Prod.pb.gw.go`

`protoc --grpc-gateway_out=logtostderr=true:../services Prod.proto`

#### 改造代码

把证书相关的代码移动到 `helper/CertHelper.go` 文件中

```go
package helper

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

// GetServerCreds 获取服务端证书配置
func GetServerCreds() credentials.TransportCredentials {
	cert, _ := tls.LoadX509KeyPair("cert/server.pem", "cert/server.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert}, //服务端证书
		ClientAuth:   tls.VerifyClientCertIfGiven,
		ClientCAs:    certPool,
	})
	return creds
}

// GetClientCreds 获取客户端证书配置
func GetClientCreds() credentials.TransportCredentials {
	cert, _ := tls.LoadX509KeyPair("cert/client.pem", "cert/client.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert}, //客户端证书
		ServerName:   "localhost",
		RootCAs:      certPool,
	})
	return creds
}
```

#### 基于 grpc-gatway 创建 http server

新建文件 `gateway/httpserver.go`

```go
package main

import (
	"context"
	"gin-grpc/helper"
	"gin-grpc/services"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

func main() {
	gwmux := runtime.NewServeMux() // 创建路由                                               
	opt := []grpc.DialOption{grpc.WithTransportCredentials(helper.GetClientCreds())} // 指定客户端请求时使用的证书
	err := services.RegisterProdServiceHandlerFromEndpoint(
		context.Background(), gwmux, "localhost:8081", opt)
	if err != nil {  //////// 路由 //// grpc 的端口 ////////////
		log.Fatal(err)
	}
	httpServer := &http.Server{
		Addr:    ":8080", // 对外提供的访问端口
		Handler: gwmux,
	}
	httpServer.ListenAndServe()
}
```

- 第1步：启动 `grpc` 服务端
- 第2步：启动 客户端，可以看到控制台输出
- 第2步：启动 `gateway` 访问浏览器 http://localhost:8080/v1/prod/3 可以看到  `{ "prodStock": 28 }`

这样就提供了内部 `grpc` 访问，第三方系统接入使用 `api` 访问。

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/9451ed43c13ffd19d42a742b130544c193089cd7#diff-dc576b33b5093f4c968f2943df65b7a64afda74e81f771e62d310a3c77e525a5L2)

### 08. 语法速学(1):返回商品”数组”、repeated修饰符

#### 第1步：写 `.proto` 文件

之前实现的是 传入一个商品 ID `ProdRequest` 获取一个商品库存 `ProdResponse`。

如果需要获取 一堆商品的库存列表呢?

```protobuf
syntax = "proto3";
package services;
option go_package = ".;services"; // .代表当前文件夹，分号后面是生成go文件引入的包名
import "google/api/annotations.proto";

message  ProdRequest {
  int32 prod_id = 1;   // 传入的商品ID
}
message ProdResponse{
  int32 prod_stock = 1; // 商品库存
}
message QuerySize {
  int32 size = 1; // 页尺寸
}
message ProdResponseList { // 使用修饰符返回商品库存列表
  repeated ProdResponse prodres = 1;
} // 修饰符  类名          变量名   顺序
service ProdService {
  rpc GetProdStock (ProdRequest) returns (ProdResponse){
    option (google.api.http) = {
      get: "/v1/prod/{prod_id}"
    };
  }

  rpc GetProdStocks(QuerySize)returns (ProdResponseList) {}
}
```

`Repeated`:是一个修饰符,返回字段可以重复任意多次(包括0次)，可以认为就是一个数组(切片)。

#### 第2步：生成 `.pb.go` 文件

`protoc --go_out=plugins=grpc:../services Prod.proto`

#### 第3步：在 `services/ProdService.go` 中实现

```go
package services

import "context"

type ProdService struct{}

func (this *ProdService) GetProdStock(ctx context.Context, request *ProdRequest) (*ProdResponse, error) {
	return &ProdResponse{ProdStock: 28}, nil
}

func (this *ProdService) GetProdStocks(context.Context, *QuerySize) (*ProdResponseList, error) {
	Prodres := []*ProdResponse{
		&ProdResponse{ProdStock: 28},
		&ProdResponse{ProdStock: 29},
		&ProdResponse{ProdStock: 30},
		&ProdResponse{ProdStock: 31},
	}
	return &ProdResponseList{Prodres: Prodres}, nil
}
```

完成服务端代码

#### 第4步：拷贝 `.pd.go` 文件到客户端

#### 第5步：修改客户端代码

```go
func main() {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(helper.GetClientCreds()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	prodClient := services.NewProdServiceClient(conn)
	//prodRes, err := prodClient.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 12})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Info(prodRes.ProdStock)
	res, err := prodClient.GetProdStocks(context.Background(), &services.QuerySize{Size: 10})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Prodres)
	fmt.Println(res.Prodres[2].ProdStock)
}
```

启动服务端和客户端，可以看到控制台信息

```bash
[prod_stock:28 prod_stock:29 prod_stock:30 prod_stock:31]
30
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/26133a2e3ffb48a2989a61520fe0871bfbb4ff07#diff-dc576b33b5093f4c968f2943df65b7a64afda74e81f771e62d310a3c77e525a5L2)

### 09. 语法速学(2): 使用枚举、获取分区商品库存

**创建枚举类型，支持分区枚举参数**

商品有区域之分譬如

-------------------------/  A 区有 10 个

 ID 为101 的商品  -  B 区有 12 个

--------------------------\ C 区有 20 个

加入枚举类型

```protobuf
enum ProdAreas{
    A=0;
    B=1;
    C=2;
}
```

修改 `pbfile/Prod.proto`

```protobuf
syntax = "proto3";
package services;
option go_package = ".;services"; // .代表当前文件夹，分号后面是生成go文件引入的包名
import "google/api/annotations.proto";

enum ProdAreas {
  A = 0; // 第一个必须是 0 表示默认值
  B = 1;
  C = 2;
}

message  ProdRequest {
  int32 prod_id = 1;   // 传入的商品ID
  ProdAreas prod_area = 2; // 传入商品区域
}
message ProdResponse {
  int32 prod_stock = 1; // 商品库存
}
message QuerySize {
  int32 size = 1; // 页尺寸
}
message ProdResponseList {// 使用修饰符返回商品库存列表
  repeated ProdResponse prodres = 1;
} // 修饰符  类名          变量名   顺序
service ProdService {
  rpc GetProdStock (ProdRequest) returns (ProdResponse){
    option (google.api.http) = {
      get: "/v1/prod/{prod_id}"
    };
  }

  rpc GetProdStocks(QuerySize) returns (ProdResponseList) {}
}
```

修改实现函数 `services/ProdService.go`

```go
package services

import "context"

type ProdService struct{}

func (this *ProdService) GetProdStock(ctx context.Context, request *ProdRequest) (*ProdResponse, error) {
	var stock int32 = 0
	if request.ProdArea == ProdAreas_A {
		stock = 39
	} else if request.ProdArea == ProdAreas_B {
		stock = 41
	} else {
		stock = 20
	}
	return &ProdResponse{ProdStock: stock}, nil
}

func (this *ProdService) GetProdStocks(context.Context, *QuerySize) (*ProdResponseList, error) {
	Prodres := []*ProdResponse{
		&ProdResponse{ProdStock: 28},
		&ProdResponse{ProdStock: 29},
		&ProdResponse{ProdStock: 30},
		&ProdResponse{ProdStock: 31},
	}
	return &ProdResponseList{Prodres: Prodres}, nil
}
```

拷贝新生成的 `Prod.pb.go` 文件到客户端

```go
func main() {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(helper.GetClientCreds()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	prodClient := services.NewProdServiceClient(conn)
	prodRes, err := prodClient.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 12, ProdArea: services.ProdAreas_B})
	if err != nil {
		log.Fatal(err)
	}
	log.Info(prodRes.ProdStock)
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/91b9a48350105a3b8ddef0a8ad248601fbb9e8bf#diff-dc576b33b5093f4c968f2943df65b7a64afda74e81f771e62d310a3c77e525a5L2)

### 10. 语法速学(3): 导入外部Proto、获取商品信息

新建文件专门存放实体 `pbfile/Models.proto`

```protobuf
syntax = "proto3";
package services; // 可以相同的包，也可以不同
option go_package = ".;services";
message ProdModel { // 商品模型
  int32 prod_id = 1;
  string prod_name = 2;
  float prod_price = 3;
}
```

外部引用

```go
import "Models.proto";
...
	rpc GetProdInfo(ProdRequest) returns (ProdModel) {}
```

生成 `.pb.go` 文件

`protoc --go_out=plugins=grpc:../services Prod.proto`

`protoc --go_out=plugins=grpc:../services Models.proto`

在 `services/ProdService.go` 文件中实现

```go
func (this *ProdService) GetProdInfo(ctx context.Context, in *ProdRequest) (*ProdModel, error) {
	ret := ProdModel{
		ProdId:    101,
		ProdName:  "测试商品",
		ProdPrice: 20.5,
	}
	return &ret, nil
}
```

拷贝两个新生成的 `.pb.go` 文件到客户端

修改客户端代码

```go
func main() {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(helper.GetClientCreds()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	prodClient := services.NewProdServiceClient(conn)
	// 获取商品库存
	//prodRes, err := prodClient.GetProdStock(context.Background(), &services.ProdRequest{ProdId: 12, ProdArea: services.ProdAreas_B})
	prodRes, err := prodClient.GetProdInfo(context.Background(), &services.ProdRequest{ProdId: 12})
	if err != nil {
		log.Fatal(err)
	}
	log.Info(prodRes)
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/f787294c152bcc155bd05a2b1c10e9a795539b63#diff-dc576b33b5093f4c968f2943df65b7a64afda74e81f771e62d310a3c77e525a5L17)

### 11. 语法学习(5)日期类型、创建主订单模型(上)

创建主订单模型 `pbfile/Models.go`

```protobuf
message OrderMain{ //主订单模型
    int32 order_id=1;//订单ID，数字自增
    string order_no=2; //订单号
    int32 user_id=3; //购买者ID
    float order_money=4;//商品金额 
   //这里还需要一个订单时间
}
```

使用TimeStamp

```protobuf
import   "google/protobuf/timestamp.proto";

message OrderMain{ //主订单模型
    int32 order_id=1;//订单ID，数字自增
    string order_no=2; //订单号
    int32 user_id=3; //购买者ID
    float order_money=4;//商品金额
    google.protobuf.Timestamp order_time=5;
}
```

新建文件 `pbfile/Orders.proto` 设定一个方法

```protobuf
syntax = "proto3";
package services;
import "Models.proto";
option go_package = ".;services";

message OrderResponse {
  string status = 1;
  string message = 2;
}

service OrderService {
  rpc NewOrder(OrderMain) returns (OrderResponse) {}
}
```

生成 `.pb.go` 文件

`protoc --go_out=plugins=grpc:../services Orders.proto`
`protoc --go_out=plugins=grpc:../services Models.proto`

在 `services/OrdersService.go` 实现

```go
package services

import (
	"context"
	"fmt"
)

type OrderService struct{}

func (this *OrderService) NewOrder(ctx context.Context, orderMain *OrderMain) (*OrderResponse, error) {
	fmt.Println(orderMain)
	return &OrderResponse{
		Status:  "ok",
		Message: "success",
	}, nil
}
```

在服务端 rpc 启动服务里注册订单服务

```go
func main() {
	rpcServer := grpc.NewServer(grpc.Creds(helper.GetServerCreds()))
	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))   // 商品服务
	services.RegisterOrderServiceServer(rpcServer, new(services.OrderService)) // 订单服务

	listen, _ := net.Listen("tcp", ":8081")
	rpcServer.Serve(listen)
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/8909dcd7ee33c221b980949641e78a73d8218356#diff-84d5cadaa82f6eeca6dc3d619649d2c2a9176ca7137584962e2235e93101175dL1)

### 12. 语法学习(5)日期类型、创建主订单模型(下)

上面创建好了 订单服务的服务端，下面看下客户端怎么调用，

先把上面新生成的 `Models.pb.go` 文件和 `Orders.pb.go` 文件拷贝到客户端。

然后在客户端中创建新的 `OrdersService`

```go
func main() {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(helper.GetClientCreds()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()
	t := timestamp.Timestamp{Seconds: time.Now().Unix()}
	orderClient := services.NewOrderServiceClient(conn)
	res, _ := orderClient.NewOrder(ctx, &services.OrderMain{
		OrderId:    1001,
		OrderNo:    "20201118",
		OrderMoney: 90,
		OrderTime:  &t,
	})
	fmt.Println(res)
}
```

运行客户端，可以看到服务端返回的 `status:"ok"  message:"success"`

在服务端，可以看到客户端的请求数据

 `order_id:1001 order_no:"20201118" order_money:90 order_time:{seconds:1605707319}`

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/dfcfcd65cc9c4caa0d22feaeae355dda8b36d02e#diff-dc576b33b5093f4c968f2943df65b7a64afda74e81f771e62d310a3c77e525a5L2)

### 13. 场景练习(1):POST提交主订单数据(gateway实现http api)

在 `pbfile/Orders.proto` 里增加 `POST` 请求

```protobuf
syntax = "proto3";
package services;
option go_package = ".;services";
import "Models.proto";
import "google/api/annotations.proto";

message OrderRequest {
  OrderMain order_main = 1;
}

message OrderResponse {
  string status = 1;
  string message = 2;
}

service OrderService {
  rpc NewOrder(OrderRequest) returns (OrderResponse) {
    option(google.api.http) = {
      post: "/v1/order"
      body: "order_main" // 提交参数
    } ;
  }
}
```

使用 `grpc-gateway` 编译出两种 `.pb.go` 文件

`protoc --go_out=plugins=grpc:../services Orders.proto`

`protoc --grpc-gateway_out=logtostderr=true:../services Orders.proto`

`protoc --grpc-gateway_out=logtostderr=true:../services Prod.proto `

修改 `services/OrdersService.go`

```go
package services

import (
	"context"
	"fmt"
)

type OrderService struct{}

func (this *OrderService) NewOrder(ctx context.Context, orderRequest *OrderRequest) (*OrderResponse, error) {
	fmt.Println(orderRequest.OrderMain)
	return &OrderResponse{
		Status:  "ok",
		Message: "success",
	}, nil
}
```

修改客户端调用的请求参数

```go
func main() {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(helper.GetClientCreds()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()
	t := timestamp.Timestamp{Seconds: time.Now().Unix()}
	orderClient := services.NewOrderServiceClient(conn)
	res, _ := orderClient.NewOrder(ctx, &services.OrderRequest{
		OrderMain: &services.OrderMain{
			OrderId:    1001,
			OrderNo:    "20201118",
			OrderMoney: 90,
			OrderTime:  &t,
		},
	})
	fmt.Println(res)
}
```

新增 `OrderService` 网关，修改代码 `gateway/httpserver.go`

```go
func main() {
	gwmux := runtime.NewServeMux() // 创建路由
	gRpcEndPoint := "localhost:8081"
	opt := []grpc.DialOption{grpc.WithTransportCredentials(helper.GetClientCreds())} // 指定客户端请求时使用的证书
	err := services.RegisterProdServiceHandlerFromEndpoint(
		context.Background(), gwmux, gRpcEndPoint, opt)
	if err != nil {
		log.Fatal(err)
	}
	err = services.RegisterOrderServiceHandlerFromEndpoint(
		context.Background(), gwmux, gRpcEndPoint, opt)
	if err != nil {
		log.Fatal(err)
	}
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: gwmux,
	}
	httpServer.ListenAndServe()
}
```

先启动 `grpc` 的 `server` ，再启动 `gateway/httpservver.go`

打开浏览器访问 http://localhost:8080/v1/prod/123 可以看到 '

```json
{
	"prodStock": 39
}
```

发送 `POST` 请求

![](../imgs/25_grpc_gateway.jpg)

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/d88bdbd4487e1e3354e909407412de195847178e#diff-dc576b33b5093f4c968f2943df65b7a64afda74e81f771e62d310a3c77e525a5L22)

### 14. 场景练习(2):POST提交主子订单、嵌套模型

上面设定了主订单模型 `pbfile/Models.proto`

```protobuf
message OrderMain{//主订单模型
  int32 order_id = 1;//订单ID，数字自增
  string order_no = 2; //订单号
  int32 user_id = 3; //购买者ID
  float order_money = 4;//商品金额
  google.protobuf.Timestamp order_time = 5; // 订单时间
}
```

下面加入子订单模型 

```protobuf
message OrderDetail{// 子订单模型
  int32 detail_id = 1;
  string order_no = 2;
  int32 prod_id = 3;
  float prod_price = 4;
  int32 prod_num = 5;
}
```

修改主订单，嵌套子订单

```protobuf
message OrderMain{//主订单模型
  int32 order_id = 1;//订单ID，数字自增
  string order_no = 2; //订单号
  int32 user_id = 3; //购买者ID
  float order_money = 4;//商品金额
  google.protobuf.Timestamp order_time = 5; // 订单时间
  repeated OrderDetail order_detail = 6; // 嵌套子订单
}
```

执行 `gen proto` 生成 `.pb.go`

`protoc --go_out=plugins=grpc:../services Models.proto`

重新运行 `grpc` 服务端可 `gateway httpserver`，可以和之前一样只提交主订单数据

![](../imgs/25_grpc_gateway.jpg)

发送嵌套订单

![](../imgs/26_grpc_gateway.jpg)

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/c26046a7e8b7cc643479df33f9322fc42f4b83fd#diff-84d5cadaa82f6eeca6dc3d619649d2c2a9176ca7137584962e2235e93101175dL14)

### 15. 穿插知识点：使用第三方库进行字段验证

1. 比如是 `gin` 和 `grpc` ，提交是统一走 `gin web` 框架，可以在 `gin` 中验证。验证通过了调用 `grpc` 方法。
2. 比如像现在内部是 `grpc` 外部有个网关 `gateway`，在网管部分验证非常麻烦，可以放到  `grpc` ，即实际实现部分，比如 `services/OrdersService.go` 中的 `orderRequest`，可以手动完成，

使用第三方库  https://github.com/envoyproxy/protoc-gen-validate

安装 `go get -u github.com/envoyproxy/protoc-gen-validate`

使用这个插件在生成 `.pb.go` 时就生成了验证规则。

虽然写起来比较麻烦，但是好处是写在 `.proto` 中间文件中，可以自动生成不同语言的验证。

修改订单金额必须大于1。

```protobuf
syntax = "proto3";
package services;
option go_package = ".;services";
import   "google/protobuf/timestamp.proto";
import "validate/validate.proto";

message ProdModel {// 商品模型
  int32 prod_id = 1;
  string prod_name = 2;
  float prod_price = 3;
}

message OrderMain{//主订单模型
  int32 order_id = 1;//订单ID，数字自增
  string order_no = 2; //订单号
  int32 user_id = 3; //购买者ID
  float order_money = 4 [(validate.rules).float.gt = 1];//商品金额
  google.protobuf.Timestamp order_time = 5; // 订单时间
  repeated OrderDetail order_detail = 6; // 嵌套子订单
}

message OrderDetail{// 子订单模型
  int32 detail_id = 1;
  string order_no = 2;
  int32 prod_id = 3;
  float prod_price = 4;
  int32 prod_num = 5;
}
```

生成文件 

```bash
protoc \
  -I . \
  -I ${GOPATH}/src \
  -I ${GOPATH}/src/github.com/protoc-gen-validate \
  --go_out=plugins=grpc:../services \
  --validate_out=lang=go:../services \
  Models.proto
```

然后可以看到 `services` 目录新生成了 `Models.pb.validate.go` 文件

```go
// Validate checks the field values on OrderMain with the rules defined in the
// proto definition for this message. If any rules are violated, an error is returned.
func (m *OrderMain) Validate() error {
	if m == nil {
		return nil
	}

	// no validation rules for OrderId

	// no validation rules for OrderNo

	// no validation rules for UserId

	if m.GetOrderMoney() <= 1 {
		return OrderMainValidationError{
			field:  "OrderMoney",
			reason: "value must be greater than 1",
		}
	}
  ...
```

修改代码 `services/OrdersService.go`

```go
package services

import (
	"context"
	"fmt"
)

type OrderService struct{}

func (this *OrderService) NewOrder(ctx context.Context, orderRequest *OrderRequest) (*OrderResponse, error) {
	err := orderRequest.OrderMain.Validate()
	if err != nil {
		return &OrderResponse{
			Status:  "error",
			Message: err.Error(),
		}, nil
	}
	fmt.Println(orderRequest.OrderMain)
	return &OrderResponse{
		Status:  "ok",
		Message: "success",
	}, nil
}
```

运行 `grpc` 服务端和 `gateway httpserver` 网关，发送请求

![](../imgs/27_grpc_gateway.jpg)

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/8357d9be080df5fccf27a04be571ab2e6b1ea463#diff-84d5cadaa82f6eeca6dc3d619649d2c2a9176ca7137584962e2235e93101175dL2)

### 16. 流模式入门(上)、场景：批量查询用户积分

为何要用流模式

前面的例子，仅仅是传输比较小的数据， 基本模式是客户端请求----服务端响应

如果是传输较大数据呢？会带来

1、数据包过大导致压力陡增

2、需要等待客户端包全部发送，才能处理以及响应

举例

假设我要从库里取 一批(x万到几十万)。批量查询用户的积分

先创建用户模型 `pbfiles/Models.go`

```protobuf
 message UserInfo{
  int32 user_id=1;
  int32 user_score=2;
}
```

创建Users.proto

```protobuf
syntax="proto3";
package services;
option go_package = ".;services";
import "Models.proto";
message UserScoreRequest{
    repeated UserInfo users=1;
}
message UserScoreResponse{
    repeated UserInfo users=1;
}
service UserService{
    rpc GetUserScore(UserScoreRequest) returns (UserScoreResponse);
}
```

生成部分

```bash
protoc --go_out=plugins=grpc:../services   Users.proto

protoc  --grpc-gateway_out=logtostderr=true:../services Users.proto
```

编写 `UserService`

```go
package services

import "context"

type UserService struct{}

func (*UserService) GetUserScore(ctx context.Context, in *UserScoreRequest) (*UserScoreResponse, error) {
	var score int32 = 101
	users := make([]*UserInfo, 0)
	for _, user := range in.Users {
		user.UserScore = score
		score++
		users = append(users, user)
	}
	return &UserScoreResponse{Users: users}, nil
}
```

`grpc` 无服务端注册用户积分服务

```go
package main

import (
	"gin-grpc/helper"
	"gin-grpc/services"
	"google.golang.org/grpc"
	"net"
)

func main() {
	rpcServer := grpc.NewServer(grpc.Creds(helper.GetServerCreds()))
	services.RegisterProdServiceServer(rpcServer, new(services.ProdService))   // 商品服务
	services.RegisterOrderServiceServer(rpcServer, new(services.OrderService)) // 订单服务
	services.RegisterUserServiceServer(rpcServer, new(services.UserService)) // 用户服务

	listen, _ := net.Listen("tcp", ":8081")
	rpcServer.Serve(listen)
}
```

客户端调用

```go
func main() {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(helper.GetClientCreds()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()
	userClient := services.NewUserServiceClient(conn)
	var i int32
	req := services.UserScoreRequest{}
	req.Users = make([]*services.UserInfo, 0)

	for i = 1; i < 20; i++ {
		req.Users = append(req.Users, &services.UserInfo{UserId: i})
	}
	res, err := userClient.GetUserScore(ctx, &req)
	fmt.Println(res.Users)
}
```

先运行服务端，再运行客户端，可以看到控制台

```bash
[user_id:1  user_score:101 user_id:2  user_score:102 user_id:3  user_score:103 user_id:4  user_score:104 user_id:5  user_score:105 user_id:6  user_score:106 user_id:7  user_score:107 user_id:8  user_score:108 user_id:9  user_score:109 user_id:10  user_score:110 user_id:11  user_score:111 user_id:12  user_score:112 user_id:13  user_score:113 user_id:14  user_score:114 user_id:15  user_score:115 user_id:16  user_score:116 user_id:17  user_score:117 user_id:18  user_score:118 user_id:19  user_score:119]
```

如果同时查询成千上万个，返回的速度会很慢，如果6个6个的查询。

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/f93c0214aa936cf7175e7e57a8221db81b296457#diff-84d5cadaa82f6eeca6dc3d619649d2c2a9176ca7137584962e2235e93101175dL25)

### 17. 服务端流模式、场景：分批发送查询结果

上面实现的是由客户端发起用户积分查询，如果有成千上个用户积分查询，

```go
func (*UserService) GetUserScore(ctx context.Context, in *UserScoreRequest) (*UserScoreResponse, error) {
	var score int32 = 101
	users := make([]*UserInfo, 0)
	for _, user := range in.Users {
		user.UserScore = score
		score++
		users = append(users, user)
	}
	return &UserScoreResponse{Users: users}, nil
}
```

之前是一次性发送 `return &UserScoreResponse{Users: users}, nil`

```go
syntax = "proto3";
package services;
option go_package = ".;services";
import "Models.proto";

message UserScoreRequest {
  repeated UserInfo users = 1;
}

message UserScoreResponse {
  repeated UserInfo users = 1;
}

service UserService {
  rpc GetUserScore(UserScoreRequest) returns (UserScoreResponse) {}
  rpc GetUserScoreByServerStream(UserScoreRequest) returns (stream UserScoreResponse) {} // 服务端流 分批处理和发送
}
```

生成 `.pb.go` 文件

```bash
protoc \
  -I . \
  -I ${GOPATH}/src \
  -I ${GOPATH}/src/github.com/protoc-gen-validate \
  --go_out=plugins=grpc:../services \
  --validate_out=lang=go:../services \
  Users.proto
```

```go
// UserServiceServer is the server API for UserService service.
type UserServiceServer interface {
	GetUserScore(context.Context, *UserScoreRequest) (*UserScoreResponse, error)
	GetUserScoreByServerStream(*UserScoreRequest, UserService_GetUserScoreByServerStreamServer) error
}
```

更新服务端 `services/UserService.go`

```go
// GetUserScoreByServerStream 服务端流
func (this *UserService) GetUserScoreByServerStream(in *UserScoreRequest, stream UserService_GetUserScoreByServerStreamServer) error {
	var score int32 = 101
	users := make([]*UserInfo, 0)
	for index, user := range in.Users {
		user.UserScore = score
		score++
		users = append(users, user)

		if (index+1)%2 == 0 && index > 0 {
			// 每隔2条发送
			err := stream.Send(&UserScoreResponse{Users: users})
			if err != nil {
				return err
			}
			users = (users)[0:0] // 发送完成之后清空切片，方便下次发送
		}
		time.Sleep(time.Second * 1) // 模拟这里处理比较耗时
	}
	if len(users) > 0 { // 发送剩余残留的数据
		err := stream.Send(&UserScoreResponse{Users: users})
		if err != nil {
			return err
		}
	}
	return nil
}
```

拷贝新生成的 `Users.pb.go` 到客户端，

```go
func main() {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(helper.GetClientCreds()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()
	userClient := services.NewUserServiceClient(conn)
	var i int32
	req := services.UserScoreRequest{}
	req.Users = make([]*services.UserInfo, 0)

	for i = 1; i < 20; i++ {
		req.Users = append(req.Users, &services.UserInfo{UserId: i})
	}
	stream, err := userClient.GetUserScoreByServerStream(ctx, &req)
	if err != nil {
		log.Fatal(err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break // 读取到结尾就 break
		}
		if err != nil { // 读取失败，就停止程序运行
			log.Fatal(err)
		}
		fmt.Println(res.Users)
	}
}
```

服务端流，在客户端可以一边接受，一边开协程处理已经接受的数据。

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/60787545fd8c224e5f8a8a5345fd83cc796d790c#diff-dc576b33b5093f4c968f2943df65b7a64afda74e81f771e62d310a3c77e525a5L7)

### 18. 流模式入门：客户端流模式、场景：分批发送请求

上面的是服务端流

![](../imgs/28_stream.jpg)

场景：  客户端批量查询用户积分

1、客户端一次性把用户**列表**发送过去（不是很多,获取列表很快）

2、服务端查询积分比较耗时， 因此查到一部分，就返回一部分。而不是全部查完再返回给客户端。

现在的场景反转：客户端批量查询用户积分

1、客户端一次性把用户列表发送过去（客户端获取列表比较慢）

2、服务端查询积分比较快，此时可以使用 **客户端流模式**

#### 第1步：修改 `Users.proto`

`rpc GetUserScoreByClientStream(stream UserScoreRequest) returns (UserScoreResponse); `

```protobuf
syntax = "proto3";
package services;
option go_package = ".;services";
import "Models.proto";

message UserScoreRequest {
  repeated UserInfo users = 1;
}

message UserScoreResponse {
  repeated UserInfo users = 1;
}

service UserService {
  rpc GetUserScore(UserScoreRequest) returns (UserScoreResponse) {}
  rpc GetUserScoreByServerStream(UserScoreRequest) returns (stream UserScoreResponse) {} // 服务端流 分批处理和发送
  rpc GetUserScoreByClientStream(stream UserScoreRequest) returns (UserScoreResponse); // 客户端流
}
```

#### 第2步：生成 `.pb.go`

```bash
protoc \
  -I . \
  -I ${GOPATH}/src \
  -I ${GOPATH}/src/github.com/protoc-gen-validate \
  --go_out=plugins=grpc:../services \
  --validate_out=lang=go:../services \
  Users.proto
```

```go
// UserServiceServer is the server API for UserService service.
type UserServiceServer interface {
	GetUserScore(context.Context, *UserScoreRequest) (*UserScoreResponse, error)
	GetUserScoreByServerStream(*UserScoreRequest, UserService_GetUserScoreByServerStreamServer) error
	GetUserScoreByClientStream(UserService_GetUserScoreByClientStreamServer) error
}
```

#### 第3步：增加处理方法 `services/UserService.go`

```go
// GetUserScoreByClientStream 客户端流
func (this *UserService) GetUserScoreByClientStream(stream UserService_GetUserScoreByClientStreamServer) error {
	var score int32 = 101
	users := make([]*UserInfo, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF { // 接受结束
			return stream.SendAndClose(&UserScoreResponse{Users: users}) // 发送并关闭
		}
		if err != nil { // 接受出错
			return err
		}
		for _, user := range req.Users {
			user.UserScore = score // 服务端做的业务处理
			score++
			users = append(users, user) // 把处理的结果放入 users
		}
	}
}
```

#### 第4步：拷贝新生成的 `.pb.go` 到客户端

#### 第5步： 客户端调用

```go
func main() {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(helper.GetClientCreds()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()
	userClient := services.NewUserServiceClient(conn)
	stream, err := userClient.GetUserScoreByClientStream(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for j := 1; j <= 3; j++ { // 模拟客户端发送耗时过程
		var i int32 // 发3次，每次发 5 条数据
		req := services.UserScoreRequest{}
		req.Users = make([]*services.UserInfo, 0)

		for i = 1; i < 6; i++ { // 假设是一个耗时的过程
			req.Users = append(req.Users, &services.UserInfo{UserId: i})
		}
		err := stream.Send(&req)
		if err != nil {
			log.Println(err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil { // 读取失败，就停止程序运行
		log.Fatal(err)
	}
	fmt.Println(res.Users)
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/034ab679fcbd7d59b5ebf3c312660f15cc850108#diff-dc576b33b5093f4c968f2943df65b7a64afda74e81f771e62d310a3c77e525a5L7)

### 19.流模式入门(5)：双向流模式的基本套路

![](../imgs/29_stream.jpg)

场景：  客户端批量查询用户积分

1、客户端分批把用户列表发送过去（**客户端获取列表比较慢**）

2、服务端查询积分**也很慢**,所以分批发送过去

此时应该使用 双向流模式

实际业务中应该使用协程方法进行处理。这里只是简单的实现。

#### 第1步： 修改 `Users.proto`

` rpc GetUserScoreByTWS(stream UserScoreRequest) returns (stream UserScoreResponse);`

```protobuf
syntax = "proto3";
package services;
option go_package = ".;services";
import "Models.proto";

message UserScoreRequest {
  repeated UserInfo users = 1;
}

message UserScoreResponse {
  repeated UserInfo users = 1;
}

service UserService {
  rpc GetUserScore(UserScoreRequest) returns (UserScoreResponse) {}
  rpc GetUserScoreByServerStream(UserScoreRequest) returns (stream UserScoreResponse) {} // 服务端流 分批处理和发送
  rpc GetUserScoreByClientStream(stream UserScoreRequest) returns (UserScoreResponse); // 客户端流
  rpc GetUserScoreByTWS(stream UserScoreRequest) returns (stream UserScoreResponse); // 双向流
}
```

#### 第2步：生成 `.pb.go` 

```bash
 protoc \
  -I . \
  -I ${GOPATH}/src \
  -I ${GOPATH}/src/github.com/protoc-gen-validate \
  --go_out=plugins=grpc:../services \
  --validate_out=lang=go:../services \
  Users.proto
```

看到 `Users.pb.go` 文件中

```go
// UserServiceServer is the server API for UserService service.
type UserServiceServer interface {
	GetUserScore(context.Context, *UserScoreRequest) (*UserScoreResponse, error)
	GetUserScoreByServerStream(*UserScoreRequest, UserService_GetUserScoreByServerStreamServer) error
	GetUserScoreByClientStream(UserService_GetUserScoreByClientStreamServer) error
	GetUserScoreByTWS(UserService_GetUserScoreByTWSServer) error
}
```

#### 第3步：服务端实现

```go
// GetUserScoreByTWS 双向流
func (this *UserService) GetUserScoreByTWS(stream UserService_GetUserScoreByTWSServer) error {
	var score int32 = 101
	users := make([]*UserInfo, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF { //接收完了
			return nil
		}
		if err != nil {
			return err
		}
		for _, user := range req.Users {
			user.UserScore = score //这里好比是服务端做的业务处理
			score++
			users = append(users, user)
		}
		err = stream.Send(&UserScoreResponse{Users: users})
		if err != nil {
			fmt.Println(err)
		}
		users = (users)[0:0]
	}
}
```

完整的 `services/UserService.go` 代码

```go
package services

import (
	"context"
	"fmt"
	"io"
	"time"
)

type UserService struct{}

// GetUserScore 普通方法一次性返回
func (*UserService) GetUserScore(ctx context.Context, in *UserScoreRequest) (*UserScoreResponse, error) {
	var score int32 = 101
	users := make([]*UserInfo, 0)
	for _, user := range in.Users {
		user.UserScore = score
		score++
		users = append(users, user)
	}
	return &UserScoreResponse{Users: users}, nil
}

// GetUserScoreByServerStream 服务端流
func (this *UserService) GetUserScoreByServerStream(in *UserScoreRequest, stream UserService_GetUserScoreByServerStreamServer) error {
	var score int32 = 101
	users := make([]*UserInfo, 0)
	for index, user := range in.Users {
		user.UserScore = score
		score++
		users = append(users, user)

		if (index+1)%2 == 0 && index > 0 {
			// 每隔2条发送
			err := stream.Send(&UserScoreResponse{Users: users})
			if err != nil {
				return err
			}
			users = (users)[0:0] // 发送完成之后清空切片，方便下次发送
		}
		time.Sleep(time.Second * 1) // 模拟这里处理比较耗时
	}
	if len(users) > 0 { // 发送剩余残留的数据
		err := stream.Send(&UserScoreResponse{Users: users})
		if err != nil {
			return err
		}
	}
	return nil
}

// GetUserScoreByClientStream 客户端流
func (this *UserService) GetUserScoreByClientStream(stream UserService_GetUserScoreByClientStreamServer) error {
	var score int32 = 101
	users := make([]*UserInfo, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF { // 接受结束
			return stream.SendAndClose(&UserScoreResponse{Users: users}) // 发送并关闭
		}
		if err != nil { // 接受出错
			return err
		}
		for _, user := range req.Users {
			user.UserScore = score // 服务端做的业务处理
			score++
			users = append(users, user) // 把处理的结果放入 users
		}
	}
}

// GetUserScoreByTWS 双向流
func (this *UserService) GetUserScoreByTWS(stream UserService_GetUserScoreByTWSServer) error {
	var score int32 = 101
	users := make([]*UserInfo, 0)
	for {
		req, err := stream.Recv()
		if err == io.EOF { //接收完了
			return nil
		}
		if err != nil {
			return err
		}
		for _, user := range req.Users {
			user.UserScore = score //这里好比是服务端做的业务处理
			score++
			users = append(users, user)
		}
		err = stream.Send(&UserScoreResponse{Users: users})
		if err != nil {
			fmt.Println(err)
		}
		users = (users)[0:0]
	}
}
```



#### 第4步：客户端拷贝新生成的 `.pb.go` 文件

#### 第5步： 客户端调用

```go
func main() {
	conn, err := grpc.Dial(":8081", grpc.WithTransportCredentials(helper.GetClientCreds()))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx := context.Background()
	userClient := services.NewUserServiceClient(conn)
	stream, err := userClient.GetUserScoreByTWS(ctx)
	if err != nil {
		log.Fatal(err)
	}
	var uid int32 = 1         // 发3次，每次发 5 条数据
	for j := 1; j <= 3; j++ { // 模拟客户端发送耗时过程
		req := services.UserScoreRequest{}
		req.Users = make([]*services.UserInfo, 0)

		for i := 1; i < 6; i++ { // 假设是一个耗时的过程
			req.Users = append(req.Users, &services.UserInfo{UserId: uid})
			uid++
		}
		err := stream.Send(&req)
		if err != nil {
			log.Println(err)
		}
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}
		fmt.Println(res.Users)
	}
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/347c3ef7cea466ddbd7db397cd77720ac545754d#diff-dc576b33b5093f4c968f2943df65b7a64afda74e81f771e62d310a3c77e525a5L7)



