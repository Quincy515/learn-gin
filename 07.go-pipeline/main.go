package main

import (
	"fmt"
	"pipeline/getdata"
	"pipeline/v1"
	"pipeline/v2"
	"time"
)

func test(v string) {
	nums := []int{2, 3, 6, 12, 22, 16, 4, 9, 23, 64, 62}
	start := time.Now().Unix()
	if v == "v1" {
		v1.Test(nums)
	} else {
		v2.Test(nums)
	}
	end := time.Now().Unix()
	fmt.Printf("%s测试--用时:%d秒\r\n", v, end-start)
}

func testData() {
	start := time.Now().Unix()
	getdata.ReadData()
	end := time.Now().Unix()
	fmt.Printf("测试--用时:%d秒\r\n", end-start)
}

func main() {
	testData()
}
