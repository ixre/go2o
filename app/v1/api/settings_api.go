package api

import (
	"context"
	"github.com/ixre/go2o/core/service"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/api"
)

var _ api.Handler = new(SettingsApi)

type SettingsApi struct {
}

func NewSettingsApi() *SettingsApi {
	return &SettingsApi{}
}

func (a SettingsApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"register_settings": a.registerSettings,
	})
}

/**
 * @api {post} /settings/register_settings 获取注册Token
 * @apiName register_settings
 * @apiGroup settings
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"code":1,"message":"api not defined"}
 */
func (a SettingsApi) registerSettings(ctx api.Context) interface{} {
	trans, cli, err := service.RegistryServiceClient()
	if err == nil {
		defer trans.Close()
		mp, _ := cli.FindRegistries(context.TODO(), &proto.String{Value: "member_register"})
		return mp.Value
	}
	return api.ResponseWithCode(1, "no register settings")
}
