// *****************************************************************************
// 作者: lgdz
// 创建时间: 2025/2/7
// 描述：
// *****************************************************************************

package pgsql

import (
	"database/sql"
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"strings"
)

type DbApi struct {
	DB     *gorm.DB
	Config Config
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
func (s *DbApi) Updates(model any, column string, columns ...any) *gorm.DB {
	return s.DB.Select(column, columns...).Updates(model)
}

// Delete 删除数据记录
func (s *DbApi) Delete(model any, conds ...any) *gorm.DB {
	return s.DB.Delete(model, conds...)
}

func (s *DbApi) Save(value any) *gorm.DB {
	return s.DB.Save(value)
}

func (s *DbApi) Debug() *gorm.DB {
	return s.DB.Debug()
}

func (s *DbApi) Unscoped() *gorm.DB {
	return s.DB.Unscoped()
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

// 关键词组装
func (s *DbApi) SqlLikeRight(keyword string) string {
	return fmt.Sprintf("%v%%", strings.Trim(keyword, " "))
}

// like模糊查询
func (s *DbApi) LikeOr(db *gorm.DB, keyword string, isRight bool, column ...string) *gorm.DB {
	if keyword != "" {
		var text []string
		for _, item := range column {
			text = append(text, fmt.Sprintf("%v like @text", item))
		}
		var value string
		if isRight {
			value = s.SqlLikeRight(keyword)
		} else {
			value = s.SqlLike(keyword)
		}
		db = db.Where(strings.Join(text, " OR "), sql.Named("text", value))
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

// query参数必须是指针切片类型，如：*[]int|*[]uint|*[]string
func (s *DbApi) QueryWhereIn(db *gorm.DB, query any, column string) *gorm.DB {
	if query != nil {
		db = db.Where(fmt.Sprintf("%v in(?)", column), query)
	}
	return db
}

func (s *DbApi) QueryWhereInString(db *gorm.DB, query vingo.TextString, column string) *gorm.DB {
	if query != "" {
		db = db.Where(fmt.Sprintf("%v in(?)", column), query.ToSlice())
	}
	return db
}

func (s *DbApi) QueryWhereNotInString(db *gorm.DB, query vingo.TextString, column string) *gorm.DB {
	if query != "" {
		db = db.Where(fmt.Sprintf("%v not in(?)", column), query.ToSlice())
	}
	return db
}

func (s *DbApi) QueryWhereInBool(db *gorm.DB, query vingo.BoolString, column string) *gorm.DB {
	if query != "" {
		db = db.Where(fmt.Sprintf("%v in(?)", column), query.ToSlice())
	}
	return db
}

func (s *DbApi) QueryWhereNotInBool(db *gorm.DB, query vingo.BoolString, column string) *gorm.DB {
	if query != "" {
		db = db.Where(fmt.Sprintf("%v not in(?)", column), query.ToSlice())
	}
	return db
}

func (s *DbApi) QueryWhereInInt(db *gorm.DB, query vingo.IntString, column string) *gorm.DB {
	if query != "" {
		db = db.Where(fmt.Sprintf("%v in(?)", column), query.ToSlice())
	}
	return db
}

func (s *DbApi) QueryWhereNotInInt(db *gorm.DB, query vingo.IntString, column string) *gorm.DB {
	if query != "" {
		db = db.Where(fmt.Sprintf("%v not in(?)", column), query.ToSlice())
	}
	return db
}

func (s *DbApi) QueryWhereDateRange(db *gorm.DB, query *vingo.DateRange, column string) *gorm.DB {
	if query != nil {
		db = s.TimeBetween(db, column, vingo.DateAt{query.Start.Format(vingo.DatetimeFormat), query.End.Format(vingo.DatetimeFormat)})
	}
	return db
}

func (s *DbApi) QueryWhereDateAt(db *gorm.DB, query *vingo.DateAt, column string) *gorm.DB {
	if query != nil {
		db = s.TimeBetween(db, column, *query)
	}
	return db
}

func (s *DbApi) QueryWhereDateAtString(db *gorm.DB, query *string, column string) *gorm.DB {
	if query != nil && *query != "" {
		arr := strings.Split(*query, ",")
		db = s.TimeBetween(db, column, vingo.DateAt{arr[0], arr[1]})
	}
	return db
}

func (s *DbApi) FindInSetBuild(column string, value string) string {
	return fmt.Sprintf("'%v' = ANY(string_to_array(%v, ','))", value, column)
}

func (s *DbApi) FindInSet(db *gorm.DB, query any, column string) *gorm.DB {
	db = db.Where(fmt.Sprintf("? = ANY(string_to_array(%v, ','))", column), vingo.ToString(query))
	return db
}

func (s *DbApi) QueryWhereFindInSetInt(db *gorm.DB, query *int, column string) *gorm.DB {
	if query != nil {
		db = s.FindInSet(db, *query, column)
	}
	return db
}

func (s *DbApi) QueryWhereFindInSetUint(db *gorm.DB, query *uint, column string) *gorm.DB {
	if query != nil {
		db = s.FindInSet(db, *query, column)
	}
	return db
}

func (s *DbApi) QueryWhereFindInSetString(db *gorm.DB, query *string, column string) *gorm.DB {
	if query != nil {
		db = s.FindInSet(db, *query, column)
	}
	return db
}

func (s *DbApi) QueryWhereFindInSetUints(db *gorm.DB, query *[]uint, column string) *gorm.DB {
	if query != nil {
		var text []string
		for _, value := range *query {
			text = append(text, fmt.Sprintf("%v = ANY(string_to_array(%v, ','))", vingo.ToString(value), column))
		}
		db = db.Where(strings.Join(text, " OR "))
	}
	return db
}

func (s *DbApi) QueryWhereFindInSetStrings(db *gorm.DB, query *[]string, column string) *gorm.DB {
	if query != nil {
		var text []string
		for _, value := range *query {
			text = append(text, fmt.Sprintf("'%v' = ANY(string_to_array(%v, ','))", vingo.ToString(value), column))
		}
		db = db.Where(strings.Join(text, " OR "))
	}
	return db
}

func (s *DbApi) QueryWhereFindInSetIntString(db *gorm.DB, query vingo.IntString, column string) *gorm.DB {
	if query != "" {
		var text []string
		for _, value := range query.ToUintSlice() {
			text = append(text, fmt.Sprintf("'%v' = ANY(string_to_array(%v, ','))", vingo.ToString(value), column))
		}
		db = db.Where(strings.Join(text, " OR "))
	}
	return db
}

func (s *DbApi) QueryWhereFindInSetTextString(db *gorm.DB, query vingo.TextString, column string) *gorm.DB {
	if query != "" {
		var text []string
		for _, value := range query.ToSlice() {
			text = append(text, fmt.Sprintf("'%v' = ANY(string_to_array(%v, ','))", vingo.ToString(value), column))
		}
		db = db.Where(strings.Join(text, " OR "))
	}
	return db
}

func (s *DbApi) QueryWhereLike(db *gorm.DB, query string, column ...string) *gorm.DB {
	if query != "" {
		db = s.LikeOr(db, query, false, column...)
	}
	return db
}

func (s *DbApi) QueryWhereLikeRight(db *gorm.DB, query string, column ...string) *gorm.DB {
	if query != "" {
		db = s.LikeOr(db, query, true, column...)
	}
	return db
}

func (s *DbApi) QueryWherePath(db *gorm.DB, query string, column string) *gorm.DB {
	if query != "" {
		db = db.Where(fmt.Sprintf("%v='%v' OR %v LIKE '%v,%%'", column, query, column, query))
	}
	return db
}

func (s *DbApi) QueryWherePaths(db *gorm.DB, query []string, column string) *gorm.DB {
	if len(query) > 0 {
		var text []string
		for _, value := range query {
			text = append(text, fmt.Sprintf("(%v='%v' OR %v LIKE '%v,%%')", column, value, column, value))
		}
		db = db.Where(strings.Join(text, " OR "))
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

// 指定字段第一个汉字按A-Z排序
func (s *DbApi) ChineseSortString(column string) string {
	return fmt.Sprintf("CONVERT(SUBSTR(%v, 1, 1) USING gbk)", column)
}

func (s *DbApi) Exists(model any, condition ...any) bool {
	err := s.DB.First(model, condition...).Error
	if err == gorm.ErrRecordNotFound {
		return false
	} else if err != nil {
		panic(err.Error())
	}
	return true
}

func (s *DbApi) TxExists(tx *gorm.DB, model any, condition ...any) bool {
	err := tx.First(model, condition...).Error
	if err == gorm.ErrRecordNotFound {
		return false
	} else if err != nil {
		panic(err.Error())
	}
	return true
}

// 记录不存在时抛出错误
func (s *DbApi) NotExistsErr(model any, condition ...any) {
	err := s.DB.First(model, condition...).Error
	if err == gorm.ErrRecordNotFound {
		panic(err.Error())
	} else if err != nil {
		panic(err.Error())
	}
}

// 记录不存在时抛出错误
func (s *DbApi) NotExistsErrMsg(msg string, model any, condition ...any) {
	err := s.DB.First(model, condition...).Error
	if err == gorm.ErrRecordNotFound {
		panic(msg)
	} else if err != nil {
		panic(err.Error())
	}
}

// 记录不存在时抛出错误(事务内)
func (s *DbApi) TXNotExistsErr(tx *gorm.DB, model any, condition ...any) {
	err := tx.First(model, condition...).Error
	if err == gorm.ErrRecordNotFound {
		panic(err.Error())
	} else if err != nil {
		panic(err.Error())
	}
}

func (s *DbApi) CheckHasChild(model any, id uint) {
	err := s.DB.First(model, "pid=?", id)
	if err.Error != gorm.ErrRecordNotFound {
		panic("记录有子项，删除失败")
	}
}

// 数据库事务自动提交
func (s *DbApi) AutoCommit(tx *gorm.DB, callback ...func()) {
	if r := recover(); r != nil {
		//fmt.Printf("%T\n%v\n", r, r)
		tx.Rollback()
		if len(callback) > 0 && callback[0] != nil {
			callback[0]()
		}
		panic(r)
	} else if err := tx.Statement.Error; err != nil {
		//fmt.Println("数据库异常事务回滚")
		tx.Rollback()
		if len(callback) > 0 && callback[0] != nil {
			callback[0]()
		}
		panic(err.Error())
	} else {
		//fmt.Println("事务提交")
		tx.Commit()
		if len(callback) > 1 && callback[1] != nil {
			callback[1]()
		}
	}
}

type TableColumn struct {
	Column     string `gorm:"column:COLUMN_NAME" json:"column"`
	ColumnType string `gorm:"column:COLUMN_TYPE" json:"-"`
	Type       string `gorm:"column:DATA_TYPE" json:"type"`
	Comment    string `gorm:"column:COLUMN_COMMENT" json:"comment"`
}

// 获取表字段
func (s *DbApi) GetTableColumn(tableName string) []TableColumn {
	var columns = make([]TableColumn, 0)
	s.DB.Raw(`
		SELECT * FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = ?
		ORDER BY ordinal_position ASC
	`, tableName).Scan(&columns)

	for index, item := range columns {
		if slice.Contain([]string{"int", "tinyint", "bigint", "float", "decimal"}, item.Type) {
			columns[index].Type = "number"
		} else if slice.Contain([]string{"char", "varchar", "text", "longtext"}, item.Type) {
			columns[index].Type = "string"
		} else if item.ColumnType == "datetime" {
			columns[index].Type = "datetime"
		}
	}
	return columns
}

type PathOption struct {
	DbApi       *DbApi
	Tx          *gorm.DB
	AppendField []string
	RootAppend  func(s reflect.Value)
	ChildAppend func(s reflect.Value, p reflect.Value)
}

// 设置数据路径，上下级数据结构包含（path、len）字段使用
// model传入必须是指针类型
func SetPath[T any](model *T, parent *T, option PathOption) {
	s := reflect.ValueOf(model).Elem()
	pid := s.FieldByName("Pid").Uint()
	if pid > 0 {
		if parent == nil {
			option.DbApi.TXNotExistsErr(option.Tx, &parent, pid)
		}
		parentValue := reflect.ValueOf(parent).Elem()
		s.FieldByName("Path").SetString(fmt.Sprintf("%v,%d", parentValue.FieldByName("Path").String(), s.FieldByName("Id").Uint()))
		s.FieldByName("Len").SetUint(parentValue.FieldByName("Len").Uint() + 1)
		if option.ChildAppend != nil {
			option.ChildAppend(s, parentValue)
		}
	} else {
		s.FieldByName("Path").SetString(strconv.Itoa(int(s.FieldByName("Id").Uint())))
		s.FieldByName("Len").SetUint(1)
		if option.RootAppend != nil {
			option.RootAppend(s)
		}
	}
	selectFiled := []string{"path", "len"}
	selectFiled = append(selectFiled, option.AppendField...)
	option.Tx.Model(model).Select(selectFiled).UpdateColumns(s.Interface())
}

// 设置所有子级路径，一般在更新pid时使用
func SetPathChild[T any](model *T, option PathOption) {
	s := reflect.ValueOf(model).Elem()
	var rows []T
	option.Tx.Find(&rows, "pid=?", s.FieldByName("Id").Uint())
	for _, row := range rows {
		SetPath[T](&row, model, option)
		SetPathChild[T](&row, option)
	}
}

// 设置自身path和所有子级path
func SetPathAndChildPath[T any](model *T, option PathOption) {
	SetPath[T](model, nil, option)
	SetPathChild[T](model, option)
}

// 事务函数
func (s *DbApi) Commit(handler func(*gorm.DB)) {
	tx := s.Begin()
	defer s.AutoCommit(tx)
	handler(tx)
}

// 通过条件获取单条记录
func Fetch[T any](tx *gorm.DB, condition ...any) (row T) {
	err := tx.First(&row, condition...).Error
	if err == gorm.ErrRecordNotFound {
		panic(&vingo.BackException{Message: err.Error()})
	} else if err != nil {
		panic(&vingo.BackException{Message: err.Error()})
	}
	return
}

// 通过主键id获取记录
func FetchById[T any](tx *gorm.DB, id any, diff ...bool) (row T) {
	row = Fetch[T](tx, "id=?", id)

	// 初始化diff
	if len(diff) > 0 && diff[0] {
		rowVal := reflect.ValueOf(&row).Elem()
		diffField := rowVal.FieldByName("Diff")
		if diffField.IsValid() && diffField.CanSet() {
			diffField.Set(reflect.ValueOf(&vingo.DiffBox{Old: row}))
		}
	}
	return
}
