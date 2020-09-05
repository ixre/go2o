package api

import (
	"encoding/json"
	"fmt"
	"github.com/ixre/gof"
	"github.com/ixre/gof/api"
	"go2o/core/service/thrift"
	"go2o/core/service/thrift/auto_gen/rpc/ttype"
	"go2o/core/service/thrift/rsi"
)

/**
 * Copyright 2009-2019 @ to2.net
 * name : shop.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-12-09 11:27
 * description :
 * history :
 */

var _ api.Handler = new(shopApi)

type shopApi struct {
	utils
}

func NewShopApi() *shopApi {
	return &shopApi{}
}

func (s shopApi) Process(fn string, ctx api.Context) *api.Response {
	return api.HandleMultiFunc(fn, ctx, map[string]api.HandlerFunc{
		"category":  s.shopCat,
		"favorite":  s.Favorite,
		"my_addrss": s.addressList,
	})
}

/**
 * @api {post} /shop/category 获取商城的商品分类
 * @apiGroup shop
 * @apiParam {int} parent_id 上级分类编号
 * @apiParam {int} shop_id 店铺编号
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (s shopApi) shopCat(ctx api.Context) interface{} {
	parentId := ctx.Form().GetInt("parent_id")
	shopId := ctx.Form().GetInt("shop_id")
	var list []*ttype.SCategory
	key := fmt.Sprintf("go2o:repo:cat:%d:json:%d", shopId, parentId)
	sto := gof.CurrentApp.Storage()
	if err := sto.Get(key, &list); err != nil {
		if parentId == 0 {
			list = rsi.ProductService.GetBigCategories(int32(shopId))
		} else {
			list = rsi.ProductService.GetChildCategories(int32(shopId), int32(parentId))
		}
		var d []byte
		d, err = json.Marshal(list)
		sto.SetExpire(key, string(d), 3600*24)
	}
	return list
}

/**
 * @api {post} /shop/favorite 收藏店铺
 * @apiGroup shop
 * @apiParam {int} item_id 店铺编号
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (s shopApi) Favorite(ctx api.Context) interface{} {
	memberId := getMemberId(ctx)
	id := ctx.Form().GetInt("shop_id")
	err := rsi.MemberService.FavoriteShop(int64(memberId), int32(id))
	if err != nil {
		return api.ResponseWithCode(1, err.Error())
	}
	return api.NewResponse(nil)
}

/*

// 登陆状态
func (s *serviceC) LoginState(c *echox.Context) error {
	mp := make(map[string]interface{})
	mobileReq := ut.MobileRequest(c.Request())
	mPrefix := util.BoolExt.TString(mobileReq, consts.DOMAIN_PREFIX_M_PASSPORT,
		consts.DOMAIN_PREFIX_PASSPORT)
	pstUrl := fmt.Sprintf("//%s%s", mPrefix, variable.Domain)
	memberId := getMemberId(c)
	if memberId <= 0 {
		registry, _ := rsi.RegistryService.GetRegistries(thrift.Context,
			[]string{"PlatformName"})
		mp["PtName"] = registry["PlatformName"]
		mp["LoginUrl"] = pstUrl + "/auth/login"
		mp["RegisterUrl"] = pstUrl + "/register"
		mp["login"] = 0
	} else {
		mmUrl := fmt.Sprintf("//%s%s",
			consts.DOMAIN_PREFIX_MEMBER, variable.Domain)
		m, _ := rsi.MemberService.GetProfile(thrift.Context, int64(memberId))
		mp["MMName"] = m.Name
		mp["LogoutUrl"] = pstUrl + "/auth/logout"
		mp["MMUrl"] = mmUrl
		mp["login"] = 1
	}
	return c.JSONP(http.StatusOK, c.QueryParam("callback"), mp)
}
*/

//

/**
 * @api {post} /shop/my_address 收货地址列表
 * @apiGroup shop
 * @apiSuccessExample Success-Response
 * {}
 * @apiSuccessExample Error-Response
 * {"err_code":1,"err_msg":"access denied"}
 */
func (s shopApi) addressList(ctx api.Context) interface{} {
	memberId := getMemberId(ctx)
	address, _ := rsi.MemberService.GetAddressList(thrift.Context, int64(memberId))
	return address
}
