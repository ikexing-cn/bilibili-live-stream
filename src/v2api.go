package bili_live_stream

import (
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
)

const V2API string = "https://api.live.bilibili.com/xlive/web-room/v2/index/getRoomPlayInfo"

func V2Initialization() {
	realRoomID := GetRealRoomID()
	if realRoomID == -1 {
		V2FormatInit()
	}
	param := map[string]string{"platform": "h5", "protocol": "0", "format": "0,1,2", "codec": "0", "room_id": strconv.FormatInt(realRoomID, 10)}
	quality := GetChooseQuality(param, "data.playurl_info.playurl.g_qn_desc", V2API)
	V2HandlerQualityUrl(quality, param)
}

func V2HandlerQualityUrl(quality int64, param map[string]string) {
	param["qn"] = strconv.FormatInt(quality, 10)
	result := GetRequest(V2API, param)
	temp := gjson.Get(result, "data.playurl_info.playurl.stream.0.format.0.codec.0").String()
	baseUrl := gjson.Get(temp, "base_url").String()
	host := gjson.Get(temp, "url_info.0.host").String()
	extra := gjson.Get(temp, "url_info.0.extra").String()

	realUrl := host + baseUrl + extra

	fmt.Println("视频地址如下：")
	fmt.Println(realUrl)

	IsOutput(realUrl)
}

func V2FormatInit() {
	fmt.Println()
	V2Initialization()
}
