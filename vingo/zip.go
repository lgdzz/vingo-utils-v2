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
	"os"
	"path/filepath"
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
	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("Error downloading file: %v", err))
	}
	defer resp.Body.Close()

	if fileName == "" {
		fileName = url[strings.LastIndex(url, "/")+1:]
	}
	fullPath := buildZipPath(fileName, folders...)

	fileInZip, err := s.ZipWriter.Create(fullPath)
	if err != nil {
		panic(fmt.Sprintf("Error creating file in zip: %v", err))
	}

	_, err = io.Copy(fileInZip, resp.Body)
	if err != nil {
		panic(fmt.Sprintf("Error copying file to zip: %v", err))
	}
}

// 添加本地文件到 zip 中
func (s *ZipObject) AddLocalFileToZip(localPath string, fileName string, folders ...string) {
	file, err := os.Open(localPath)
	if err != nil {
		panic(fmt.Sprintf("Error opening local file: %v", err))
	}
	defer file.Close()

	if fileName == "" {
		fileName = filepath.Base(localPath)
	}
	fullPath := buildZipPath(fileName, folders...)

	fileInfo, err := file.Stat()
	if err != nil {
		panic(fmt.Sprintf("Error getting file info: %v", err))
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		panic(fmt.Sprintf("Error creating zip header: %v", err))
	}
	header.Name = fullPath
	header.Method = zip.Deflate

	writer, err := s.ZipWriter.CreateHeader(header)
	if err != nil {
		panic(fmt.Sprintf("Error creating file in zip: %v", err))
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		panic(fmt.Sprintf("Error writing file to zip: %v", err))
	}
}

// 关闭 zip writer（重要：写入完成后需要关闭）
func (s *ZipObject) Close() {
	err := s.ZipWriter.Close()
	if err != nil {
		panic(fmt.Sprintf("Error closing zip writer: %v", err))
	}
}

// 工具函数：拼接 zip 内路径
func buildZipPath(fileName string, folders ...string) string {
	if len(folders) == 0 {
		return fileName
	}
	return strings.Join(folders, "/") + "/" + fileName
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
