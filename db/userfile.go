package db

import (
	"log"
	mydb "stashx/db/mysql"
	"time"
)

// UserFile: 用户文件表结构体
type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    *time.Time
	LastUpdated *time.Time
}

// OnUserFileUploadFinished: 更新用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into stashx_user_file (`user_name`,`file_sha1`,`file_name`,`file_size`) values (?,?,?,?)",
	)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, filehash, filename, filesize)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		log.Println("没有影响的数量")
		if rf <= 0 {
			return false
		}
		return true
	}
	return false
}

// QueryUserFileMetas: 批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_name,file_size,upload_at,last_update from stashx_user_file where user_name=? limit ?",
	)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(
			&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			log.Println(err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}
	log.Println("userFiles: ", userFiles)
	return userFiles, nil
}
