package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lgdzz/vingo-utils-v2/db/redis"
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"time"
)

type JwtApi[T any] struct {
	Secret   string
	RedisApi *redis.RedisApi
}

type JwtTicket struct {
	Key string `json:"key"`
	TK  string `json:"tk"`
}

type JwtBody[T any] struct {
	ID       string     `json:"id"`
	Day      uint       `json:"day"` // 默认有效期90天
	Business T          `json:"business"`
	CheckTK  bool       `json:"checkTk"`
	Ticket   *JwtTicket `json:"ticket"`
}

func NewJwt[T any](secret string, redisApi *redis.RedisApi) *JwtApi[T] {
	return &JwtApi[T]{
		Secret:   secret,
		RedisApi: redisApi,
	}
}

// 生成token
// JwtIssued(jwt.JwtBody[Business]{}, "123456")
// Business是声明body中business字段类型
func (s *JwtApi[T]) Issued(body JwtBody[T]) (string, int64) {
	if body.Day == 0 {
		body.Day = 90
	}
	day := 3600 * 24 * int64(body.Day)
	exp := time.Now().Unix() + day
	if body.CheckTK {
		body.Ticket = &JwtTicket{Key: vingo.MD5(fmt.Sprintf("%v%v", s.Secret, body.ID)), TK: vingo.RandomString(50)}
		s.RedisApi.Set(body.Ticket.Key, body.Ticket.TK, time.Second*time.Duration(day))
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": body.ID, "checkTk": body.CheckTK, "ticket": body.Ticket, "business": body.Business, "exp": exp}).SignedString([]byte(s.Secret))
	if err != nil {
		panic(err)
	}
	return token, exp
}

// 验证token
// JwtCheck[Business](token, "123456")
// Business是声明body中business字段类型
func (s *JwtApi[T]) Check(token string) JwtBody[T] {
	claims, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Secret), nil
	})
	if err != nil {
		panic(&vingo.AuthException{Message: err.Error()})
	}
	var body JwtBody[T]
	vingo.CustomOutput(claims.Claims, &body)
	if body.CheckTK {
		var tk string
		if !s.RedisApi.Get(body.Ticket.Key, &tk) || tk != body.Ticket.TK {
			panic(&vingo.AuthException{Message: "登录已失效"})
		}
	}
	return body
}
