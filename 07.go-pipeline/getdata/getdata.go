package getdata

import "fmt"

type Book struct {
	BookId   int    `gorm:"column:book_id"`
	BookName string `gorm:"column:book_name"`
}

func (this *Book) String() string {
	return fmt.Sprintf("bookid:%d,book_name:%s\n", this.BookId, this.BookName)
}

type BookList struct {
	Data []*Book
	Page int
}

// Result 管道结果集
type Result struct {
	Page int
	Err  error
}
