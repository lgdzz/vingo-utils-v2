package pgsql

import (
	"database/sql"
	"fmt"
	"github.com/lgdzz/vingo-utils-v2/vingo"
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
      margin: 0 50px;
      font-size: 14px;
      padding-bottom: 50px;
    }

    table {
      border-collapse: collapse;
      width: 100%;
    }

    th,
    td {
      border: 1px solid #ddd;
      padding: 8px;
      text-align: left;
    }

    th {
      background-color: #f2f2f2;
    }

    .main {
      display: flex;
      height: 85vh;
    }

    .menu {
      margin-right: 50px;
      height: 100%;
      overflow: auto;
    }

    .menu a {
      display: flex;
      color: #2196f3;
      font-size: 12px;
      text-decoration: inherit;
    }

    .menu a div:nth-child(1) {
      flex: 1;
    }

    .menu a div:nth-child(2) {
      color: #ccc;
      margin: 0 10px;
    }

    .table {
      flex: 1;
      height: 100%;
      overflow: auto;
    }
  </style>
</head>
<body>
  <h1>{{ .Name }} 数据字典<span style="float:right">{{ .ReleaseTime }}</span></h1>
  <div class="main">
	  <div class="menu">
	  {{ range .Tables }}
	  <a href="#{{ .Name }}">
      	<div>{{ .Name }}</div>
        <div>{{ .Comment }}</div>
      </a>
	  {{ end }}
	  </div>

  	  <div class="table">
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
	  </div>
  </div>
</body>
</html>
`

func (s *DbApi) BuildBook() error {
	var tables []TableItem
	var dbName string

	// PostgreSQL 获取当前数据库名
	err := s.DB.Raw("SELECT current_database()").Row().Scan(&dbName)
	if err != nil {
		return err
	}

	vingo.Mkdir("dbbook")
	var outputFilePath = filepath.Join("dbbook", fmt.Sprintf("%v_%v.html", dbName, time.Now().Format("20060102")))

	// PostgreSQL 获取所有表名及注释
	rows, err := s.DB.Raw(`
		SELECT c.relname AS table_name, obj_description(c.oid) AS comment
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE c.relkind = 'r' AND n.nspname = 'public'
	`).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		var comment sql.NullString

		if err := rows.Scan(&tableName, &comment); err != nil {
			return err
		}

		tableComment := ""
		if comment.Valid {
			tableComment = comment.String
		}

		// 获取列信息
		columns, err := s.getTableColumns(tableName)
		if err != nil {
			return err
		}

		// 字段排序
		fieldOrder := s.getFieldOrder(tableName)
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

	database := Database{
		Name:        dbName,
		Tables:      tables,
		ReleaseTime: time.Now().Format("2006年01月02日"),
	}

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
func (s *DbApi) getFieldOrder(tableName string) []string {
	var fields []string

	rows, err := s.DB.Raw(`
		SELECT column_name
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = ?
		ORDER BY ordinal_position
	`, tableName).Rows()
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

// 获取表字段信息（PostgreSQL）
func (s *DbApi) getTableColumns(tableName string) ([]Column, error) {
	var columns []Column

	rows, err := s.DB.Raw(`
		SELECT
			a.attname AS column_name,
			format_type(a.atttypid, a.atttypmod) AS data_type,
			NOT a.attnotnull AS is_nullable,
			CASE WHEN ct.contype = 'p' THEN 'PRI' ELSE '' END AS column_key,
			COALESCE(pg_get_expr(d.adbin, d.adrelid), '') AS column_default,
			col_description(a.attrelid, a.attnum) AS column_comment
		FROM
			pg_attribute a
		JOIN pg_class c ON a.attrelid = c.oid
		JOIN pg_namespace n ON c.relnamespace = n.oid
		LEFT JOIN pg_attrdef d ON d.adrelid = c.oid AND d.adnum = a.attnum
		LEFT JOIN pg_constraint ct ON ct.conrelid = c.oid AND a.attnum = ANY(ct.conkey)
		WHERE
			a.attnum > 0
			AND NOT a.attisdropped
			AND c.relname = ?
			AND n.nspname = 'public'
		ORDER BY a.attnum
	`, tableName).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var col Column
		var nullable bool
		var comment sql.NullString

		if err := rows.Scan(&col.Field, &col.Type, &nullable, &col.Key, &col.Default, &comment); err != nil {
			return nil, err
		}

		col.Null = map[bool]string{true: "YES", false: "NO"}[nullable]
		col.Extra = "" // PostgreSQL 没有 extra
		col.Comment = ""
		if comment.Valid {
			col.Comment = comment.String
		}

		columns = append(columns, col)
	}

	return columns, nil
}
