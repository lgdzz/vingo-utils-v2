package pgsql

import (
	"fmt"
	"github.com/lgdzz/vingo-utils/vingo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var Db *gorm.DB

type PgsqlConfig struct {
	Host         string `yaml:"host" json:"host"`
	Port         string `yaml:"port" json:"port"`
	Dbname       string `yaml:"dbname" json:"dbname"`
	Username     string `yaml:"username" json:"username"`
	Password     string `yaml:"password" json:"password"`
	MaxIdleConns int    `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxOpenConns int    `yaml:"maxOpenConns" json:"maxOpenConns"`
}

func InitClient(config *PgsqlConfig) {
	InitPgsqlService(config)
}

func InitPgsqlService(config *PgsqlConfig) {
	if Db != nil {
		return
	}

	if config.MaxIdleConns == 0 {
		config.MaxIdleConns = 10
	}
	if config.MaxOpenConns == 0 {
		config.MaxOpenConns = 100
	}

	dsn := "user=" + config.Username +
		" password=" + config.Password +
		" host=" + config.Host +
		" port=" + config.Port +
		" dbname=" + config.Dbname +
		" sslmode=disable TimeZone=Asia/Shanghai options='--client_encoding=UTF8'"

	//这里 gorm.Open()函数与之前版本的不一样，大家注意查看官方最新gorm版本的用法
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  logger.Warn, // 日志级别
				IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  true,        // 禁用彩色打印
			},
		),
		NowFunc: func() time.Time {
			tmp := time.Now().Local().Format(vingo.DatetimeFormat)
			now, _ := time.ParseInLocation(vingo.DatetimeFormat, tmp, time.Local)
			return now
		},
	})
	if err != nil {
		panic("Error to Db connection, err: " + err.Error())
	}

	// 连接池配置
	sqlDB, _ := db.DB()
	// 最大空闲数
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	// 最大连接数
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	// 连接最大存活时长
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	// 注册统一异常插件
	RegisterAfterQuery(db)
	RegisterAfterCreate(db)
	RegisterAfterUpdate(db)
	RegisterAfterDelete(db)

	Db = db
}

func RegisterAfterQuery(db *gorm.DB) {
	err := db.Callback().Query().After("gorm:query").Register("gormerror:after_query", func(db *gorm.DB) {
		if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
			panic(&vingo.DbException{Message: db.Error.Error()})
		}
	})
	if err != nil {
		panic(fmt.Sprintf("插件注册失败: %v", err.Error()))
	}
}

func RegisterAfterCreate(db *gorm.DB) {
	err := db.Callback().Create().After("gorm:create").Register("gormerror:after_create", func(db *gorm.DB) {
		if db.Error != nil {
			panic(&vingo.DbException{Message: db.Error.Error()})
		}
	})
	if err != nil {
		panic(fmt.Sprintf("插件注册失败: %v", err.Error()))
	}
}

func RegisterAfterUpdate(db *gorm.DB) {
	err := db.Callback().Update().After("gorm:update").Register("gormerror:after_update", func(db *gorm.DB) {
		if db.Error != nil {
			panic(&vingo.DbException{Message: db.Error.Error()})
		}
	})
	if err != nil {
		panic(fmt.Sprintf("插件注册失败: %v", err.Error()))
	}
}

func RegisterAfterDelete(db *gorm.DB) {
	err := db.Callback().Delete().After("gorm:delete").Register("gormerror:after_delete", func(db *gorm.DB) {
		if db.Error != nil {
			panic(&vingo.DbException{Message: db.Error.Error()})
		}
	})
	if err != nil {
		panic(fmt.Sprintf("插件注册失败: %v", err.Error()))
	}
}
