package handler

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path"
	"regexp"
	redisLayer "stashx/cache/redis"
	dblayer "stashx/db"
	"stashx/util"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

// 做成 bitmap

// MultipartUploadInfo : 初始化信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

func InitialMultipartUploadHandler(writer http.ResponseWriter, request *http.Request) {
	// 1. 解析用户请求参数
	request.ParseForm()
	username := request.Form.Get("username")
	filehash := request.Form.Get("filehash")
	filesize, err := strconv.Atoi(request.Form.Get("filesize"))
	if err != nil {
		writer.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}

	// 2. 获取 redis 连接池中的一个连接
	rConn := redisLayer.RedisPool().Get()
	defer rConn.Close()

	// 3. 生成分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   filesize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  5 * 1024 * 1024, // 5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	// 4. 将初始化信息写入 redis 缓存中
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)

	// 5. 将响应初始化数据返回到客户端
	writer.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())
}

// UploadPartHandler : 上传文件分块
func UploadPartHandler(writer http.ResponseWriter, request *http.Request) {
	// 1. 解析用户请求参数
	request.ParseForm()
	_ = request.Form.Get("username")
	uploadID := request.Form.Get("uploadid")
	chunkIndex := request.Form.Get("index")

	// 2. 获得 redis 连接池中的一个连接
	rConn := redisLayer.RedisPool().Get()
	defer rConn.Close()

	// 3. 获得文件句柄，用于存储分块内容
	fpath := "/temp/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create("/temp/" + uploadID + "/" + chunkIndex)
	if err != nil {
		writer.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()
	for {
		buf := make([]byte, 1024*1024)
		n, err := request.Body.Read(buf)
		if err != nil {
			break
		}
		fd.Write(buf[:n]) // 写入文件内容
	}

	// 4. 更新 redis 缓存状态, 添加一条分块信息
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)

	// 5. 返回处理结果到客户端
	writer.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// CompleteUploadHandler : 通知上传合并
func CompleteUploadHandler(writer http.ResponseWriter, request *http.Request) {
	// 1. 解析请求参数
	request.ParseForm()
	uploadID := request.Form.Get("uploadid")
	username := request.Form.Get("username")
	filehash := request.Form.Get("filehash")
	filesize := request.Form.Get("filesize")
	filename := request.Form.Get("filename")

	// 2. 获得 redis 连接池中的一个连接
	rConn := redisLayer.RedisPool().Get()
	defer rConn.Close()

	// 3. 通过 uploadid 查询 redis 并判断是否所有分块上传完成
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+uploadID))
	if err != nil {
		writer.Write(util.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
		return
	}
	total := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		// 通过 HGETALL 返回的结果，data[i] 为 key，data[i+1] 为 value
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			total, _ = strconv.Atoi(v)
		}
		if match, _ := regexp.MatchString(`^chkidx_\d+$`, k); match {
			chunkCount++
		}
	}

	if total != chunkCount {
		writer.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
		return
	}
	// 4. TODO 合并分块
	// 5. 更新唯一文件表及用户文件表
	fsize, _ := strconv.Atoi(filesize)
	dblayer.OnFileUploadFinished(filehash, filename, int64(fsize), "")
	dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))

	// 6. 响应处理结果
	writer.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// CancelUploadPartHandler : 取消分块上传
func CancelUploadPartHandler(writer http.ResponseWriter, request *http.Request) {
	// 1. 取消已存在的分块文件
	// 2. 删除redis缓存分块状态
	// 3. 更新mysql文件状态
}

// MultipartUploadStatusHandler : 查询文件分块上传初始化信息
func MultipartUploadStatusHandler(writer http.ResponseWriter, request *http.Request) {
	// 1. 检查分块上传状态是否有效
	// 2. 获取分块初始化信息
	// 3. 获取已上传的分块信息
}
