package main

import (
	"fmt"

	"./pack"
)

func main() {
	// entriesresources.Download()
	entrie := entriesresources.HandleRetries("./temp/entries.json")

	fmt.Println(len(entrie.Assets))

	for i := 0; i < len(entrie.Assets); i++ {
		fmt.Println(i, " -> ", entrie.Assets[i].AccessibilityLabel)
		fmt.Println(i, " -> ", entrie.Assets[i].Url1080H264)
		fmt.Println(i, " -> ", entrie.Assets[i].Url1080HDR)
		fmt.Println(i, " -> ", entrie.Assets[i].Url1080SDR)
		fmt.Println(i, " -> ", entrie.Assets[i].Url4KHDR)
		fmt.Println(i, " -> ", entrie.Assets[i].Url4KSDR)
	}

	entriesresources.DownloadFromURL(entrie.Assets[1].Url1080H264)
	// fmt.Printf("Hello, worldxxx\n")
}
