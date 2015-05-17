/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package session

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/web"
	"go2o/src/core"
	"go2o/src/core/infrastructure/domain"
	"net/http"
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
	loginTokenResult := domain.Md5PartnerPwd(usr, pwd)
	loginResult := loginTokenResult == ctx.App.Config().GetString("master_token")

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
	if cookie == nil || cookie.Value != ctx.App.Config().GetString("master_token") {
		return false
	}
	return true
}
