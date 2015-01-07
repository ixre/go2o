package session

import (
	"com/ording"
	"com/ording/dao"
	"com/ording/entity"
	"com/share/glob"
	"net/http"
	"github.com/newmin/gof/app"
	"github.com/newmin/gof/web"
	"strconv"
	"strings"
	"time"
)

var (
	LSession *loginSession = &loginSession{glob.CurrContext()}
)

//商户业务逻辑
type loginSession struct {
	app.Context
}

func (this *loginSession) AdministratorLogin(w http.ResponseWriter, usr string, pwd string) bool {
	loginTokenResult := ording.EncodePartnerPwd(usr, pwd)
	loginResult := loginTokenResult == glob.CurrContext().Config().GetString("master_token")

	if loginResult {

		//存入cookie
		expires := time.Now()
		expires = expires.Add(3600 * 72 * 1e9) //72H
		cookie := http.Cookie{Name: "mtkey",
			Value:   loginTokenResult,
			Expires: expires}
		http.SetCookie(w, &cookie)
		web.Seria2json(w, true, "", nil)
		return true
	}
	return false
}

func (this *loginSession) IsAdministrator(r *http.Request) bool {
	cookie, _ := r.Cookie("mtkey")
	if cookie == nil || cookie.Value != glob.CurrContext().Config().GetString("master_token") {
		return false
	}
	return true
}

//验证登陆
func (pb *loginSession) ValidLogin(usr string, pwd string) (*entity.Partner, bool, string) {
	var message string
	var result bool
	var pt *entity.Partner

	id := dao.Partner().Verify(usr, pwd)

	if id == -1 {
		result = false
		message = "用户或密码不正确！"
	} else {
		pt = dao.Partner().GetPartnerById(id)
		if pt.Expires.Before(time.Now()) {
			message = "您的账号已经过期"
			result = false
		} else {
			result = true
		}
	}
	return pt, result, message
}

//验证登陆
func (pb *loginSession) WebValidLogin(w http.ResponseWriter, usr string, pwd string) {

	pt, result, message := pb.ValidLogin(usr, pwd)

	if result {
		//存入cookie
		tokenStr := strconv.Itoa(pt.Id) + "," + pt.Pwd
		expires := time.Now()
		expires = expires.Add(3600 * 72 * 1e9) //72H
		cookie := http.Cookie{Name: "pntoken", Value: tokenStr, Expires: expires}
		http.SetCookie(w, &cookie)
		//fmt.Println(tokenStr + expires.String())
	}
	web.Seria2json(w, result, message, nil)
}

//从cookie中获取当前会话信息
func (pb *loginSession) GetCurrentSessionFromCookie(r *http.Request) (*entity.Partner, error) {
	cookie, cookieErr := r.Cookie("pntoken")
	if cookie == nil {
		return nil, cookieErr
	}
	partnerId, parseErr := strconv.Atoi(strings.Split(cookie.Value, ",")[0])
	if parseErr != nil {
		return nil, parseErr
	}

	return dao.Partner().GetPartnerById(partnerId), nil
}

//获取合作商编号
func (pb *loginSession) GetPartnerIdFromCookie(r *http.Request) (int, error) {
	cookie, cookieErr := r.Cookie("pntoken")
	if cookie == nil {
		return -1, cookieErr
	}

	partnerId, parseErr := strconv.Atoi(strings.Split(cookie.Value, ",")[0])
	if parseErr != nil {
		return -1, parseErr
	}
	return partnerId, nil
}

func (pb *loginSession) PartnerLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("pntoken")
	if err == nil {
		cookie.Expires = time.Now().Add(-5 * time.Second)
		//cookie.MaxAge = 0
		cookie.Path = "/"
		http.SetCookie(w, cookie)
	}
}
