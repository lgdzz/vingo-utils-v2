package qiniu

import (
	"context"
	"github.com/lgdzz/vingo-utils-v2/db/redis"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/objects"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"time"
)

type Config struct {
	AccessKey    string
	SecretKey    string
	Bucket       string
	Host         string
	RedisApi     *redis.RedisApi
	SignCacheKey string
}

type ClientApi struct {
	Config        Config
	mac           *credentials.Credentials
	bucketManager *objects.Bucket
}

// 在主进程中只需要执行一次
func InitClient(config Config) (api ClientApi) {
	api.Config = config
	return api
}

// 上传签名
func (s *ClientApi) Sign() string {

	// 如果配置了redis，则从缓存中读取upToken
	if s.Config.RedisApi != nil {
		var upToken string
		if s.Config.RedisApi.Get(s.Config.SignCacheKey, &upToken) {
			return upToken
		}
	}

	var expiry = 1 * time.Hour

	mac := s.NewMac()
	putPolicy, err := uptoken.NewPutPolicy(s.Config.Bucket, time.Now().Add(expiry))
	if err != nil {
		panic(err)
	}
	upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		panic(err)
	}

	// 如果配置了redis，则将新的upToken保存到缓存中
	if s.Config.RedisApi != nil {
		s.Config.RedisApi.Set(s.Config.SignCacheKey, upToken, expiry-(5*time.Minute))
	}

	return upToken
}

func (s *ClientApi) Delete(objectName string) error {
	err := s.BucketManager().Object(objectName).Delete().Call(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (s *ClientApi) NewMac() *credentials.Credentials {
	if s.mac == nil {
		s.mac = credentials.NewCredentials(s.Config.AccessKey, s.Config.SecretKey)
	}
	return s.mac
}

func (s *ClientApi) BucketManager() *objects.Bucket {
	if s.bucketManager == nil {
		objectsManager := objects.NewObjectsManager(&objects.ObjectsManagerOptions{
			Options: http_client.Options{Credentials: s.NewMac()},
		})
		s.bucketManager = objectsManager.Bucket(s.Config.Bucket)
	}
	return s.bucketManager
}
