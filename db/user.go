package db

import (
	"fmt"
	mydb "stashx/db/mysql"
)

// UserSignUp : 通过用户名及密码完成user表的注册操作
func UserSignUp(username, password string) bool {
	sqlStr := "insert ignore into user (username, password) values (?, ?)"
	stmt, err := mydb.DBConn().Prepare(
		sqlStr,
	)
	if err != nil {
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("insert user error: ", err)
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}
