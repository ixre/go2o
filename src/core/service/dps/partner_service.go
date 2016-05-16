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
	"go2o/src/core/domain/interface/ad"
	"go2o/src/core/domain/interface/merchant"
	"go2o/src/core/domain/interface/merchant/mss"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/domain/interface/valueobject"
	"go2o/src/core/query"
	"log"
	"strings"
)

type partnerService struct {
	_partnerRep merchant.IMerchantRep
	_adRep      ad.IAdvertisementRep
	_saleRep    sale.ISaleRep
	_query      *query.PartnerQuery
}

func NewPartnerService(r merchant.IMerchantRep, saleRep sale.ISaleRep,
	adRep ad.IAdvertisementRep, q *query.PartnerQuery) *partnerService {
	return &partnerService{
		_partnerRep: r,
		_query:      q,
		_saleRep:    saleRep,
		_adRep:      adRep,
	}
}

// 验证用户密码并返回编号
func (this *partnerService) Verify(usr, pwd string) int {
	usr = strings.ToLower(strings.TrimSpace(usr))
	return this._query.Verify(usr, pwd)
}

func (this *partnerService) GetMerchant(merchantId int) (*merchant.MerchantValue, error) {
	pt, err := this._partnerRep.GetMerchant(merchantId)
	if pt != nil {
		v := pt.GetValue()
		return &v, err
	}
	return nil, err
}

func (this *partnerService) SaveMerchant(merchantId int, v *merchant.MerchantValue) (int, error) {
	var pt merchant.IMerchant
	var err error
	var isCreate bool

	v.Id = merchantId

	if merchantId > 0 {
		pt, _ = this._partnerRep.GetMerchant(merchantId)
		if pt == nil {
			err = errors.New("no such partner")
		} else {
			err = pt.SetValue(v)
		}
	} else {
		isCreate = true
		pt, err = this._partnerRep.CreateMerchant(v)
	}

	if err != nil {
		return 0, err
	}

	merchantId, err = pt.Save()

	if isCreate {
		this.initializePartner(merchantId)
	}

	return merchantId, err
}

func (this *partnerService) initializePartner(merchantId int) {

	// 初始化广告
	this._adRep.GetPartnerAdvertisement(merchantId).InitInternalAdvertisements()

	// 初始化会员默认等级
	pt, _ := this._partnerRep.GetMerchant(merchantId)
	pt.LevelManager().InitDefaultLevels()

	// 保存站点设置
	pt.SaveSiteConf(&merchant.SiteConf{
		MerchantId: pt.GetAggregateRootId(),
		IndexTitle: pt.GetValue().Name,
	})

	// 保存销售设置
	pt.SaveSaleConf(&merchant.SaleConf{
		MerchantId:   pt.GetAggregateRootId(),
		RegisterMode: merchant.ModeRegisterNormal,
	})

	// 初始化销售标签
	this._saleRep.GetSale(merchantId).InitSaleTags()
}

// 根据主机查询商户编号
func (this *partnerService) GetMerchantIdByHost(host string) int {
	return this._query.QueryMerchantIdByHost(host)
}

// 获取商户的域名
func (this *partnerService) GetMerchantMajorHost(merchantId int) string {
	pt, err := this._partnerRep.GetMerchant(merchantId)
	if err != nil {
		log.Println("[ Partner][ Service]-", err.Error())
	}
	return pt.GetMajorHost()
}

func (this *partnerService) SaveSiteConf(merchantId int, v *merchant.SiteConf) error {
	pt, _ := this._partnerRep.GetMerchant(merchantId)
	return pt.SaveSiteConf(v)
}

func (this *partnerService) SaveSaleConf(merchantId int, v *merchant.SaleConf) error {
	pt, _ := this._partnerRep.GetMerchant(merchantId)
	return pt.SaveSaleConf(v)
}

func (this *partnerService) GetSaleConf(merchantId int) *merchant.SaleConf {
	pt, err := this._partnerRep.GetMerchant(merchantId)
	if err != nil {
		log.Println("[ Partner][ Service]-", err.Error())
	}
	conf := pt.GetSaleConf()
	return &conf
}

func (this *partnerService) GetSiteConf(merchantId int) *merchant.SiteConf {
	pt, err := this._partnerRep.GetMerchant(merchantId)
	if err != nil {
		log.Println("[ Partner][ Service]-", err.Error())
	}
	conf := pt.GetSiteConf()
	return &conf
}

// 检查注册权限
func (this *partnerService) CheckRegisterPerm(merchantId int, isInvitation bool) error {
	//	if conf.RegisterMode == partner.ModeRegisterClosed {
	//		return errors.New("1011:系统暂不开放注册")
	//	}
	//	if conf.RegisterMode == partner.ModeRegisterMustInvitation && len(code) == 0 {
	//		return errors.New("1011:必须使用推荐码注册")
	//	}
	//	if conf.RegisterMode == partner.ModeRegisterMustRedirect && len(code) > 0 {
	//		return errors.New("1011:系统暂不开放推荐注册")
	//	}

	pt, err := this._partnerRep.GetMerchant(merchantId)
	if err == nil {
		err = pt.RegisterPerm(isInvitation)
	}
	return err
}

func (this *partnerService) GetShopsOfMerchant(merchantId int) []*merchant.ValueShop {
	pt, err := this._partnerRep.GetMerchant(merchantId)
	if err != nil {
		log.Println("[ Partner][ Service]-", err.Error())
	}
	shops := pt.GetShops()
	sv := make([]*merchant.ValueShop, len(shops))
	for i, v := range shops {
		vv := v.GetValue()
		sv[i] = &vv
	}
	return sv
}

func (this *partnerService) GetShopValueById(merchantId, shopId int) *merchant.ValueShop {
	pt, err := this._partnerRep.GetMerchant(merchantId)
	if err != nil {

		log.Println("[ Partner][ Service]-", err.Error())
	}
	v := pt.GetShop(shopId).GetValue()
	return &v
}

func (this *partnerService) SaveShop(merchantId int, v *merchant.ValueShop) (int, error) {
	pt, err := this._partnerRep.GetMerchant(merchantId)
	if err != nil {

		log.Println("[ Partner][ Service]-", err.Error())
		return 0, err
	}
	var shop merchant.IShop
	if v.Id > 0 {
		shop = pt.GetShop(v.Id)
		if shop == nil {
			return 0, errors.New("门店不存在")
		}
	} else {
		shop = pt.CreateShop(v)
	}
	err = shop.SetValue(v)
	if err != nil {
		return 0, err
	}
	return shop.Save()
}

func (this *partnerService) DeleteShop(merchantId, shopId int) error {
	pt, err := this._partnerRep.GetMerchant(merchantId)
	if err != nil {

		log.Println("[ Partner][ Service]-", err.Error())
	}
	return pt.DeleteShop(shopId)
}

func (this *partnerService) GetPartnersId() []int {
	return this._partnerRep.GetMerchantsId()
}

// 保存API信息
func (this *partnerService) SaveApiInfo(merchantId int, d *merchant.ApiInfo) error {
	pt, _ := this._partnerRep.GetMerchant(merchantId)
	return pt.SaveApiInfo(d)
}

// 获取API接口
func (this *partnerService) GetApiInfo(merchantId int) *merchant.ApiInfo {
	pt, _ := this._partnerRep.GetMerchant(merchantId)
	v := pt.GetApiInfo()
	return &v
}

// 根据API ID获取MerchantId
func (this *partnerService) GetMerchantIdByApiId(apiId string) int {
	return this._partnerRep.GetMerchantIdByApiId(apiId)
}

// 获取所有会员等级
func (this *partnerService) GetMemberLevels(merchantId int) []*valueobject.MemberLevel {
	pt, _ := this._partnerRep.GetMerchant(merchantId)
	return pt.LevelManager().GetLevelSet()
}

// 根据编号获取会员等级信息
func (this *partnerService) GetMemberLevelById(merchantId, id int) *valueobject.MemberLevel {
	pt, _ := this._partnerRep.GetMerchant(merchantId)
	return pt.LevelManager().GetLevelById(id)
}

// 保存会员等级信息
func (this *partnerService) SaveMemberLevel(merchantId int, v *valueobject.MemberLevel) (int, error) {
	pt, _ := this._partnerRep.GetMerchant(merchantId)
	return pt.LevelManager().SaveLevel(v)
}

// 删除会员等级
func (this *partnerService) DelMemberLevel(merchantId, levelId int) error {
	pt, _ := this._partnerRep.GetMerchant(merchantId)
	return pt.LevelManager().DeleteLevel(levelId)
}

// 获取等级
func (this *partnerService) GetLevel(merchantId, level int) *valueobject.MemberLevel {
	pt, _ := this._partnerRep.GetMerchant(merchantId)
	return pt.LevelManager().GetLevelByValue(level)
}

// 获取下一个等级
func (this *partnerService) GetNextLevel(merchantId, levelValue int) *valueobject.MemberLevel {
	pt, _ := this._partnerRep.GetMerchant(merchantId)
	return pt.LevelManager().GetNextLevel(levelValue)
}

// 获取键值字典
func (this *partnerService) GetKeyMapsByKeyword(merchantId int, keyword string) map[string]string {
	pt, _ := this._partnerRep.GetMerchant(merchantId)
	return pt.KvManager().GetsByChar(keyword)
}

// 保存键值字典
func (this *partnerService) SaveKeyMaps(merchantId int, data map[string]string) error {
	pt, err := this._partnerRep.GetMerchant(merchantId)
	if err != nil {
		return err
	}
	return pt.KvManager().Sets(data)
}

// 获取邮件模版
func (this *partnerService) GetMailTemplate(merchantId int, id int) (*mss.MailTemplate, error) {
	pt, err := this._partnerRep.GetMerchant(merchantId)
	if err != nil {
		return nil, err
	}
	return pt.MssManager().GetMailTemplate(id), nil
}

// 保存邮件模板
func (this *partnerService) SaveMailTemplate(merchantId int, v *mss.MailTemplate) (int, error) {
	if v.MerchantId != merchantId {
		return 0, merchant.ErrPartnerNotMatch
	}
	pt, err := this._partnerRep.GetMerchant(merchantId)
	if err != nil {
		return 0, err
	}
	return pt.MssManager().SaveMailTemplate(v)
}

// 获取邮件模板
func (this *partnerService) GetMailTemplates(merchantId int) []*mss.MailTemplate {
	pt, err := this._partnerRep.GetMerchant(merchantId)
	if err != nil {
		return nil
	}
	return pt.MssManager().GetMailTemplates()
}

// 删除邮件模板
func (this *partnerService) DeleteMailTemplate(merchantId int, id int) error {
	pt, err := this._partnerRep.GetMerchant(merchantId)
	if err != nil {
		return err
	}
	return pt.MssManager().DeleteMailTemplate(id)
}
