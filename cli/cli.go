package cli

import (
	"flag"
	"fmt"
	"github.com/lgdzz/vingo-utils-v2/db/mysql"
	"github.com/lgdzz/vingo-utils-v2/db/pgsql"
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Options struct {
	Enable     bool
	DbApi      *mysql.DbApi
	PgSqlDbApi *pgsql.DbApi
	Register   func()
}

func InitCli(options Options) {
	if !options.Enable {
		return
	}
	model := flag.String("model", "", "生成数据库模型，支持多个表生成，格式：table1,table2")
	flag.StringVar(model, "m", "", "生成数据库模型，支持多个表生成，格式：table1,table2")

	dbbook := flag.Bool("dbbook", false, "生成数据库字典")
	flag.BoolVar(dbbook, "d", false, "生成数据库字典")

	secret := flag.Bool("secret", false, "生成字符串加解密secret")
	flag.BoolVar(secret, "s", false, "生成字符串加解密secret")

	buildDev := flag.String("build-dev", "", "打包开发版，参数：l=linux;w=windows;m=mac")
	buildProd := flag.String("build-prod", "", "打包正式版，参数：l=linux;w=windows;m=mac")

	updateVingo := flag.String("v2", "", "更新vingo-v2版本")

	if options.Register != nil {
		options.Register()
	}

	help := flag.Bool("h", false, "Show help")

	// 解析命令行参数
	flag.Parse()

	if *help {
		// 如果使用 -h 或 --help 标志，则显示帮助信息
		flag.Usage()
		os.Exit(0)
	}

	// 创建数据库字典
	if *dbbook {
		if options.DbApi != nil {
			err := options.DbApi.BuildBook()
			if err != nil {
				fmt.Println(err)
			}
		} else if options.PgSqlDbApi != nil {
			err := options.PgSqlDbApi.BuildBook()
			if err != nil {
				fmt.Println(err)
			}
		}
		os.Exit(0)
	}

	// 创建数据表模型文件
	if *model != "" {
		if options.DbApi != nil {
			_, _ = options.DbApi.CreateDbModel(strings.Split(*model, ",")...)
		} else if options.PgSqlDbApi != nil {
			_, _ = options.PgSqlDbApi.CreateDbModel(strings.Split(*model, ",")...)

		}
		os.Exit(0)
	}

	if *buildDev != "" {
		BuildProject(*buildDev, "dev")
	}
	if *buildProd != "" {
		BuildProject(*buildProd, "prod")
	}

	if *updateVingo != "" {
		cmd := exec.Command("go", "get", "-u", "github.com/lgdzz/vingo-utils-v2@"+*updateVingo)
		_ = cmd.Run()
		os.Exit(0)
	}

}

func BuildProject(value string, version string) {
	var goos string
	var osName string
	switch value {
	case "l":
		goos = "linux"
		osName = "linux"
	case "w":
		goos = "windows"
		osName = "windows"
	case "m":
		goos = "darwin"
		osName = "mac"
	}
	var err error
	if err = os.Setenv("CGO_ENABLED", "0"); err != nil {
		log.Println("设置CGO_ENABLED错误：", err.Error())
		os.Exit(0)
	}
	if err = os.Setenv("GOOS", goos); err != nil {
		log.Println("设置GOOS错误：", err.Error())
		os.Exit(0)
	}
	if err = os.Setenv("GOARCH", "amd64"); err != nil {
		log.Println("设置GOARCH错误：", err.Error())
		os.Exit(0)
	}

	var moduleName = vingo.GetModuleName()
	var outputName = fmt.Sprintf("%v.%v-%v", moduleName, version, osName)
	if osName == "windows" {
		outputName += ".exe"
	}

	// 执行打包命令
	cmd := exec.Command("go", "build", "-ldflags=-X "+moduleName+"/config.version="+version, "-o", outputName)
	err = cmd.Run()
	if err != nil {
		log.Println("执行打包命令错误：", err.Error())
		os.Exit(0)
	}

	// 获取文件信息
	fileInfo, err := os.Stat(outputName)
	if err != nil {
		log.Println("获取打包文件信息错误：", err.Error())
	}
	fileSize := fileInfo.Size()
	log.Println("文件名称：", outputName)
	log.Println("文件大小：", vingo.FormatBytes(fileSize, 2))
	os.Exit(0)
}
