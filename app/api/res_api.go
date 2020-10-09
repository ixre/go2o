package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/ixre/gof"
	"github.com/ixre/gof/api"
	"go2o/core/domain/interface/ad"
	"go2o/core/domain/interface/registry"
	"go2o/core/infrastructure/domain"
	"go2o/core/infrastructure/gen"
	"go2o/core/infrastructure/tool"
	"go2o/core/service"
	"go2o/core/service/impl"
	"go2o/core/service/proto"
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

type resApi struct {
	utils
}

func NewResApi() api.Handler {
	return &resApi{}
}
func (r resApi) Process(fn string, ctx api.Context) *api.Response {
	switch fn {
	case "ad_api":
		return r.adApi(ctx)
	case "geo_location":
		return r.geoLocation(ctx)
	case "area":
		return r.childArea(ctx)
	}
	return api.RUndefinedApi
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
func (r resApi) adApi(ctx api.Context) *api.Response {
	posKeys := ctx.Form().GetString("pos_keys")
	userId := ctx.Form().GetInt("user_id")
	namesParams := strings.TrimSpace(posKeys)
	names := strings.Split(namesParams, "|")
	as := impl.AdService
	result := make(map[string]*ad.AdDto, len(names))
	key := fmt.Sprintf("go2o:repo:ad:%d:front:%s", userId,
		domain.Md5(namesParams))
	rds := gof.CurrentApp.Storage()
	if rds == nil {
		panic("storage need redis support")
	}
	var seconds = 0
	_ = rds.RWJson(key, &result, func() interface{} {
		//从缓存中读取
		for _, n := range names {
			//分别绑定广告
			dto := as.GetAdAndDataByKey(int64(userId), n)
			if dto == nil {
				result[n] = nil
				continue
			}
			result[n] = dto
		}
		regArr := []string{registry.CacheAdMaxAge}
		trans, cli, err := service.RegistryServiceClient()
		if err == nil {
			mp, _ := cli.GetValues(context.TODO(), &proto.StringArray{Value: regArr})
			_ = trans.Close()
			seconds, _ = strconv.Atoi(mp.Value[regArr[0]])
		}
		return result
	}, int64(seconds))
	return r.utils.success(result)
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
func (r resApi) qrCode(ctx api.Context) interface{} {
	text := ctx.Form().GetString("text")
	data := gen.BuildQrCodeForUrl(text, 20)
	return base64.URLEncoding.EncodeToString(data)
}

/**
 * @api {post} /res/geo_location 返回用户的IP及其地址
 * @apiGroup res
 * @apiSuccessExample Success-Response
 * {"ip":"121.83.88.49":"address":"中国上海徐汇区"}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (r resApi) geoLocation(ctx api.Context) *api.Response {
	ip := ctx.Form().GetString("$user_addr")
	if i := strings.Index(ip, ":"); i != -1 {
		ip = ip[:i]
	}
	addr := tool.GetLocation(ip)
	mp := map[string]string{
		"ip":      ip,
		"address": addr,
	}
	return r.utils.success(mp)
}

/**
 * @api {post} /res/area 获取地区(省市区)数据
 * @apiGroup res
 * @apiParam {int} area_code(省市区)地区编码, 如获取全部省, 参数传:0,如果获取四川省下面的市,则传四川省的数字编码
 * @apiParam {int} area_type 地区类型,1:省, 2:市, 3:区
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (r resApi) childArea(ctx api.Context) *api.Response {
	code, _ := strconv.Atoi(ctx.Form().GetString("area_code"))
	areaType := ctx.Form().GetInt("area_type")
	tran, cli, err := service.FoundationServiceClient()
	var areas *proto.AreaListResponse
	if err == nil {
		areas, _ = cli.GetChildAreas(context.TODO(),
			&proto.Int32{Value: int32(code)})
		_ = tran.Close()
		if areaType == 3 {
			for i, v := range areas.Value {
				if strings.TrimSpace(v.Name) == "市辖区" {
					areas.Value = append(areas.Value[:i], areas.Value[i+1:]...)
				}
			}
		}
	}
	return r.success(areas)
}
