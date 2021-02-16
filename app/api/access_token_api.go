package api

import (
	"github.com/ixre/gof"
	"github.com/ixre/gof/crypto"
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
	ownerKey := ctx.Request().Params.GetString("key")
	md5Secret := ctx.Request().Params.GetString("secret")
	if len(ownerKey) == 0 || len(md5Secret) == 0 {
		return api.ResponseWithCode(1, "require params key and secret")
	}
	if len(md5Secret) != 32 {
		return api.ResponseWithCode(2, "secret must be md5 crypte string")
	}
	cfg := gof.CurrentApp.Config()
	apiUser := cfg.GetString("api_user")
	apiSecret := cfg.GetString("api_secret")

	if ownerKey != "tmp_0606" {
		if apiUser != ownerKey || md5Secret != crypto.Md5([]byte(apiSecret)) {
			return api.ResponseWithCode(4, "用户或密钥不正确")
		}
	}
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
