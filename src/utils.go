package bili_live_stream

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GetRealRoomID() int64 {
	fmt.Println("请输入BiliBili直播间房间号：")
	var roomID string
	_, _ = fmt.Scanln(&roomID)
	address := "https://api.live.bilibili.com/room/v1/Room/room_init"
	result := GetRequest(address, map[string]string{"id": roomID})
	realRoomID := HandlerLiveStatus(result)
	return realRoomID
}

func HandlerLiveStatus(result string) int64 {
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
	return gjson.Get(result, "data.room_id").Int()
}

func GetChooseQuality(param map[string]string, path string, api string) int64 {
	result := GetRequest(api, param)

	var qualityMap = make(map[int64]string)

	gjson.Get(result, path).ForEach(func(key, value gjson.Result) bool {
		qualityMap[key.Int()] = value.String()
		return true
	})

	fmt.Println("请选择清晰度（请尽量选择网页上存在的清晰度，不然可能会有问题）：")
	for k, v := range qualityMap {
		fmt.Println(k, gjson.Get(v, "desc").String())
	}

	var quality int64
	_, _ = fmt.Scanln(&quality)

	return gjson.Get(qualityMap[quality], "qn").Int()
}

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return false
	}
	return true
}

func IsOutput(content string) {
	fmt.Println("是否需要输出到文件？(输入任意键执行，输入n取消)")
	var isOutput string
	_, _ = fmt.Scanln(&isOutput)
	if isOutput != "n" {
		WriteString(content)
		fmt.Println("写入完成")
	}

	_, _ = fmt.Scanln()
}

func WriteString(content string) {
	fileName := "urls.txt"
	var dstFile *os.File
	if !IsExists(fileName) {
		dstFile, _ = os.Create(fileName)
	} else {
		_ = os.Remove(fileName)
		dstFile, _ = os.Create(fileName)
	}

	defer func(dstFile *os.File) {
		_ = dstFile.Close()
	}(dstFile)

	_, _ = dstFile.WriteString(content)
}

func GetRequest(address string, params map[string]string) string {
	paramsTemp := url.Values{}
	Url, _ := url.Parse(address)
	for k, v := range params {
		paramsTemp.Set(k, v)
	}

	Url.RawQuery = paramsTemp.Encode()

	client := &http.Client{}
	println(Url.String())
	req, err := http.NewRequest("GET", Url.String(), strings.NewReader(""))
	if err != nil {
		log.Println(err)
	}
	resp, _ := client.Do(req)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
