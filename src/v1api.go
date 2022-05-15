package bili_live_stream

import (
	"fmt"
	"github.com/tidwall/gjson"
	"os/exec"
	"strconv"
)

const V1API string = "https://api.live.bilibili.com/xlive/web-room/v1/playUrl/playUrl"

func V1Initialization() {
	potPlayerAddress := GetPotplayerAddress()
	realRoomID := GetRealRoomID()
	if realRoomID == -1 {
		V1FormatInit()
	}
	println(strconv.FormatInt(realRoomID, 10))
	param := map[string]string{"cid": strconv.FormatInt(realRoomID, 10), "platform": "hls"}
	V1HandlerQualityUrl(GetChooseQuality(param, "data.quality_description", V1API), param, potPlayerAddress)
}

func V1HandlerQualityUrl(qn int64, param map[string]string, potPlayerAddress string) {
	param["qn"] = strconv.FormatInt(qn, 10)
	result := GetRequest(V1API, param)

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

	cmd := exec.Command(potPlayerAddress, "potplayer://"+content)
	cmd.Run()

	IsOutput(content)
}

func V1FormatInit() {
	fmt.Println()
	V1Initialization()
}
