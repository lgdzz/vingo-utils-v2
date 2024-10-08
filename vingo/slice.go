package vingo

import (
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"reflect"
	"strconv"
)

// 通用切片去重
// s := []string{"a", "b", "c", "a", "b", "a"}
// SliceUnique(&s)
// print：["a", "b", "c"]
//
// i := []int{1, 1, 2, 3, 3, 1, 3}
// SliceUnique(&i)
// print：[1, 2, 3]
func SliceUnique(values []any) []any {
	return slice.Unique(values)
}

// 将[]string数据去重返回
func SliceUniqueString(slice []string) []string {
	uniqueMap := make(map[string]interface{})
	for _, v := range slice {
		uniqueMap[v] = nil
	}
	var uniqueSlice []string
	for k := range uniqueMap {
		uniqueSlice = append(uniqueSlice, k)
	}
	return uniqueSlice
}

// 将[]uint数据去重返回
func SliceUniqueUint(slice []uint) []uint {
	uniqueMap := make(map[uint]interface{})
	for _, v := range slice {
		uniqueMap[v] = nil
	}
	var uniqueSlice []uint
	for k := range uniqueMap {
		uniqueSlice = append(uniqueSlice, k)
	}
	return uniqueSlice
}

// 将[]int数据去重返回
func SliceUniqueInt(slice []int) []int {
	uniqueMap := make(map[int]interface{})
	for _, v := range slice {
		uniqueMap[v] = nil
	}
	var uniqueSlice []int
	for k := range uniqueMap {
		uniqueSlice = append(uniqueSlice, k)
	}
	return uniqueSlice
}

// []string转[]int
func SliceStringToInt(s []string) []int {
	slice := make([]int, 0)
	for _, v := range s {
		num, _ := strconv.Atoi(v)
		slice = append(slice, num)
	}
	return slice
}

// []string转[]uint
func SliceStringToUint(s []string) []uint {
	slice := make([]uint, 0)
	for _, v := range s {
		num, _ := strconv.Atoi(v)
		slice = append(slice, uint(num))
	}
	return slice
}

// []string转[]float64
func SliceStringToFloat64(s []string) []float64 {
	slice := make([]float64, 0)
	for _, v := range s {
		num, _ := strconv.Atoi(v)
		slice = append(slice, float64(num))
	}
	return slice
}

// []int转[]string
func SliceIntToString(s []int) []string {
	slice := make([]string, 0)
	for _, v := range s {
		slice = append(slice, strconv.Itoa(v))
	}
	return slice
}

// []uint转[]string
func SliceUintToString(s []uint) []string {
	slice := make([]string, 0)
	for _, v := range s {
		slice = append(slice, strconv.Itoa(int(v)))
	}
	return slice
}

// []float64转[]string
func SliceFloat64ToString(s []float64) []string {
	slice := make([]string, 0)
	for _, v := range s {
		slice = append(slice, fmt.Sprintf("%f", v))
	}
	return slice
}

// 判断一个节点是否在切片中，与IsInSlice函数不同，该函数支持更多场景，而IsInSlice只适合切片类型
// 判断字符串是否在字符串切片中
// 判断数字是否在整型切片中
// 判断字符串是否在字符串字典中
// 判断结构体是否在结构体切片中
func IsInSliceAny(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

// 判断参数是否是切片类型
func IsSlice(slice interface{}) bool {
	sliceType := reflect.TypeOf(slice)
	return sliceType.Kind() == reflect.Slice
}

// 切片取差集
func UintSliceDiff(slice1 []uint, slices ...[]uint) []uint {
	m := make(map[uint]bool)
	result := make([]uint, 0)

	for _, item := range slices {
		for _, val := range item {
			m[val] = true
		}
	}

	for _, item := range slice1 {
		if _, exists := m[item]; !exists {
			result = append(result, item)
		}
	}

	return result
}

// 切片取交集
func UintSliceIntersect(slices ...[]uint) []uint {
	if len(slices) == 0 {
		return nil
	}

	intersect := make([]uint, 0)
	set := make(map[uint]int)

	// 计算元素在数组中出现的次数
	for i, arr := range slices {
		for _, v := range arr {
			set[v] = i + 1
		}
	}

	// 检查元素是否在所有数组中都出现
	for k, v := range set {
		if v == len(slices) {
			intersect = append(intersect, k)
		}
	}

	return intersect
}

// 判断一个节点是否在切片中
func IsInSlice(item interface{}, items interface{}) bool {
	s := reflect.ValueOf(items)
	if s.Kind() != reflect.Slice {
		panic("not a slice")
	}

	for i := 0; i < s.Len(); i++ {
		if s.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func GetSliceElement(slice interface{}, index int) interface{} {
	value := reflect.ValueOf(slice)
	if value.Kind() != reflect.Slice {
		panic("GetSliceElement函数参数1不是切片类型")
	}
	if index >= value.Len() {
		panic(fmt.Sprintf("Index out of range: %d", index))
	}
	element := value.Index(index)
	if !element.IsValid() {
		panic(fmt.Sprintf("Element does not exist: %d", index))
	}
	return element.Interface()
}

// 在切片中搜索元素，返回索引，-1未找到
func IndexOf(item interface{}, items interface{}) int {
	value := reflect.ValueOf(items)
	if value.Kind() != reflect.Slice {
		panic("IndexOf函数参数2不是切片类型")
	}

	for i := 0; i < value.Len(); i++ {
		if reflect.DeepEqual(value.Index(i).Interface(), item) {
			return i
		}
	}
	return -1
}

// 将数组对象转字典对象
func SliceToMapSlice[T any](slice []T, column string) map[string]T {
	var result = map[string]T{}
	for _, row := range slice {
		var rowValue = reflect.ValueOf(row)
		var keyValue = rowValue.FieldByName(column)
		var keyString string
		if keyValue.Kind() != reflect.String {
			keyString = ToString(keyValue.Interface())
		} else {
			keyString = keyValue.Interface().(string)
		}
		result[keyString] = row
	}
	return result
}

func ForEach[T any, R any](collection []T, callback func(item T, index int) R) []R {
	result := make([]R, 0)

	for i, item := range collection {
		result = append(result, callback(item, i))
	}

	return result
}

// 将元素添加到切片中
func SlicePush[T any](items *[]T, item T) {
	*items = append(*items, item)
}

// 将元素从切片中移除，如出现多次则删除多个
func SliceRemove[T any](items *[]T, item T) {
	for i := 0; i < len(*items); i++ {
		if reflect.DeepEqual((*items)[i], item) {
			// 执行删除操作
			*items = append((*items)[:i], (*items)[i+1:]...)
			i-- // 调整索引以处理下一个元素
		}
	}
}
