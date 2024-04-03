// *****************************************************************************
// 作者: lgdz
// 创建时间: 2024/4/3
// 描述：
// *****************************************************************************

package mysql

import (
	"database/sql"
	"fmt"
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

type DbApi struct {
	DB *gorm.DB
}

func (s *DbApi) NewDB() *gorm.DB {
	return s.DB
}

// Begin 开始事务
func (s *DbApi) Begin(opts ...*sql.TxOptions) *gorm.DB {
	return s.DB.Begin(opts...)
}

// Create 创建数据记录
func (s *DbApi) Create(value any) *gorm.DB {
	return s.DB.Create(value)
}

// FirstOrCreate 不存在则创建
func (s *DbApi) FirstOrCreate(dest any, conds ...any) *gorm.DB {
	return s.DB.FirstOrCreate(dest, conds...)
}

// Updates 更新指定模型字段
func (s *DbApi) Updates(model *any, column string, columns ...any) *gorm.DB {
	return s.DB.Select(column, columns...).Updates(model)
}

// Delete 删除数据记录
func (s *DbApi) Delete(model *any, conds ...any) *gorm.DB {
	return s.DB.Delete(model, conds...)
}

func (s *DbApi) Save(value any) *gorm.DB {
	return s.DB.Save(value)
}

func (s *DbApi) Debug() *gorm.DB {
	return s.DB.Debug()
}

func (s *DbApi) Table(name string, args ...any) *gorm.DB {
	return s.DB.Table(name, args...)
}

func (s *DbApi) Model(value any) *gorm.DB {
	return s.DB.Model(value)
}

func (s *DbApi) Select(query any, args ...any) *gorm.DB {
	return s.DB.Select(query, args...)
}

func (s *DbApi) Where(query any, args ...any) *gorm.DB {
	return s.DB.Where(query, args...)
}

func (s *DbApi) Order(value any) *gorm.DB {
	return s.DB.Order(value)
}

func (s *DbApi) Like(db *gorm.DB, keyword string) *gorm.DB {
	if keyword != "" {
		db = db.Where("name like @text OR description like @text", sql.Named("text", s.SqlLike(keyword)))
	}
	return db
}

// 关键词组装
func (s *DbApi) SqlLike(keyword string) string {
	return fmt.Sprintf("%%%v%%", strings.Trim(keyword, " "))
}

// like模糊查询
func (s *DbApi) LikeOr(db *gorm.DB, keyword string, column ...string) *gorm.DB {
	if keyword != "" {
		var text []string
		for _, item := range column {
			text = append(text, fmt.Sprintf("%v like @text", item))
		}
		db = db.Where(strings.Join(text, " OR "), sql.Named("text", s.SqlLike(keyword)))
	}
	return db
}

// 时间范围查询
func (s *DbApi) TimeBetween(db *gorm.DB, column string, dateAt vingo.DateAt) *gorm.DB {
	return db.Where(fmt.Sprintf("%v BETWEEN ? AND ?", column), dateAt.Start(), dateAt.End())
}

func (s *DbApi) QueryWhere(db *gorm.DB, query any, column string) *gorm.DB {
	valueOf := reflect.ValueOf(query)
	typeOf := valueOf.Type()
	if typeOf.Kind() == reflect.Ptr {
		if valueOf.IsNil() {
			//fmt.Println("空指针无条件")
			return db
		} else {
			query = valueOf.Elem().Interface()
		}
	} else {
		switch v := query.(type) {
		case string:
			if v == "" {
				//fmt.Println("string无条件")
				return db
			}
		}
		query = valueOf.Interface()
	}
	if query != nil {
		db = db.Where(fmt.Sprintf("%v=?", column), query)
	}
	return db
}

func (s *DbApi) QueryWhereDateAt(db *gorm.DB, query *vingo.DateAt, column string) *gorm.DB {
	if query != nil {
		db = s.TimeBetween(db, column, *query)
	}
	return db
}

func (s *DbApi) QueryWhereLike(db *gorm.DB, query string, column ...string) *gorm.DB {
	if query != "" {
		db = s.LikeOr(db, query, column...)
	}
	return db
}

func (s *DbApi) QueryWhereBetween(db *gorm.DB, query *[2]any, column string) *gorm.DB {
	if query != nil {
		db = db.Where(fmt.Sprintf("%v BETWEEN ? AND ?", column), query[0], query[1])
	}
	return db
}

func (s *DbApi) QueryWhereDeletedAt(db *gorm.DB, column string) *gorm.DB {
	db = db.Where(fmt.Sprintf("%v IS NULL", column))
	return db
}
