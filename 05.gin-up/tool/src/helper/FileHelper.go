package Helper

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// 获取当前工作目录
func GetWorkDir() string {
	wd, _ := os.Getwd()
	return strings.Replace(wd, "\\", "/", -1)
}

// 判断api文件是否存在
func IsFilePathExist(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func IsFileExist(filepath string) bool {
	fi, err := os.Stat(filepath)
	if err != nil {
		log.Println(err)
		return false
	}
	if fi.IsDir() {
		return false
	}
	return true
}

//读取文件
func ReadFile(f string) []byte {
	file, err := os.OpenFile(f, os.O_RDWR, 0666)
	if err != nil {
		log.Println("open file err:", err)
		return nil
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("read file err:", err)
		return nil
	}
	return b
}

//遍历文件夹，把静态文件读取出来
//map["SERVICE_TPL"]=文件里面的内容
func LoadResource(dir string) map[string]string {
	ret := make(map[string]string)
	dirlist, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal("read dir err", err)
	}
	for _, fi := range dirlist {
		if fi.IsDir() {
			continue //目前只处理一级 ，我们这个项目只需要这样
		}
		//这里统一把.换成下划线 并且大写
		keyName := strings.ToUpper(strings.Replace(fi.Name(), ".", "_", -1))
		ret[keyName] = string(ReadFile(dir + "/" + fi.Name()))
	}
	return ret
}
