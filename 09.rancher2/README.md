Rancher2+k8s无脑上手 https://www.jtthink.com/course/play/2757

### 1. 部署第一个程序

1. 简单部署

   ecs + sftp、docker 部署（nginx等）

2. 容器部署

   会使用公有云的其他服务，如日志服务、redis服务等

3. 简易集群

#### 简单部署

##### 1. 写 `main.go` 文件

```go
package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"result": "pong"})
	})
	s := &http.Server{Addr: ":8080", Handler: router}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logrus.Fatal("forced to shutdown: ", err)
	}
	logrus.Println("server exiting, good bye!")
}
```

##### 2. 写 `build.sh` 文件

```sh
echo "start build (linux,amd64)"
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
go build -o build/myserver main.go
echo "complete build (linux,amd64)";
```

##### 3. 上传 `go.mod` 和 `main.go` 文件到装好 `docker` 的云服务器 `home/custer/myweb/src` 目录下。

该 `src` 目录是源码目录，编译好之后是可以删除的，可以不保留在服务器上。

在 `myweb` 目录下生成可执行程序。

### 2. 使用Golang容器来编译程序

#### 1. 在云服务器上安装 `docker`

https://docs.docker.com/engine/install/centos/

###### Uninstall old versions

```sh
sudo yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine
```

###### SET UP THE REPOSITORY

```sh
sudo yum install -y yum-utils

sudo yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo
```

###### 查看 `docker-ce` 版本、`docker-ce-cli` 版本

```sh
yum list docker-ce --showduplicates |sort -r
yum list docker-ce-cli --showduplicates |sort -r
```

######  安装指定版本.(`docker-ce-版本号` `docker-ce-cli-版本号`).

```bash
yum install docker-ce-18.09.9-3.el7 docker-ce-cli-18.09.9-3.el7 containerd.io
```

###### INSTALL DOCKER ENGINE

``` bash 
sudo yum install docker-ce docker-ce-cli containerd.io
```

###### Start Docker

```bash
sudo systemctl start docker
```

###### 查看 `docker` 是否安装成功 

```bash
docker ps
```

###### `docker sudo `

由于docker daemon需要绑定到主机的Unix socket而不是普通的TCP端口，而Unix socket的属主为root用户，所以其他用户只有在命令前添加sudo选项才能执行相关操作。

如果不想每次使用docker命令的时候还要额外的敲一下sudo，可以按照下面的方法配置一下。

1. 创建一个docker组
2. `$ sudo groupadd docker`
3. 添加当前用户到docker组
4. `$ sudo usermod -aG docker $USER`
5. 登出，重新登录shell
6. 验证docker命令是否可以运行
7. `$ docker ps`

重启 `docker`

```bash
sudo systemctl restart docker
```

#### 2. 创建编译容器

`Go` 的官方镜像 https://hub.docker.com/_/golang，选择 `alpine` 版本即可。

```bash
docker pull golang:1.15.7-alpine3.13
```

程序最后也是在 `alpine` 镜像中运行，而不是在云服务器中直接运行。

#### 3. 运行编译容器

```bash
docker run --rm -it \
-v /home/custer/myweb:/app \
-w /app/src \
-e GOPROXY=https://goproxy.cn \
golang:1.15.7-alpine3.13 \
go build -o ../myserver main.go
```

解释：

> --rm：容器停止后自动删除该容器
>
> -it：交互式，实际上也不需要
>
> -v：把当前的 myweb 目录映射到容器中的 app 目录下
>
> -w：把容器当前的工作目录设置为 app/src
>
> -e：环境变量 goproxy 下载第三方库
>
> golang：镜像名
>
> go build -o ../myserver 编译之后的程序是在 app 下命名为 myserver 的程序

直接运行 

```bash
./myserver
```

会报错

```bash
[root@localhost myweb]# ./myserver 
bash: ./myserver: /lib/ld-musl-x86_64.so.1: bad ELF interpreter: No such file or directory
```

因为 `Go` 一些内置库如 `net` 包都用到了 `CGO`。禁用它就不会去寻找对应的 `libc` 库。

这里加入禁用 `CGO` 的环境变量即可。 

**注意：**最后程序直接放在容器里运行。就不需要禁用 `CGO` 了。这里是在云服务器中运行，所以需要禁用。

删除生成的执行程序 `myserver`

```bash
rm myserver
```

运行加入禁用 `CGO` 环境变量的编译容器

```bash
docker run --rm -it \
-v /home/custer/myweb:/app \
-w /app/src \
-e CGO_ENABLED=0 \
-e GOPROXY=https://goproxy.cn \
golang:1.15.7-alpine3.13 \
go build -o ../myserver main.go
```

这里还会下载第三方库，如果把 `gopath` 映射到容器里，就不会每次构建都下载。

### 3. `Goland` 同步程序、自动远程编译

#### 映射 `gopath` 

查看容器的默认 `Gopath`

```bash
docker run --rm -it golang:1.15.7-alpine3.13 go env
```

可以看到，这个容器的默认 `gopath` 在 `/go` 下，所以上面的命令，添加映射目录就行了。

在云服务器中创建 `gopath` 文件夹，映射到容器中，每次 `go build` 就不会拉取第三方库了。

```bash
docker run --rm -it \
-v /home/custer/myweb:/app \
-v /home/custer/gopath:/go \
-w /app/src \
-e CGO_ENABLED=0 \
-e GOPROXY=https://goproxy.cn \
golang:1.15.7-alpine3.13 \
go build -o ../myserver main.go
```

#### 配置 `Goland` 自动同步

```bash
Goland 点击 Tools → Deployment → Configuration → 点击 + 号 → 选择 SFTP →
create New Server name：填写 myweb → 选择 SSH configuration → 
填写 Host、Port、User name、Authentication type：Password、Password →
Save password → TEST CONNECTION → Root path：点击 AUTODETECT →
Mappings 选择 Local path 和 Deployment path →
Excluded Paths → Add excluded path → Local Path
```

#### 远程执行

 `go get golang.org/x/crypto/ssh`

新建文件 `build.go` 通过 `ssh` 连接，远程执行

```go
package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"time"
)

func SSHConnect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	hostKeyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}
	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: hostKeyCallback,
	}
	// connect to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil
}

func main() {
	build_script := "sh /home/custer/myweb/build.sh"
	var stdOut, stdErr bytes.Buffer
	session, err := SSHConnect(" ", " ", "192.168.172.2", 22)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	session.Stdout = &stdOut
	session.Stderr = &stdErr
	err = session.Run(build_script)
	log.Println(stdOut.String())
	log.Println(stdErr.String())
	if err != nil {
		log.Fatal(err)
	}
}
```

在云服务器中新建文件 `/home/custer/myweb/build.sh`

```sh
docker run --rm \
-v /home/custer/myweb:/app \
-v /home/custer/gopath:/go \
-w /app/src \
-e CGO_ENABLED=0 \
-e GOPROXY=https://goproxy.cn \
golang:1.15.7-alpine3.13 \
go build -o ../myserver main.go
```

`build.go` 是可以不上传到云服务器的，

在本地运行 `go run build.go` ，可以看到云服务器已经生成了执行文件 `myserver`

### 4. 使用alpine镜像启动Go API

首先下载纯净的 `Alpine` 镜像 https://hub.docker.com/_/alpine

选择上面对应的编译版本，比如 3.13 版本 `docker pull alpine:3.13`

这个 `alpine` 镜像是用来运行的，`docker images` 查看大小只有 5MB

所以 `Go` 官方镜像用来编译，`alpine` 镜像用来运行程序。

程序发布

```sh
docker run --name myweb -d \
-v /home/custer/myweb:/app \
-w /app \
-p 80:8080 \
alpine:3.13 \
./myserver
```

`./myserver` 是启动容器后立即运行的程序。因为 `myserver` 是在目录 `/home/custer/myweb` 下，而这个目录又映射到了容器里的 `app` 目录，而当前容器工作目录也是 `app`，所以可以直接运行 `./myserver`。

停止和删除服务

`docker stop myweb && docker rm myweb`

新建文件 `tool/scripts/build_script.go`

```go
package scripts

const BuildScript = `
docker run --rm \
-v /home/custer/myweb:/app \
-v /home/custer/gopath:/go \
-w /app/src \
-e GOPROXY=https://goproxy.cn \
golang:1.15.7-alpine3.13 \
go build -o ../myserver main.go
`
```

这样就可以删除云服务器上的 `build.sh` 脚本 `rm build.sh`

在 `Goland` 中也可以排除掉 `tool` 文件夹

```go
Goland 点击 Tools → Deployment → Configuration → Excluded Paths → Add Local Path
```

这样就可以修改 `build.go` 中的代码

```go
package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"rancher2/tool/scripts"
	"time"
)

func SSHConnect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	hostKeyCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}
	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: hostKeyCallback,
	}
	// connect to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil
}

func main() {
	//build_script := "sh /home/custer/myweb/build.sh"
	var stdOut, stdErr bytes.Buffer
	session, err := SSHConnect("custer", "root1234", "192.168.172.2", 22)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	session.Stdout = &stdOut
	session.Stderr = &stdErr
	err = session.Run(scripts.BuildScript)
	log.Println(stdOut.String())
	log.Println(stdErr.String())
	if err != nil {
		log.Fatal(err)
	}
}
```

其实这个 `build.go` 可以做成带参数的命令行工具。

运行 `go run build.go` 这样就可以通过 `ssh` 编译容器。

代码变动 [git commit]()

### 5. 使用Go调用Docker API

