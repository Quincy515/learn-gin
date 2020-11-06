package src

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
)

// MustLogin 必须登录
func MustLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, status := c.GetQuery("token"); !status {
			c.String(http.StatusUnauthorized, "缺少 token 参数")
			c.Abort()
		} else {
			c.Next()
		}
	}
}
func GetTopicDetail(c *gin.Context) {
	tid := c.Param("topic_id")
	topics := Topic{}
	//DBHelper.Find(&topics, tid)
	conn := RedisDefaultPool.Get() // 获取连接池
	defer conn.Close()             // 不是关闭是还给连接池
	redisKey := "topic_" + tid
	ret, err := redis.Bytes(conn.Do("get", redisKey))
	if err != nil { // 缓存里没有
		DBHelper.Find(&topics, tid)
		retData, _ := json.Marshal(topics)
		if topics.TopicID == 0 { // 表示从数据库没有匹配到
			conn.Do("setex", redisKey, 20, retData)
		} else { // 正常数据，50秒缓存
			conn.Do("setex", redisKey, 50, retData)
		}
		c.JSON(200, topics)
		log.Println("从数据库总读取")
	} else {
		json.Unmarshal(ret, &topics)
		c.JSON(200, topics)
		log.Println("从 redis 读取")
	}
}

// NewTopic 单条帖子新增
func NewTopic(c *gin.Context) {
	topic := Topic{}
	err := c.BindJSON(&topic)
	if err != nil {
		c.String(400, "参数错误:%s", err.Error())
	} else {
		c.JSON(200, topic)
	}
}

// NewTopics 多条帖子批量新增
func NewTopics(c *gin.Context) {
	topics := Topics{}
	err := c.BindJSON(&topics)
	if err != nil {
		c.String(400, "参数错误:%s", err.Error())
	} else {
		c.JSON(200, topics)
	}
}

func DeleteTopic(c *gin.Context) {
	// 判断登录
	c.String(200, "删除帖子")
}

// GetTopicList 获取帖子列表
func GetTopicList(c *gin.Context) {
	query := TopicQuery{}
	err := c.BindQuery(&query)
	if err != nil {
		c.String(400, "参数错误:%s", err.Error())
	} else {
		c.JSON(200, query)
	}
}
