// *****************************************************************************
// 作者: lgdz
// 创建时间: 2025/6/9
// 描述：
// *****************************************************************************

package excel_v2

import (
	"github.com/xuri/excelize/v2"
	"mime/multipart"
)

func ReadDataByFromFile(fileHeader *multipart.FileHeader, startRowIndex int, rowFunc ...func(int, []string)) {
	file, err := fileHeader.Open()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 从流中读取
	f, err := excelize.OpenReader(file)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	// 遍历第一个 sheet（你可以按需处理多个）
	sheets := f.GetSheetList()
	for sheetIndex, fn := range rowFunc {
		if sheetIndex >= len(sheets) {
			continue
		}
		rows, err := f.GetRows(sheets[sheetIndex])
		if err != nil {
			panic(err)
		}
		for i := startRowIndex; i < len(rows); i++ {
			row := rows[i]
			if len(row) == 0 || row[0] == "" {
				continue
			}
			fn(i, row)
		}
	}
}
