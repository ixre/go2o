package api

import (
	"encoding/base64"
	"fmt"
	"github.com/ixre/gof"
	"github.com/ixre/gof/api"
	"go2o/core/domain/interface/ad"
	"go2o/core/domain/interface/registry"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/gen"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift"
	"strconv"
	"strings"
)

/**
 * Copyright 2009-2019 @ to2.net
 * name : res_api.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-11-26 17:07
 * description :
 * history :
 */

var _ api.Handler = new(resApi)
type resApi struct{
	utils
}

func NewResApi() api.Handler {
	return &resApi{}
}
func (a resApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"ad_api": a.adApi,
	})
}


/**
 * @api {post} /res/ad_api 获取广告数据
 * @apiGroup res
 * @apiParam {string} pos_keys 广告位KEYS, 多个广告位用"|"拼接
 * @apiParam {int} user_id 广告用户编号, 默认: 0
 * @apiSuccess {json} data 返回广告数据,根据广告位的type进行区分,不同的广告,广告的数据格式不一样.
 * @apiSuccessExample Success-Response
 * {"mobi-index-scroller":{"type":3,"data":[{"AdId":40,"Title":"1","LinkUrl":"","ImageUrl":"https://img30.360buyimg.com/da/jfs/t1/87253/21/3038/151384/5ddb9821Efa7cc1bf/39f7d4ba7fb76f77.jpg!q90.webp","SortNum":0,"Enabled":1},
 * {"AdId":40,"Title":"2","LinkUrl":"http://","ImageUrl":"https://img30.360buyimg.com/da/jfs/t1/48966/19/16642/115982/5ddc9f0dE82a8c4ed/9d144383a1790b3b.jpg!q90.webp","SortNum":0,"Enabled":1},
 * {"AdId":40,"Title":"3","LinkUrl":"http://","ImageUrl":"https://img10.360buyimg.com/da/jfs/t1/45172/26/15880/100131/5dce6481E630eb605/2c36d6bc89d7dd29.jpg","SortNum":0,"Enabled":1}]}}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (a resApi) adApi(ctx api.Context) interface{} {
	posKeys :=ctx.Form().GetString("pos_keys")
	userId := ctx.Form().GetInt("user_id")
	namesParams := strings.TrimSpace(posKeys)
	names := strings.Split(namesParams, "|")
	as := rsi.AdService
	result := make(map[string]*ad.AdDto, len(names))
	key := fmt.Sprintf("go2o:repo:ad:%d:front:%s", userId,
		domain.Md5(namesParams))
	rds := gof.CurrentApp.Storage()
	if rds == nil {
		panic("storage need redis support")
	}
	var seconds int = 0
	rds.RWJson(key, &result, func() interface{} {
		//从缓存中读取
		for _, n := range names {
			//分别绑定广告
			dto := as.GetAdAndDataByKey(int32(userId), n)
			if dto == nil {
				result[n] = nil
				continue
			}
			result[n] = dto
		}
		regArr := []string{registry.CacheAdMaxAge}
		trans, cli, err := thrift.RegistryServeClient()
		if err == nil {
			defer trans.Close()
			mp, _ := cli.GetRegistries(thrift.Context, regArr)
			seconds, _ = strconv.Atoi(mp[regArr[0]])
		}
		return result
	}, int64(seconds))
	return result
}



/**
 * @api {post} /res/qr_code 生成二维码
 * @apiGroup res
 * @apiParam {string} text 文本内容
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (a resApi) qrCode(ctx api.Context) interface{} {
	text := ctx.Form().GetString("text")
	data := gen.BuildQrCodeForUrl(text, 20)
	return base64.URLEncoding.EncodeToString(data)
}

