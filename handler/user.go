package handler

import (
	"fmt"
	"net/http"
	dblayer "stashx/db"
	"stashx/util"
	"time"
)

const (
	pwdSalt = "*#890"
)

// SignupHandler : 处理用户注册请求
func SignupHandler(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	if method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	request.ParseForm()
	username := request.Form.Get("username")
	password := request.Form.Get("password")
	fmt.Println(username, password)

	if len(username) < 3 || len(password) < 5 {
		writer.Write([]byte("Invalid parameter"))
		return
	}

	encPassword := util.Sha1([]byte(password + pwdSalt))
	suc := dblayer.UserSignUp(username, encPassword)
	if suc {
		writer.Write([]byte("SUCCESS"))
	} else {
		writer.Write([]byte("FAILED"))
	}
}

// SigninHandler : 登录接口
func SigninHandler(writer http.ResponseWriter, request *http.Request) {
	// 1. 校验用户名和密码
	request.ParseForm()
	username := request.Form.Get("username")
	password := request.Form.Get("password")
	encpwd := util.Sha1([]byte(password + pwdSalt))
	pwdChecked := dblayer.UserSignin(username, encpwd)
	if !pwdChecked {
		writer.Write([]byte("FAILED"))
		return
	}
	// 2. 生成访问凭证(token)
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		writer.Write([]byte("FAILED"))
		return
	}
	// 3. 登录成功后重定向到首页
	writer.Write([]byte("http://" + request.Host + "/static/view/home.html"))
}

func GenToken(username string) string {
	// 40位字符 md5(uesrname+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}
