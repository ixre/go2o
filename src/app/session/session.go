/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package session

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/infrastructure/domain"
	"go2o/src/core/service/dps"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var _lsession *loginSession

func GetLSession() *loginSession {
	if _lsession == nil {
		_lsession = &loginSession{core.GlobalApp}
	}
	return _lsession
}

//商户业务逻辑
type loginSession struct {
	gof.App
}

func (this *loginSession) AdministratorLogin(w http.ResponseWriter, usr string, pwd string) bool {
	loginTokenResult := domain.EncodePartnerPwd(usr, pwd)
	loginResult := loginTokenResult == this.App.Config().GetString("master_token")

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
	if cookie == nil || cookie.Value != this.App.Config().GetString("master_token") {
		return false
	}
	return true
}

//验证登陆
func (pb *loginSession) ValidLogin(usr string, pwd string) (*partner.ValuePartner, bool, string) {
	var message string
	var result bool
	var pt *partner.ValuePartner
	var err error

	id := dps.PartnerService.Verify(usr, pwd)

	if id == -1 {
		result = false
		message = "用户或密码不正确！"
	} else {
		pt, err = dps.PartnerService.GetPartner(id)
		if err != nil {
			message = err.Error()
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
func (pb *loginSession) GetCurrentSessionFromCookie(r *http.Request) (*partner.ValuePartner, error) {
	cookie, cookieErr := r.Cookie("pntoken")
	if cookie == nil {
		return nil, cookieErr
	}
	partnerId, parseErr := strconv.Atoi(strings.Split(cookie.Value, ",")[0])
	if parseErr != nil {
		return nil, parseErr
	}

	return dps.PartnerService.GetPartner(partnerId)
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

func (pb *loginSession) PartnerLogout(ctx *web.Context) {
	r, w := ctx.Request, ctx.ResponseWriter
	cookie, err := r.Cookie("pntoken")
	if err == nil {
		cookie.Expires = time.Now().Add(-5 * time.Second)
		//cookie.MaxAge = 0
		cookie.Path = "/"
		http.SetCookie(w, cookie)
	}
}
