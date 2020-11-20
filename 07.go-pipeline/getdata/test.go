package getdata

import (
	"fmt"
	"pipeline/AppInit"
	"time"

	"log"
)

const sql = "select * from books order by book_id limit ? offset ?"

// GetPage 获取页码
func GetPage(args ...interface{}) InChan {
	in := make(InChan)
	go func() {
		defer close(in)
		for i := 1; i <= 80; i++ { // 第 1 页到 80 页
			in <- i
		}
	}()
	return in
}

// GetData 取数据源的普通方法
func GetData(in InChan) OutChan {
	out := make(OutChan)
	go func() {
		defer close(out)
		for d := range in {
			page := d.(int)
			pagesize := 1000
			booklist := &BookList{make([]*Book, 0), page}
			db := AppInit.GetDB().Raw(sql, pagesize, (page-1)*pagesize).Find(&booklist.Data)
			if db.Error != nil {
				log.Println(db.Error)
			}
			out <- booklist.Data
		}
	}()
	return out
}

// DoData 模拟处理数据，必须写成管道函数
func DoData(in InChan) OutChan {
	out := make(OutChan)
	go func() {
		defer close(out)
		for d := range in {
			v := d.([]*Book)
			time.Sleep(time.Second * 1)
			out <- fmt.Sprintf("处理了%d条数据,%d\n", len(v), time.Now().Unix())
		}
	}()
	return out
}
func PipeTest() {
	p1 := NewPipe()           // 新建一个管道
	p1.SetCmd(GetPage)        // 获取原始数据
	p1.SetPipeCmd(GetData, 5) // 5 表示多路复用
	out := p1.Exec()          // 执行管道

	//for item := range out {
	//	v := item.([]*Book)
	//	fmt.Println(v)
	//}

	p2 := NewPipe()
	p2.SetCmd(func(args ...interface{}) InChan {
		return InChan(out)
	})
	p2.SetPipeCmd(DoData, 2)
	out2 := p2.Exec()

	for item := range out2 {
		fmt.Println(item)
	}
}
