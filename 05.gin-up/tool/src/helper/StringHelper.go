package Helper

import (
	"strings"
)

// 首字母大写
func Ucfirst(str string) string {
	if str == "" {
		return str
	}
	clist := []rune(str)
	if len(clist) == 1 {
		return strings.ToUpper(str)
	}
	if clist[0] >= 97 && clist[0] <= 122 {
		clist[0] -= 32
		return string(clist[0]) + string(clist[1:])
	} else {
		return str
	}
}

// 下划线模式
func SnakeCase(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// 驼峰模式
func CamelCase(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
