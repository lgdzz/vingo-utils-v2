package vingo

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

func TreeBuildString(list *[]map[string]any, id string, pidName string) (result []map[string]any) {
	for _, row := range *list {

		if row[pidName] != id {
			continue
		}

		children := TreeBuildString(list, row["id"].(string), pidName)

		if len(children) > 0 {
			row["hasChild"] = true
			row["children"] = children
		} else {
			row["hasChild"] = false
		}
		row["id"] = GetUUID()
		result = append(result, row)
	}
	return
}

func TreeBuild(list *[]map[string]any, id uint, pidName string, already *[]uint) (result []map[string]any) {

	for _, row := range *list {

		if ToUint(row[pidName]) != id {
			continue
		}

		*already = append(*already, id)

		children := TreeBuild(list, ToUint(row["id"]), pidName, already)

		if len(children) > 0 {
			row["hasChild"] = true
			row["children"] = children
		} else {
			row["hasChild"] = false
		}
		result = append(result, row)
	}
	return
}

func TreeBuilds(list *[]map[string]any, ids []uint, pidName string) []map[string]any {
	result := make([]map[string]any, 0)
	already := make([]uint, 0)
	for _, id := range ids {
		if IsInSlice(id, already) {
			continue
		}
		result = append(result, TreeBuild(list, id, pidName, &already)...)
	}
	return result
}

// enable：如果为true时，则过滤掉禁用的数据
func Tree[T any](rows *[]T, pidName string, enable bool) []map[string]any {
	var hideIds = make([]uint, 0)
	var ids = make([]uint, 0)
	var newRows = make([]T, 0)
	for _, row := range *rows {
		rowValue := reflect.ValueOf(row)

		if enable {
			var status = uint(rowValue.FieldByName("Status").Uint())
			var currentId = uint(rowValue.FieldByName("Id").Uint())
			if status != 1 {
				hideIds = append(hideIds, currentId)
				continue
			}
			var path = strings.Split(rowValue.FieldByName("Path").String(), ",")
			var cHide bool
			for _, p := range path {
				var currentP = ToUint(p)
				if IsInSlice(currentP, hideIds) && !IsInSlice(currentId, hideIds) {
					cHide = true
					break
				}
			}
			if cHide {
				hideIds = append(hideIds, currentId)
				continue
			}

			newRows = append(newRows, row)
		}
		ids = append(ids, uint(rowValue.FieldByName("Pid").Uint()))
	}
	if enable {
		rows = &newRows
	}
	var list []map[string]any
	CustomOutput(rows, &list)
	return TreeBuilds(&list, ids, pidName)
}

func CallStructFunc(obj any, method string, param map[string]any) any {
	t := reflect.TypeOf(obj)
	_func, ok := t.MethodByName(method)
	if !ok {
		panic(fmt.Sprintf("%v方法不存在", method))
	}

	_param := make([]reflect.Value, 0)
	_param = append(_param, reflect.ValueOf(obj))
	for _, value := range param {
		_param = append(_param, reflect.ValueOf(value))
	}
	res := _func.Func.Call(_param)
	return res[0].Interface()
}

func CallStructFuncNoResult(obj any, method string, param map[string]any) {
	t := reflect.TypeOf(obj)
	_func, ok := t.MethodByName(method)
	if !ok {
		panic(fmt.Sprintf("%v方法不存在", method))
	}

	_param := make([]reflect.Value, 0)
	_param = append(_param, reflect.ValueOf(obj))
	for _, value := range param {
		_param = append(_param, reflect.ValueOf(value))
	}
	_func.Func.Call(_param)
}

func CheckErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// Deprecated: This function is no longer recommended for use.
// Suggested: Please use mysql.SqlLike instead.
func SqlLike(keyword string) string {
	return fmt.Sprintf("%%%v%%", strings.Trim(keyword, " "))
}

/**
 * 增长率计算
 * @param float64 $now 现在值
 * @param float64 $prev 过去值
 * @return string
 */
func ComputeGrowRate(now float64, prev float64) string {
	if now == prev {
		return "0.00"
	} else if prev == 0 {
		return "-"
	} else {
		return fmt.Sprintf("%.2f", ((now - prev) / prev * 100))
	}
}

/**
 * 根据起点坐标和终点坐标测距离
 * @param [2]float64 $from [起点坐标(经纬度),例如:[2]float64{118.012951,36.810024}]
 * @param [2]float64 $to [终点坐标(经纬度)]
 * @param bool $km 是否以公里为单位 false:米 true:公里(千米)
 * @param int $decimal 精度 保留小数位数
 * @return float  距离数值
 */
func Distance(from Location, to Location, km bool, decimal int) float64 {
	EARTH_RADIUS := 6370.996 // 地球半径系数
	fromSorted := from
	toSorted := to
	if from.Lng > to.Lng {
		fromSorted = to
		toSorted = from
	}

	dLat := (toSorted.Lng - fromSorted.Lng) * math.Pi / 180
	dLon := (toSorted.Lat - fromSorted.Lat) * math.Pi / 180

	fromLat := fromSorted.Lng * math.Pi / 180
	toLat := toSorted.Lng * math.Pi / 180

	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(fromLat)*math.Cos(toLat)*math.Pow(math.Sin(dLon/2), 2)
	c := 2 * math.Asin(math.Sqrt(a))

	distance := EARTH_RADIUS * c * 1000

	if km {
		distance = distance / 1000
	}

	return math.Round(distance*math.Pow10(decimal)) / math.Pow10(decimal)
}

// 密码加密
func PasswordToCipher(text string, salt string) string {
	return MD5(MD5(text) + salt)
}

// 密码强度验证
// level-2：中等密码，任意两种字符组合
// level-3：复杂密码，必须包含四种字符组合
func PasswordStrength(password string, level int) {
	if len(password) < 6 || len(password) > 18 {
		panic("密码长度需符合6-18个字符长度要求")
	}
	if level == 2 {
		// 中等密码，任意两种字符组合
		hasDigit := false
		hasUpper := false
		hasLower := false
		hasSpecial := false

		for _, ch := range password {
			if unicode.IsDigit(ch) {
				hasDigit = true
			} else if unicode.IsUpper(ch) {
				hasUpper = true
			} else if unicode.IsLower(ch) {
				hasLower = true
			} else if unicode.IsPunct(ch) || unicode.IsSymbol(ch) {
				hasSpecial = true
			}
		}
		if !(hasDigit && hasUpper) && !(hasDigit && hasLower) && !(hasDigit && hasSpecial) &&
			!(hasUpper && hasLower) && !(hasUpper && hasSpecial) && !(hasLower && hasSpecial) {
			panic("密码需满足两种以上的字符组合（数字、大写字母、小写字母、特殊符号）")
		}
	} else if level == 3 {
		// 复杂密码，必须包含四种字符组合
		hasDigit := false
		hasUpper := false
		hasLower := false
		hasSpecial := false

		for _, ch := range password {
			if unicode.IsDigit(ch) {
				hasDigit = true
			} else if unicode.IsUpper(ch) {
				hasUpper = true
			} else if unicode.IsLower(ch) {
				hasLower = true
			} else if unicode.IsPunct(ch) || unicode.IsSymbol(ch) {
				hasSpecial = true
			}
		}
		if !hasDigit || !hasUpper || !hasLower || !hasSpecial {
			panic("密码需满足四种字符组合（数字、大写字母、小写字母、特殊符号）")
		}
	}
}

// 返回传入参数的指针
func Of[T any](v T) *T {
	return &v
}

// 返回传入参数的指针
// Deprecated: This function is no longer recommended for use.
// Suggested: Please use Of() instead.
func StringPointer(text string) *string {
	return &text
}

// Deprecated: This function is no longer recommended for use.
// Suggested: Please use Of() instead.
func Int64Pointer(v int64) *int64 {
	return &v
}

// Deprecated: This function is no longer recommended for use.
// Suggested: Please use Of() instead.
func IntPointer(v int) *int {
	return &v
}

// Deprecated: This function is no longer recommended for use.
// Suggested: Please use Of() instead.
func UintPointer(v uint) *uint {
	return &v
}

// Deprecated: This function is no longer recommended for use.
// Suggested: Please use Of() instead.
func Float64Pointer(v float64) *float64 {
	return &v
}

// Deprecated: This function is no longer recommended for use.
// Suggested: Please use Of() instead.
func BoolPointer(v bool) *bool {
	return &v
}

// 版本号自增
func IncrementVersion(version string) (string, error) {
	if version == "" {
		return "1.0.0", nil
	}
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("Invalid version format. Expected format: major.minor.patch")
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", fmt.Errorf("Invalid major version: %v", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("Invalid minor version: %v", err)
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", fmt.Errorf("Invalid patch version: %v", err)
	}

	patch++
	if patch >= 10 {
		patch = 0
		minor++
		if minor >= 10 {
			minor = 0
			major++
		}
	}

	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}

// 获取当前项目模块名称(mod-name)
func GetModuleName() (name string) {
	// 获取当前项目的根目录路径
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Println("无法获取当前目录路径：", err)
		return
	}

	// 执行go mod命令获取模块名称
	cmd := exec.Command("go", "list", "-m")
	cmd.Dir = rootDir
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("无法获取模块名称：", err)
		return
	}

	// 解析输出结果，获取模块名称
	name = strings.TrimSpace(string(output))

	return
}

// 三元运算
func SY[T any](condition bool, trueValue T, falseValue T) T {
	if condition {
		return trueValue
	} else {
		return falseValue
	}
}

// 获取当前函数名
func GetCurrentFunctionName() string {
	pc, _, _, _ := runtime.Caller(1)
	currentFunction := runtime.FuncForPC(pc).Name()
	return currentFunction
}
