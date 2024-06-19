/**
 * Copyright 2015 @ 56x.net.
 * name : express
 * author : jarryliu
 * date : 2016-07-05 14:30
 * description :
 * history :
 */
package express

import "github.com/ixre/go2o/core/infrastructure/domain"

var (
	ErrExpressTemplateName = domain.NewError(
		"err_express_template_name", "运费模板名称不能为空")
	ErrUserNotMatch = domain.NewError(
		"err_express_user_not_match", "运费模板用户不匹配")
	ErrExpressBasis = domain.NewError(
		"err_express_basis", "运费计价方式不正确")
	ErrAddFee = domain.NewError(
		"err_express_add_fee", "续重(件)费用必须大于零")
	ErrFirstUnitNotSet = domain.NewError(
		"err_express_first_unit_not_set", "首重(件)单位数量未填写")
	ErrAddUnitNotSet = domain.NewError(
		"err_express_add_unit_not_set", "续重(件)单位数量未填写")
	ErrExpressTemplateMissingAreaCode = domain.NewError(
		"err_express_template_missing_area_code", "运费模板未指定地区")
	ErrExistsAreaTemplateSet = domain.NewError(
		"err_express_exists_area_template_set", "地区已存在运费模板设置")
	ErrNoSuchTemplate = domain.NewError(
		"err_express_no_such_template", "运费模板不存在")
	ErrTemplateNotEnabled = domain.NewError(
		"err_express_template_not_enabled", "运费模板未启用")
	ErrNotSupportProvider = domain.NewError(
		"err_express_no_support_provider", "不支持该物流服务商")
)

const (
	//根据件数计算运费,通常大件物品,可以按件收费
	BasisByNumber = 1
	//根据重量计算运费
	BasisByWeight = 2
	//按体积(容积)来计算运费,比如饮料
	BasisByVolume = 3
)

var (
	//todo: 选择一些主流的快递
	// 系统支持的快递服务商
	SupportedExpressProvider = []*Provider{
		NewExpressProvider("安能快递", "A-E", "ANE66", "ANE66"),
		NewExpressProvider("百世汇通", "常用,A-E", "HTKY", "HTKY"),
		NewExpressProvider("CCES", "A-E", "CCES", "CCES"),
		NewExpressProvider("大田物流", "A-E", "DTWL", "DTWL"),
		NewExpressProvider("德邦物流", "常用,A-E", "DBL", "DBL"),
		NewExpressProvider("EMS", "常用,A-E", "EMS", "EMS"),

		NewExpressProvider("飞远配送", "F-J", "GZLT", "GZLT"),
		NewExpressProvider("港中能达", "F-J", "NEDA", "NEDA"),
		NewExpressProvider("龙邦物流", "F-J", "LB", "LB"),
		NewExpressProvider("联邦快递", "F-J", "FEDEX", "FEDEX"),
		NewExpressProvider("联昊通物流", "F-J", "LHT", "LHT"),
		NewExpressProvider("国通快递", "F-J", "GTO", "GTO"),
		NewExpressProvider("海航天天快递", "F-J", "HHTT", "HHTT"),

		NewExpressProvider("全峰快递", "常用,P-T", "QFKD", "QFKD"),
		NewExpressProvider("全一快递", "P-T", "UAPEX", "UAPEX"),
		NewExpressProvider("全日通快递", "P-T", "QRT", "QRT"),

		NewExpressProvider("顺丰快递", "常用,U-Z", "SF", "SF"),
		NewExpressProvider("申通物流", "常用,U-Z", "STO", "STO"),
		NewExpressProvider("圆通速递", "常用,U-Z", "YTO", "YTO"),
		NewExpressProvider("中通速递", "常用,U-Z", "ZTO", "ZTO"),
		NewExpressProvider("韵达快运", "常用,U-Z", "YD", "YD"),
		NewExpressProvider("优速物流", "常用,U-Z", "UC", "UC"),
		NewExpressProvider("宅急送", "常用,U-Z", "ZJS", "ZJS"),
		NewExpressProvider("新邦物流", "U-Z", "XBWL", "XBWL"),
		NewExpressProvider("邮政平邮", "常用,U-Z", "YZPY", "YZPY"),
		NewExpressProvider("中铁快运", "U-Z", "ZTKY", "ZTKY"),
		NewExpressProvider("亚风速递", "U-Z", "YFSD", "YFSD"),
		NewExpressProvider("其它", "U-Z", "OTHER", "OTHER"),
	}
)

type (
	// 运费计算器
	IExpressCalculator interface {
		// 添加计算项,tplId为运费模板的编号
		Add(tplId int, unit int) error

		// 计算运费
		Calculate(areaCode string)

		// 获取累计运费
		Total() int64

		// 获取运费模板编号与费用的集合
		Fee() map[int]int64
	}

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

		// 创建运费计算器
		CreateCalculator() IExpressCalculator
	}

	// 快递模板
	IExpressTemplate interface {
		// 获取领域对象编号
		GetDomainId() int

		// 获取快递模板数据
		Value() ExpressTemplate

		// 设置地区的快递模板
		Set(v *ExpressTemplate) error

		// 设置地区运费
		SetDistrictExpress(arr *[]DistrictExpressTemplate) error

		// 地区运费设置
		DistrictExpress() []DistrictExpressTemplate

		// 是否启用
		Enabled() bool

		// 保存
		Save() (int, error)

		// 根据地区编码获取运费模板
		GetAreaExpressTemplateByAreaCode(areaCode string) *DistrictExpressTemplate
	}

	IExpressRepo interface {
		// 获取所有快递公司
		GetExpressProviders() []*Provider

		// 获取快递公司
		GetExpressProvider(id int32) *Provider

		// 保存快递公司
		SaveExpressProvider(v *Provider) (int32, error)

		// 获取用户的快递
		GetUserExpress(userId int) IUserExpress

		// 获取用户的快递模板
		GetUserAllTemplate(userId int) []*ExpressTemplate

		// 删除快递模板
		DeleteExpressTemplate(userId int, templateId int) error

		// 保存快递模板
		SaveExpressTemplate(value *ExpressTemplate) (int, error)

		// 获取模板的所有地区设置
		GetExpressTemplateAllAreaSet(templateId int) []DistrictExpressTemplate

		// 保存模板的地区设置
		SaveExpressTemplateAreaSet(t *DistrictExpressTemplate) (int, error)

		// 删除模板的地区设置
		DeleteAreaExpressTemplate(templateId int, areaSetId int) error
	}

	// 快递服务商
	Provider struct {
		// 快递公司编号
		Id int32 `db:"id" pk:"yes" auto:"yes"`
		// 快递名称
		Name string `db:"name"`
		// 首字母，用于索引分组
		FirstLetter string `db:"-"` //`db:"letter"`
		// 分组,多个组,用","隔开
		GroupFlag string `db:"group_flag"`
		// 快递公司编码
		Code string `db:"code"`
		// 接口编码
		ApiCode string `db:"api_code"`
		// 是否启用
		Enabled int `db:"enabled"`
	}

	// 快递模板
	ExpressTemplate struct {
		// 编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 运营商编号
		VendorId int `db:"vendor_id"`
		// 运费模板名称
		Name string `db:"name"`
		// 是否卖价承担运费
		IsFree int `db:"is_free"`
		// 运费计价依据
		Basis int `db:"basis"`
		// 首次计价单位,如首重为2kg
		FirstUnit int `db:"first_unit"`
		// 首次计价单价(元),如续重1kg
		FirstFee int64 `db:"first_fee"`
		// 超过首次计价计算单位,如续重1kg
		AddUnit int `db:"add_unit"`
		// 超过首次计价单价(元)，如续重1kg
		AddFee int64 `db:"add_fee"`
		// 是否启用
		Enabled int `db:"enabled"`
	}

	// 快递地区模板
	DistrictExpressTemplate struct {
		// 模板编号
		Id int `db:"id" pk:"yes" auto:"yes"`
		// 运费模板编号
		TemplateId int `db:"template_id"`
		// 地区编号列表，通常精确到省即可
		CodeList string `db:"code_list"`
		// 地区名称列表
		NameList string `db:"name_list"`
		// 首次数值，如 首重为2kg
		FirstUnit int32 `db:"first_unit"`
		// 首次金额，如首重10元
		FirstFee int64 `db:"first_fee"`
		// 增加数值，如续重1kg
		AddUnit int32 `db:"add_unit"`
		// 增加产生费用，如续重1kg 10元
		AddFee int64 `db:"add_fee"`
	}
)

func NewExpressProvider(name, group, code, apiCode string) *Provider {
	return &Provider{
		Name:      name,
		GroupFlag: group,
		Code:      code,
		ApiCode:   apiCode,
		Enabled:   1,
	}
}
