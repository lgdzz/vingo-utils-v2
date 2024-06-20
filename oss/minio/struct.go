// *****************************************************************************
// 作者: lgdz
// 创建时间: 2024/6/20
// 描述：
// *****************************************************************************

package minio

import "github.com/minio/minio-go/v7"

// 连接池对象
type MinIOApi struct {
	Config Config
	Client *minio.Client
}

// 配置
type Config struct {
	Endpoint        string `yaml:"endpoint" json:"endpoint"`
	AccessKeyId     string `yaml:"accessKeyId" json:"accessKeyId"`
	SecretAccessKey string `yaml:"secretAccessKey" json:"secretAccessKey"`
	Bucket          string `yaml:"bucket" json:"bucket"`
	Location        string `yaml:"location" json:"location"`
	UseSSL          bool   `yaml:"useSSL" json:"useSSL"`
	Domain          string `yaml:"domain" json:"domain"`
}

// 存储对象
type Object struct {
	Name        string
	ContentType string
}

// put签名
type PutSign struct {
	Key string `json:"key"`
	Url string `json:"url"`
}

// post签名
type PostSign struct {
	Policy map[string]string `json:"policy"`
	Url    string            `json:"url"`
}
