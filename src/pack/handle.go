package entriesresources

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Entries struct {
	Assets            []Asset
	InitialAssetCount int
	Version           int
}

type Asset struct {
	Id                 string
	AccessibilityLabel string
	Url1080H264        string `json:"url-1080-H264"`
	Url1080HDR         string `json:"url-1080-HDR"`
	Url1080SDR         string `json:"url-1080-SDR"`
	Url4KHDR           string `json:"url-4K-HDR"`
	Url4KSDR           string `json:"url-4K-SDR"`
}

func HandleRetries(path string) {
	// fmt.Printf(path)
	// entriesresources.Test()
	file, _ := ioutil.ReadFile(path)

	data := Entries{}

	// JSON 读取
	_ = json.Unmarshal([]byte(file), &data)
	fmt.Println(len(data.Assets))

	fmt.Println(data.Assets[0].AccessibilityLabel)
	fmt.Println(data.Assets[0].Url1080H264)
	fmt.Println(data.Assets[0].Url1080HDR)
	fmt.Println(data.Assets[0].Url1080SDR)
	fmt.Println(data.Assets[0].Url4KHDR)
	fmt.Println(data.Assets[0].Url4KSDR)

	// 测试打印所有下载链接
	// for i := 0; i < len(data.Assets); i++ {
	// 	fmt.Println(i, " -> ", data.Assets[i].AccessibilityLabel)
	// 	fmt.Println(i, " -> ", data.Assets[i].Url1080H264)
	// 	fmt.Println(i, " -> ", data.Assets[i].Url1080HDR)
	// 	fmt.Println(i, " -> ", data.Assets[i].Url1080SDR)
	// 	fmt.Println(i, " -> ", data.Assets[i].Url4KHDR)
	// 	fmt.Println(i, " -> ", data.Assets[i].Url4KSDR)
	// }
}
