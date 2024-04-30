package redis

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/lgdzz/vingo-utils-v2/config"
	"time"
)

type Config struct {
	config.Config
	Host         string `yaml:"host" json:"host"`
	Port         string `yaml:"port" json:"port"`
	Select       int    `yaml:"select" json:"select"`
	Password     string `yaml:"password" json:"password"`
	PoolSize     int    `yaml:"poolSize" json:"poolSize"`
	MinIdleConns int    `yaml:"minIdleConns" json:"minIdleConns"`
	Prefix       string `yaml:"prefix" json:"prefix"`
}

type RedisApi struct {
	Client *redis.Client
	Config Config
}

func (s *RedisApi) BuildKey(key string) string {
	return s.Config.Prefix + key
}

func (s *RedisApi) Get(key string, value any) (exist bool) {
	text, err := s.Client.Get(s.BuildKey(key)).Result()
	if err == redis.Nil {
		return
	} else if err != nil {
		panic(err)
	} else {
		err = json.Unmarshal([]byte(text), value)
		if err != nil {
			panic(err.Error())
		}
		exist = true
		return
	}
}

func (s *RedisApi) Set(key string, value any, expiration time.Duration) string {
	v, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	result, err := s.Client.Set(s.BuildKey(key), v, expiration).Result()
	if err != nil {
		panic(err)
	}
	return result
}

func (s *RedisApi) HSet(key string, field string, value any) bool {
	v, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	result, err := s.Client.HSet(s.BuildKey(key), field, v).Result()
	if err != nil {
		panic(err)
	}
	return result
}

func (s *RedisApi) HGet(key string, field string, value any) (exist bool) {
	text, err := s.Client.HGet(s.BuildKey(key), field).Result()
	if err == redis.Nil {
		return
	} else if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(text), value)
	if err != nil {
		panic(err.Error())
	}
	exist = true
	return
}

func (s *RedisApi) Del(key ...string) int64 {
	var keys = make([]string, 0)
	for _, item := range key {
		keys = append(keys, s.BuildKey(item))
	}
	result, err := s.Client.Del(keys...).Result()
	if err != nil {
		panic(err)
	}
	return result
}

// 新建一个redis连接池
func NewRedis(config Config) *RedisApi {
	config.StringValue(&config.Host, "127.0.0.1")
	config.StringValue(&config.Port, "6379")
	config.StringValue(&config.Prefix, "")
	config.IntValue(&config.PoolSize, 20)
	config.IntValue(&config.MinIdleConns, 10)

	var redisApi = RedisApi{
		Config: config,
	}

	redisApi.Client = redis.NewClient(&redis.Options{
		//连接信息
		Network:  "tcp",                                          //网络类型，tcp or unix，默认tcp
		Addr:     fmt.Sprintf("%v:%v", config.Host, config.Port), //主机名+冒号+端口，默认localhost:6379
		Password: config.Password,                                //密码
		DB:       config.Select,                                  // redis数据库index

		//连接池容量及闲置连接数量
		PoolSize:     config.PoolSize,     // 连接池最大socket连接数，应该设置为服务器CPU核心数的两倍
		MinIdleConns: config.MinIdleConns, // 在启动阶段创建指定数量的Idle连接，一般来说，可以将其设置为PoolSize的一半
	})
	// 测试连接是否正常
	_, err := redisApi.Client.Ping().Result()
	if err != nil {
		fmt.Println(fmt.Sprintf("Redis连接异常：%v", err.Error()))
	}

	return &redisApi
}
