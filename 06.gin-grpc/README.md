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

#### 第4步：启动服务

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

代码变动 [git commit]()



















