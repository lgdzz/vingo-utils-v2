package mysql

import (
	"fmt"
	"github.com/lgdzz/vingo-utils/vingo"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitService(config *Config) *gorm.DB {
	config.StringValue(&config.Host, "127.0.0.1")
	config.StringValue(&config.Port, "3306")
	config.StringValue(&config.Username, "root")
	config.StringValue(&config.Password, "123456789")
	config.StringValue(&config.Charset, "utf8mb4")
	config.IntValue(&config.MaxIdleConns, 10)
	config.IntValue(&config.MaxOpenConns, 100)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Dbname,
		config.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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
			tmp := time.Now().Local().Format("2006-01-02 15:04:05")
			now, _ := time.ParseInLocation("2006-01-02 15:04:05", tmp, time.Local)
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

	return db
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
