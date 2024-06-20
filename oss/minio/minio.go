package minio

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"strings"
	"time"
)

// 创建新的客户端（全局只需要执行一次）
func NewClient(config Config) *MinIOApi {
	var api = MinIOApi{
		Config: config,
	}

	// Initialize minio client object.
	client, err := minio.New(api.Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(api.Config.AccessKeyId, api.Config.SecretAccessKey, ""),
		Secure: api.Config.UseSSL,
	})
	if err != nil {
		msg := fmt.Sprintf("MinIO初始化异常：%v", err.Error())
		vingo.Log(msg)
		fmt.Println(msg)
	}

	api.Client = client
	return &api
}

// 完整文件url
func (s *MinIOApi) GetObjectUrl(objectName string) string {
	return s.Config.Domain + s.Config.Bucket + "/" + objectName
}

// 上传本地文件
func (s *MinIOApi) UploadFileOfLocal(object Object, localFilePath string) minio.UploadInfo {
	info, err := s.Client.FPutObject(context.Background(), s.Config.Bucket, object.Name, localFilePath, minio.PutObjectOptions{ContentType: object.ContentType})
	if err != nil {
		panic(err.Error())
	}
	return info
}

// 上传base64文件
func (s *MinIOApi) UploadFileOfBase64(object Object, fileBase64 string) minio.UploadInfo {
	// 解码Base64图像数据
	var fileBase64Array = strings.Split(fileBase64, ",")
	if len(fileBase64Array) > 1 {
		fileBase64 = fileBase64Array[1]
	}
	imageData, err := base64.StdEncoding.DecodeString(fileBase64)
	if err != nil {
		panic(err.Error())
	}
	info, err := s.Client.PutObject(context.Background(), s.Config.Bucket, object.Name, strings.NewReader(string(imageData)), int64(len(imageData)), minio.PutObjectOptions{ContentType: object.ContentType})
	if err != nil {
		panic(err.Error())
	}
	return info
}

// 删除文件
func (s *MinIOApi) DeleteFile(objectName string) error {
	err := s.Client.RemoveObject(context.Background(), s.Config.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// 获取文件授权访问地址(默认24小时有效)
func (s *MinIOApi) GetObjectSignUrl(objectName string, expires *time.Duration) string {
	if expires == nil {
		*expires = time.Hour * 24
	}
	url, err := s.Client.PresignedGetObject(context.Background(), s.Config.Bucket, objectName, *expires, nil)
	if err != nil {
		return err.Error()
	}
	return url.String()
}

// 获取对象put上传签名
func (s *MinIOApi) GetObjectPutSign(objectName string) PutSign {
	url, err := s.Client.PresignedPutObject(context.Background(), s.Config.Bucket, objectName, time.Minute*10)
	if err != nil {
		panic(err.Error())
	}
	return PutSign{
		Key: objectName,
		Url: url.String(),
	}
}

// 获取对象post上传签名
func (s *MinIOApi) GetObjectPostSign(objectName string) PostSign {
	policy := minio.NewPostPolicy()
	_ = policy.SetExpires(time.Now().Add(time.Minute * 10))
	_ = policy.SetKey(objectName)
	_ = policy.SetBucket(s.Config.Bucket)
	url, formData, err := s.Client.PresignedPostPolicy(context.Background(), policy)
	if err != nil {
		panic(err.Error())
	}
	return PostSign{
		Policy: formData,
		Url:    url.String(),
	}
}
