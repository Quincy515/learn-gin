package goft

import (
	"gopkg.in/yaml.v2"
	"log"
)

type UserConfig map[interface{}]interface{}

// GetConfigValue 递归读取用户配置文件
func GetConfigValue(m UserConfig, prefix []string, index int) interface{} {
	key := prefix[index]
	if v, ok := m[key]; ok {
		if index == len(prefix)-1 {
			return v
		} else {
			index = index + 1
			if mv, ok := v.(UserConfig); ok {
				return GetConfigValue(mv, prefix, index)
			} else {
				return nil
			}
		}
	}
	return nil
}

type ServerConfig struct {
	Port int32
	Name string
}

// SysConfig 系统配置
type SysConfig struct {
	Server *ServerConfig
	Config UserConfig
}

// NewSysConfig 初始化默认配置
func NewSysConfig() *SysConfig {
	return &SysConfig{Server: &ServerConfig{Port: 8080, Name: "goft"}}
}

func InitConfig() *SysConfig {
	config := NewSysConfig()             // 如果没有设定配置文件，使用默认配置
	if b := LoadConfigFile(); b != nil { // 如果设定了配置文件
		err := yaml.Unmarshal(b, config) // 把字符串类型的配置文件映射到 struct
		if err != nil {
			log.Fatal(err)
		}
	}
	return config
}
