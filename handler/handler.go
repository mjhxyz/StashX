package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"stashx/meta"
	"stashx/util"
	"time"
)

// FileDeleteHandler : 删除文件
func FileDeleteHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fileSha1 := request.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)
	writer.WriteHeader(http.StatusOK)
}

// FileUploadHandler : 处理文件上传
func MetaUpdateHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	opType := request.Form.Get("op")
	fileSha1 := request.Form.Get("filehash")
	newFileName := request.Form.Get("filename")

	if opType != "0" {
		return
	}

	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		fmt.Printf("序列化失败:%v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}

// DownloadHandler : 下载文件
func DownloadHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fsha1 := request.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	file, err := os.Open(fm.Location)
	if err != nil {
		fmt.Printf("打开文件失败:%v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 只适合小文件
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("读取文件失败:%v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/octect-stream")
	writer.Header().Set("Content-Disposition", "attachment;filename=\""+fm.FileName+"\"")
	writer.Write(data)
}

// ListFileMetaHandler : 获取文件元信息列表
func ListFileMetaHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fileMetas := meta.GetFileMetaList()
	data, err := json.Marshal(fileMetas)
	if err != nil {
		fmt.Printf("序列化失败:%v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Write(data)
}

// GetFileMetaHandler : 获取文件元信息
func GetFileMetaHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	filehash := request.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)

	data, err := json.Marshal(fMeta)
	if err != nil {
		fmt.Printf("序列化失败:%v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Write(data)
}

func HandleUpload(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	file, header, err := request.FormFile("file")
	if err != nil {
		fmt.Printf("获取文件失败:%v\n", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileMeta := meta.FileMeta{
		FileName: header.Filename,
		Location: "/tmp/" + header.Filename,
		UploadAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	newFile, err := os.Create(fileMeta.Location)
	if err != nil {
		fmt.Printf("创建文件失败:%v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	fileMeta.FileSize, err = io.Copy(newFile, file)
	if err != nil {
		fmt.Printf("复制文件失败:%v\n", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 计算 hash
	newFile.Seek(0, 0)
	fileMeta.FileSha1 = util.FileSha1(newFile)
	fmt.Printf("文件hash:%s\n", fileMeta.FileSha1)

	meta.UpdateFileMeta(fileMeta)

	writer.Write([]byte("上传成功!"))
}
