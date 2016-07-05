/**
 * Copyright 2015 @ z3q.net.
 * name : express
 * author : jarryliu
 * date : 2016-07-05 14:30
 * description :
 * history :
 */
package express

import "go2o/core/infrastructure/domain"

var (
	ErrNotFullExpressTemplate *domain.DomainError = domain.NewDomainError(
		"err_not_full_express_template", "不完整的运费模板")
	ErrExpressTemplateMissingAreaCode *domain.DomainError = domain.NewDomainError(
		"err_express_template_missing_area_code", "运费模板未指定地区")
	ErrExistsAreaTemplateSet *domain.DomainError = domain.NewDomainError(
		"err_express_exists_area_template_set", "地区已存在运费模板设置")
)

const (
	//根据重量计算面积
	BasisByWeight = 1
	//根据件数计算运费
	BasisByNumber = 2
	//根据面积计算运费
	BasisBySpace = 3
)

var (
	//todo: 选择一些主流的快递
	// 系统支持的快递服务商
	SupportedExpressProvider = []*ExpressProvider{
		NewExpressProvider("顺丰快递", "S", "SF", "SF"),
		NewExpressProvider("圆通速递", "Y", "YTO", "YTO"),
		NewExpressProvider("中通速递", "Z", "ZTO", "ZTO"),
		NewExpressProvider("韵达快运", "Y", "YD", "YD"),
		NewExpressProvider("海航天天快递", "H", "HHTT", "HHTT"),
		NewExpressProvider("全峰快递", "Q", "QFKD", "QFKD"),
		NewExpressProvider("EMS", "E", "EMS", "EMS"),
		NewExpressProvider("优速物流", "Y", "UC", "UC"),
		NewExpressProvider("宅急送", "Z", "ZJS", "ZJS"),
		NewExpressProvider("全一快递", "Q", "UAPEX", "UAPEX"),
		NewExpressProvider("联邦快递", "L", "FEDEX", "FEDEX"),
		NewExpressProvider("汇通快运", "H", "HTKY", "HTKY"),
		NewExpressProvider("德邦物流", "D", "DBL", "DBL"),
		NewExpressProvider("中铁快运", "Z", "ZTKY", "ZTKY"),
		NewExpressProvider("CCES", "C", "CCES", "CCES"),
		NewExpressProvider("联昊通物流", "L", "LHT", "LHT"),
		NewExpressProvider("申通物流", "S", "STO", "STO"),
		NewExpressProvider("龙邦物流", "L", "LB", "LB"),
		NewExpressProvider("新邦物流", "X", "XBWL", "XBWL"),
		NewExpressProvider("港中能达", "G", "NEDA", "NEDA"),
		NewExpressProvider("全日通快递", "Q", "QRT", "QRT"),
		NewExpressProvider("邮政平邮", "Y", "YZPY", "YZPY"),
		NewExpressProvider("亚风", "Y", "YFSD", "YFSD"),
		NewExpressProvider("大田物流", "D", "DTWL", "DTWL"),
		NewExpressProvider("其它", "O", "OTHER", "OTHER"),
	}
)

type (
	// 物流快递
	IUserExpress interface {
		// 获取聚合根编号
		GetAggregateRootId() int

		// 创建快递模板
		CreateTemplate(t *ExpressTemplate) IExpressTemplate

		// 获取快递模板
		GetTemplate(id int) IExpressTemplate

		// 获取所有的快递模板
		GetAllTemplate() []IExpressTemplate

		// 删除模板
		DeleteTemplate(id int) error

		// 获取快递费,传入地区编码，根据单位值，如总重量。
		GetExpressFee(templateId int, areaCode string, unit int) float32
	}

	// 快递模板
	IExpressTemplate interface {
		// 获取领域对象编号
		GetDomainId() int

		// 获取快递模板数据
		Value() ExpressTemplate

		// 设置地区的快递模板
		Set(v *ExpressTemplate) error

		// 保存
		Save() (int, error)

		// 根据地区编码获取运费模板
		GetAreaExpressTemplateByAreaCode(areaCode string) *ExpressAreaTemplate

		// 根据编号获取地区的运费模板
		GetAreaExpressTemplate(id int) *ExpressAreaTemplate

		// 保存地区快递模板
		SaveAreaTemplate(t *ExpressAreaTemplate) (int, error)

		// 获取所有的地区快递模板
		GetAllAreaTemplate() []ExpressAreaTemplate
	}

	IExpressRep interface {
		// 获取所有快递公司
		GetExpressProviders() []*ExpressProvider

		// 获取快递公司
		GetExpressProvider(id int) *ExpressProvider

		// 保存快递公司
		SaveExpressProvider(v *ExpressProvider) (int, error)

		// 获取用户的快递
		GetUserExpress(userId int) IUserExpress

		// 获取用户的快递模板
		GetUserAllTemplate(userId int) []*ExpressTemplate

		// 删除快递模板
		DeleteExpressTemplate(userId int, templateId int) error

		// 保存快递模板
		SaveExpressTemplate(value *ExpressTemplate) (int, error)

		// 获取模板的所有地区设置
		GetExpressTemplateAllAreaSet(templateId int) []ExpressAreaTemplate

		// 保存模板的地区设置
		SaveExpressTemplateAreaSet(t *ExpressAreaTemplate) (int, error)
	}

	// 快递服务商
	ExpressProvider struct {
		// 快递公司编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 快递名称
		Name string `db:"name"`
		// 首字母，用于索引分组
		FirstLetter string `db:"letter"`
		// 快递公司编码
		Code string `db:"code"`
		// 接口编码
		ApiCode string `db:"api_code"`
		// 是否启用
		Enabled int `db:"enabled"`
	}

	// 快递模板
	ExpressTemplate struct {
		// 快递模板编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 用户编号
		UserId int `db:"user_id"`
		// 快递模板名称
		Name string `db:"name"`
		// 是否卖家承担运费，0为顾客承担
		IsFree int `db:"is_free"`
		// 计价方式:1:按重量;2:按数量;3:按体积
		Basis int `db:"basis"`
		// 首次数值，如 首重为2kg
		FirstUnit int `db:"first_unit"`
		// 首次金额，如首重10元
		FirstFee float32 `db:"first_fee"`
		// 增加数值，如续重1kg
		AddUnit int `db:"add_unit"`
		// 增加产生费用，如续重1kg 10元
		AddFee float32 `db:"add_fee"`
		// 是否启用
		Enabled int `db:"enabled"`
	}

	// 快递地区模板
	ExpressAreaTemplate struct {
		// 模板编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 运费模板编号
		TemplateId int `db:"template_id"`
		// 地区编号列表，通常精确到省即可
		CodeList string `db:"code_list"`
		// 地区名称列表
		NameList string `db:"name_list"`
		// 首次数值，如 首重为2kg
		FirstUnit int `db:"first_unit"`
		// 首次金额，如首重10元
		FirstFee float32 `db:"first_fee"`
		// 增加数值，如续重1kg
		AddUnit int `db:"add_unit"`
		// 增加产生费用，如续重1kg 10元
		AddFee float32 `db:"add_fee"`
	}
)

func NewExpressProvider(name, letter, code, apiCode string) *ExpressProvider {
	return &ExpressProvider{
		Name:        name,
		FirstLetter: letter,
		Code:        code,
		ApiCode:     apiCode,
		Enabled:     1,
	}
}
