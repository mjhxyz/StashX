package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:000000@tcp(192.168.60.100:3306)/stashx?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Printf("open mysql failed, err:%v\n", err)
		panic(err)
	}
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var results []map[string]interface{}
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			panic(err)
		}
		row := make(map[string]interface{})
		for i, col := range values {
			if col != nil {
				row[columns[i]] = col
			}
		}
		results = append(results, row)
	}
	return results
}

// DBConn 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}
