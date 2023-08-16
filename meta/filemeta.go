package meta

import "stashx/db"

// FileMeta : 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// GetFileMetaListDB : 从mysql获取批量的文件元信息
func GetFileMetaListDB(limit int) ([]FileMeta, error) {
	tfiles, err := db.GetFileMetaList(limit)
	if err != nil {
		return nil, err
	}
	fmeta := make([]FileMeta, 0)
	for _, tfile := range tfiles {
		fmeta = append(fmeta, FileMeta{
			FileSha1: tfile.FileHash,
			FileName: tfile.FileName.String,
			FileSize: tfile.FileSize.Int64,
			Location: tfile.FileAddr.String,
		})
	}
	return fmeta, nil
}

// GetFileMetaDB : 从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (*FileMeta, error) {
	tfile, err := db.GetFileMeta(fileSha1)
	if err != nil {
		return nil, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return &fmeta, nil
}

// UpdateFileMetaDB : 新增/更新文件元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return db.OnFileUploadFinished(
		fmeta.FileSha1, fmeta.FileName,
		fmeta.FileSize, fmeta.Location,
	)
}

// UpdateFileMeta : 新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// GetFileMetaList : 获取批量文件元信息列表
func GetFileMetaList() []FileMeta {
	fileMetaList := make([]FileMeta, 0)
	for _, v := range fileMetas {
		fileMetaList = append(fileMetaList, v)
	}
	return fileMetaList
}

// GetFileMeta : 通过sha1值获取文件元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}

func RemoveFileMetaDB(fileSha1 string) bool {
	return db.RemoveFileMeta(fileSha1)
}
