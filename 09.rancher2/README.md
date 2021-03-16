Rancher2+k8s无脑上手 https://www.jtthink.com/course/play/2757

[toc]

Centos7 虚拟机安装之后的配置

```bash
su // 切换到 root 用户
sudo vi /etc/sudoers // 添加操作用户的 sudo 权限
sudo vi /etc/sysconfig/network-scripts/ifcfg-ens33 // 修改虚拟机网络
ONBOOT=no  改为  ONBOOT=yes
sudo service network restart // 重启网络服务
ip addr // 查看ip地址
hostnamectl set-hostname k8s.rancher2.com // 修改主机名
hostname // 查看主机名
```



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

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/30f883a85ea8095a7f2a85466713c27ffcb1b20f)

### 12. Rancher编排容器(5):多节点下nginx配置、SSL等

之前学习了单节点 nginx 配置，上面使用了内置的负载均衡，首先停止上面使用的内置的负载均衡。

然后部署 nginx，使用域名的方式，并使用证书。

对照下之前手动 docker run 创建的 nginx

```dockerfile
docker run --name nginx -d \
-v /home/custer/webconfig/myweb.conf:/etc/nginx/conf.d/default.conf \
-p 80:80 \
nginx:1.19.6-alpine
```

这里需要加上证书的位置

```dockerfile
/home/custer/webconfig/certs:/certs
```

对照下，通过 rancher 来创建

<img src="../imgs/49.rancher-mynginx-1.jpg" style="zoom:150%;" />

<img src="../imgs/50.rancher-mynginx-2.jpg" style="zoom:150%;" />

修改 nginx 配置文件 

```nginx
upstream goapi {
  server 10.42.19.222;
  server 10.42.55.71;
}

server {
  listen 80;
  listen [::]:80;

  server_name localhost;

  location /api/ {

    proxy_pass http://goapi/;

    proxy_set_header HOST $host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Real-IP $remote_addr;
  }
}
```

<img src="../imgs/51.rancher-mynginx-3.jpg" style="zoom:150%;" />

访问部署的那个容器的宿主机ip地址，或者添加域名了，直接访问域名即可。

配置证书

```nginx
upstream goapi {
  server 10.42.19.222;
  server 10.42.55.71;
}

server {
  listen 80;
  server_name www.xxx.com;
  charset     utf-8;
  return 301 https://$server_name$request_uri;
}

server {
  listen 443 ssl;
  listen [::]:443 ssl;
  ssl_certificate "/certs/1_www.xxx.com_bundle.crt";
  ssl_certificate_key "/certs/2_www.xxx.com.key";
  ssl_session_timeout 5m;
  ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
  ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:HIGH:!aNULL:!MD5:!RC4:!DHE;
  ssl_prefer_server_ciphers on;

  server_name www.xxx.com;

  location /api/ {

    proxy_pass http://goapi/;

    proxy_set_header HOST $host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Real-IP $remote_addr;
  }
}
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/24c67308c3db06dc9694cf94196b897bd12d5822)

### 13. Rancher编排容器(6):前后端分离简单部署

前后端分离，前端使用 vue + iview 构建一个特别简单的程序

For Vue 3, you should use Vue CLI v4.5 available on `npm` as `@vue/cli`. To upgrade, you need to reinstall the latest version of `@vue/cli` globally:

```bash
yarn global add @vue/cli
# OR
npm install -g @vue/cli
```

Then in the Vue projects, run

```bash
vue upgrade --next
```

Vue projects can quickly be set up with Vite by running the following commands in your terminal.

With npm:

```bash
npm init @vitejs/app <project-name>
cd <project-name>
npm install
npm run dev
```

Or with Yarn:

```bash
yarn create @vitejs/app <project-name>
cd <project-name>
yarn
yarn dev
```

It might occur, that when your username has a space in it like 'Mike Baker' that vite cannot succeed. Have a try with

```bash
create-vite-app <project-name>
```

创建 vue-cli 项目

1. `vue create myui`
2. `cd myui` 打开文件夹
3. `npm install iview --save` 安装 iview 框架
4. `npm install --save vue-router` 安装 view 路由

iview 文档  http://iview.talkingdata.com/#/components/guide/start

5. 启动服务 `npm run serve`
6. 打包服务 `npm run build`

生成编译好的 `dist` 文件，需要上传到服务器，在服务器上新建目录 myui ，上传dist 目录下的内容。

<img src="../imgs/52.myui.jpg" style="zoom:150%;" />

修改配置文件

```nginx
upstream goapi {
  server 10.42.19.222;
  server 10.42.55.71;
}

# localhost
server {
  listen 80;
  listen [::]:80;

  server_name localhost;

  set $base /usr/share/nginx;
  root $base/html;

  location /api/ {

    proxy_pass http://goapi/;

    proxy_set_header HOST $host;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_set_header X-Real-IP $remote_addr;
  }
}
```

注意这里 

```nginx
	set $base /usr/share/nginx;
  root $base/html;
```

只需要把前端文件目录  `myui` 映射到容器里的 `/usr/share/nginx/html` 即可。

更新 rancher mynginx 的配置

<img src="../imgs/53.rancher-myui-1.jpg" style="zoom:100%;" />

<img src="../imgs/54.rancher-myui-2.jpg" style="zoom:100%;" />

点击 upgrade

<img src="../imgs/55.rancher-myui-3.jpg" style="zoom:100%;" />

前后端分离部署，使用 nginx。

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/03e9e259f5764b7d17b082cd682fe212f910313f)

### 14. Rancher编排容器(7):单机部署Redis

https://hub.docker.com/_/redis  

```dockerfile
docker pull redis:5-alpine
```

在第二台服务器 192.168.172.3 上按照 redis 的镜像。

接下来，创建一个文件夹 redis，分别包含

1. data
2. logs（用来存放日志）
3. redis.conf 配置文件

```redis
bind 0.0.0.0

protected-mode yes
 
port 6379

# TCP连接最大积压数
#
# 在大量客户端连接的情况下，应该提高该值，以免客户端连接慢。
# 但该值受系统内核参数的限制，包括 somaxconn 和 tcp_max_syn_backlog。
tcp-backlog 511

# Unix socket.
#
# 指定将用于监听传入连接的Unix套接字的路径。 没有默认值，因此Redis在未指定时不会侦听unix套接字。
#
# unixsocket /tmp/redis.sock
# unixsocketperm 700

# 当连接的客户端连续空闲指定时间后，就断开该连接。指定值为0时禁用超时机制。
timeout 0

# TCP keepalive.
# 周期性检测客户端是否可用
# 如果非零，则在没有通信的情况下使用SO_KEEPALIVE向客户端发送TCP ACK。
# 此选项的合理值为300秒
tcp-keepalive 300

################################# GENERAL #####################################

# 设定是否以守护进程启动服务（默认是no），守护进程会生成 PID 文件 /var/run/redis_6379.pid。
daemonize no

 
supervised no

# 启用守护进程模式时，会生成该文件。
pidfile /var/run/redis_6379.pid

# 指定日志级别
# 日志级别有以下选项:
# debug (适用于开发/测试)
# verbose (很少但很有用的信息)
# notice (信息适中，推荐选项)
# warning (只记录非常重要/关键的消息)
loglevel notice

# 指定保存日志的文件。请注意，如果您使用标准输出进行日志记录，守护进程情况下，日志将发送到/dev/null
logfile /logs/redis.log

# 要启用日志记录到系统记录器，只需将“syslog-enabled”设置为yes，并可选择更新其他syslog参数以满足您的需要。
# syslog-enabled no

# 指定syslog标识。
# syslog-ident redis

# 指定syslog工具。 必须是USER或LOCAL0-LOCAL7之间。
# syslog-facility local0

# 设置数据库数量，默认为16. 默认数据库是 DB 0, 你可以使用 SELECT <dbid> 选择使用的数据库。
# 数据库编号在 0 到 'databases'-1
databases 16

# 启动日志中是否显示redis logo，默认是开启的
always-show-logo yes

################################ SNAPSHOTTING  ################################
#
# 数据持久化:
#
#   save <seconds> <changes>
#
#   指定时间间隔后，如果数据变化达到指定次数，则导出生成快照文件
#
#   示例如下：
#  
#   900 秒(15 分钟)内至少有1个key被修改
#   300 秒(5分钟)内至少有10个key被修改
#   60 秒(1分钟)内至少有10000个key被修改
#
#
#   如果指定 save ""，则相当于清除前面指定的所有 save 设置
#
#   save ""

save 900 1
save 300 10
save 60 10000

# 在启用快照的情况下(指定了有效的 save)，如果遇到某次快照生成失败(比如目录无权限)，
# 之后的数据修改就会被禁止。这有利于用户及早发现快照保存失败，以免更多的数据不能持久化而丢失的风险。
# 当快照恢复正常后，数据的修改会自动开启。
# 如果你有其他的持久化监控，你可以关闭本机制。
stop-writes-on-bgsave-error yes

# 快照中字符串值是否压缩
rdbcompression yes

# 如果开启，校验和会被放在文件尾部。这将使快照数据更可靠，但会在快照生成与加载时降低大约 10% 的性能，追求高性能时可关闭该功能。
rdbchecksum yes

# 指定保存快照文件的名称
dbfilename dump.rdb

# 指定保存快照文件的目录，AOF(Append Only File) 文件也会生成到该目录
dir /data

################################# REPLICATION #################################

# 主从复制。 使用 replicaof 使Redis实例成为另一台Redis服务器的副本。
#
#   +------------------+      +---------------+
#   |      Master      | ---> |    Replica    |
#   | (receive writes) |      |  (exact copy) |
#   +------------------+      +---------------+
#
# 1) Redis复制是异步的，但是如果master与一定数量的副本无法连接，则可以将主服务器配置为停止接受写入。
# 2) 如果再较短时间内与副本失去了连接，当Redis副本与master重新连接时可以执行部分重新同步。因此就要求配置一个合理的 backlog 值。
# 3) 当副本节点重新连接到master时，重新同步复制时自动的，不需要用户干预。
#
# replicaof <masterip> <masterport>

# 如果主服务器受密码保护（使用下面的“requirepass”配置指令），则可以在启动复制同步过程之前告知副本服务器进行身份验证，否则主服务器将拒绝副本服务器请求。
#
# masterauth <master-password>

# 当从库与主库连接中断，或者主从同步正在进行时，如果有客户端向从库读取数据：
# - yes: 从库答复现有数据，可能是旧数据(初始从未修改的值则为空值)
# - no: 从库报错“正在从主库同步”
replica-serve-stale-data yes

# 从库只允许读取
replica-read-only yes

# 无盘同步
#
# -------------------------------------------------------
# WARNING: DISKLESS REPLICATION IS EXPERIMENTAL CURRENTLY
# -------------------------------------------------------

# 新连接(包括连接中断后重连)的从库不可采用增量同步，只能采用全量同步(RDB文件由主库传给从库)，有两种传递方式：
# - 磁盘形式：主库创建子进程，子进程写入磁盘 RDB 文件，再由父进程立即传给从库；
# - 无磁盘形式：主库创建子进程，子进程把 RDB 文件直接写入从库的 SOCKET 连接。
repl-diskless-sync no

# 无盘同步传输间隔（秒）
repl-diskless-sync-delay 5

# 从库向主库PING的间隔（秒）
#
# repl-ping-replica-period 10

# 以下选项设置复制超时：
#
# 1) 从副本的角度来看，在SYNC期间批量传输I / O.
# 2) 从副本（data，ping）的角度来看master超时。
# 3) 从主服务器的角度来看副本超时（REPLCONF ACK ping）。
#
# 确保此值大于为repl-ping-replica-period指定的值非常重要，否则每次主服务器和副本服务器之间的流量较低时都会检测到超时。
#
# repl-timeout 60

# 在SYNC之后禁用副本套接字上的TCP_NODELAY？
#
# 如果选择“yes”，Redis将使用较少数量的TCP数据包和较少的带宽将数据发送到副本。 但这可能会增加数据在副本端出现的延迟，使用默认配置的Linux内核最多可达40毫秒。
#
# 如果选择“no”，则副本上显示的数据延迟将减少，但将使用更多带宽进行复制。
#
# 默认情况下，我们针对低延迟进行优化，但是在非常高的流量条件下，或者当主节点和副本很多时，将其设置为 yes 或许是较好的选择
repl-disable-tcp-nodelay no

# 设置复制积压大小（backlog）。 积压是一个缓冲区，当副本断开连接一段时间后会累积副本数据，因此当副本想要再次重新连接时，通常不需要完全重新同步，只需要部分重新同步就足够了
#
# 复制backlog越大，副本可以断开连接的时间越长。
#
# repl-backlog-size 1mb

# 当master与副本节点断开时间超过指定时间后，将释放复制积压缓冲区（backlog）
#
# 如果设置为0，表示一直不释放复制积压缓冲区
#
# repl-backlog-ttl 3600

# 副本优先级，哨兵模式下，如果主服务器不再正常工作，Redis Sentinel 将优先使用它来选择要升级为主服务器的副本。
#
# 值越低，优先级越高
#
# 优先级为0会将副本标记为无法担任master的角色，因此Redis Sentinel永远不会选择优先级为0的副本进行升级。
#
# 默认情况下，优先级为100。
replica-priority 100

# 如果可用连接的副本数少于N个，并且延迟小于或等于M秒，则master停止接受写入。
#
# 以秒为单位的延迟（必须<=指定值）是根据从副本接收的最后一次ping计算的，通常每秒发送一次。
#
# 例如，要求至少3个在线且滞后时间<= 10秒的副本：
#
# min-replicas-to-write 3
# min-replicas-max-lag 10
#
# 以上两个属性，任意一个设置为0，都会禁用该功能。
#
# 默认情况下，min-replicas-to-write设置为0（功能已禁用），min-replicas-max-lag设置为10。


# 当使用端口转发或网络地址转换（NAT）时，实际上可以通过不同的IP和端口对副本进行访问。
# 副本可以使用以下两个选项，向其主服务器报告一组特定的IP和端口。
#
# 如果只需要覆盖端口或IP地址，则无需使用这两个选项。
#
# replica-announce-ip 5.5.5.5
# replica-announce-port 1234

################################## SECURITY ###################################

# 设置redis访问密码
#
# requirepass foobared

# 命令重命名.
# 对于一些敏感的命令，不希望任意客户端都可以执行，可以改掉默认的名字，新名字只告知特定的客户端来执行。
# 可以命令改名：rename-command CONFIG b840fc02d524045429941cc15f59e41cb7be6c52
# 可以禁用命令：rename-command CONFIG ""，即新名称为空串。
# 需要注意的是，命令改名保存至 AOF 文件或传输至从库，可能导致问题。
# rename-command CONFIG ""


################################### CLIENTS ####################################

# 同一时刻最多可以接纳的客户端数目(Redis 服务要占用其中的大约 32 个文件描述符)。
# 如果客户端连接数达到该上限，新来客户端将被告知“已达到最大客户端连接数”。
#
# maxclients 10000

############################## MEMORY MANAGEMENT ################################

# 内存使用上限
#
# 当内存达到上限时，Redis 将使用指定的策略清除其他键值。
# 如果 Redis 无法清除(或者策略不允许清除键值)，将对占用内存的命令报错，但对只读的命令正常服务。
#
# maxmemory <bytes>


# - volatile-lru: 针对到期的键值，采取 LRU 策略；
# - volatile-lfu: 针对到期的键值，采取 LFU 策略；
# - volatile-random: 针对到期的键值，采取随机策略；
# - allkeys-lru: 针对所有键值，采取 LRU 策略；
# - allkeys-lfu: 针对所有键值，采取 LFU 策略；
# - allkeys-random: 针对所有键值，采取随机策略；
# - volatile-ttl: 删除最近到期的key（次要TTL）
# - noeviction: 不清除任何内容，只是在写入操作时报错。
#
# LRU表示最近最少使用
# LFU意味着最少使用
#
# LRU，LFU和volatile-ttl都是使用近似随机算法实现的。
#
# 默认值是：noeviction
#
# maxmemory-policy noeviction

# 清除键值时取样数量
# LRU，LFU和最小TTL算法不是精确的算法，而是近似算法（为了节省内存），因此您可以调整它以获得速度或精度。
# 默认情况下，Redis将检查五个键并选择最近使用的键，您可以使用以下配置指令更改样本大小。
# 默认值为5会产生足够好的结果。 10:近似非常接近真实的LRU但成本更高的CPU。 3:更快但不是很准确。
#
# maxmemory-samples 5

# 从Redis 5开始，默认情况下，副本将忽略其maxmemory设置（除非在故障转移后或手动将其提升为主设备）。 
# 这意味着key的清除将由主服务器处理，当主服务器中的key清除时，将DEL命令发送到副本。
#
# 此行为可确保主服务器和副本服务器保持一致，但是如果您的副本服务器是可写的，或者您希望副本服务器具有不同的内存设置，
# 并且您确定对副本服务器执行的所有写操作都是幂等的， 然后你可以改变这个默认值（但一定要明白你在做什么）。
#
# replica-ignore-maxmemory yes

############################# LAZY FREEING ####################################

# lazy free 为惰性删除或延迟释放；
# 当删除键的时候,redis提供异步延时释放key内存的功能，
# 把key释放操作放在bio(Background I/O)单独的子线程处理中，
# 减少删除big key对redis主线程的阻塞。有效地避免删除big key带来的性能和可用性问题。

# lazy free的使用分为2类：第一类是与DEL命令对应的主动删除，第二类是过期key删除。

# 针对redis内存使用达到maxmeory，并设置有淘汰策略时；在被动淘汰键时，是否采用lazy free机制；
lazyfree-lazy-eviction no

# 针对设置有TTL的键，达到过期后，被redis清理删除时是否采用lazy free机制；
lazyfree-lazy-expire no

# 针对有些指令在处理已存在的键时，会带有一个隐式的DEL键的操作。如rename命令，当目标键已存在,redis会先删除目标键，如果这些目标键是一个big key,那就会引入阻塞删除的性能问题
lazyfree-lazy-server-del no

# 针对slave进行全量数据同步，slave在加载master的RDB文件前，会运行flushall来清理自己的数据场景，
replica-lazy-flush no

############################## APPEND ONLY MODE ###############################

# 可以同时启用AOF和RDB持久性而不会出现问题。 如果在启动时检查到启用了AOF，Redis将优先加载AOF。
# AOF 持久化机制默认是关闭的
# 
appendonly no

# AOF持久化文件名称默认为 appendonly.aof
appendfilename "appendonly.aof"

# fsync() 调用会告诉操作系统将缓冲区的数据同步到磁盘，可取三种值：always、everysec和no。
# always：实时会极大消弱Redis的性能，因为这种模式下每次write后都会调用fsync。
# no：write后不会有fsync调用，由操作系统自动调度刷磁盘，性能是最好的。
# everysec：每秒调用一次fsync（默认）

# appendfsync always
appendfsync everysec
# appendfsync no

# 在AOF文件 rewrite期间,是否对aof新记录的append暂缓使用文件同步策略,主要考虑磁盘IO开支和请求阻塞时间。默认为no,表示"不暂缓",新的aof记录仍然会被立即同步
no-appendfsync-on-rewrite no

# 当AOF增长超过指定比例时，重写AOF文件，设置为0表示不自动重写AOF文件，重写是为了使aof体积保持最小，而确保保存最完整的数据。
# 这里表示增长一倍
auto-aof-rewrite-percentage 100

#触发aof rewrite的最小文件大小，这里表示，文件大小最小64mb才会触发重写机制
auto-aof-rewrite-min-size 64mb

# AOF文件可能在尾部是不完整的。那redis重启时load进内存的时候就有问题了。
# 
# 如果选择的是yes，当截断的aof文件被导入的时候，会自动发布一个log给客户端然后load。如果是no，用户必须手动redis-check-aof修复AOF文件才可以。默认值为 yes。
aof-load-truncated yes

# 开启混合持久化
# redis保证RDB转储跟AOF重写不会同时进行。
# 当redis启动时，即便RDB和AOF持久化同时启用且AOF，RDB文件都存在，则redis总是会先加载AOF文件，这是因为AOF文件被认为能够更好的保证数据一致性，
# 当加载AOF文件时，如果启用了混合持久化，那么redis将首先检查AOF文件的前缀，如果前缀字符是REDIS，那么该AOF文件就是混合格式的，redis服务器会先加载RDB部分，然后再加载AOF部分。
aof-use-rdb-preamble yes

################################ LUA SCRIPTING  ###############################

# Lua脚本执行超时时间
#
# 设置成0或者负值表示不限时
lua-time-limit 5000

################################ REDIS CLUSTER  ###############################
#
# 开启集群功能，此redis实例作为集群的一个节点
#
# cluster-enabled yes

# 集群配置文件
# 此配置文件不能人工编辑，它是集群节点自动维护的文件，主要用于记录集群中有哪些节点、他们的状态以及一些持久化参数等，方便在重启时恢复这些状态。通常是在收到请求之后这个文件就会被更新
# cluster-config-file nodes-6379.conf

# 集群中的节点能够失联的最大时间，超过这个时间，该节点就会被认为故障。如果主节点超过这个时间还是不可达，则用它的从节点将启动故障迁移，升级成主节点
#
# cluster-node-timeout 15000

# 如果设置成０，则无论从节点与主节点失联多久，从节点都会尝试升级成主节点。
# 如果设置成正数，则cluster-node-timeout*cluster-slave-validity-factor得到的时间，是从节点与主节点失联后，
# 此从节点数据有效的最长时间，超过这个时间，从节点不会启动故障迁移。
# 假设cluster-node-timeout=5，cluster-slave-validity-factor=10，则如果从节点跟主节点失联超过50秒，此从节点不能成为主节点。
# 注意，如果此参数配置为非0，将可能出现由于某主节点失联却没有从节点能顶上的情况，从而导致集群不能正常工作，
# 在这种情况下，只有等到原来的主节点重新回归到集群，集群才恢复运作。
#
# cluster-replica-validity-factor 10

# 主节点需要的最小从节点数，只有达到这个数，主节点失败时，从节点才会进行迁移。
#
# cluster-migration-barrier 1

# 在部分key所在的节点不可用时，如果此参数设置为"yes"(默认值), 则整个集群停止接受操作；
# 如果此参数设置为”no”，则集群依然为可达节点上的key提供读操作。
#
# cluster-require-full-coverage yes

# 在主节点失效期间,从节点不允许对master失效转移
# cluster-replica-no-failover no


########################## CLUSTER DOCKER/NAT support  ########################

#默认情况下，Redis会自动检测自己的IP和从配置中获取绑定的PORT，告诉客户端或者是其他节点。
#而在Docker环境中，如果使用的不是host网络模式，在容器内部的IP和PORT都是隔离的，那么客户端和其他节点无法通过节点公布的IP和PORT建立连接。
#如果开启以下配置，Redis节点会将配置中的这些IP和PORT告知客户端或其他节点。而这些IP和PORT是通过Docker转发到容器内的临时IP和PORT的。
#
# Example:
#
# cluster-announce-ip 10.1.1.5
# cluster-announce-port 6379
# cluster-announce-bus-port 6380

################################## SLOW LOG ###################################

# 执行时间大于slowlog-log-slower-than的才会定义成慢查询，才会被slow-log记录
# 这里的单位是微秒，默认是 10ms 
slowlog-log-slower-than 10000

# 慢查询最大的条数，当slowlog超过设定的最大值后，会将最早的slowlog删除，是个FIFO队列
slowlog-max-len 128

################################ LATENCY MONITOR ##############################

# Redis延迟监视子系统在运行时对不同的操作进行采样，以便收集可能导致延时的数据根源。
#
# 通过LATENCY命令，可以打印图表并获取报告。
#
# 系统仅记录在等于或大于 latency-monitor-threshold 指定的毫秒数的时间内执行的操作。 当其值设置为0时，将关闭延迟监视器。
#
# 默认情况下，延迟监视被禁用，因为如果您没有延迟问题，则通常不需要延迟监视，并且收集数据会对性能产生影响，虽然非常小。 
# 如果需要，可以使用命令“CONFIG SET latency-monitor-threshold <milliseconds>”在运行时轻松启用延迟监视。
latency-monitor-threshold 0

############################# EVENT NOTIFICATION ##############################

# Redis可以向Pub / Sub客户端通知键空间发生的事件。
#
# 例如，如果启用了键空间事件通知，并且客户端对存储在数据库0中的键 foo 执行DEL操作，则将通过Pub / Sub发布两条消息:
#
# PUBLISH __keyspace@0__:foo del
# PUBLISH __keyevent@0__:del foo
# 以 keyspace 为前缀的频道被称为键空间通知（key-space notification）， 而以 keyevent 为前缀的频道则被称为键事件通知（key-event notification）。
# It is possible to select the events that Redis will notify among a set
# of classes. Every class is identified by a single character:
#
#  K     键空间通知，所有通知以 __keyspace@<db>__ 为前缀.
#  E     键事件通知，所有通知以 __keyevent@<db>__ 为前缀
#  g     DEL 、 EXPIRE 、 RENAME 等类型无关的通用命令的通知
#  $     字符串命令的通知
#  l     列表命令的通知
#  s     集合命令的通知
#  h     哈希命令的通知
#  z     有序集合命令的通知
#  x     过期事件：每当有过期键被删除时发送
#  e     驱逐(evict)事件：每当有键因为 maxmemory 策略而被删除时发送
#  A     参数 g$lshzxe 的别名
#
# 输入的参数中至少要有一个 K 或者 E ， 否则的话， 不管其余的参数是什么， 都不会有任何通知被分发。
# 如果只想订阅键空间中和列表相关的通知， 那么参数就应该设为 Kl。将参数设为字符串 "AKE" 表示发送所有类型的通知。
notify-keyspace-events ""

############################### ADVANCED CONFIG ###############################

# hash类型的数据结构在编码上可以使用ziplist和hashtable。
# ziplist的特点就是文件存储(以及内存存储)所需的空间较小,在内容较小时,性能和hashtable几乎一样。
# 因此redis对hash类型默认采取ziplist。如果hash中条目个数或者value长度达到阀值,内部编码将使用hashtable。
# 这个参数指的是ziplist中允许存储的最大条目个数，默认为512，建议为128
hash-max-ziplist-entries 512
# ziplist中允许条目value值最大字节数，默认为64，建议为1024
hash-max-ziplist-value 64

# 当取正值的时候，表示按照数据项个数来限定每个quicklist节点上的ziplist长度。比如，当这个参数配置成5的时候，表示每个quicklist节点的ziplist最多包含5个数据项。
# 当取负值的时候，表示按照占用字节数来限定每个quicklist节点上的ziplist长度。这时，它只能取-1到-5这五个值
# -5: max size: 64 Kb  <-- not recommended for normal workloads
# -4: max size: 32 Kb  <-- not recommended
# -3: max size: 16 Kb  <-- probably not recommended
# -2: max size: 8 Kb   <-- good
# -1: max size: 4 Kb   <-- good
# 性能最高的选项通常为-2（8 Kb大小）或-1（4 Kb大小）。
list-max-ziplist-size -2

# 一个quicklist两端不被压缩的节点个数
# 参数list-compress-depth的取值含义如下：
# 0: 表示都不压缩。这是Redis的默认值
# 1: 表示quicklist两端各有1个节点不压缩，中间的节点压缩。
#    So: [head]->node->node->...->node->[tail]
#    [head], [tail] 不压缩; 内部节点将被压缩.
# 2: [head]->[next]->node->node->...->node->[prev]->[tail]
#    2：表示quicklist两端各有2个节点不压缩，中间的节点压缩
# 3: [head]->[next]->[next]->node->node->...->node->[prev]->[prev]->[tail]
# 3: 表示quicklist两端各有3个节点不压缩，中间的节点压缩。
# etc.
list-compress-depth 0

# 数据量小于等于512用intset，大于512用set
set-max-intset-entries 512

# 数据量小于等于zset-max-ziplist-entries用ziplist，大于zset-max-ziplist-entries用zset
zset-max-ziplist-entries 128
zset-max-ziplist-value 64

# value大小小于等于hll-sparse-max-bytes使用稀疏数据结构（sparse）
# 大于hll-sparse-max-bytes使用稠密的数据结构（dense），一个比16000大的value是几乎没用的，
# 建议的value大概为3000。如果对CPU要求不高，对空间要求较高的，建议设置到10000左右
hll-sparse-max-bytes 3000

#Streams宏节点最大大小。流数据结构是基数编码内部多个项目的大节点树。使用此配置
#可以配置单个节点的字节数，以及切换到新节点之前可能包含的最大项目数
#追加新的流条目。如果以下任何设置设置为0，忽略限制，因此例如可以设置一个
#大入口限制将max-bytes设置为0，将max-entries设置为所需的值
stream-node-max-bytes 4096
stream-node-max-entries 100

# 主动重新散列每100毫秒CPU时间使用1毫秒，以帮助重新散列主Redis散列表（将顶级键映射到值）。 
# Redis使用的散列表实现（请参阅dict.c）执行延迟重新散列：您在重新散列的散列表中运行的操作越多，执行的重复“步骤”就越多，
# 因此如果服务器处于空闲状态，则重新散列将永远不会完成 哈希表使用了一些内存。
activerehashing yes

# 对客户端输出缓冲进行限制可以强迫那些不从服务器读取数据的客户端断开连接，用来强制关闭传输缓慢的客户端。
# 对于normal client，第一个0表示取消hard limit，第二个0和第三个0表示取消soft limit，normal client默认取消限制
client-output-buffer-limit normal 0 0 0

# 对于slave client和MONITER client，如果client-output-buffer一旦超过256mb，又或者超过64mb持续60秒，那么服务器就会立即断开客户端连接。
client-output-buffer-limit replica 256mb 64mb 60

# 对于pubsub client，如果client-output-buffer一旦超过32mb，又或者超过8mb持续60秒，那么服务器就会立即断开客户端连接。
client-output-buffer-limit pubsub 32mb 8mb 60

# 客户端查询缓冲区累积新命令。 默认情况下，它被限制为固定数量，以避免协议失步（例如由于客户端中的错误）将导致查询缓冲区中的未绑定内存使用。 
# 但是，如果您有非常特殊的需求，可以在此配置它，例如我们巨大执行请求。
#
# client-query-buffer-limit 1gb

# 在Redis协议中，批量请求（即表示单个字符串的元素）通常限制为512 MB。 但是，您可以在此更改此限制。
#
# proto-max-bulk-len 512mb

# Redis调用内部函数来执行许多后台任务，例如在超时时关闭客户端的连接，清除从未请求过期的过期密钥等等。
#
# 并非所有任务都以相同的频率执行，但Redis会根据指定的“hz”值检查要执行的任务。
#
# 默认情况下，hz设置为10.提高值时，在Redis处于空闲状态下，将使用更多CPU
# 但同时，当有很多键同时到期时，Redis会响应更快，并且可以更精确地处理超时。
#  
# 范围介于1到500之间，但超过100的值通常不是一个好主意。 大多数用户应使用默认值10，除非仅在需要非常低延迟的环境中将此值提高到100。
hz 10

# 通常，推荐使HZ的值与连接的客户端数量成比例。这有助于避免为每个后台任务调用处理太多客户端，以避免延迟峰值。
#
# 默认情况下默认的HZ值为10。Redis 提供并启用自适应HZ值的功能，当有很多连接的客户端时，该值会临时增加。
#
# 启用动态HZ时，实际配置的HZ将用作基线，但是一旦连接了更多客户端，将根据实际需要使用配置的HZ值的倍数。 
# 通过这种方式，空闲实例将使用非常少的CPU时间，而繁忙的实例将更具响应性。
dynamic-hz yes

# 当一个子进程重写AOF文件时，如果启用下面的选项，则文件每生成32M数据会被同步。
aof-rewrite-incremental-fsync yes

# 当redis保存RDB文件时，如果启用了以下选项，则每生成32 MB数据将对文件进行fsync。 这对于以递增方式将文件提交到磁盘并避免大延迟峰值非常有用。
rdb-save-incremental-fsync yes

# 可以调整Redis LFU（参见maxmemory设置）。 但是，最好使用默认设置，仅在调查如何改进性能以及LFU如何随时间变化后更改它们，这可以通过OBJECT FREQ命令进行检查。
#
# Redis LFU实现中有两个可调参数：计数器对数因子和计数器衰减时间。 在更改它们之前，了解这两个参数的含义非常重要。
#
# LFU计数器每个键只有8位，它的最大值是255，因此Redis使用具有对数行为的概率增量。 给定旧计数器的值，当访问密钥时，计数器以这种方式递增：
#
# 1. A random number R between 0 and 1 is extracted.
# 2. A probability P is calculated as 1/(old_value*lfu_log_factor+1).
# 3. The counter is incremented only if R < P.
#
# The default lfu-log-factor is 10. This is a table of how the frequency
# counter changes with a different number of accesses with different
# logarithmic factors:
#
# +--------+------------+------------+------------+------------+------------+
# | factor | 100 hits   | 1000 hits  | 100K hits  | 1M hits    | 10M hits   |
# +--------+------------+------------+------------+------------+------------+
# | 0      | 104        | 255        | 255        | 255        | 255        |
# +--------+------------+------------+------------+------------+------------+
# | 1      | 18         | 49         | 255        | 255        | 255        |
# +--------+------------+------------+------------+------------+------------+
# | 10     | 10         | 18         | 142        | 255        | 255        |
# +--------+------------+------------+------------+------------+------------+
# | 100    | 8          | 11         | 49         | 143        | 255        |
# +--------+------------+------------+------------+------------+------------+
#
# NOTE: The above table was obtained by running the following commands:
#
#   redis-benchmark -n 1000000 incr foo
#   redis-cli object freq foo
#
# NOTE 2: The counter initial value is 5 in order to give new objects a chance
# to accumulate hits.
#
# The counter decay time is the time, in minutes, that must elapse in order
# for the key counter to be divided by two (or decremented if it has a value
# less <= 10).
#
# The default value for the lfu-decay-time is 1. A Special value of 0 means to
# decay the counter every time it happens to be scanned.
#
# lfu-log-factor 10
# lfu-decay-time 1
 

# 启用主动碎片整理
# activedefrag yes

# 启动活动碎片整理的最小碎片浪费量
# active-defrag-ignore-bytes 100mb

# 启动碎片整理的最小碎片百分比
# active-defrag-threshold-lower 10

# 使用最大消耗时的最大碎片百分比
# active-defrag-threshold-upper 100

# 在CPU百分比中进行碎片整理的最小消耗
# active-defrag-cycle-min 5

# 在CPU百分比达到最大值时，进行碎片整理
# active-defrag-cycle-max 75

# 从set / hash / zset / list 扫描的最大字段数
# active-defrag-max-scan-fields 1000
```

即可以使用之前创建的 STACKS myweb 也可以新建一个

<img src="../imgs/56.rancher-myredis-1.jpg" style="zoom:100%;" />

<img src="../imgs/57.rancher-myredis-2.jpg" style="zoom:100%;" />

<img src="../imgs/58.rancher-myredis-3.jpg" style="zoom:100%;" />

<img src="../imgs/59.rancher-myredis-4.jpg" style="zoom:100%;" />

<img src="../imgs/60.rancher-myredis-5.jpg" style="zoom:100%;" />

进入命令行

```bash
$ docker exec -it df3f6a147d3f redis-cli
127.0.0.1:6379> keys *
(empty list or set)
127.0.0.1:6379> set name custer
OK
127.0.0.1:6379>
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/0c5e5fd3b23c6922f444b0655b50086875a67d17)

### 15. Rancher编排容器(8):单机部署mysql5.7

https://hub.docker.com/_/mysql 

```dockerfile
docker pull mysql:5.7
```

查看文档，默认配置文件在 `/etc/mysql/my.cnf`

接下来创建一个 mysql 文件夹，包含

1. data 目录，存放数据
2. logs 目录，存放日志
3. my.cnf 主配置文件

```mysql
[client]
port=3306
default-character-set=utf8mb4
[mysql]
default-character-set=utf8mb4
 
[mysqld]
log-error=/mysql/logs/error.log
slow_query_log = on
long_query_time=2
slow-query-log-file =/mysql/logs/slow.log

secure-file-priv=''
character-set-client-handshake = FALSE 
#服务器端的端口号
port=3306
 
#MySQL数据库数据文件的目录
datadir=/mysql/data
 
character-set-server = utf8mb4 
collation-server = utf8mb4_general_ci
init_connect='SET NAMES utf8mb4'
 
#MySQL软件的存储引擎
default-storage-engine=INNODB
# Set the SQL mode to strict
#MySQL软件的SQL模式
sql-mode="STRICT_TRANS_TABLES,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION"
#MySQL软件的最大连接数
max_connections=200
#MySQL软件的查询缓存
query_cache_size=0
 
#MySQL软件内存中可以存储临时表的最大值
tmp_table_size=11M
#MySQL软件中可以保留的客户端连接线程数
thread_cache_size=8
#MySQL软件重建索引时允许的最大临时文件的大小
myisam_max_sort_file_size=100G
#MySQL软件重建索引时允许的最大缓存大小
myisam_sort_buffer_size=22M
#MySQL软件中最大关键字缓存大小
key_buffer_size=10M
#MySQL软件全扫描MyISAM表时的缓存大小
read_buffer_size=64K
#MySQL软件可以插入排序好数据的缓存大小
read_rnd_buffer_size=256K
#MySQL软件用户排序时缓存大小
sort_buffer_size=256K
#*** INNODB Specific options ***
#关于INNODB存储引擎参数设置
 
#关于提交日志的时机
innodb_flush_log_at_trx_commit=1
#存储日志数据的缓存区的大小
innodb_log_buffer_size=1M
#缓存池中缓存区大小
innodb_buffer_pool_size=52M
#日志文件的大小
innodb_log_file_size=26M
#允许线程的最大数
innodb_thread_concurrency=9
```

我们在目标主机 192.168.172.3 上操作

```bash
mkdir mysql && cd mysql
mkdir data  logs
vi my.cnf
```

手工 docker run 启动命令

```dockerfile
docker run --name mysql -d \
-v /home/custer/mysql:/mysql \
-v /home/custer/mysql/my.cnf:/etc/mysql/my.cnf \
-e MYSQL_ROOT_PASSWORD=123456 \
mysql:5.7
```

使用 rancher 图形界面，首先 Add Stack

<img src="../imgs/61.rancher-mysql-1.jpg" style="zoom:100%;" />

添加 Add Service

<img src="../imgs/62.rancher-mysql-2.jpg" style="zoom:100%;" />

<img src="../imgs/63.rancher-mysql-3.jpg" style="zoom:100%;" />

<img src="../imgs/64.rancher-mysql-4.jpg" style="zoom:100%;" />

此时 logs 文件夹没有权限处理。会启动错误

处理 logs 文件夹的权限问题

```dockerfile
docker run -it --rm --entrypoint="/bin/bash" mysql:5.7 -c "cat /etc/group"
```

找到 logs 文件夹 

```bash
sudo chown -R xxx logs
```

然后重新启动即可。具体操作如下

```bash
[custer@localhost ~]$ docker run -it --rm --entrypoint="/bin/bash" mysql:5.7 -c "cat /etc/group"
root:x:0:
daemon:x:1:
bin:x:2:
sys:x:3:
adm:x:4:
tty:x:5:
disk:x:6:
lp:x:7:
mail:x:8:
news:x:9:
uucp:x:10:
man:x:12:
proxy:x:13:
kmem:x:15:
dialout:x:20:
fax:x:21:
voice:x:22:
cdrom:x:24:
floppy:x:25:
tape:x:26:
sudo:x:27:
audio:x:29:
dip:x:30:
www-data:x:33:
backup:x:34:
operator:x:37:
list:x:38:
irc:x:39:
src:x:40:
gnats:x:41:
shadow:x:42:
utmp:x:43:
video:x:44:
sasl:x:45:
plugdev:x:46:
staff:x:50:
games:x:60:
users:x:100:
nogroup:x:65534:
mysql:x:999:
[custer@localhost ~]$ id
uid=1000(custer) gid=1000(custer) 组=1000(custer),994(docker) 环境=unconfined_u:unconfined_r:unconfined_t:s0-s0:c0.c1023
[custer@localhost ~]$ cd mysql/
[custer@localhost mysql]$ ls -la
总用量 4
drwxrwxr-x. 4 custer  custer   44 3月   8 06:59 .
drwx------. 6 custer  custer  139 3月   8 06:58 ..
drwxrwxr-x. 2 polkitd custer    6 3月   8 07:16 data
drwxrwxr-x. 2 custer  custer    6 3月   8 06:59 logs
-rw-rw-r--. 1 custer  custer 1723 3月   8 06:59 my.cnf
[custer@localhost mysql]$ sudo chown -R 999 logs
[sudo] custer 的密码：
[custer@localhost mysql]$ docker exec -it 95f58e2a0972 /bin/bash
root@95f58e2a0972:/# mysql -uroot -p
Enter password:
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 2
Server version: 5.7.33-log MySQL Community Server (GPL)

Copyright (c) 2000, 2021, Oracle and/or its affiliates.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql>
```

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/0236e9b926b9ee81d09964c81d92fa7a4a54ac87)

### 16. 快速部署Rancher2和K8s集群

文档 https://www.kubernetes.org.cn/k8s 

Kubernetes（通常称为 K8s）是开源容器集群管理系统，用于自动部署，扩展和管理容器化应用程序。

准备工作

1. 非必要

```bash
docker stop $(docker ps -aq) // stop 停止所有容器
docker rm $(docker ps -aq) // remove 删除所有容器
```

2. 关闭防火墙

```bash
systemctl stop firewalld && systemctl disable firewalld
```

3. 关闭 SELinux

```bash
sudo setenforce 0
sudo sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
```

4. 关闭 swap

```bash
sudo swapoff -a
```

重启docker

```bash
sudo systemctl daemon-reload
sudo systemctl restart docker
```

在主机创建空的文件夹

```bash
[custer@k8s01 ~]$ mkdir rancher
```

启动rancher，安装到主节点或者专门的 rancher 服务器上，一般有三台服务器，

1. rancher 服务器
2. master 主节点 
3. work 工作节点

```bash
sudo docker run -d --restart=unless-stopped -p 8080:80 -p 8443:443 -v /home/custer/rancher:/var/lib/rancher/ rancher/rancher:stable
```

如果是把rancher安装在主节点上，访问主节点的公网ip:8080

**docker 启动rancher，无法访问，报错**

>[custer@k8s01 ~]$ docker logs naughty_ishizaka
>ERROR: Rancher must be ran with the --privileged flag when running outside of Kubernetes

问题解决: 启动时添加 --privileged *作用其实就是启动的 container内的root拥有真正的root权限！！！*

```bash
sudo docker run --privileged -d --restart=unless-stopped -p 8080:80 -p 8443:443 -v /home/custer/rancher:/var/lib/rancher/ rancher/rancher:stable
```

<img src="../imgs/65.rancher2.jpg" style="zoom:100%;" />

<img src="../imgs/66.rancher2-url.jpg" style="zoom:100%;" />

点击添加集群

<img src="../imgs/67.rancher2-mycluster-1.jpg" style="zoom:100%;" />

<img src="../imgs/68.rancher2-mycluster-2.jpg" style="zoom:100%;" />

<img src="../imgs/69.rancher2-mycluer-3.jpg" style="zoom:100%;" />

<img src="../imgs/70.rancher2-mycluster-4.jpg" style="zoom:100%;" />

和rancher1.6类似，需要到主机上执行，针对master服务器这几个都要勾选。

<img src="../imgs/71.rancher-mycluster-5.jpg" style="zoom:100%;" />

针对work节点，只需要 worker。

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/a35e08a7d98d7bdce84914e8c24de8527071878f)

### 17. 创建项目、namespace、初步部署nginx、nodeport

Namespace

> 是对一组资源和对象的抽象集合， **用来将系统内部的对象划分为不同的项目组或用户组**
>
> 常用来隔离不同的用户，比如Kubernetes自带的服务一般运行在kube-system namespace中。

<img src="../imgs/73.k8s-namespace-1.jpg" style="zoom:100%;" />

<img src="../imgs/74.k8s-namespace-2.jpg" style="zoom:100%;" />

<img src="../imgs/75.k8s-namespace-3.jpg" style="zoom:100%;" />

<img src="../imgs/76.k8s-namespace-4.jpg" style="zoom:100%;" />

<img src="../imgs/77.k8s-namespace-5.jpg" style="zoom:100%;" />

workload

> Pod是所有业务类型的基础，也是K8S管理的最小单位级，
>
> 可以理解为  **它是一个或多个容器的组合**

<img src="../imgs/72.k8s.png" style="zoom:50%;" />

<img src="../imgs/78.k8s-workload-1.jpg" style="zoom:100%;" />

NodePort

> 在所有节点（虚拟机）上开放一个特定端口，任何发送到该端口的流量都被转发到对应服务
>
> 端口范围 30000-32767

Nginx    [nginx:1.18-alpine](https://hub.docker.com/_/nginx)

<img src="../imgs/79.k8s-nodeport-1.jpg" style="zoom:100%;" />

<img src="../imgs/79.k8s-nodeport-2.jpg" style="zoom:100%;" />

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/e70b13994838ba16f423061673f9e0dbab00aad4)

### 18. 部署一个go api到k8s集群中(初级)、Hostport

<img src="../imgs/80.k8s-goapi-1.jpg" style="zoom:100%;" />

HostPort 

> 直接将容器的端口与所调度的节点上的端口进行映射。

之前使用的单机部署

```dockerfile
docker pull alpine:3.13

docker run -d --name myweb \
-v /home/custer/myweb:/app \
-w /app \
-p 8081:80 \
alpine:3.13 \
./myserver
```

<img src="../imgs/81.k8s-goapi-2.jpeg" style="zoom:100%;" />

<img src="../imgs/82.k8s-goapi-3.jpg" style="zoom:100%;" />

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/35dc05640e44d1822c915cf2296a4ac2eb5e8cca)

### 19. 2个go api进行负载均衡(ingress)

目前环境 

rancher服务器 master 主机 192.168.172.4

worker 主机 192.168.172.5 

现在在这两个服务器上都上传go的源代码文件，并使用 docker 进行 build

```dockerfile
docker run --rm -it \
-v /home/custer/myweb:/app \
-w /app/src \
-e GOPROXY=https://goproxy.cn \
golang:1.15.7-alpine3.13 \
go build -o ../myserver main.go
```

<img src="../imgs/84.k8s-goapi-4.jpg" style="zoom:100%;" />

之前是把 go api 部署在单节点，现在进行修改，

<img src="../imgs/85.k8s-goapi-5.jpg" style="zoom:100%;" />

修改为

<img src="../imgs/86.k8s-goapi-6.jpg" style="zoom:100%;" />

删除残余的出错节点

<img src="../imgs/87.k8s-goapi-7.jpg" style="zoom:100%;" />

添加负载均衡

<img src="../imgs/88.k8s-lb-1.jpg" style="zoom:100%;" />

Ingress 

> 相当与一个7层负载均衡器，理解为进行反代并定义规则的一个api 对象
>
> ingressController通过监听 ingress api 转化为各自的配置(常用的有nginx-ingress，trafik-ingress)

```bash
                                    ingress
                             /          |             \
                     xxx1.com        xxx2.com            xxx3.com
                     /                  |                   \  
            workload/service     workload/service       workload/service
             /        \              /      \               /        \
           pod         pod         pod       pod           pod       pod
```

<img src="../imgs/89.k8s-lb-2.jpg" style="zoom:100%;" />

<img src="../imgs/90.k8s-lb-3.jpg" style="zoom:100%;" />

添加了负载均衡，就可以在外网访问。

<img src="../imgs/91.k8s-lb-4.jpg" style="zoom:100%;" />

相比之前的就是80是外网访问，32602/tcp是内网访问，如果之前在添加rancher时使用的是内网ip地址，则32602是不能被外网访问的。

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/92e0f4104a6ca73169a90838151f74a46f6f5a3e)

### 20. ClusterIP模式、服务发现基本入门和调用

<img src="../imgs/92.k8s-clusterip-1.jpg" style="zoom:100%;" />

<img src="../imgs/93.k8s-clusterip-2.jpg" style="zoom:100%;" />

把 NodePort 切换成 ClusterIP

> ClusterIP，创建集群内的服务，应用只要在集群内都可以访问，外部（如公网）无法访问。

<img src="../imgs/94.k8s-clusterip-3.jpg" style="zoom:100%;" />

上传test.go文件到master服务器 `/home/custer/mytest/src`

```go
package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

func main(){

	r:=gin.Default()
	r.Handle("GET","/", func(context *gin.Context) {
		 host:=context.Query("host")
		 if host==""{
			 context.JSON(400,gin.H{"error":"no host!"})
			 return
		 }
		 // http://mygo
		 rsp,err:=http.Get("http://"+host)
		 if err!=nil{
		 	context.JSON(400,gin.H{"error":err})
		 }else{
		 	b,err:=ioutil.ReadAll(rsp.Body)
		 	if err!=nil{
				context.JSON(400,gin.H{"error":err})
			}else{
				context.JSON(200,gin.H{"message":string(b)})
			}

		 }
	})


	err:=r.Run(":80")
	if err!=nil{
		log.Fatal(err)
	}
}
```

新建文件 `go.mod`

```go
module test
go 1.15
```

使用docker编译go代码

```dockerfile
docker run --rm -it \
-v /home/custer/mytest:/app \
-w /app/src \
-e GOPROXY=https://goproxy.cn \
golang:1.15.7-alpine3.13 \
go build -o ../myserver test.go
```

然后使用 rancher 部署mytest服务

手工部署

```dockerfile
docker run -d --name myweb \
-v /home/custer/mytest:/app \
-w /app \
-p 8081:80 \
alpine:3.13 \
./myserver
```

<img src="../imgs/95.k8s-mytest-1.jpeg" style="zoom:100%;" />

这样就部署好了 mytest服务，下面如何访问设置为 clusterIP 的 mygo 服务。

<img src="../imgs/96.k8s-mytest-2.jpg" style="zoom:100%;" />

选择服务发现 Service Discovery

<img src="../imgs/97.k8s-mytest-3.jpg" style="zoom:100%;" />

服务发现简单机制

Rancher2.4 使用 k8s-coredns 作为服务发现基础

- 在同一个命名空间内：可以通过 service_name 直接解析
- 在不同命名空间内：service_name.namespace_name

在上创建 workloads 时的 mygo 和 mytest 名称，会自动添加一条解析记录，

在容器内可以直接使用该名称进行访问。

访问 http://192.168.172.4:8081/?host=mygo 可以自动解析 mygo 访问到内容。

这是访问同一个namespace 下的 workload，现在再部署一个 mytest 在不同的 namespace 下。

<img src="../imgs/98.k8s-mytest-4.jpeg" style="zoom:100%;" />

不在同一个命名空间，进行访问 http://192.168.172.4:8082/?host=mygo.myweb

<img src="../imgs/99.k8s-mytest-5.jpg" style="zoom:100%;" />

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/c58fb3b2dba77a404a85258729664d8c1e587aab)

### 21. 补充：部署nfs服务进行跨主机文件共享

在前面部署服务的时候，如果需要在不同的主机上部署相同的api，进行负载均衡的设置。

在前面是使用文件挂载，使用 -v 进行映射。

在每一台机器上上传相同的目录进行挂载。如果节点特别多，都需要对每台服务器进行上传，

使用 nfs 服务进行跨主机文件共享。

Network File System

> 通过网络，让不同机器、不同的操作系统可以共享彼此的文件。

每台主机都是 centos7 系统，使用相同的用户名，正规来说应该再准备一个 nfs 服务器。

这里使用 work 服务器 192.168.172.5 部署 nfs。

安装和配置：

```bash
sudo yum -y install nfs-utils
sudo vi /etc/sysconfig/nfs

// 加入 
LOCKD_TCPPORT=30001#TCP锁使用的端口
LOCKD_UDPPORT=30002#UDP锁使用的端口
MOUNTD_PORT=30003#挂载使用的端口
STATD_PORT=30004#状态使用的端口
```

因为这里部署k8s就把所有防火墙都关闭了。固定端口。

启动/重启服务：

```bash
sudo systemctl restart rpcbind.service
sudo systemctl restart nfs-server.service
```

开机启动：

```bash
sudo systemctl enable rpcbind.service
sudo systemctl enable nfs-server.service
```

新建共享目录：

```bash
[custer@k8s02 ~]$ mkdir goapi
[custer@k8s02 ~]$ ls
goapi  myweb
[custer@k8s02 ~]$ cd goapi/
[custer@k8s02 goapi]$ echo "abc" >> test.text
[custer@k8s02 goapi]$ ls
test.text
[custer@k8s02 goapi]$ cat test.text
abc
```

编辑共享目录：

```bash
sudo vi /etc/exports
```

写入如下内容-共享目录 内部IP地址 配置

```bash
/home/custer/goapi 192.168.172.5/24(rw,async)
```

| 参数           | 作用                                                         |
| -------------- | ------------------------------------------------------------ |
| ro             | 只读                                                         |
| rw             | 读写                                                         |
| root_squash    | 当NFS客户端以root管理员访问时，映射为NFS服务器的匿名用户     |
| no_root_squash | 当NFS客户端以root管理员访问时，映射为NFS服务器的root管理员   |
| all_squash     | 无论NFS客户端用什么账户访问，均映射为NFS服务器的匿名用户     |
| sync           | 同时将数据写入到内存与硬盘中，保证不丢失数据                 |
| async          | 优先将数据保存到内存，然后再写入硬盘；这样效率更高，但可能会丢失数据 |

查看挂载

```bash
showmount -e localhost		
```

会发现没有，于是重启 nfs 服务

```bash
sudo systemctl restart nfs-server.service
```

这样 nfs 服务器就启动成功了，可以在其他服务器上进行共享。

来到另一台服务器上

```bash
sudo yum -y install nfs-utils
```

这样就好了，不需要启动nfs服务，直接执行如下操作：

```bash
showmount -e 192.168.172.5
```

尝试进行挂载

```bash
sudo mount -t nfs 192.168.172.5:/home/custer/goapi /home/custer/goapi
```

卸载只要 

```bash
sudo umount /home/custer/goapi
```

查看挂载 `df -h`

```bash
[custer@k8s01 ~]$ mkdir goapi
[custer@k8s01 ~]$ sudo mount -t nfs 192.168.172.5:/home/custer/goapi /home/custer/goapi
[custer@k8s01 ~]$ cat goapi/test.text
abc
[custer@k8s01 ~]$ df -h
文件系统                          容量  已用  可用 已用% 挂载点
devtmpfs                          1.9G     0  1.9G    0% /dev
tmpfs                             1.9G     0  1.9G    0% /dev/shm
tmpfs                             1.9G   14M  1.9G    1% /run
tmpfs                             1.9G     0  1.9G    0% /sys/fs/cgroup
/dev/mapper/centos-root            17G  9.0G  8.1G   53% /
/dev/sda1                        1014M  151M  864M   15% /boot
tmpfs                             378M     0  378M    0% /run/user/1000
192.168.172.5:/home/custer/goapi   17G  5.5G   12G   33% /home/custer/goapi
```

这样只需要在一台服务器上进行文件的上传，其他服务器只要挂载该目录，就都有了。

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/cc3d437cc3e77ebd3bf09f116c5526d9cac4a9d5)

### 22. 使用Rancher创建PV和PVC、运行GoAPI

之前创建goapi的方式，是在每台主机上创建文件夹，把程序文件分别上传到每台主机，

使用 nfs 就可以挂载 nfs 服务器上的文件夹。

我们把 nfs 服务器放到了 work 机器，192.168.172.5。

先把我们的 goapi 程序上传到 nfs 服务器

```bahs
[custer@k8s02 goapi]$ ls
myserver  src
```

然后把之前在其他服务器上的挂载 `umount`

```bash
[custer@k8s01 src]$ df -h
文件系统                          容量  已用  可用 已用% 挂载点
devtmpfs                          1.9G     0  1.9G    0% /dev
tmpfs                             1.9G     0  1.9G    0% /dev/shm
tmpfs                             1.9G   14M  1.9G    1% /run
tmpfs                             1.9G     0  1.9G    0% /sys/fs/cgroup
/dev/mapper/centos-root            17G  9.1G  8.0G   54% /
/dev/sda1                        1014M  151M  864M   15% /boot
tmpfs                             378M     0  378M    0% /run/user/1000
192.168.172.5:/home/custer/goapi   17G  5.5G   12G   33% /home/custer/goapi
[custer@k8s01 src]$ sudo umount /home/custer/goapi
umount.nfs: /home/custer/goapi: device is busy
[custer@k8s01 src]$ cd
[custer@k8s01 ~]$ sudo umount /home/custer/goapi
[custer@k8s01 ~]$ df -h
文件系统                 容量  已用  可用 已用% 挂载点
devtmpfs                 1.9G     0  1.9G    0% /dev
tmpfs                    1.9G     0  1.9G    0% /dev/shm
tmpfs                    1.9G   14M  1.9G    1% /run
tmpfs                    1.9G     0  1.9G    0% /sys/fs/cgroup
/dev/mapper/centos-root   17G  9.1G  8.0G   54% /
/dev/sda1               1014M  151M  864M   15% /boot
tmpfs                    378M     0  378M    0% /run/user/1000
[custer@k8s01 ~]$
```

实际上在正式上线部署，也不能在每台机器上都放一个 myserver 可执行程序。

下面先把前面的内容做几个改动

```bash
sudo vi /etc/exports
```

修改成

```bash
/home/custer/goapi 192.168.172.5/24(rw,async,insecure,no_root_squash)
```

然后执行 `exportfs -a` 重新加载配置

```bash
sudo exportfs -a
```

参数说明：

> Root_squash(默认): 将来访的root用户映射为匿名用户或用户组
>
> no_root_squash: 来访的root用户保持root账号权限
>
> no_all_squash(默认): 访问用户先与本机用户匹配，匹配失败后再映射为匿名用户或用户组
>
> all_squash: 将来访的所有用户映射为匿名用户或用户组
>
> secure(默认): 限制客户端之鞥呢从小于1024的tcp/ip端口连接服务器
>
> insecure: 允许客户端从大于1024的tcp/ip端口连接服务器
>
> anonuid: 匿名用户的uid值，通常是nobody或nfsnobody，可以在此处自行设定
>
> anongid: 匿名用户的gid值
>
> no_subtree_check: 如果 nfs 输出的是一个子目录，则无需检查其父目录的权限(可提高效率)

PV (Persistent Volume持久卷) PVC(Persistent Volume Claim)

> PV: 定义 Volume 的类型，挂载目录，远程存储服务器等
>
> PVC：定义 Pod 想要使用的持久化属性，比如存储大小、读写权限等

StorageClass 存储类

>  PV的模板，自动为PVC创建PV。

<img src="../imgs/100.k8s-pv-1.jpg" style="zoom:100%;" />

PV 全局的集群资源，不针对某一个项目，或命名空间。

添加 PV

<img src="../imgs/101.k8s-pv-2.jpg" style="zoom:100%;" />

<img src="../imgs/102.k8s-pvc-1.jpg" style="zoom:100%;" />

<img src="../imgs/103.k8s-pvc-2.jpg" style="zoom:100%;" />

<img src="../imgs/104.k8s-pvc-3.jpg" style="zoom:100%;" />

这样 PV 就可以 PVC 绑定了，PVC 也定义了

<img src="../imgs/104.k8s-pvc-3.jpg" style="zoom:100%;" />

清空了之前的工作服在 workload，部署新的工作负载

<img src="../imgs/105.k8s-workload-1.jpeg" style="zoom:100%;" />

<img src="../imgs/106.k8s-workload-2.jpg" style="zoom:50%;" />

可以发现，在k8s01上是没有程序运行文件的，但是通过PV和PVC可以部署。

代码变动 [git commit](https://github.com/custer-go/learn-gin/commit/ab0fd40258f97b54494d8d1252f819e8f7ed3be7)

### 23. k8s负载均衡加域名和路径重写

<img src="../imgs/107.k8s-lb-1.jpeg" style="zoom:100%;" />

添加负载均衡

<img src="../imgs/108.k8s-lb-2.jpg" style="zoom:100%;" />

<img src="../imgs/109.k8s-lb-3.jpg" style="zoom:100%;" />

<img src="../imgs/110.k8s-lb-4.jpg" style="zoom:100%;" />

假设想像nginx一样，配置不同的路径访问api

<img src="../imgs/111-k8s-lb-5.jpg" style="zoom:100%;" />

现在访问 [mylb.myweb.192.168.172.4.xip.io/api](http://mylb.myweb.192.168.172.4.xip.io/api) 显示 404 page not found，

因为 nginx 配置是自动配置路径，在 rancher 里或 k8s 里需要配置，文档：

[在 Minikube 环境中使用 NGINX Ingress 控制器配置 Ingress | Kubernetes](https://kubernetes.io/zh/docs/tasks/access-application-cluster/ingress-minikube/)

<img src="../imgs/112.k8s-lb-6.jpeg" style="zoom:100%;" />

这时访问 [mylb.myweb.192.168.172.4.xip.io/api/ping](http://mylb.myweb.192.168.172.4.xip.io/api/ping) 就可以看到显示内容了。

代码变动 [git commit]()
