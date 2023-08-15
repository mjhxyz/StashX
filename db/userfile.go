package db

import (
	"database/sql"
	mydb "stashx/db/mysql"
)

// UserFile: 用户文件表结构体
type TableUserFile struct {
	UserName    string
	FileHash    string
	FileName    sql.NullString
	FileSize    sql.NullInt64
	UploadAt    sql.NullTime
	LastUpdated sql.NullTime
}

// OnUserFileUploadFinished: 更新用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into stashx_user_file (`user_name`,`file_sha1`,`file_name`,`file_size`) values (?,?,?,?)",
	)
	if err != nil {
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, filehash, filename, filesize)
	if err != nil {
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			return false
		}
		return true
	}
	return false
}
