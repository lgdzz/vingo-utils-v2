package vingo

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

var Valid *validator.Validate

func init() {
	if Valid == nil {
		Valid = validator.New()
	}
}

type Context struct {
	*gin.Context
}

// 动态添加查询条件
func (c *Context) AddQuery(kv ...KeyValue) {
	var query = c.Request.URL.Query()
	for _, item := range kv {
		query.Add(item.Key, item.Value)
	}
	c.Request.URL.RawQuery = query.Encode()
}

// url解码
func (c *Context) UrlDecode() (dStr string) {
	var (
		err  error
		eStr = c.Request.RequestURI
	)
	dStr, err = url.QueryUnescape(eStr)
	if err != nil {
		dStr = eStr
	}
	return
}

// 获取客户端真实IP
func (c *Context) GetRealClientIP() string {
	var ip string
	if ip = c.Request.Header.Get("X-Forwarded-For"); ip == "" {
		ip = c.Request.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = c.Request.RemoteAddr
	} else {
		ips := strings.Split(ip, ", ")
		ip = ips[len(ips)-1]
	}
	return ip
}

// 请求成功
func (c *Context) Response(d *ResponseData) {
	c.Set("clientIp", c.GetRealClientIP())
	if d.Message == "" {
		d.Message = "Success"
	}
	if d.Status == 0 {
		d.Status = 200
	}
	uuid := c.GetString("requestUUID")
	c.JSON(d.Status, gin.H{
		"uuid":      uuid,
		"error":     d.Error,
		"message":   d.Message,
		"data":      d.Data,
		"timestamp": time.Now().Unix(),
	})

	if !d.NoLog {
		// 记录请求日志
		go func(context *Context, uuid string, d *ResponseData) {
			startTime := context.GetTime("requestStart")
			endTime := time.Now()
			latency := endTime.Sub(startTime)
			millisecond := float64(latency.Nanoseconds()) / float64(time.Millisecond)
			duration := fmt.Sprintf("%.3fms", millisecond)
			if millisecond > 300 {
				duration += ":慢接口"
			}

			var err string
			if d.Error == 1 {
				err = d.Message
			}

			if context.Request.Method == "GET" {
				LogRequest(duration, fmt.Sprintf("{\"uuid\":\"%v\",\"method\":\"%v\",\"url\":\"%v\",\"err\":\"%v\",\"errType\":\"%v\",\"userAgent\":\"%v\",\"clientIP\":\"%v\",\"user\":\"%v\"}", uuid, context.Request.Method, context.UrlDecode(), err, d.ErrorType, c.GetHeader("User-Agent"), c.GetString("clientIp"), c.GetString("user")))
			} else {
				body := context.GetString("requestBody")
				if body == "" {
					body = "\"\""
				}
				LogRequest(duration, fmt.Sprintf("{\"uuid\":\"%v\",\"method\":\"%v\",\"url\":\"%v\",\"body\":%v,\"err\":\"%v\",\"errType\":\"%v\",\"userAgent\":\"%v\",\"clientIP\":\"%v\",\"user\":\"%v\"}", uuid, context.Request.Method, context.Request.RequestURI, body, err, d.ErrorType, c.GetHeader("User-Agent"), c.GetString("clientIp"), c.GetString("user")))
			}
		}(c, uuid, d)
	}
}

// 请求成功，带data数据
func (c *Context) ResponseBody(data any) {
	c.Response(&ResponseData{Data: data})
}

// 请求成功，默认
func (c *Context) ResponseSuccess(data ...any) {
	if data == nil {
		c.Response(&ResponseData{})
	} else {
		c.Response(&ResponseData{Data: data})
	}
}

// 注册get路由
func RoutesGet(g *gin.RouterGroup, path string, handler func(*Context)) {
	g.GET(path, func(c *gin.Context) {
		handler(&Context{Context: c})
	})
}

// 注册post路由
func RoutesPost(g *gin.RouterGroup, path string, handler func(*Context)) {
	g.POST(path, func(c *gin.Context) {
		handler(&Context{Context: c})
	})
}

// 注册put路由
func RoutesPut(g *gin.RouterGroup, path string, handler func(*Context)) {
	g.PUT(path, func(c *gin.Context) {
		handler(&Context{Context: c})
	})
}

// 注册patch路由
func RoutesPatch(g *gin.RouterGroup, path string, handler func(*Context)) {
	g.PATCH(path, func(c *gin.Context) {
		handler(&Context{Context: c})
	})
}

// 注册delete路由
func RoutesDelete(g *gin.RouterGroup, path string, handler func(*Context)) {
	g.DELETE(path, func(c *gin.Context) {
		handler(&Context{Context: c})
	})
}

// 注册websocket路由
func RoutesWebsocket(g *gin.RouterGroup, path string, handler func(*Context)) {
	RoutesGet(g, path, handler)
}

// 屏蔽搜索引擎爬虫
func ShieldRobots(r *gin.Engine) {
	r.GET("/robots.txt", func(c *gin.Context) {
		c.String(200, `User-agent: *
Disallow: /`)
	})
}

// 设置允许跨域访问的域名或IP地址
func AllowCrossDomain(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
}

// 输出接口地址
func ApiAddress(port uint) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, item := range addr {
		if ipNet, ok := item.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				fmt.Println(fmt.Sprintf("+ 接口地址：http://%v:%d", ipNet.IP.String(), port))
			}
		}
	}
}

func (c *Context) GetUserId() uint {
	return c.GetUint("userId")
}

func (c *Context) GetAccId() uint {
	return c.GetUint("accId")
}

func (c *Context) GetOrgId() uint {
	return c.GetUint("orgId")
}

func (c *Context) GetRoleId() UintIds {
	id, exists := c.Get("roleId")
	if !exists {
		id = UintIds{}
	}
	return id.(UintIds)
}

func (c *Context) GetRealName() string {
	return c.GetString("realName")
}

func (c *Context) GetOrgName() string {
	return c.GetString("orgName")
}

// 获取请求body
// GetRequestBody[结构体类型](c)
func GetRequestBody[T any](c *Context) T {
	var body T
	if err := c.ShouldBindJSON(&body); err != nil {
		panic(err.Error())
	}

	if err := Valid.Struct(body); err != nil {
		// handle validation error
		panic(err)
	}

	if data, err := json.Marshal(body); err != nil {
		panic(err.Error())
	} else {
		c.Set("requestBody", string(data))
	}
	return body
}

// 获取请求query
// GetRequestQuery[结构体类型](c)
func GetRequestQuery[T any](c *Context) T {
	var query T
	if err := c.ShouldBindQuery(&query); err != nil {
		panic(err.Error())
	}
	return query
}

type ResponseData struct {
	Status    int    // 状态
	Error     int    // 0-无错误|1-有错误
	ErrorType string // 错误类型
	Message   string // 消息
	Data      any    // 返回数据内容
	NoLog     bool   // true时不记录日志
}
