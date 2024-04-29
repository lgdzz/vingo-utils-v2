// *****************************************************************************
// 作者: lgdz
// 创建时间: 2024/4/29
// 描述：zip压缩
// *****************************************************************************

package vingo

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ZipObject struct {
	ZipBuffer *bytes.Buffer
	ZipWriter *zip.Writer
}

// 创建新的压缩包
func NewZip() ZipObject {
	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)
	return ZipObject{
		ZipBuffer: zipBuffer,
		ZipWriter: zipWriter,
	}
}

// 将url文件添加到压缩包
func (s *ZipObject) AddUrlFileToZip(url string, fileName string, folders ...string) {
	// 发起 HTTP 请求获取远程文件内容
	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("Error downloading file: %v", err.Error()))
	}
	defer resp.Body.Close()
	// 创建 zip 文件中的条目
	var folder string
	if len(folders) > 0 {
		folder = strings.Join(folders, "/") + "/"
	}
	if fileName == "" {
		fileName = url[strings.LastIndex(url, "/")+1:] // 提取文件名
	}
	fileInZip, err := s.ZipWriter.Create(folder + fileName)
	if err != nil {
		panic(fmt.Sprintf("Error creating file in zip: %v", err.Error()))
	}
	// 将远程文件内容复制到 zip 文件中
	_, err = io.Copy(fileInZip, resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Error copying file to zip: %v", err.Error()))
	}
}

func (s *ZipObject) Download(c *Context, filename string) {
	// 关闭 zip.Writer
	err := s.ZipWriter.Close()
	if err != nil {
		panic(fmt.Sprintf("Error closing zip writer: %v", err.Error()))
	}
	// 设置响应头，指定为 zip 文件
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%v.zip", url.QueryEscape(filename)))
	// 将 zipBuffer 中的内容作为文件流返回给客户端
	c.Data(http.StatusOK, "application/zip", s.ZipBuffer.Bytes())
}
