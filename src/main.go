package main

import "fmt"

func main() {
	fmt.Println("请输入API类型: ")
	fmt.Println("1 V1API")
	fmt.Println("2 V2API")
	var apiType int
	_, _ = fmt.Scanln(&apiType)
	if apiType == 1 {
		V1Initialization()
	} else {
		V2Initialization()
	}
}
