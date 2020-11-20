package getdata

import (
	"encoding/csv"
	"fmt"
	"os"
	"pipeline/AppInit"
	"strconv"
	"sync"
	"time"
)

type Book struct {
	BookId   int    `gorm:"column:book_id"`
	BookName string `gorm:"column:book_name"`
}

type BookList struct {
	Data []*Book
	Page int
}

// InChan 管道入参
type InChan chan *BookList

// Result 管道结果集
type Result struct {
	Page int
	Err  error
}

// OutChan 管道数据输出
type OutChan chan *Result

// DataCmd 读取数据源的方法类型
type DataCmd func() InChan

// DataPipeCmd 管道方法的类型
type DataPipeCmd func(in InChan) OutChan

// Pipe 管道函数
func Pipe(c1 DataCmd, cs ...DataPipeCmd) OutChan {
	in := c1()
	out := make(OutChan)
	wg := sync.WaitGroup{}
	for _, c := range cs {
		getChan := c(in)
		wg.Add(1)
		go func(input OutChan) {
			defer wg.Done()
			for v := range input {
				out <- v
			}
		}(getChan)
	}
	go func() {
		defer close(out)
		wg.Wait()
	}()
	return out
}

const sql = "SELECT * FROM books ORDER BY book_id LIMIT ? OFFSET ?"

// ReadData 读取数据源
func ReadData() InChan {
	page := 1
	pageSize := 1000
	in := make(InChan)
	go func() {
		defer close(in)
		for {
			bookList := &BookList{make([]*Book, 0), page}
			db := AppInit.GetDB().Raw(sql, pageSize, (page-1)*pageSize).Find(&bookList.Data)
			if db.Error != nil || db.RowsAffected == 0 {
				break
			}
			in <- bookList
			page++
		}
	}()
	return in
}

// WriteData 执行管道的函数
func WriteData(in InChan) OutChan {
	out := make(OutChan)
	go func() {
		defer close(out)
		for d := range in {
			err := SaveData(d)
			out <- &Result{Page: d.Page, Err: err}
		}
	}()
	return out
}

// SaveData 写入到 csv 文件
func SaveData(data *BookList) error {
	time.Sleep(time.Millisecond * 500)
	file := fmt.Sprintf("./csv/%d.csv", data.Page)
	csvFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer csvFile.Close()
	w := csv.NewWriter(csvFile) // 创建一个新的写入文件流
	header := []string{"book_id", "book_name"}
	export := [][]string{
		header,
	}
	for _, d := range data.Data {
		cnt := []string{
			strconv.Itoa(d.BookId),
			d.BookName,
		}
		export = append(export, cnt)
	}
	err = w.WriteAll(export)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

func Test() {
	out := Pipe(ReadData, WriteData, WriteData, WriteData, WriteData, WriteData)
	for o := range out {
		fmt.Printf("%d.csv文件执行完成,结果:%v\n", o.Page, o.Err)
	}
}
