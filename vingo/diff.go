package vingo

import (
	"fmt"
	"reflect"
)

type DiffBox struct {
	Old    any
	New    any
	Result *map[string]DiffItem
}

type DiffItem struct {
	Column   string
	OldValue any
	NewValue any
	Message  string
}

func (s *DiffItem) SetMessage() {
	s.Message = fmt.Sprintf("将%v的值[%v]变更为[%v]；", s.Column, s.OldValue, s.NewValue)
}

// 设置新值
func (s *DiffBox) SetNew(newValue any) {
	s.New = newValue
}

// 设置新值并且执行比对
func (s *DiffBox) SetNewAndCompare(newValue any, result func(diff *DiffBox)) {
	s.New = newValue
	s.Compare()
	if result != nil {
		result(s)
	}
}

// 比较
func (s *DiffBox) Compare() {
	result := map[string]DiffItem{}
	oldVal := reflect.ValueOf(s.Old)
	newVal := reflect.ValueOf(s.New)
	if oldVal.Kind() != reflect.Struct || newVal.Kind() != reflect.Struct {
		return
	}
	oldType := oldVal.Type()
	newType := newVal.Type()
	if oldType != newType {
		return
	}
	for i := 0; i < oldVal.NumField(); i++ {
		oldField := oldVal.Field(i)
		newField := newVal.Field(i)

		if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
			name := oldType.Field(i).Name
			if IsInSlice(name, []string{"CreatedAt", "UpdatedAt", "DeletedAt"}) {
				continue
			}
			diffItem := DiffItem{
				Column:   name,
				OldValue: oldField.Interface(),
				NewValue: newField.Interface(),
			}
			diffItem.SetMessage()
			result[diffItem.Column] = diffItem
		}
	}
	s.Result = &result
}

// 判断指定字段是否被修改，被修改返回true
func (s *DiffBox) IsChange(column string) bool {
	if s.Result == nil {
		s.Compare()
	}
	_, ok := (*s.Result)[column]
	return ok
}

// 判断指定字段是否被修改，被修改则进行相应处理
func (s *DiffBox) IsModify(column string, callback func()) {
	if s.Result == nil {
		s.Compare()
	}
	_, ok := (*s.Result)[column]
	if ok && callback != nil {
		callback()
	}
}

// 批量或判断（只要有一个被修改则返回真）
func (s *DiffBox) IsChangeOr(column ...string) bool {
	for _, item := range column {
		if s.IsChange(item) {
			return true
		}
	}
	return false
}

// 批量且判断（只要有一个未修改则返回假）
func (s *DiffBox) IsChangeAnd(column ...string) bool {
	for _, item := range column {
		if !s.IsChange(item) {
			return false
		}
	}
	return false
}

// 返回字符串结果
func (s *DiffBox) ResultContent() string {
	var text string
	for _, item := range *s.Result {
		text += item.Message
	}
	if text == "" {
		return "无修改"
	}
	return text
}
