package getdata

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"pipeline/AppInit"
	"strconv"
)

type Book struct {
	BookId   int    `gorm:"column:book_id"`
	BookName string `gorm:"column:book_name"`
}

type BookList struct {
	Data []*Book
	Page int
}

const sql = "SELECT * FROM books ORDER BY book_id LIMIT ? OFFSET ?"

func ReadData() {
	page := 1
	pageSize := 1000
	for {
		bookList := &BookList{make([]*Book, 0), page}
		db := AppInit.GetDB().Raw(sql, pageSize, (page-1)*pageSize).Find(&bookList.Data)
		if db.Error != nil || db.RowsAffected == 0 {
			break
		}
		err := SaveData(bookList)
		if err != nil {
			log.Println(err)
		}
		page++
	}
}

// SaveData 写入到 csv 文件
func SaveData(data *BookList) error {
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
