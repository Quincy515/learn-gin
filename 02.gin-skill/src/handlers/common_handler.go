package handlers

import (
	"github.com/gin-gonic/gin"
	"sync"
)

type JSONResult struct {
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Result  interface{} `json:"result"`
}

func NewJSONResult(message string, code string, result interface{}) *JSONResult {
	return &JSONResult{Message: message, Code: code, Result: result}
}

// 每次使用都需要初始化 JSONResult 实例，这个是没有必要的，所以把相关内容放入到 临时对象池
// 临时对象池
var ResultPool *sync.Pool

func init() {
	ResultPool = &sync.Pool{
		New: func() interface{} {
			return NewJSONResult("", "", nil)
		},
	}
}

type ResultFunc func(message string, code string, result interface{})

// 定义函数对 {message: "xxx", code: "10001", result: nil} 进行封装
func OK(c *gin.Context) ResultFunc {
	return func(message string, code string, result interface{}) {
		r := ResultPool.Get().(*JSONResult)
		defer ResultPool.Put(r)
		r.Message = message
		r.Code = code
		r.Result = result
		c.JSON(200, r)
	}
}
