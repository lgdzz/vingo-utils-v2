package pgsql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lgdzz/vingo-utils/vingo"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type TableItem struct {
	Name    string
	Comment string
	Columns []Column
}

type Database struct {
	Name        string
	ReleaseTime string
	Tables      []TableItem
}

const booktpl = `
<!DOCTYPE html>
<html>
<head>
  <title>{{ .Name }} 数据字典</title>
  <style>
    body {
        width: 1000px;
    	margin: 0 auto;
		font-size: 14px;
		padding-bottom: 50px;
    }
    table {
      border-collapse: collapse;
      width: 100%;
    }

    th, td {
      border: 1px solid #ddd;
      padding: 8px;
      text-align: left;
    }

    th {
      background-color: #f2f2f2;
    }
    
    .menu {
      position: fixed;
      left: 100px;
    }

    .menu a {
      display: block;
      color: #2196f3;
    }
  </style>
</head>
<body>
  <h1>{{ .Name }} 数据字典<span style="float:right">{{ .ReleaseTime }}</span></h1>

  <div class="menu">
  {{ range .Tables }}
  <a href="#{{ .Name }}">{{ .Name }} {{ .Comment }}</a>
  {{ end }}
  </div>

  {{ range .Tables }}
  <h2 id="{{ .Name }}">{{ .Name }} {{ .Comment }}</h2>

  <table>
    <tr>
      <th>字段名</th>
      <th>数据类型</th>
      <th>允许空值</th>
      <th>键</th>
      <th>默认值</th>
      <th>备注</th>
    </tr>
    {{ range .Columns }}
    <tr>
      <td>{{ .Field }}</td>
      <td>{{ .Type }}</td>
      <td>{{ .Null }}</td>
      <td>{{ .Key }}</td>
      <td>{{ .Default }}</td>
      <td>{{ .Comment }}</td>
    </tr>
    {{ end }}
  </table>

  {{ end }}
</body>
</html>
`

// 生成数据库字典
func BuildBook() error {
	var tables []TableItem
	var dbName string
	err := Db.Raw("SELECT DATABASE()").Row().Scan(&dbName)
	if err != nil {
		return err
	}

	vingo.Mkdir("dbbook")
	var outputFilePath = filepath.Join(".", "dbbook", fmt.Sprintf("%v_%v.html", dbName, time.Now().Format("20060102")))

	// 查询所有表的信息
	rows, err := Db.Raw(`SELECT TABLE_NAME, TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?`, dbName).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, tableComment string
		if err := rows.Scan(&tableName, &tableComment); err != nil {
			return err
		}

		// 查询每张表的列信息并按字段顺序排序
		columns, err := getTableColumns(dbName, tableName)
		if err != nil {
			return err
		}

		// 获取字段顺序
		fieldOrder := getFieldOrder(dbName, tableName)

		// 根据字段顺序排序
		sortedColumns := make([]Column, len(columns))
		for i, field := range fieldOrder {
			for _, col := range columns {
				if col.Field == field {
					sortedColumns[i] = col
					break
				}
			}
		}

		tables = append(tables, TableItem{
			Name:    tableName,
			Comment: tableComment,
			Columns: sortedColumns,
		})
	}

	// 构造 Database 对象
	database := Database{
		Name:        dbName,
		Tables:      tables,
		ReleaseTime: time.Now().Format("2006年01月02日"),
	}

	// 渲染模板
	t, err := template.New("tpl").Parse(booktpl)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	if err := t.Execute(outputFile, database); err != nil {
		return err
	}

	return nil
}

// 获取字段顺序
func getFieldOrder(dbName string, tableName string) []string {
	var fields []string

	rows, err := Db.Raw(`SELECT COLUMN_NAME FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION`, dbName, tableName).Rows()
	if err != nil {
		return fields
	}
	defer rows.Close()

	for rows.Next() {
		var field string
		if err := rows.Scan(&field); err != nil {
			continue
		}
		fields = append(fields, field)
	}

	return fields
}

func getTableColumns(dbName string, tableName string) ([]Column, error) {
	var columns []Column
	rows, err := Db.Raw(`SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, COLUMN_KEY, COLUMN_DEFAULT, EXTRA, COLUMN_COMMENT FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?`, dbName, tableName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var column Column
		if err := rows.Scan(&column.Field, &column.Type, &column.Null, &column.Key, &column.Default, &column.Extra, &column.Comment); err != nil {
			return nil, err
		}

		columns = append(columns, column)
	}

	return columns, nil
}
