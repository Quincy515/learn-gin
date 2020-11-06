package src

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"log"
)

// 缓存装饰器
func CacheDecorator(h gin.HandlerFunc, parm string, redKeyPattern string, empty interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		getID := c.Param(parm)                        // 得到 ID 值
		redisKey := fmt.Sprintf(redKeyPattern, getID) // 拼接 redisKey
		conn := RedisDefaultPool.Get()                // 获取连接池
		defer conn.Close()                            // 不是关闭是还给连接池
		ret, err := redis.Bytes(conn.Do("get", redisKey))
		if err != nil { // 缓存里没有
			h(c) // 执行目标方法
			dbResult, exists := c.Get("dbResult")
			if !exists {
				dbResult = empty // 空对象
			}
			retData, _ := json.Marshal(dbResult)
			conn.Do("setex", redisKey, 20, retData)
			c.JSON(200, dbResult)
			log.Println("从数据库总读取")
		} else { // 缓存里有需要获取的数据
			json.Unmarshal(ret, &empty)
			c.JSON(200, empty)
			log.Println("从 redis 读取")
		}
	}
}
