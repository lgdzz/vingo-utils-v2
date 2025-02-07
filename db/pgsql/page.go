package pgsql

import (
	"fmt"
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"gorm.io/gorm"
	"strings"
)

type PageResult struct {
	Page  int   `json:"page"`
	Size  int   `json:"size"`
	Total int64 `json:"total"` // 总的记录数
	Items any   `json:"items"` // 查询数据列表
}

type PageLimit struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func (s *PageLimit) GetPage() int {
	if s.Page > 0 {
		return s.Page
	} else {
		return 1
	}
}

func (s *PageLimit) GetSize() int {
	if s.Size > 0 {
		return s.Size
	} else {
		return 10
	}
}

func (s *PageLimit) Offset() int64 {
	if s.Page > 0 {
		return int64((s.Page - 1) * s.Size)
	} else {
		return 0
	}
}

type PageOrder struct {
	Column string `form:"sortField"`
	Sort   string `form:"sortOrder"`
}

type PageQuery struct {
	Limit PageLimit
	Order *PageOrder
}

func (s *PageOrder) HandleColumn() string {
	var sort = strings.ToLower(s.Sort)
	if sort != "asc" && sort != "desc" {
		panic("存在sql注入的风险")
	}
	var items = strings.Split(s.Column, ".")
	for index, item := range items {
		items[index] = "`" + item + "`"
	}
	return fmt.Sprintf("%v %v", strings.Join(items, "."), s.Sort)
}

type PageOption[T any] struct {
	Db       *gorm.DB     // 必须
	Query    PageQuery    // 必须
	DefOrder *PageOrder   // 默认排序
	Orders   *[]PageOrder // 服务端指定多个排序条件
	Handle   func(T) any
}

// 创建一个新的分页查询
func NewPage[T any](option PageOption[T]) (result PageResult) {

	if option.Query.Order == nil && option.Orders == nil {
		option.Query.Order = option.DefOrder
	}

	var count int64
	var items = make([]T, 0)
	option.Db.Count(&count)
	result.Total = count
	result.Page = option.Query.Limit.GetPage()
	result.Size = option.Query.Limit.GetSize()
	if count > 0 {
		option.Db = option.Db.Order(option.BuildOrderString())
		option.Db.Limit(option.Query.Limit.GetSize()).Offset(int(option.Query.Limit.Offset())).Scan(&items)

		if option.Handle != nil {
			result.Items = vingo.ForEach(items, func(item T, index int) any {
				return option.Handle(item)
			})
			return
		}
	}
	result.Items = items
	return
}

func (s *PageOption[T]) BuildOrderString() string {
	// 默认排序
	if s.Query.Order == nil && s.Orders == nil {
		return "`id` desc"
	}

	if s.Query.Order != nil {
		s.Orders = &[]PageOrder{*s.Query.Order}
	}

	var orders = make([]string, 0)
	for _, item := range *s.Orders {
		orders = append(orders, item.HandleColumn())
	}
	return strings.Join(orders, ",")
}
