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
		GetExpressFee(areaCode string, unit int) float32
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
		SaveAreaTemplate(t *ExpressAreaTemplate) error

		// 获取所有的地区快递模板
		GetAllAreaTemplate() []ExpressTemplate
	}

	IExpressRep interface {
		// 获取所有快递公司
		GetExpressProviders() []ExpressProvider

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
		Id int
		// 首字母，用于索引分组
		FirstLetter string
		// 快递公司编码
		Code string
		// 接口编码
		ApiCode string
	}

	// 快递模板
	ExpressTemplate struct {
		// 快递模板编号
		Id int
		// 用户编号
		UserId int
		// 快递模板名称
		Name string
		// 是否卖家承担运费，0为顾客承担
		IsFree int
		// 计价方式:1:按重量;2:按数量;3:按体积
		Basis int
		// 首次数值，如 首重为2kg
		FirstUnit int
		// 首次金额，如首重10元
		FirstFee float32
		// 增加数值，如续重1kg
		AddUnit int
		// 增加产生费用，如续重1kg 10元
		AddFee int
		// 是否启用
		Enabled int
	}

	// 快递地区模板
	ExpressAreaTemplate struct {
		// 模板编号
		Id int
		// 运费模板编号
		TemplateId int
		// 地区编号列表，通常精确到省即可
		CodeList string
		// 地区名称列表
		NameList string
		// 首次数值，如 首重为2kg
		FirstUnit int
		// 首次金额，如首重10元
		FirstFee float32
		// 增加数值，如续重1kg
		AddUnit int
		// 增加产生费用，如续重1kg 10元
		AddFee int
	}
)
