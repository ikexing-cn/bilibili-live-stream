package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
)

const URL string = "https://api.live.bilibili.com/xlive/web-room/v1/playUrl/playUrl"

func main() {
	Initialization()
}

func Initialization() {
	fmt.Println("请输入BiliBili直播间房间号：")
	var roomID string
	_, _ = fmt.Scanln(&roomID)
	address := "https://api.live.bilibili.com/room/v1/Room/room_init"
	result := GetRequest(address, map[string]string{"id": roomID})
	realRoomID := HandlerLiveStatus(result)
	if realRoomID == -1 {
		FormatInit()
	}
	HandlerQualityUrl(GetChooseQuality(realRoomID), realRoomID)
}

func GetChooseQuality(realRoomID int) int64 {
	param := map[string]string{"cid": strconv.Itoa(realRoomID), "platform": "h5"}
	result := GetRequest(URL, param)

	var qualityMap = make(map[int64]string)

	gjson.Get(result, "data.quality_description").ForEach(func(key, value gjson.Result) bool {
		qualityMap[key.Int()] = value.String()
		return true
	})

	fmt.Println("请选择清晰度：")
	for k, v := range qualityMap {
		fmt.Println(k, gjson.Get(v, "desc").String())
	}

	var quality int64
	_, _ = fmt.Scanln(&quality)

	return gjson.Get(qualityMap[quality], "qn").Int()
}

func HandlerQualityUrl(qn int64, realRoomID int) {
	param := map[string]string{"cid": strconv.Itoa(realRoomID), "platform": "h5", "qn": strconv.FormatInt(qn, 10)}
	result := GetRequest(URL, param)

	var urls []string

	gjson.Get(result, "data.durl").ForEach(func(key, value gjson.Result) bool {
		value.Get("url").ForEach(func(key, value gjson.Result) bool {
			urls = append(urls, value.String())
			return true
		})
		return true
	})

	fmt.Println("视频地址如下(包含全部线路)：")

	var content string

	for url := range urls {
		fmt.Println(urls[url])
		content += urls[url] + "\n"
	}

	fmt.Println("是否需要输出到文件？(输入任意键执行，输入n取消)")
	var isOutput string
	_, _ = fmt.Scanln(&isOutput)
	if isOutput != "n" {
		WriteString(content)
		fmt.Println("写入完成")
	}

}

func HandlerLiveStatus(result string) int {
	code := gjson.Get(result, "code").Int()
	if code == 60004 {
		fmt.Println("直播间不存在")
		return -1
	}
	if code == 0 {
		liveStatus := gjson.Get(result, "data.live_status").Int()
		if liveStatus != 1 {
			fmt.Println("直播间未开播")
			return -1
		}
	}
	return int(gjson.Get(result, "data.room_id").Int())
}
