package main

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func errPrintln(err error) {
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

// Untar will decompress a tar archive, moving all files and folders
// within the tar file (parameter 1) to an output directory (parameter 2).
func untar(src string) ([]string, error) {

	var filenames []string

	// 将 tar 包打开
	fr, err := os.Open(src)
	errPrintln(err)
	defer fr.Close()

	// 通过 fr 创建一个 tar.*Reader 结构，然后将 tr 遍历，并将数据保存到磁盘中
	tr := tar.NewReader(fr)

	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		// 处理 err ！= nil 的情况
		errPrintln(err)
		// 获取文件信息
		fi := hdr.FileInfo()

		// 创建一个空文件，用来写入解包后的数据
		fw, err := os.Create(fi.Name())
		errPrintln(err)

		// 将 tr 写入到 fw
		n, err := io.Copy(fw, tr)
		errPrintln(err)
		log.Printf("解包： %s 到 %s ，共处理了 %d 个字符的数据。", srcFile, fi.Name(), n)

		// 设置文件权限，这样可以保证和原始文件权限相同，如果不设置，会根据当前系统的 umask 来设置。
		os.Chmod(fi.Name(), fi.Mode().Perm())

		// 注意，因为是在循环中，所以就没有使用 defer 关闭文件
		// 如果想使用 defer 的话，可以将文件写入的步骤单独封装在一个函数中即可
		fw.Close()
	}
	return filenames, nil
}

func downloadFromURL(url string) string {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	fmt.Println("file exixxxxst")
	// 判断文件是否已经存在
	if fileExists(fileName) {
		// 文件已经存在，不进行任何操作
		fmt.Println("file exist")
		return fileName
	} else {
		// 文件不存在的情况下开始下载
		fmt.Println("Downloading", url, "to", fileName)
	}

	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return fileName
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return fileName
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return fileName
	}

	fmt.Println(n, "bytes downloaded.")
	return fileName
}

func main() {
	// 下载实体文件
	var url = "https://sylvan.apple.com/Aerials/resources.tar"
	srcfile := downloadFromURL(url)
	fmt.Println(srcfile)
	// untar(srcfile)
}
