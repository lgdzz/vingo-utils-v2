package mysql

import (
	"database/sql"
	"fmt"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/lgdzz/vingo-utils/vingo"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

const tpl = `// *****************************************************************************
// 作者: lgdz
// 创建时间: {{ .Date }}
// 描述：
// *****************************************************************************

package model

type {{ .ModelName }} struct {
	{{ range .TableColumns }}{{ .DataName }}   {{ .DataType }}  ` + "`gorm:\"{{ if eq .Key \"PRI\" }}primaryKey;{{ end }}column:{{ .Field }}\" json:\"{{ .JsonName }}\"`" + ` {{ if .Comment }}// {{ .Comment }}{{ end }}
	{{ end }}
}

func (s *{{ .ModelName }}) TableName() string {
	return "{{ .TableName }}"
}

type {{ .ModelName }}Query struct {
	page.Limit
	Keyword string ` + "`form:\"keyword\"`" + `
}

type {{ .ModelName }}Body struct {
	{{ .ModelName }}
}
`

type TableData struct {
	TableName    string
	ModelName    string
	TableComment string
	TableColumns []Column
	Date         string
}

type Column struct {
	Field    string
	Type     string
	Null     string
	Key      string
	Default  sql.NullString
	Extra    string
	Comment  string
	DataName string
	DataType string
	JsonName string
}

func (s *DbApi) CreateDbModel(tableNames ...string) (bool, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("请检查数据库是否正常连接")
			fmt.Println(r)
		}
	}()
	vingo.Mkdir("model")
	for _, tableName := range tableNames {

		var modelPath = filepath.Join(".", "model", tableName+".go")
		// 模型文件存在则不创建
		if vingo.FileExists(modelPath) {
			continue
		}

		var columns []Column
		s.DB.Raw("SHOW FULL COLUMNS FROM " + tableName).Select("Field,Type,Comment").Scan(&columns)
		columns = vingo.ForEach[Column](columns, func(item Column, index int) Column {
			if vingo.StringStartsWith(item.Type, []string{"int", "tinyint", "smallint"}) {
				if vingo.StringContainsAnd(item.Type, "unsigned") {
					item.DataType = "uint"
				} else {
					item.DataType = "int"
				}
			} else if vingo.StringStartsWith(item.Type, []string{"decimal"}) {
				item.DataType = "float64"
			} else if item.Field == "deleted_at" {
				item.DataType = "gorm.DeletedAt"
			} else if vingo.StringStartsWith(item.Type, []string{"datetime"}) {
				item.DataType = "*vingo.LocalTime"
			} else {
				item.DataType = "string"
			}
			item.JsonName = strutil.CamelCase(item.Field)
			item.DataName = strutil.UpperFirst(item.JsonName)
			return item
		})

		// 渲染模板
		t, err := template.New("tpl").Parse(tpl)
		if err != nil {
			fmt.Println(err)
			return false, err
		}

		outputFile, err := os.Create(modelPath)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		defer outputFile.Close()

		if err = t.Execute(outputFile, TableData{
			TableName:    tableName,
			ModelName:    strutil.UpperFirst(strutil.CamelCase(tableName)),
			TableColumns: columns,
			Date:         time.Now().Format("2006/01/02"),
		}); err != nil {
			fmt.Println(err)
			return false, err
		}
	}
	return true, nil
}
