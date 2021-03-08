Rancher2+k8s无脑上手 https://www.jtthink.com/course/play/2757

[toc]

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

如果执行后报错Peer's Certificate has expired

首先安装时间同步软件 `sudo yum install ntp -y` 同步时间 `sudo ntpdate -u 0.centos.pool.ntp.org`

###### 查看 `docker-ce` 版本、`docker-ce-cli` 版本

```sh
sudo yum list docker-ce --showduplicates |sort -r
sudo yum list docker-ce-cli --showduplicates |sort -r
```

######  安装指定版本.(`docker-ce-版本号` `docker-ce-cli-版本号`).

```bash
sudo yum install docker-ce-18.09.9-3.el7 docker-ce-cli-18.09.9-3.el7 containerd.io
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

```dockerfile
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

```dockerfile
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

```dockerfile
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

```dockerfile
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

```dockerfile
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

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/ae9930ea466046a01701f1af96e7c9778bff907c#diff-3d390161cc978954fc3a4e858bce12ad7e05cd1d90a2416b8f916745fa894a69L6)

### 5. 使用Go调用Docker API

上面完成了业务代码的上传和程序自动调用 `ssh` 在容器中进行编译打包镜像。

在云服务器上运行 `docker ps` 是查看本机的 `docker` 服务，

如果在另外一台云服务器，想要查看这台的 `docker` 服务，默认是不可以的。

首先 `docker` 开放 `tcp` 连接，对于 `centos7` 文件在 

`/usr/lib/systemd/system/docker.service`

**1.打开编辑**：

```bash
sudo vi /usr/lib/systemd/system/docker.service
```

**2.注释原有的：**

```bash
ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock
```

**3.添加新的：**

```bash
ExecStart=/usr/bin/dockerd -H unix:///var/run/docker.sock -H tcp://0.0.0.0:2345
```

 -H代表指定docker的监听方式，这里是socket文件文件位置，也就是socket方式，2345就是tcp端口

**4.保存并退出**

**5.重新加载系统服务配置文件（包含刚刚修改的文件）**

```bash
systemctl daemon-reload
```

**6.重启docker服务**

```bash
systemctl restart docker
```

**7.查看端口是否被docker监听**

```bash  
ss -tnl | grep 2345
```

**8.在内部系统上测试端口是否可以使用**

```bash
[custer@localhost ~]$ docker -H tcp://localhost:2345 ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
[custer@localhost ~]$ docker -H tcp://localhost:2346 ps
Cannot connect to the Docker daemon at tcp://localhost:2346. Is the docker daemon running?
```

**8.查看防火墙是否开放2375端口**

```bash  
sudo firewall-cmd --zone=public --query-port=2345/tcp
```

**9.防火墙添加开放2375端口** 

```bash  
sudo firewall-cmd --zone=public --add-port=2345/tcp --permanent
```

**10.重启防火墙** 

 ```bash 
sudo firewall-cmd --reload
 ```

**11.在外部系统上测试端口是否可以使用**

```bash  
telnet 192.168.172.2 2345
```

使用 `docker` `sdk` 的官方文档 https://docs.docker.com/engine/api/sdk/

```bash
go get github.com/docker/docker/client
```

新建文件 `docker.go` 文件

```go
package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
)

func main() {
	cli, err := client.NewClient("tcp://192.168.172.2:2345", "v1.39", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, image := range images {
		fmt.Println(image.ID, image.Labels)
	}
}
```

注意这里的  `version` 需要查看云服务器的 `docker` 版本

```bash
[custer@localhost ~]$ docker version
Client:
 Version:           18.09.9
 API version:       1.39
 Go version:        go1.11.13
 Git commit:        039a7df9ba
 Built:             Wed Sep  4 16:51:21 2019
 OS/Arch:           linux/amd64
 Experimental:      false

Server: Docker Engine - Community
 Engine:
  Version:          18.09.9
  API version:      1.39 (minimum version 1.12)
  Go version:       go1.11.13
  Git commit:       039a7df
  Built:            Wed Sep  4 16:22:32 2019
  OS/Arch:          linux/amd64
  Experimental:     false
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/b5d91e56aaabe7f449596045a8d80cc61e979581#diff-d9b46c18eb099b2ccb6ffd4d1c1301939fab3dee37a5f3e13ff937352719e94dR1)

### 6. Go调用Docker API：启动容器

文档：https://docs.docker.com/engine/api/v1.41/#operation/ContainerCreate

```go
package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
)

type Empty struct{}

func main() {
	cli, err := client.NewClient("tcp://192.168.172.2:2345", "v1.39", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	config := &container.Config{
		WorkingDir: "/app",
		ExposedPorts: map[nat.Port]struct{}{
			"80/tcp": Empty{},
		},
		Image: "alpine:3.13",
		Cmd:   []string{"./myserver"},
	}
	hostConfig := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			"80/tcp": []nat.PortBinding{
				nat.PortBinding{HostPort: "8080"}, //宿主机的端口
			},
		},
		Binds: []string{"/home/custer/myweb:/app"},
	}
	ret, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil,"myweb")
	if err != nil {
		log.Fatal(err)
	}
	err = cli.ContainerStart(ctx, ret.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("容器启动成功,ID是:", ret.ID)
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/e901582fa76d4b2bb886495e1542c2de6c85710a#diff-d9b46c18eb099b2ccb6ffd4d1c1301939fab3dee37a5f3e13ff937352719e94dL4)

### 7. NGINX反代Go API

一般网站的架构会有 nginx 代理 80端口，映射到 api 服务 8080 或 8081。

首先 nginx 容器 https://hub.docker.com/_/nginx 使用 alpine 版本

```bash
docker pull nginx:1.19.6-alpine
```

通过文档可以得知，其配置文件在 

```bash
/etc/nginx/nginx.conf
```

默认 HTML 目录在 

```bash
/usr/share/nginx/html
```

配置文件工具，可以自动生成 https://www.digitalocean.com/community/tools/nginx

这里还有个https://it.baiked.com/tools/nginxconfig.html

启动容器(先启动 go api)

```dockerfile
docker run --name myweb -d \
-v /home/custer/myweb:/app \
-w /app \
alpine:3.13 \
./myserver
```

注意：不需要映射端口，也就是浏览器永远无法直接访问 `myweb`，必须通过 `nginx`，所以这里不需要端口的映射。但是容器和容器之间是可以访问的。

```bash
docker inspect myweb

可以看到 
"IPAddress": "172.17.0.2",
```

找到这个 ip 地址，等下需要使用。

<img src="../imgs/33.nginx.jpg" style="zoom:150%;" />

在云服务器根目录创建文件夹 

```bash
cd && mkdir webconfig && cd webconfig && vi nginx.conf
```

并把自动生成的配置文件复制粘贴

启动 nginx 容器

```dockerfile
docker run --name nginx -d \
-v /home/custer/webconfig/nginx.conf:/etc/nginx/nginx.conf \
-p 80:80 \
nginx:1.19.6-alpine
```

查看日志 

```bash
docker logs nginx
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/78b211b2b63b4120e962a8c9c5bb4d9c38881b3b#diff-d9b46c18eb099b2ccb6ffd4d1c1301939fab3dee37a5f3e13ff937352719e94dL1)

**注意**

> 如要进入alpine容器，命令是（后面的路径不是/bin/bash）：
>
> ` docker exec -it 容器id /bin/sh`

如果可以访问到 nginx 服务器，在 nginx 容器中也可以 curl 到 web 容器，但是在页面访问 web，404报错，应该是 nginx 容器中的 `/etc/nginx/conf.d/default.conf` 文件优先匹配了 `location \ {}` , 需要在下方添加匹配规则 

```nginx
server {
  ...

  location /api/ {

    proxy_pass http://172.17.0.2/;

    proxy_set_header HOST $host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Real-IP $remote_addr;
  }
}
```

也可以新建文件 `webconfig/myweb.conf` 

```nginx
# localhost
server {
  listen 80;
  listen [::]:80;

  server_name localhost;

  location /api/ {

    proxy_pass http://172.17.0.2/;

    proxy_set_header HOST $host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Real-IP $remote_addr;
  }
}
```

使用如下方法启动 nginx 容器

```dockerfile
docker run --name nginx -d \
-v /home/custer/webconfig/myweb.conf:/etc/nginx/conf.d/default.conf \
-p 80:80 \
nginx:1.19.6-alpine
```

### 8. 使用Rancher 1.x 来编排容器(1):基本操作

上面之所以使用 `docker run` 没有使用 `docker-compose` 工具来编排容器，

是因为我们要学习使用 `rancher` 来编排容器。

[rancher/server](https://hub.docker.com/r/rancher/server/)  This image is for Rancher 1.x.

**For Rancher 2.x, see the [rancher/rancher](https://hub.docker.com/r/rancher/rancher/) image.**

Rancher 是个非常强大好用的容器管理（编排）工具。

1.x 主攻容器编排，2.x 则全面转向 k8s 集群管理。

这里使用的是1.6 文档：https://rancher.com/docs/rancher/v1.6/zh/

```bash
docker pull rancher/server

sudo docker run -d --restart=unless-stopped \
--name rc \
-p 8080:8080 rancher/server
```

需要开启8080端口

**查看防火墙是否开放8080端口**

```bash  
sudo firewall-cmd --zone=public --query-port=8080/tcp
```

**防火墙添加开放8080端口** 

```bash  
sudo firewall-cmd --zone=public --add-port=8080/tcp --permanent
```

**重启防火墙** 

 ```bash 
sudo firewall-cmd --reload
 ```

**在外部系统上测试端口是否可以使用**

```bash  
telnet 192.168.172.2 8080
```

访问浏览器查看 http://192.168.172.2:8080/ 

首先点击菜单栏 ADMIN 选择 Access Control 添加 Local Authentication，设定 Login Username 和密码。

然后选择菜单栏 infrastructure 点击 host 添加 Custom，关键在于第5步，复制，在服务器中执行 

```shell
sudo docker run --rm --privileged -v /var/run/docker.sock:/var/run/docker.sock -v /var/lib/rancher:/var/lib/rancher rancher/agent:v1.2.11 http://192.168.172.2:8080/v1/scripts/5892C8EDE6C3284D3270:1609372800000:qhpXECF0cHxC08hKJHBf76jLuEE
```

执行结果 

```shell
[custer@localhost ~]$ sudo docker run --rm --privileged -v /var/run/docker.sock:/var/run/docker.sock -v /var/lib/rancher:/var/lib/rancher rancher/agent:v1.2.11 http://192.168.172.2:8080/v1/scripts/5892C8EDE6C3284D3270:1609372800000:qhpXECF0cHxC08hKJHBf76jLuEE
[sudo] custer 的密码：
Unable to find image 'rancher/agent:v1.2.11' locally
v1.2.11: Pulling from rancher/agent
b3e1c725a85f: Pull complete
6a710864a9fc: Pull complete
d0ac3b234321: Pull complete
87f567b5cf58: Pull complete
063e24b217c4: Pull complete
d0a3f58caef0: Pull complete
16914729cfd3: Pull complete
bbad862633b9: Pull complete
3cf9849d7f3c: Pull complete
Digest: sha256:0fba3fb10108f7821596dc5ad4bfa30e93426d034cd3471f6ccd3afb5f87a963
Status: Downloaded newer image for rancher/agent:v1.2.11

INFO: Running Agent Registration Process, CATTLE_URL=http://192.168.172.2:8080/v1
INFO: Attempting to connect to: http://192.168.172.2:8080/v1
INFO: http://192.168.172.2:8080/v1 is accessible
INFO: Configured Host Registration URL info: CATTLE_URL=http://192.168.172.2:8080/v1 ENV_URL=http://192.168.172.2:8080/v1
INFO: Inspecting host capabilities
INFO: Boot2Docker: false
INFO: Host writable: true
INFO: Token: xxxxxxxx
INFO: Running registration
INFO: Printing Environment
INFO: ENV: CATTLE_ACCESS_KEY=8ED985F780AD9E78AF7F
INFO: ENV: CATTLE_HOME=/var/lib/cattle
INFO: ENV: CATTLE_REGISTRATION_ACCESS_KEY=registrationToken
INFO: ENV: CATTLE_REGISTRATION_SECRET_KEY=xxxxxxx
INFO: ENV: CATTLE_SECRET_KEY=xxxxxxx
INFO: ENV: CATTLE_URL=http://192.168.172.2:8080/v1
INFO: ENV: DETECTED_CATTLE_AGENT_IP=172.17.0.1
INFO: ENV: RANCHER_AGENT_IMAGE=rancher/agent:v1.2.11
INFO: Launched Rancher Agent: 4d89d7cda516c978864f0d7e4747f5cfc60818148fc4ad4ecce5e46873d56395
[custer@localhost ~]$
```

在浏览器中点击 close 可以查看到 rancher 管理的容器，点击 add container 可以创建容器

首先创建 go api 容器，手动执行docker 语句是

```dockerfile
docker run --name myweb -d \
-v /home/custer/myweb:/app \
-w /app \
alpine:3.13 \
./myserver
```

在 rancher 中创建容器，网络使用内置的 Bridge，点击 Create

<img src="../imgs/34.rancher-create-1.jpg" alt="34.rancher-create-1" style="zoom:150%;" />

<img src="../imgs/35.rancher-create-2.jpg" style="zoom:150%;" />

<img src="../imgs/36.rancher-create-3.jpg" style="zoom:150%;" />

创建成功之后可以在 infrastructure -> host 里可以看到在 Standalone Containers 里有 myweb ，点击进入可以看到 

<img src="../imgs/37.rancher-myweb.jpg" style="zoom:150%;" />

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/6f6885a2646396818c06e8366ea889c3a1fe5d0b)

### 9.使用Rancher编排容器(2):nginx反代Go API

对比手动部署 nginx 容器

```dockerfile
docker run --name nginx -d \
-v /home/custer/webconfig/nginx.conf:/etc/nginx/nginx.conf \
-p 80:80 \
nginx:1.19.6-alpine
```

<img src="../imgs/38.rancher-nginx-1.jpg" style="zoom:150%;" />

<img src="../imgs/39.rancher-nginx-2.jpg" style="zoom:150%;" />

<img src="../imgs/40.rancher-nginx-3.jpg" style="zoom:150%;" />

点击 create 创建 nginx 容器

<img src="../imgs/41.rancher-nginx.jpg" style="zoom:150%;" />

在浏览器访问 http://192.168.172.2/api/ping 可以看到 **{**"result": "pong"**}**

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/e14737d74de924b49238a3518bfe6aea898890a3)

### 10. Rancher编排容器(3):多主机启动API

前面是单主机部署 go api 和 nginx 容器，首先再添加一个服务器或虚拟机，

该虚拟机ip地址是 192.168.172.3 用户名是 custer 密码是 root1234

在浏览器页面中删除上面创建的 infrastructure 中的 host ，这里从新创建 host

两台虚拟主机的公网 ip 分别是 192.168.172.2 和 192.168.172.3，

如果是在公有云中使用两台机器的内网 ip ，两个主机放在同一安全组，这样两台主机的端口都能互通。

添加第一个主机

```dockerfile
sudo docker run --rm -e CATTLE_AGENT_IP=192.168.172.2  --privileged -v /var/run/docker.sock:/var/run/docker.sock -v /var/lib/rancher:/var/lib/rancher rancher/agent:v1.2.11 http://192.168.172.2:8080/v1/scripts/5892C8EDE6C3284D3270:1609372800000:qhpXECF0cHxC08hKJHBf76jLuEE
```

Cattle 是 Rancher 自己内置的缺省的编排环境 (其他的还有 k8s、swarm 等等)

点击 close 可以查看已经创建了第一个主机，点击 add host 创建第二个主机

<img src="../imgs/42.rancher-add-host-1.jpg" style="zoom:90%;" />

在第二台主机中运行 

```dockerfile
sudo docker run --rm --privileged -v /var/run/docker.sock:/var/run/docker.sock -v /var/lib/rancher:/var/lib/rancher rancher/agent:v1.2.11 http://192.168.172.2:8080/v1/scripts/5892C8EDE6C3284D3270:1609372800000:qhpXECF0cHxC08hKJHBf76jLuEE
```

<img src="../imgs/43.rancher-add-host-2.jpg" style="zoom:90%;" />

下面创建之前的 Go 服务，点击菜单栏 STACKS 的 ALL

对照手动 的 docker run 命令

```dockerfile
docker run --name myweb -d \
-v /home/custer/myweb:/app \
-w /app \
alpine:3.13 \
./myserver

docker run --name nginx -d \
-v /home/custer/webconfig/myweb.conf:/etc/nginx/conf.d/default.conf \
-p 80:80 \
nginx:1.19.6-alpine
```

注意这里的第二台服务器也需要提前准备好 go 代码和 alpine 和 ngixn 镜像

<img src="../imgs/44.rancher-stack-myweb-1.jpg" style="zoom:150%;" />

![](../imgs/45.rancher-stack-myweb-2.jpg)

<img src="../imgs/45.rancher-stack-myweb-2.jpg" style="zoom:150%;" />

**有防火墙注意需要开放这两个端口 500 和 4500**

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/0764265b90b86658c0ae954578bf106de013bf07)

### 11. Rancher编排容器(4):创建简单多节点负载均衡

<img src="../imgs/46.rancher-stacks-load-balancer.jpg" style="zoom:150%;" />

<img src="../imgs/47.rancher-load-balance.jpg" style="zoom:150%;" />

点击 create 可以看到 

<img src="../imgs/48.rancher-goapi-lb.jpg" style="zoom:150%;" />

点击 Ports: 80/tcp 可以在浏览器上看到访问的 ip，可以判断是部署在哪个服务器上的。

如果上面的部署是使用内网IP，则输入对应的服务器的公网ip即可。

如果两个服务器中输出代码不同，可以查看到轮询的效果。

代码变动 [git commit]()