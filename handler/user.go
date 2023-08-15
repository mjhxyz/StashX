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
	// writer.Write([]byte("http://" + request.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + request.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	writer.Write(resp.JSONBytes())
}

func UserInfoHandler(writer http.ResponseWriter, request *http.Request) {
	// 1. 解析请求参数
	request.ParseForm()
	username := request.Form.Get("username")
	token := request.Form.Get("token")

	// 2. 验证token是否有效
	isValidToken := IsTokenValid(token)
	if !isValidToken {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	// 3. 查询用户信息
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		writer.WriteHeader(http.StatusForbidden)
		return
	}

	// 4. 组装并且响应用户数据
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	writer.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	// 40位字符 md5(uesrname+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// IsTokenValid : token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// 1. token是否过期
	// 2. token是否存在于数据库
	// 3. token是否与数据库中一致
	return true
}
