package db

import (
	"fmt"
	mydb "stashx/db/mysql"
)

// UserSignin : 判断密码是否一致
func UserSignin(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"select * from stashx_user where user_name=? limit 1",
	)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found: ", username)
		return false
	}

	pRows := mydb.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		fmt.Println("success")
		return true
	} else {
		fmt.Println("failed")
		return false
	}
}

// UserSignUp : 通过用户名及密码完成user表的注册操作
func UserSignUp(username, password string) bool {
	sqlStr := "insert ignore into stashx_user (user_name, user_pwd) values (?, ?)"
	stmt, err := mydb.DBConn().Prepare(
		sqlStr,
	)
	if err != nil {
		fmt.Println("insert user error: ", err)
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

// UpdateToken : 刷新用户登录的token
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into stashx_user_token (`user_name`, `user_token`) values (?, ?)",
	)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
