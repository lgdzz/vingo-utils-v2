package vingo

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lgdzz/vingo-utils-v2/db/redis"
	"time"
)

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

// 生成token
// JwtIssued(jwt.JwtBody[Business]{}, "123456")
// Business是声明body中business字段类型
func JwtIssued[T any](body JwtBody[T], signingKey string) string {
	if body.Day == 0 {
		body.Day = 90
	}
	day := 3600 * 24 * int64(body.Day)
	exp := time.Now().Unix() + day
	if body.CheckTK {
		body.Ticket = &JwtTicket{Key: MD5(fmt.Sprintf("%v%v", signingKey, body.ID)), TK: RandomString(50)}
		redis.Set(body.Ticket.Key, body.Ticket.TK, time.Second*time.Duration(day))
	}
	signedString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": body.ID, "checkTk": body.CheckTK, "ticket": body.Ticket, "business": body.Business, "exp": exp}).SignedString([]byte(signingKey))
	if err != nil {
		panic(err)
	}
	return signedString
}

// 验证token
// JwtCheck[Business](token, "123456")
// Business是声明body中business字段类型
func JwtCheck[T any](token string, signingKey string) JwtBody[T] {
	claims, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		panic(&AuthException{Message: err.Error()})
	}
	var body JwtBody[T]
	CustomOutput(claims.Claims, &body)
	if body.CheckTK {
		tkPointer := redis.Get[string](body.Ticket.Key)
		if tkPointer == nil || body.Ticket.TK != *tkPointer {
			panic(&AuthException{Message: "登录已失效"})
		}
	}
	return body
}
