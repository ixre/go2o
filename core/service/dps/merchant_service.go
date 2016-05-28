/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-19 22:49
 * description :
 * history :
 */

package dps

import (
	"errors"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/mss"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/domain/interface/sale"
	"go2o/core/query"
	"log"
	"strings"
)

type merchantService struct {
	_mchRep  merchant.IMerchantRep
	_saleRep sale.ISaleRep
	_query   *query.MerchantQuery
}

func NewMerchantService(r merchant.IMerchantRep, saleRep sale.ISaleRep,
	q *query.MerchantQuery) *merchantService {
	return &merchantService{
		_mchRep:  r,
		_query:   q,
		_saleRep: saleRep,
	}
}

// 验证用户密码并返回编号
func (this *merchantService) Verify(usr, pwd string) int {
	usr = strings.ToLower(strings.TrimSpace(usr))
	return this._query.Verify(usr, pwd)
}

// 获取企业信息
func (this *merchantService) GetEnterpriseInfo(mchId int,
	reviewed bool) merchant.EnterpriseInfo {
	mch, _ := this._mchRep.GetMerchant(mchId)
	if reviewed {
		return mch.ProfileManager().GetReviewedEnterpriseInfo()
	}
	return mch.ProfileManager().GetEnterpriseInfo()
}

// 保存企业信息
func (this *merchantService) SaveEnterpriseInfo(mchId int,
	e *merchant.EnterpriseInfo) (int, error) {
	mch, _ := this._mchRep.GetMerchant(mchId)
	return mch.ProfileManager().SaveEnterpriseInfo(e)
}

// 审核企业信息
func (this *merchantService) ReviewEnterpriseInfo(mchId int, pass bool, remark string) error {
	mch, _ := this._mchRep.GetMerchant(mchId)
	return mch.ProfileManager().ReviewEnterpriseInfo(pass, remark)
}

func (this *merchantService) GetMerchant(merchantId int) (*merchant.Merchant, error) {
	pt, err := this._mchRep.GetMerchant(merchantId)
	if pt != nil {
		v := pt.GetValue()
		return &v, err
	}
	return nil, err
}

func (this *merchantService) SaveMerchant(merchantId int, v *merchant.Merchant) (int, error) {
	var pt merchant.IMerchant
	var err error
	var isCreate bool

	v.Id = merchantId

	if merchantId > 0 {
		pt, _ = this._mchRep.GetMerchant(merchantId)
		if pt == nil {
			err = errors.New("no such partner")
		} else {
			err = pt.SetValue(v)
		}
	} else {
		isCreate = true
		pt, err = this._mchRep.CreateMerchant(v)
	}

	if err != nil {
		return 0, err
	}

	merchantId, err = pt.Save()

	if isCreate {
		this.initializeMerchant(merchantId)
	}

	return merchantId, err
}

func (this *merchantService) initializeMerchant(merchantId int) {

	// 初始化会员默认等级
	mch, _ := this._mchRep.GetMerchant(merchantId)

	// 保存站点设置
	mch.ShopManager().SaveSiteConf(&shop.ShopSiteConf{
		MerchantId: mch.GetAggregateRootId(),
		IndexTitle: mch.GetValue().Name,
	})

	conf := merchant.DefaultSaleConf
	conf.MerchantId = mch.GetAggregateRootId()
	// 保存销售设置
	mch.ConfManager().SaveSaleConf(&conf)

	// 初始化销售标签
	this._saleRep.GetSale(merchantId).InitSaleTags()
}

// 获取商户的状态
func (this *merchantService) Stat(merchantId int) error {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		return err
	}
	return mch.Stat()
}

// 根据主机查询商户编号
func (this *merchantService) GetMerchantIdByHost(host string) int {
	return this._query.QueryMerchantIdByHost(host)
}

// 获取商户的域名
func (this *merchantService) GetMerchantMajorHost(merchantId int) string {
	pt, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		log.Println("[ Merchant][ Service]-", err.Error())
	}
	return pt.GetMajorHost()
}

func (this *merchantService) SaveSiteConf(merchantId int, v *shop.ShopSiteConf) error {
	mch, _ := this._mchRep.GetMerchant(merchantId)
	return mch.ShopManager().SaveSiteConf(v)
}

func (this *merchantService) SaveSaleConf(merchantId int, v *merchant.SaleConf) error {
	mch, _ := this._mchRep.GetMerchant(merchantId)
	return mch.ConfManager().SaveSaleConf(v)
}

func (this *merchantService) GetSaleConf(merchantId int) *merchant.SaleConf {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		log.Println("[ Merchant][ Service]-", err.Error())
	}
	conf := mch.ConfManager().GetSaleConf()
	return &conf
}

func (this *merchantService) GetSiteConf(merchantId int) *shop.ShopSiteConf {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		log.Println("[ Merchant][ Service]-", err.Error())
	}
	conf := mch.ShopManager().GetSiteConf()
	return &conf
}

func (this *merchantService) GetShopsOfMerchant(merchantId int) []*shop.Shop {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		log.Println("[ Merchant][ Service]-", err.Error())
	}
	shops := mch.ShopManager().GetShops()
	sv := make([]*shop.Shop, len(shops))
	for i, v := range shops {
		vv := v.GetValue()
		sv[i] = &vv
	}
	return sv
}

func (this *merchantService) GetShopValueById(merchantId, shopId int) *shop.Shop {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {

		log.Println("[ Merchant][ Service]-", err.Error())
	}
	v := mch.ShopManager().GetShop(shopId).GetValue()
	return &v
}

func (this *merchantService) SaveShop(merchantId int, v *shop.Shop) (int, error) {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {

		log.Println("[ Merchant][ Service]-", err.Error())
		return 0, err
	}
	var shop shop.IShop
	if v.Id > 0 {
		shop = mch.ShopManager().GetShop(v.Id)
		if shop == nil {
			return 0, errors.New("门店不存在")
		}
	} else {
		shop = mch.ShopManager().CreateShop(v)
	}
	err = shop.SetValue(v)
	if err != nil {
		return 0, err
	}
	return shop.Save()
}

func (this *merchantService) DeleteShop(merchantId, shopId int) error {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {

		log.Println("[ Merchant][ Service]-", err.Error())
	}
	return mch.ShopManager().DeleteShop(shopId)
}

func (this *merchantService) GetMerchantsId() []int {
	return this._mchRep.GetMerchantsId()
}

// 保存API信息
func (this *merchantService) SaveApiInfo(merchantId int, d *merchant.ApiInfo) error {
	pt, _ := this._mchRep.GetMerchant(merchantId)
	return pt.ApiManager().SaveApiInfo(d)
}

// 获取API接口
func (this *merchantService) GetApiInfo(merchantId int) *merchant.ApiInfo {
	pt, _ := this._mchRep.GetMerchant(merchantId)
	v := pt.ApiManager().GetApiInfo()
	return &v
}

// 启用/停用接口权限
func (this *merchantService) ApiPerm(merchantId int, enabled bool) error {
	pt, _ := this._mchRep.GetMerchant(merchantId)
	if enabled {
		return pt.ApiManager().EnableApiPerm()
	}
	return pt.ApiManager().DisableApiPerm()
}

// 根据API ID获取MerchantId
func (this *merchantService) GetMerchantIdByApiId(apiId string) int {
	return this._mchRep.GetMerchantIdByApiId(apiId)
}

// 获取所有会员等级
func (this *merchantService) GetMemberLevels(merchantId int) []*merchant.MemberLevel {
	pt, _ := this._mchRep.GetMerchant(merchantId)
	return pt.LevelManager().GetLevelSet()
}

// 根据编号获取会员等级信息
func (this *merchantService) GetMemberLevelById(merchantId, id int) *merchant.MemberLevel {
	pt, _ := this._mchRep.GetMerchant(merchantId)
	return pt.LevelManager().GetLevelById(id)
}

// 保存会员等级信息
func (this *merchantService) SaveMemberLevel(merchantId int, v *merchant.MemberLevel) (int, error) {
	pt, _ := this._mchRep.GetMerchant(merchantId)
	return pt.LevelManager().SaveLevel(v)
}

// 删除会员等级
func (this *merchantService) DelMemberLevel(merchantId, levelId int) error {
	pt, _ := this._mchRep.GetMerchant(merchantId)
	return pt.LevelManager().DeleteLevel(levelId)
}

// 获取等级
func (this *merchantService) GetLevel(merchantId, level int) *merchant.MemberLevel {
	pt, _ := this._mchRep.GetMerchant(merchantId)
	return pt.LevelManager().GetLevelByValue(level)
}

// 获取下一个等级
func (this *merchantService) GetNextLevel(merchantId, levelValue int) *merchant.MemberLevel {
	pt, _ := this._mchRep.GetMerchant(merchantId)
	return pt.LevelManager().GetNextLevel(levelValue)
}

// 获取键值字典
func (this *merchantService) GetKeyMapsByKeyword(merchantId int, keyword string) map[string]string {
	pt, _ := this._mchRep.GetMerchant(merchantId)
	return pt.KvManager().GetsByChar(keyword)
}

// 保存键值字典
func (this *merchantService) SaveKeyMaps(merchantId int, data map[string]string) error {
	pt, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		return err
	}
	return pt.KvManager().Sets(data)
}

// 获取邮件模版
func (this *merchantService) GetMailTemplate(merchantId int, id int) (*mss.MailTemplate, error) {
	pt, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		return nil, err
	}
	return pt.MssManager().GetMailTemplate(id), nil
}

// 保存邮件模板
func (this *merchantService) SaveMailTemplate(merchantId int, v *mss.MailTemplate) (int, error) {
	if v.MerchantId != merchantId {
		return 0, merchant.ErrMerchantNotMatch
	}
	pt, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		return 0, err
	}
	return pt.MssManager().SaveMailTemplate(v)
}

// 获取邮件模板
func (this *merchantService) GetMailTemplates(merchantId int) []*mss.MailTemplate {
	pt, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		return nil
	}
	return pt.MssManager().GetMailTemplates()
}

// 删除邮件模板
func (this *merchantService) DeleteMailTemplate(merchantId int, id int) error {
	pt, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		return err
	}
	return pt.MssManager().DeleteMailTemplate(id)
}
