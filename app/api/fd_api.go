package api

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/ixre/go2o/app/api/util"
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/gen"
	"github.com/ixre/go2o/core/infrastructure/tool"
	"github.com/ixre/go2o/core/service"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof"
	"github.com/ixre/gof/jwt-api"
	"strconv"
	"strings"
)

/**
 * Copyright 2009-2019 @ 56x.net
 * name : res_api.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-11-26 17:07
 * description :
 * history :
 */

var _ api.Handler = new(fdApi)

// 基础接口
type fdApi struct {
	util.Utils
}

func (r fdApi) Group() string {
	return "fd"
}

func (r fdApi) Process(fn string, ctx api.Context) *api.Response {
	switch fn {
	case "ad_api":
		return r.adApi(ctx)
	case "geo_location":
		return r.geoLocation(ctx)
	case "area":
		return r.childArea(ctx)
	case "check_sensitive":
		return r.checkSensitive(ctx)
	case "replace_sensitive":
		return r.replaceSensitive(ctx)
	}
	return nil
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
func (r fdApi) adApi(ctx api.Context) *api.Response {
	posKeys := ctx.Request().Params.GetString("pos_keys")
	userId := ctx.Request().Params.GetInt("user_id")
	namesParams := strings.TrimSpace(posKeys)
	names := strings.Split(namesParams, "|")
	result := make(map[string]*proto.SAdvertisementDto, len(names))
	key := fmt.Sprintf("go2o:repo:ad:%d:front:%s", userId,
		domain.Md5(namesParams))
	rds := gof.CurrentApp.Storage()
	if rds == nil {
		panic("storage need redis support")
	}
	var seconds = 0
	_ = rds.RWJson(key, &result, func() interface{} {
		trans, cli, _ := service.AdvertisementServiceClient()
		defer trans.Close()
		//从缓存中读取
		for _, n := range names {
			//分别绑定广告
			dto, _ := cli.GetAdvertisement(context.TODO(),
				&proto.AdIdRequest{
					AdUserId:   int64(userId),
					AdKey:      n,
					ReturnData: true,
				})
			if dto == nil {
				result[n] = nil
				continue
			}
			result[n] = dto.Data
		}
		regArr := []string{registry.CacheAdMaxAge}
		trans2, cli2, _ := service.RegistryServiceClient()
		mp, _ := cli2.GetValues(context.TODO(), &proto.StringArray{Value: regArr})
		_ = trans2.Close()
		seconds, _ = strconv.Atoi(mp.Value[regArr[0]])
		return result
	}, int64(seconds))
	return r.Utils.Success(result)
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
func (r fdApi) qrCode(ctx api.Context) interface{} {
	text := ctx.Request().Params.GetString("text")
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
func (r fdApi) geoLocation(ctx api.Context) *api.Response {
	ip := ctx.Request().Params.GetString("$user_addr")
	if i := strings.Index(ip, ":"); i != -1 {
		ip = ip[:i]
	}
	addr := tool.GetLocation(ip)
	mp := map[string]string{
		"ip":      ip,
		"address": addr,
	}
	return r.Utils.Success(mp)
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
func (r fdApi) childArea(ctx api.Context) *api.Response {
	code, _ := strconv.Atoi(ctx.Request().Params.GetString("area_code"))
	areaType := ctx.Request().Params.GetInt("area_type")
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
	return r.Success(areas)
}

// 检查是否包含敏感词
func (r fdApi) checkSensitive(ctx api.Context) *api.Response {
	text := ctx.Request().Params.GetString("text")
	trans, cli, err := service.FoundationServiceClient()
	if err == nil {
		defer trans.Close()
		ret, _ := cli.CheckSensitive(context.TODO(), &proto.String{Value: text})
		return r.Success(ret.Value)
	}
	return r.Error(err)
}

func (r fdApi) replaceSensitive(ctx api.Context) *api.Response {
	mp := map[string]string{}
	ctx.Request().Bind(&mp)
	text := mp["text"]
	replacement := mp["replacement"]
	trans, cli, err := service.FoundationServiceClient()
	if err == nil {
		defer trans.Close()
		ret, _ := cli.ReplaceSensitive(context.TODO(), &proto.ReplaceSensitiveRequest{
			Text:        text,
			Replacement: replacement,
		})
		return r.Success(ret.Value)
	}
	return r.Error(err)
}
