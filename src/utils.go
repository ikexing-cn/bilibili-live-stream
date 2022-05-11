package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func FormatInit() {
	fmt.Println()
	Initialization()
}

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return false
	}
	return true
}

func WriteString(content string) {
	fileName := "urls.txt"
	var dstFile *os.File
	if !IsExists(fileName) {
		dstFile, _ = os.Create(fileName)
	} else {
		_ = os.Remove(fileName)
		dstFile, _ = os.Create(fileName) // easy way to io use
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
