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

// DBConn 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}
