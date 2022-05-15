package main

import (
	biliLiveStream "bilibili-live-stream/src"
	"fmt"
)

func main() {
	fmt.Println("请输入API类型: ")
	fmt.Println("1 V1API")
	fmt.Println("2 V2API")
	var apiType int
	_, _ = fmt.Scanln(&apiType)
	if apiType == 1 {
		biliLiveStream.V1Initialization()
	} else {
		biliLiveStream.V2Initialization()
	}
}
