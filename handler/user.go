package handler

import (
	"net/http"
)

func SignupHandler(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// db.UserSignUp(writer, request)
}
