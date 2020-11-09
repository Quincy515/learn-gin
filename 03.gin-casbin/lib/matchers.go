package lib

import (
	"fmt"
	"strings"
)

var ADMIN = []string{"admin", "root"}

func init() {
	E.AddFunction("methodMatch", func(arguments ...interface{}) (i interface{}, e error) {
		if len(arguments) == 2 {
			k1, k2 := arguments[0].(string), arguments[1].(string)
			return MethodMatch(k1, k2), nil
		}
		return nil, fmt.Errorf("methodMatch failed.")
	})

	E.AddFunction("isSuper", func(arguments ...interface{}) (i interface{}, e error) {
		if len(arguments) == 1 {
			user := arguments[0].(string)
			return IsSuperAdmin(user), nil
		}
		return nil, fmt.Errorf("superMatch error")
	})
}

func IsSuperAdmin(userName string) bool {
	for _, user := range ADMIN {
		if user == userName {
			return true
		}
	}
	return false
}

// 自定义匹配函数 匹配空格
func MethodMatch(key1 string, key2 string) bool {
	ks := strings.Split(key2, " ")
	for _, s := range ks {
		if s == key1 {
			return true
		}
	}
	return false
}
