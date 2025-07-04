package pgsql

import (
	"fmt"
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// 新建一个数据库连接池
func NewPgSql(config Config) *DbApi {
	config.StringValue(&config.Host, "127.0.0.1")
	config.StringValue(&config.Port, "54321")
	config.StringValue(&config.Username, "system")
	config.StringValue(&config.Password, "123456")
	config.StringValue(&config.Charset, "utf8mb4")
	config.IntValue(&config.MaxIdleConns, 10)
	config.IntValue(&config.MaxOpenConns, 100)

	var dbApi = DbApi{
		Config: config,
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", config.Host, config.Port, config.Username, config.Password, config.Dbname)
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
			loc, _ := time.LoadLocation("Asia/Shanghai")
			return time.Now().In(loc)
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

	dbApi.DB = db
	return &dbApi
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
