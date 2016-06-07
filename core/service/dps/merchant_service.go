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
	"go2o/core/infrastructure/domain"
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
	mch, err := this._mchRep.GetMerchant(mchId)
	if err != nil {
		domain.HandleError(err)
		return merchant.EnterpriseInfo{}
	}
	if reviewed {
		return mch.ProfileManager().GetReviewedEnterpriseInfo()
	}
	return mch.ProfileManager().GetEnterpriseInfo()
}

// 保存企业信息
func (this *merchantService) SaveEnterpriseInfo(mchId int,
	e *merchant.EnterpriseInfo) (int, error) {
	mch, err := this._mchRep.GetMerchant(mchId)
	if err != nil {
		return -1, err
	}
	return mch.ProfileManager().SaveEnterpriseInfo(e)
}

// 审核企业信息
func (this *merchantService) ReviewEnterpriseInfo(mchId int, pass bool, remark string) error {
	mch, err := this._mchRep.GetMerchant(mchId)
	if err != nil {
		return err
	}
	return mch.ProfileManager().ReviewEnterpriseInfo(pass, remark)
}

func (this *merchantService) GetMerchant(merchantId int) (*merchant.Merchant, error) {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if mch != nil {
		v := mch.GetValue()
		return &v, err
	}
	return nil, err
}

func (this *merchantService) SaveMerchant(merchantId int, v *merchant.Merchant) (int, error) {
	var mch merchant.IMerchant
	var err error
	var isCreate bool

	v.Id = merchantId

	if merchantId > 0 {
		mch, _ = this._mchRep.GetMerchant(merchantId)
		if mch == nil {
			err = errors.New("no such partner")
		} else {
			err = mch.SetValue(v)
		}
	} else {
		isCreate = true
		mch, err = this._mchRep.CreateMerchant(v)
	}

	if err != nil {
		return 0, err
	}

	merchantId, err = mch.Save()

	if isCreate {
		this.initializeMerchant(merchantId)
	}

	return merchantId, err
}

func (this *merchantService) initializeMerchant(merchantId int) {

	// 初始化会员默认等级
	// this._mchRep.GetMerchant(merchantId)

	//conf := merchant.DefaultSaleConf
	//conf.MerchantId = mch.GetAggregateRootId()
	// 保存销售设置
	//mch.ConfManager().SaveSaleConf(&conf)

	// 初始化销售标签
	this._saleRep.GetSale(merchantId).LabelManager().InitSaleLabels()
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
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		log.Println("[ Merchant][ Service]-", err.Error())
	}
	return mch.GetMajorHost()
}

func (this *merchantService) SaveSaleConf(merchantId int, v *merchant.SaleConf) error {
	mch, _ := this._mchRep.GetMerchant(merchantId)
	return mch.ConfManager().SaveSaleConf(v)
}

func (this *merchantService) GetSaleConf(merchantId int) *merchant.SaleConf {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		log.Println("[ Merchant][ Service]-", err.Error(), ",ID:", merchantId)
	}
	conf := mch.ConfManager().GetSaleConf()
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

func (this *merchantService) GetMerchantsId() []int {
	return this._mchRep.GetMerchantsId()
}

// 保存API信息
func (this *merchantService) SaveApiInfo(merchantId int, d *merchant.ApiInfo) error {
	mch, _ := this._mchRep.GetMerchant(merchantId)
	return mch.ApiManager().SaveApiInfo(d)
}

// 获取API接口
func (this *merchantService) GetApiInfo(merchantId int) *merchant.ApiInfo {
	mch, _ := this._mchRep.GetMerchant(merchantId)
	v := mch.ApiManager().GetApiInfo()
	return &v
}

// 启用/停用接口权限
func (this *merchantService) ApiPerm(merchantId int, enabled bool) error {
	mch, _ := this._mchRep.GetMerchant(merchantId)
	if enabled {
		return mch.ApiManager().EnableApiPerm()
	}
	return mch.ApiManager().DisableApiPerm()
}

// 根据API ID获取MerchantId
func (this *merchantService) GetMerchantIdByApiId(apiId string) int {
	return this._mchRep.GetMerchantIdByApiId(apiId)
}

// 获取所有会员等级
func (this *merchantService) GetMemberLevels(merchantId int) []*merchant.MemberLevel {
	mch, _ := this._mchRep.GetMerchant(merchantId)
	return mch.LevelManager().GetLevelSet()
}

// 根据编号获取会员等级信息
func (this *merchantService) GetMemberLevelById(merchantId, id int) *merchant.MemberLevel {
	mch, _ := this._mchRep.GetMerchant(merchantId)
	return mch.LevelManager().GetLevelById(id)
}

// 保存会员等级信息
func (this *merchantService) SaveMemberLevel(merchantId int, v *merchant.MemberLevel) (int, error) {
	mch, _ := this._mchRep.GetMerchant(merchantId)
	return mch.LevelManager().SaveLevel(v)
}

// 删除会员等级
func (this *merchantService) DelMemberLevel(merchantId, levelId int) error {
	mch, _ := this._mchRep.GetMerchant(merchantId)
	return mch.LevelManager().DeleteLevel(levelId)
}

// 获取等级
func (this *merchantService) GetLevel(merchantId, level int) *merchant.MemberLevel {
	mch, _ := this._mchRep.GetMerchant(merchantId)
	return mch.LevelManager().GetLevelByValue(level)
}

// 获取下一个等级
func (this *merchantService) GetNextLevel(merchantId, levelValue int) *merchant.MemberLevel {
	mch, _ := this._mchRep.GetMerchant(merchantId)
	return mch.LevelManager().GetNextLevel(levelValue)
}

// 获取键值字典
func (this *merchantService) GetKeyMapsByKeyword(merchantId int, keyword string) map[string]string {
	mch, _ := this._mchRep.GetMerchant(merchantId)
	return mch.KvManager().GetsByChar(keyword)
}

// 保存键值字典
func (this *merchantService) SaveKeyMaps(merchantId int, data map[string]string) error {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		return err
	}
	return mch.KvManager().Sets(data)
}

// 获取邮件模版
func (this *merchantService) GetMailTemplate(merchantId int, id int) (*mss.MailTemplate, error) {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		return nil, err
	}
	return mch.MssManager().GetMailTemplate(id), nil
}

// 保存邮件模板
func (this *merchantService) SaveMailTemplate(merchantId int, v *mss.MailTemplate) (int, error) {
	if v.MerchantId != merchantId {
		return 0, merchant.ErrMerchantNotMatch
	}
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		return 0, err
	}
	return mch.MssManager().SaveMailTemplate(v)
}

// 获取邮件模板
func (this *merchantService) GetMailTemplates(merchantId int) []*mss.MailTemplate {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		return nil
	}
	return mch.MssManager().GetMailTemplates()
}

// 删除邮件模板
func (this *merchantService) DeleteMailTemplate(merchantId int, id int) error {
	mch, err := this._mchRep.GetMerchant(merchantId)
	if err != nil {
		return err
	}
	return mch.MssManager().DeleteMailTemplate(id)
}
