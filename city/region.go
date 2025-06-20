// *****************************************************************************
// 作者: lgdz
// 创建时间: 2025/6/20
// 描述：
// *****************************************************************************

package city

import "github.com/duke-git/lancet/v2/slice"

type RegionJson struct {
	Code     string       `json:"code"`
	Name     string       `json:"name"`
	Children []RegionJson `json:"children"`
}

// 将字符串区域编码转成数组
func ParseAreaCode(code string) []string {
	splitPoints := []int{2, 4, 6, 9}
	var result []string

	for _, p := range splitPoints {
		if len(code) >= p {
			result = append(result, code[:p])
		}
	}
	return result
}

func FindChildrenByCode(regions []RegionJson, code string) []RegionJson {
	path := ParseAreaCode(code)
	current := regions

	for _, p := range path {
		found := false
		for _, region := range current {
			if region.Code == p {
				current = region.Children
				found = true
				break
			}
		}
		if !found {
			return make([]RegionJson, 0)
		}
	}
	return current
}

func RegionNames(regions []RegionJson) []string {
	var result []string
	slice.ForEach(regions, func(index int, item RegionJson) {
		result = append(result, item.Name)
	})
	return result
}

func RegionCodes(regions []RegionJson) []string {
	var result []string
	slice.ForEach(regions, func(index int, item RegionJson) {
		result = append(result, item.Name)
	})
	return result
}

func GetRegionChildNameByCode(code string, regions []RegionJson) []string {
	return RegionNames(FindChildrenByCode(regions, code))
}

func GetRegionChildCodeByCode(code string, regions []RegionJson) []string {
	return RegionCodes(FindChildrenByCode(regions, code))
}
