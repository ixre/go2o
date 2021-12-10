package api

import (
	api "github.com/ixre/gof/jwt-api"
	"time"
)

var _ api.Handler = new(AccessTokenApi)

type AccessTokenApi struct {
}

func (a AccessTokenApi) Group() string {
	return "access_token"
}

func (a AccessTokenApi) Process(fn string, ctx api.Context) *api.Response {
	return a.createAccessToken(ctx)
}

func (a AccessTokenApi) createAccessToken(ctx api.Context) *api.Response {
	// 创建token并返回
	claims := api.CreateClaims("0", "go2o",
		"go2o-api-jwt", time.Now().Unix()+7200).(api.MapClaims)
	claims["global"] = true
	token, err := api.AccessToken(claims, getJWTSecret())
	if err != nil {
		return api.ResponseWithCode(4, err.Error())
	}
	return api.NewResponse(token)
}
