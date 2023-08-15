package main

import (
	"fmt"
	"net/http"
	"stashx/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(http.Dir("")))

	mux.HandleFunc("/file/upload", handler.HandleUpload)
	mux.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	mux.HandleFunc("/file/download", handler.DownloadHandler)
	mux.HandleFunc("/file/update", handler.MetaUpdateHandler)
	mux.HandleFunc("/file/delete", handler.FileDeleteHandler)
	mux.HandleFunc("/file/list", handler.ListFileMetaHandler)
	mux.HandleFunc("/file/query", handler.FileQueryHandler)

	mux.HandleFunc("/user/signup", handler.SignupHandler)
	mux.HandleFunc("/user/signin", handler.SigninHandler)
	mux.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	fmt.Println("server start")
	server := &http.Server{
		Addr:    ":8888",
		Handler: mux,
	}
	server.ListenAndServe()
}
