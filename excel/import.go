package excel

import (
	"github.com/tealeg/xlsx"
	"os"
)

// 读取excel表格中的数据
// 示例：
//
//	type Person struct {
//		Name  string
//		Phone string
//	}
//
//	excel.ReadData("test.xlsx", func(cells []*xlsx.Cell) {
//			person = append(person, Person{
//				Name:  cells[0].String(),
//				Phone: cells[1].String(),
//			})
//		})
//
// Deprecated: This function is no longer recommended for use.
// Suggested: Please use ReadDatas() instead.
func ReadData(excelPath string, rowFunc ...func([]*xlsx.Cell)) {
	// 打开Excel文件
	xlFile, err := xlsx.OpenFile(excelPath)
	if err != nil {
		panic(err.Error())
	}

	for sheetIndex := range rowFunc {
		sheet := xlFile.Sheets[sheetIndex]

		// 遍历每一行（忽略表头）
		for _, row := range sheet.Rows[1:] {
			// 读取单元格数据
			rowFunc[sheetIndex](row.Cells)
		}
	}

	// 删除Excel文件
	err = os.Remove(excelPath)
	if err != nil {
		panic(err.Error())
	}
}

func ReadDatas(file string, startRowIndex int, rowFunc ...func(int, []*xlsx.Cell)) {
	// 打开Excel文件
	xlFile, err := xlsx.OpenFile(file)
	if err != nil {
		panic(err.Error())
	}

	for sheetIndex, f := range rowFunc {
		currentIndex := sheetIndex
		if currentIndex < len(xlFile.Sheets) {
			sheet := xlFile.Sheets[sheetIndex]
			for rowIndex, row := range sheet.Rows[startRowIndex:] {
				if len(row.Cells) == 0 {
					continue
				} else if row.Cells[0].String() == "" {
					// 如果第一列单元格为空，则认定该行为空数据
					continue
				}
				f(rowIndex, row.Cells)
			}
		}
	}

	// 删除Excel文件
	err = os.Remove(file)
	if err != nil {
		panic(err.Error())
	}
}

func ReadCell(cells []*xlsx.Cell, index int) *xlsx.Cell {
	endIndex := len(cells) - 1
	if index >= 0 && index <= endIndex {
		return cells[index]
	} else {
		return &xlsx.Cell{}
	}
}
