package pgsql

import (
	"database/sql"
	"fmt"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

const tpl = `// *****************************************************************************
// 作者: lgdz
// 创建时间: {{ .Date }}
// 描述：{{ .TableComment }}
// *****************************************************************************

package model

import(
	"github.com/lgdzz/vingo-utils-v2/vingo"
	"github.com/lgdzz/vingo-utils-v2/db/page"
	pgsql "github.com/lgdzz/vingo-utils-v2/db/pgsql"
	"gorm.io/gorm"
)

type {{ .ModelName }} struct {
	{{ range .TableColumns }}{{ .DataName }}   {{ .DataType }}  ` + "`gorm:\"{{ if eq .Key \"PRI\" }}primaryKey;{{ end }}column:{{ .Field }}\" json:\"{{ .JsonName }}\"`" + ` {{ if .Comment }}// {{ .Comment }}{{ end }}
    {{ end }}
}

func (s *{{ .ModelName }}) TableName() string {
	return "{{ .TableName }}"
}

type {{ .ModelName }}Query struct {
	*pgsql.PageQuery
	CreatedAt *string ` + "`form:\"createdAt\"`" + `
	Keyword string ` + "`form:\"keyword\"`" + `
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
		if vingo.FileExists(modelPath) {
			continue
		}

		var tableComment sql.NullString
		queryTableComment := `
	SELECT obj_description(c.oid, 'pg_class') 
	FROM pg_class c
	WHERE relname = ?
`
		s.DB.Raw(queryTableComment, tableName).Scan(&tableComment)

		commentStr := ""
		if tableComment.Valid {
			commentStr = tableComment.String
		}

		var columns []Column
		queryColumn := `
			SELECT 
				a.attname AS field,
				format_type(a.atttypid, a.atttypmod) AS type,
				col_description(a.attrelid, a.attnum) AS comment,
				CASE WHEN a.attnotnull THEN 'NO' ELSE 'YES' END AS null
			FROM 
				pg_attribute a
			JOIN 
				pg_class c ON a.attrelid = c.oid
			JOIN 
				pg_namespace n ON c.relnamespace = n.oid
			WHERE 
				c.relname = ? 
				AND a.attnum > 0 
				AND NOT a.attisdropped
		`
		s.DB.Raw(queryColumn, tableName).Scan(&columns)

		columns = vingo.ForEach[Column](columns, func(item Column, index int) Column {
			typ := item.Type
			if vingo.StringStartsWith(typ, []string{"integer", "int4"}) {
				item.DataType = "uint"
			} else if vingo.StringStartsWith(typ, []string{"bigint", "int8"}) {
				item.DataType = "uint"
			} else if vingo.StringStartsWith(typ, []string{"smallint", "int2"}) {
				item.DataType = "uint"
			} else if vingo.StringStartsWith(typ, []string{"numeric", "decimal", "double precision", "real"}) {
				item.DataType = "float64"
			} else if vingo.StringStartsWith(typ, []string{"boolean"}) {
				item.DataType = "bool"
			} else if vingo.StringStartsWith(typ, []string{"timestamp", "date", "time"}) {
				item.DataType = "*vingo.LocalTime"
			} else {
				item.DataType = "string"
			}
			item.JsonName = strutil.CamelCase(item.Field)
			item.DataName = strutil.UpperFirst(item.JsonName)

			// 判断主键
			var isPrimaryKey string
			s.DB.Raw(`
				SELECT a.attname 
				FROM pg_index i 
				JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = ANY(i.indkey)
				WHERE i.indrelid = ?::regclass AND i.indisprimary
			`, tableName).Scan(&isPrimaryKey)

			if isPrimaryKey == item.Field {
				item.Key = "PRI"
			}
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
			TableComment: commentStr,
			TableColumns: columns,
			Date:         time.Now().Format("2006/01/02"),
		}); err != nil {
			fmt.Println(err)
			return false, err
		}
	}
	return true, nil
}
