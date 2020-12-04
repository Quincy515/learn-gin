package Helper

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"log"

	"text/template"
)

//给模板使用的函数集
func NewTplFunction() template.FuncMap {
	fm := make(map[string]interface{})
	fm["CamelCase"] = CamelCase
	fm["SnakeCase"] = SnakeCase
	fm["Ucfirst"] = Ucfirst
	fm["Gzip"] = Gzip

	return fm
}

func Gzip(str string) string {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write([]byte(str))
	if err != nil {
		log.Println(err)
		return ""
	}
	err = gz.Close() //这里要关掉，否则取不到数据  也可手工flush.但依然要关掉gz
	if err != nil {
		log.Println(err)
		return ""
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func UnGzip(str string) string {
	dbytes, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Println(err)
		return ""
	}
	read_data := bytes.NewReader(dbytes)
	reader, err := gzip.NewReader(read_data)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer reader.Close()
	ret, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Println("read gzip error:", err)
		return ""
	}
	return string(ret)
}
