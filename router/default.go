package router

import (
	"fmt"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"github.com/lgdzz/vingo-utils-v2/db/mysql"
	"github.com/lgdzz/vingo-utils-v2/db/redis"
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"net/http"
	"time"
)

type Hook struct {
	Option         HookOption
	RegisterRouter func(r *gin.Engine)
	BaseMiddle     func(c *gin.Context)
	LoadWeb        []WebItem // 通过此配置将前端项目打包到项目中
}

type HookOption struct {
	Name      string
	Port      uint
	Copyright string
	Debug     bool
	Database  mysql.Config
	Redis     redis.Config
}

type WebItem struct {
	Route string
	FS    *assetfs.AssetFS // go-bindata-assetfs -pkg {admin} -o router/{admin}/bindata.go dist/...
}

// 初始化路由
func InitRouter(hook *Hook) {
	var option = hook.Option
	if option.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	vingo.GinDebug = option.Debug

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	_ = r.SetTrustedProxies(nil)

	// 加载web前端
	for _, item := range hook.LoadWeb {
		currentItem := item
		r.GET(currentItem.Route+"/*filepath", func(c *gin.Context) {
			http.FileServer(currentItem.FS).ServeHTTP(c.Writer, c.Request)
		})
	}

	// 屏蔽搜索引擎爬虫
	vingo.ShieldRobots(r)

	// 404捕捉
	r.NoRoute(func(c *gin.Context) {
		context := vingo.Context{Context: c}
		context.Response(&vingo.ResponseData{Message: "404:Not Found", Error: 1})
	})

	// 注册异常处理、基础中间件
	r.Use(vingo.ExceptionHandler, BaseMiddle(hook))

	// 注册路由
	hook.RegisterRouter(r)

	fmt.Println("+------------------------------------------------------------+")
	fmt.Println(fmt.Sprintf("+ 项目名称：%v", option.Name))
	fmt.Println(fmt.Sprintf("+ 服务端口：%d", option.Port))
	fmt.Println(fmt.Sprintf("+ 调试模式：%v", option.Debug))
	if option.Database.Host != "" {
		fmt.Println(fmt.Sprintf("+ Mysql：%v:%v db:%v", option.Database.Host, option.Database.Port, option.Database.Dbname))
	}
	if option.Redis.Host != "" {
		fmt.Println(fmt.Sprintf("+ Redis：%v:%v db:%v", option.Redis.Host, option.Redis.Port, option.Redis.Select))
	}
	vingo.ApiAddress(option.Port)
	fmt.Println(fmt.Sprintf("+ 启动时间：%v", time.Now().Format(vingo.DatetimeFormatChinese)))
	fmt.Println(fmt.Sprintf("+ 技术支持：%v", option.Copyright))
	fmt.Println("+------------------------------------------------------------+")

	// 开启服务
	_ = r.Run(fmt.Sprintf(":%d", option.Port))
}
