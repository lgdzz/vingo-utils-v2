package aliyun

import (
	aliyun "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/lgdzz/vingo-utils-v2/oss"
)

type Config struct {
	Endpoint        string // 必填
	Bucket          string // 必填
	AccessKeyId     string // 必填
	SecretAccessKey string // 必填
}

type ClientApi struct {
	Config Config
	Client *aliyun.Client
}

// 在主进程中只需要执行一次
func InitClient(config Config) (api ClientApi) {
	api.Config = config

	client, err := aliyun.New(api.Config.Endpoint, api.Config.AccessKeyId, api.Config.SecretAccessKey)
	if err != nil {
		panic(err.Error())
	}

	api.Client = client

	return api
}

func (s *ClientApi) Upload(object oss.Object, localFilePath string) *oss.UploadRes {
	bucket, err := s.Client.Bucket(s.Config.Bucket)
	if err != nil {
		panic(err.Error())
	}

	err = bucket.PutObjectFromFile(object.Name, localFilePath)
	if err != nil {
		panic(err.Error())
	}
	return &oss.UploadRes{Key: object.Name}
}

func (s *ClientApi) Delete(objectName string) error {
	bucket, err := s.Client.Bucket(s.Config.Bucket)
	if err != nil {
		return err
	}

	err = bucket.DeleteObject(objectName)
	if err != nil {
		return err
	}
	return nil
}
