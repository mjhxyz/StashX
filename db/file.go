package db

import (
	"database/sql"
	"fmt"
	mydb "stashx/db/mysql"
)

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

func RemoveFileMeta(filehash string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"delete from stashx_file where file_sha1=?",
	)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(filehash)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("No file with hash:%s\n", filehash)
		}
		return true
	}
	return false
}

// GetFileMetaList : 从mysql获取批量的文件元信息
func GetFileMetaList(limit int) ([]TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_name,file_size,file_addr from stashx_file limit ?",
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(limit)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer rows.Close()
	var tfiles []TableFile
	for rows.Next() {
		tfile := TableFile{}
		err = rows.Scan(
			&tfile.FileHash, &tfile.FileName, &tfile.FileSize, &tfile.FileAddr)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		tfiles = append(tfiles, tfile)
	}
	return tfiles, nil
}

// GetFileMeta : 从mysql获取文件元信息
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1,file_name,file_size," +
			"file_addr from stashx_file where file_sha1=? and status=1 limit 1",
	)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()
	tfile := TableFile{}
	// 按照 select 顺序，将查询结果依次赋值给 tfile
	err = stmt.QueryRow(filehash).Scan(
		&tfile.FileHash, &tfile.FileName, &tfile.FileSize, &tfile.FileAddr)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &tfile, nil
}

// OnFileUploadFinished : 文件上传完成，保存meta
func OnFileUploadFinished(
	filehash string, filename string,
	filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore stashx_file(file_sha1,file_name,file_size,file_addr,status)" +
			" values(?, ?, ?, ?, 1)",
	)
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Printf("Failed to insert, err:" + err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}
