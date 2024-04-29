package minio

import (
	"context"
	"github.com/lgdzz/vingo-utils-v2/oss"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type Config struct {
	Endpoint        string // 必填
	Bucket          string // 必填
	AccessKeyId     string // 必填
	SecretAccessKey string // 必填
	Location        string
	UseSSL          bool
}

type ClientApi struct {
	Config Config
	Client *minio.Client
}

// 在主进程中只需要执行一次
func InitClient(config Config) (api ClientApi) {
	api.Config = config

	// Initialize minio client object.
	Client, err := minio.New(api.Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(api.Config.AccessKeyId, api.Config.SecretAccessKey, ""),
		Secure: api.Config.UseSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	api.Client = Client

	return api
}

// 上传文件
func (s *ClientApi) Upload(object oss.Object, localFilePath string) *oss.UploadRes {
	ctx := context.Background()
	err := s.Client.MakeBucket(ctx, s.Config.Bucket, minio.MakeBucketOptions{Region: s.Config.Location})
	if err != nil {
		exists, errBucketExists := s.Client.BucketExists(ctx, s.Config.Bucket)
		if errBucketExists != nil || !exists {
			panic(err.Error())
		}
	}
	info, err := s.Client.FPutObject(ctx, s.Config.Bucket, object.Name, localFilePath, minio.PutObjectOptions{ContentType: object.ContentType})
	if err != nil {
		panic(err.Error())
	}
	return &oss.UploadRes{
		Key:  info.Key,
		Info: info,
	}
}

// 删除文件
func (s *ClientApi) Delete(objectName string) error {
	ctx := context.Background()
	err := s.Client.RemoveObject(ctx, s.Config.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
