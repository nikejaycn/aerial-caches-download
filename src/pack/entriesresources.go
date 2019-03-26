package entriesresources

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

func unDirTar(dst, src string) (err error) {
	// 打开文件
	srcFileReader, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcFileReader.Close()

	// 获取压缩文件的目录情况
	tarFileReader := tar.NewReader(srcFileReader)
	if err != nil {
		return
	}

	for {
		hdr, err := tarFileReader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case hdr == nil:
			continue
		}

		// 处理保存目录，将要保存的目录加上 header 中的 Name
		// 这个变量保存的有肯能是目录，也有可能是文件
		dstFileDir := filepath.Join(dst, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir: // 目录
			// 判断目录是否存在，不存在就创建
			if b := existDir(dstFileDir); !b {
				// MkdirAll = (mkdir -p)
				err := os.MkdirAll(dstFileDir, 0775)
				if err != nil {
					return err
				}
			}
		case tar.TypeReg: // 文件
			// 创建一个可以读写的文件，权限就使用 header 中记录的权限
			// 因为操作系统的 FileMode 是 int32 类型的，hdr 中的是 int64，所以转换下
			file, err := os.OpenFile(dstFileDir, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}

			n, err := io.Copy(file, tarFileReader)
			if err != nil {
				return err
			}

			// 将解压结果输出显示
			fmt.Printf("成功解压： %s , 共处理了 %d 个字符\n", dstFileDir, n)

			// 不要忘记关闭打开的文件，因为它是在 for 循环中，不能使用 defer
			// 如果想使用 defer 就放在一个单独的函数中
			file.Close()
		}
	}

	return nil
}

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type WriteCounter struct {
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\rDownloading... %d MB complete", wc.Total/1024/1024)
}

// 判断目录是否存在
func existDir(dirname string) bool {
	fileinfo, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fileinfo.IsDir()
}

// 获取文件的执行目录
func getCurrentPath(src string) (string, error) {
	path, err := filepath.Abs(src)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

// 根据链接下载视频
func DownloadFromURL(url string) string {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

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

	// 实现下载过程中显示进度
	counter := &WriteCounter{}
	_, err = io.Copy(output, io.TeeReader(response.Body, counter))
	if err != nil {
		fmt.Println("Error while TeeReader", fileName, "-", err)
		return fileName
	}

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return fileName
	}

	fmt.Println(n, "bytes downloaded.")
	return fileName
}

func Download() {

	// 下载实体文件
	var url = "https://sylvan.apple.com/Aerials/resources.tar"
	srcfile := DownloadFromURL(url)
	// fmt.Println(filepath.Abs(srcfile))

	// 获取当前目录
	path, err := getCurrentPath(srcfile)
	if err != nil {
		fmt.Println(err)
	}
	dst := path + "temp/"
	if b := existDir(dst); !b {
		os.MkdirAll(dst, 0775)
	}
	result := unDirTar(dst, srcfile)
	fmt.Println(result)
}
