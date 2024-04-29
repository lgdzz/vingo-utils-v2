package sqlite

import (
	"fmt"
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// Db 连接池句柄
var Db *gorm.DB

func InitSqliteService() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
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
