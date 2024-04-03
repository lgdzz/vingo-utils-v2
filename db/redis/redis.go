package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Option struct {
	Host         string `yaml:"host" json:"host"`
	Port         string `yaml:"port" json:"port"`
	Select       int    `yaml:"select" json:"select"`
	Password     string `yaml:"password" json:"password"`
	PoolSize     int    `yaml:"poolSize" json:"poolSize"`
	MinIdleConns int    `yaml:"minIdleConns" json:"minIdleConns"`
	Prefix       string `yaml:"prefix" json:"prefix"`
}

var Client *redis.Client
var KeyPrefix string

// 默认配置
func DefaultConfig(option *Option) {
	if option.Host == "" {
		option.Host = "127.0.0.1"
	}

	if option.Port == "" {
		option.Port = "6379"
	}

	if option.PoolSize == 0 {
		option.PoolSize = 20
	}

	if option.MinIdleConns == 0 {
		option.PoolSize = 10
	}

	KeyPrefix = option.Prefix
}

// redis初始化
func InitClient(option *Option) {

	if Client != nil {
		return
	}

	DefaultConfig(option)

	Client = redis.NewClient(&redis.Options{
		//连接信息
		Network:  "tcp",                                          //网络类型，tcp or unix，默认tcp
		Addr:     fmt.Sprintf("%v:%v", option.Host, option.Port), //主机名+冒号+端口，默认localhost:6379
		Password: option.Password,                                //密码
		DB:       option.Select,                                  // redis数据库index

		//连接池容量及闲置连接数量
		PoolSize:     option.PoolSize,     // 连接池最大socket连接数，应该设置为服务器CPU核心数的两倍
		MinIdleConns: option.MinIdleConns, // 在启动阶段创建指定数量的Idle连接，一般来说，可以将其设置为PoolSize的一半
	})
	// 测试连接是否正常
	_, err := Client.Ping().Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis连接异常：%v", err.Error()))
	}
}

func RedisResult(cmd *redis.StringCmd) string {
	result, err := cmd.Result()
	if err == redis.Nil {
		return ""
	} else if err != nil {
		panic(err.Error())
	} else {
		return result
	}
}

func RedisSaveResult(err error) {
	if err != nil {
		panic(err.Error())
	}
}
