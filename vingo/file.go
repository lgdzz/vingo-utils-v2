package vingo

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 创建目录
func Mkdir(dirPath string, args ...os.FileMode) string {

	var perm os.FileMode
	if len(args) > 0 {
		perm = args[0]
	} else {
		perm = 0777
	}

	// 将路径中的反斜杠替换为正斜杠，以支持 Windows 目录
	dirPath = filepath.ToSlash(dirPath)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 目录不存在，创建目录
		if err = os.MkdirAll(dirPath, perm); err != nil {
			panic(fmt.Sprintf("创建目录失败：%v", err.Error()))
		} else {
			return dirPath
		}
	} else if err != nil {
		// 其他错误
		panic(fmt.Sprintf("判断目录是否存在时发生错误：%v", err.Error()))
	} else {
		// 目录存在
		return dirPath
	}
}

// 保存文件
// 将字节数据保存为文件
// dirPath 保存文件所在目录
// fileName 保存文件的名称
// data 保存文件的内容
func SaveFile(dirPath string, fileName string, data []byte) {
	targetFile := filepath.Join(dirPath, fileName)
	if err := os.WriteFile(targetFile, data, 0644); err != nil {
		panic(fmt.Sprintf("保存文件失败：%v", err.Error()))
	}
}

func SaveFileSetMode(dirPath string, fileName string, data []byte, perm os.FileMode) {
	targetFile := filepath.Join(dirPath, fileName)
	if err := os.WriteFile(targetFile, data, perm); err != nil {
		panic(fmt.Sprintf("保存文件失败：%v", err.Error()))
	}
}

// 读取文件
func ReadFile(filename string) []byte {
	// 读取文件内容
	data, err := os.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}
	return data
}

// 读取文件，并返回字符串
func ReadFileString(filename string) string {
	return string(ReadFile(filename))
}

// 判断文件是否存在
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// 判断目录是否有读写权限
func HasDirReadWritePermission(dirPath string) bool {
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		LogError(fmt.Sprintf("目录 %v 不存在或无法访问: %v", dirPath, err.Error()))
		return false
	}

	mode := fileInfo.Mode()
	if !mode.IsDir() {
		LogError(fmt.Sprintf("%v 不是一个目录", dirPath))
		return false
	}

	perm := mode.Perm()
	if perm&(1<<(uint(7))) == 0 || perm&(1<<(uint(6))) == 0 {
		LogError(fmt.Sprintf("目录 %v 没有读和写权限", dirPath))
		return false
	}

	return true
}

// 接收上传的文件，返回文件基础信息
func FileUpload(path string, request *http.Request) *FileInfo {
	var (
		requestFile multipart.File
		header      *multipart.FileHeader
		err         error
	)
	requestFile, header, err = request.FormFile("file")
	if err != nil {
		panic(err.Error())
	}
	defer requestFile.Close()

	// 获取文件大小
	fileSize := header.Size

	// 获取文件名称、类型、后缀
	fileName := header.Filename
	fileType := header.Header.Get("Content-Type")
	fileSuffix := filepath.Ext(fileName)

	// 获取当前日期，用于存储文件
	dateString := time.Now().Format(DateFormat)

	// 指定存储目录，如果不存在则创建
	dirPath := Mkdir(filepath.Join(path, dateString))

	// 创建文件
	filePath := filepath.Join(dirPath, fmt.Sprintf("%v%v", GetUUID(), fileSuffix))
	newFile, err := os.Create(filePath)
	if err != nil {
		panic(err.Error())
	}
	defer newFile.Close()

	// 将文件内容拷贝到新文件中
	if _, err = io.Copy(newFile, requestFile); err != nil {
		panic(err.Error())
	}

	// 返回结果
	return &FileInfo{
		Name:      fileName,
		Mimetype:  fileType,
		Extension: fileSuffix,
		Size:      fileSize,
		Realpath:  strings.Replace(filePath, "\\", "/", -1),
	}
}

func FileUploadSetName(path string, name string, request *http.Request, args ...os.FileMode) *FileInfo {
	var (
		requestFile multipart.File
		header      *multipart.FileHeader
		err         error
	)
	requestFile, header, err = request.FormFile("file")
	if err != nil {
		panic(err.Error())
	}
	defer requestFile.Close()

	// 获取文件大小
	fileSize := header.Size

	// 获取文件名称、类型、后缀
	fileName := header.Filename
	fileType := header.Header.Get("Content-Type")
	fileSuffix := filepath.Ext(fileName)

	// 创建文件
	filePath := filepath.Join(path, fmt.Sprintf("%v%v", name, fileSuffix))
	newFile, err := os.Create(filePath)
	if err != nil {
		panic(err.Error())
	}
	defer newFile.Close()

	// 将文件内容拷贝到新文件中
	if _, err = io.Copy(newFile, requestFile); err != nil {
		panic(err.Error())
	}

	var perm os.FileMode
	if len(args) > 0 {
		perm = args[0]
	} else {
		perm = 0644
	}

	// 设置目标文件的权限
	err = os.Chmod(filePath, perm)
	if err != nil {
		panic(err.Error())
	}

	contentType := header.Header.Get("Content-Type")
	err = os.Setenv("Content-Type", contentType)
	if err != nil {
		panic(err.Error())
	}

	// 返回结果
	return &FileInfo{
		Name:      fileName,
		Mimetype:  fileType,
		Extension: fileSuffix,
		Size:      fileSize,
		Realpath:  strings.Replace(filePath, "\\", "/", -1),
	}
}

// 复制文件
// src 文件位置
// dstDir 要复制到的位置
func FileCopy(src, dstDir string) string {
	// Open the source file for reading.
	srcFile, err := os.Open(src)
	if err != nil {
		panic(err.Error())
	}
	defer srcFile.Close()

	// 获取当前日期，用于存储文件
	dateString := time.Now().Format(DateFormat)

	// 指定存储目录，如果不存在则创建
	dstDir = filepath.Join(dstDir, dateString)
	if _, err = os.Stat(dstDir); os.IsNotExist(err) {
		if err = os.MkdirAll(dstDir, 0755); err != nil {
			panic(err.Error())
		}
	}

	// Generate a unique file name with the same extension as the source file.
	srcExt := filepath.Ext(src)
	dstFileName := uuid.New().String() + srcExt

	// Create the destination file with the generated file name.
	dstFile, err := os.Create(filepath.Join(dstDir, dstFileName))
	if err != nil {
		panic(err.Error())
	}
	defer dstFile.Close()

	// Copy the contents of the source file to the destination file.
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		panic(err.Error())
	}

	// Return the path to the copied file.
	return filepath.Join(dstDir, dstFileName)
}

func FileCopySetMode(src, dst string, perm os.FileMode) {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	// 确保目标目录存在
	dstDir := filepath.Dir(dst)
	err = os.MkdirAll(dstDir, perm)
	if err != nil {
		panic(err)
	}

	// 创建目标文件并设置权限
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()

	// 使用 io.Copy 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		panic(err)
	}

	// 确保数据写入磁盘
	err = dstFile.Sync()
	if err != nil {
		panic(err)
	}
}

// 删除文件
func FileDelete(path string, showErr bool) {
	// 文件存在并且路径是安全的
	if FileExists(path) && CheckFilePath(path) {
		// 删除文件
		if err := os.Remove(path); err != nil {
			if showErr {
				panic(fmt.Sprintf("删除文件失败：%v", err.Error()))
			}
		}
	}
}

// 检查路径安全（安全路径必须非/开头，路径中不包含..）
// true-安全|false-不安全
func CheckFilePath(path string) bool {
	// 判断路径中是否包含".."
	if strings.Contains(path, "..") {
		return false
	} else if strings.HasPrefix(path, "/") {
		return false
	}
	return true
}

// 修改文件路径扩展名，如：test.docx 修改为 test.pdf
// filePath 文件路径
// newExt 新的扩展名，不带点，如：pdf
func ModifyPathExtName(filePath string, newExt string) string {
	fileName := filepath.Base(filePath)
	fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	return filepath.Join(filepath.Dir(filePath), fileNameWithoutExt+"."+newExt)
}

// 从文件路径中获取文件名（带文件后缀）
func GetFileNameByPathExt(filePath string) string {
	return filepath.Base(filePath)
}

// 从文件路径中获取文件名（不带文件后缀）
func GetFileNameByPath(filePath string) string {
	fileName := filepath.Base(filePath)
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

// 替换文件路径扩展
// newExt为新的扩展名，包含"."，如：".jpeg"
func ReplaceFilePathExt(path string, newExt string) string {
	return strings.Replace(path, filepath.Ext(path), newExt, 1)
}
