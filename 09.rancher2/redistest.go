package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	r := gin.New()
	r.Handle("GET", "/", func(c *gin.Context) {
		key := c.Query("key")
		ret, err := rdb.Get(c, key).Result()
		if err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, ret)
		}
	})
	r.Run(":8082")
}
