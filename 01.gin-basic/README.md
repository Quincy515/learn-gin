### 01.最简单的服务启动

分几步
1、随便创建一个文件夹（不要有中文、空格和奇怪的字符串），比如 gin-basic

2、在 gin-basic下创建topic ，代表 话题相关API

3、cd进入topic执行:     `go mod init topic` 

4、在当前目录下执行    `go get  github.com/gin-gonic/gin`

5、用 goland 打开 topic 目录

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.Writer.Write([]byte("Hello"))
	})
	router.Run() // 8080
}
```

返回 JSON 格式

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"username": "custer",
		})
	})
	router.Run() // 8080
}
```

也可以这样修改

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	m := make(map[string]interface{})
	m["username"] = "custer"
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, m)
	})
	router.Run() // 8080
}
```

或这样实现

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Topic struct {
	TopicID    int
	TopicTitle string
}

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, Topic{101, "话题标题"})
	})
	router.Run() // 8080
}
```

也可以使用内置的 `gin.H{}`

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"TopicID": 101, "TopicTitle": "话题标题"})
	})
	router.Run() // 8080
}
```

### 02. API 的 URL 规则设计、带参数的路由

传统的一些API路径设计方式(仔细看看行不行)

- GET /topic/{topic_id} 获取帖子明细
- GET /topic/{user_name} 获取用户发布的帖子列表
- GET /topic/top 获取最热帖子列表

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/topic/:topic_id", func(c *gin.Context) {
		c.String(http.StatusOK, "获取帖子ID为:%s", c.Param("topic_id"))
	})
	router.Run() // 8080
}
```

访问 http://localhost:8080/topic/12 可以看到 

>  获取帖子ID为:12

gin的路由使用的是httprouter库（请自行github一下）,性能好，相对功能够用。

但是目前不支持正则，也不支持 固定路径和参数路径共存。比如：

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/topic/:topic_id", func(c *gin.Context) {
		c.String(http.StatusOK, "获取帖子ID为:%s", c.Param("topic_id"))
	})
	router.GET("/topic/top", func(c *gin.Context) {
		c.String(http.StatusOK, "获取最热帖子列表")
	})
	router.Run() // 8080
}
```

会出现错误

> panic: 'top' in new path '/topic/top' conflicts with existing wildcard ':topic_id' in existing prefix '/:topic_id'

`router.GET("/topic/:id", xxxoo)`
`router.GET("/topic/user", xxxoo)`

甚至 `"/topic/user/:username"`  也会冲突 

#### 重新设计API规则

1、api有版本信息
譬如
`/v1/xxxoo`
`/v2/xxxoo`

2、尽可能使用复数，且含义明确。名词最佳
  `/v1/topics`
  `/v1/users`
  `/v1/getusers`  //不推荐

3、 使用GET参数规划数据展现规则 
`/v1/users` //显示全部或默认条数
`/v1/users?limit=10`  //只显示10条

 `/v1/topics?username=xxxoo` //显示xxoo的帖子

```go
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/v1/topics", func(c *gin.Context) {
		if c.Query("username") == "" {
			c.String(200, "获取帖子列表")
		} else {
			c.String(200, "获取用户=%s的帖子列表", c.Query("username"))
		}
	})
	router.GET("/v1/topics/:topic_id", func(c *gin.Context) {
		c.String(200, "获取topicid=%s的帖子", c.Param("topic_id"))
	})

	router.Run() // 8080
}
```

- 访问 http://localhost:8080/v1/topics 可以看到 `获取帖子列表`
- 访问 http://localhost:8080/v1/topics?username=custer 可以看到 `获取用户=custer的帖子列表`
- 访问 http://localhost:8080/v1/topics/12 可以看到 `获取topicid=12的帖子`

获取参数可以使用 `c.Query()` , 没有就返回为空，可以使用 `c.DefaultQuery("abc", 1)` 设置默认返回值。



















