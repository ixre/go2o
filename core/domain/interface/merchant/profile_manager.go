/**
 * Copyright 2015 @ 56x.net.
 * name : profilemanager.go
 * author : jarryliu
 * date : 2016-05-26 21:19
 * description :
 * history :
 */
package merchant

type (
	// Authenticate 商户认证信息
	Authenticate struct {
		// Id
		Id int `json:"id" db:"id" gorm:"column:id" pk:"yes" auto:"yes" bson:"id"`
		// 商户编号
		MchId int `json:"mchId" db:"mch_id" gorm:"column:mch_id" bson:"mchId"`
		// 商户名称
		MchName string `json:"mchName" db:"mch_name" gorm:"column:mch_name" bson:"mchName"`
		// 省
		Province int `json:"province" db:"province" gorm:"column:province" bson:"province"`
		// 市
		City int `json:"city" db:"city" gorm:"column:city" bson:"city"`
		// 区
		District int `json:"district" db:"district" gorm:"column:district" bson:"district"`
		// 公司名称
		OrgName string `json:"orgName" db:"org_name" gorm:"column:org_name" bson:"orgName"`
		// 营业执照编号
		LicenceNo string `json:"orgNo" db:"org_no" gorm:"column:org_no" bson:"orgNo"`
		// 公司地址
		OrgAddress string `json:"orgAddress" db:"org_address" gorm:"column:org_address" bson:"orgAddress"`
		// 营业执照照片
		LicencePic string `json:"orgPic" db:"org_pic" gorm:"column:org_pic" bson:"orgPic"`
		// 办公地
		WorkCity int `json:"workCity" db:"work_city" gorm:"column:work_city" bson:"workCity"`
		// 资质图片
		QualificationPic string `json:"qualificationPic" db:"qualification_pic" gorm:"column:qualification_pic" bson:"qualificationPic"`
		// 法人身份证号
		PersonId string `json:"personId" db:"person_id" gorm:"column:person_id" bson:"personId"`
		// 法人姓名
		PersonName string `json:"personName" db:"person_name" gorm:"column:person_name" bson:"personName"`
		// 法人身份证照片(正反面)
		PersonFrontPic string `json:"personFrontPic" db:"person_front_pic" gorm:"column:person_front_pic" bson:"personFrontPic"`
		// 联系人手机
		PersonPhone string `json:"personPhone" db:"person_phone" gorm:"column:person_phone" bson:"personPhone"`
		// 授权书
		AuthorityPic string `json:"authorityPic" db:"authority_pic" gorm:"column:authority_pic" bson:"authorityPic"`
		// 开户银行
		BankName string `json:"bankName" db:"bank_name" gorm:"column:bank_name" bson:"bankName"`
		// 开户户名
		BankAccount string `json:"bankAccount" db:"bank_account" gorm:"column:bank_account" bson:"bankAccount"`
		// 开户账号
		BankNo string `json:"bankNo" db:"bank_no" gorm:"column:bank_no" bson:"bankNo"`
		// 扩展数据
		ExtraData string `json:"extraData" db:"extra_data" gorm:"column:extra_data" bson:"extraData"`
		// 审核时间
		ReviewTime int `json:"reviewTime" db:"review_time" gorm:"column:review_time" bson:"reviewTime"`
		// 审核状态
		ReviewStatus int `json:"reviewStatus" db:"review_status" gorm:"column:review_status" bson:"reviewStatus"`
		// 审核备注
		ReviewRemark string `json:"reviewRemark" db:"review_remark" gorm:"column:review_remark" bson:"reviewRemark"`
		// 版本号: 0: 待审核 1: 已审核
		Version int `json:"version" db:"version" gorm:"column:version" bson:"version"`
		// 更新时间
		UpdateTime int `json:"updateTime" db:"update_time" gorm:"column:update_time" bson:"updateTime"`
		// 身份证背面照片
		PersonBackPic string `json:"personBackPic" db:"person_back_pic" gorm:"column:person_back_pic" bson:"personBackPic"`
		// 银行账户信息表(企业)/银行卡(个体)
		BankAccountPic string `json:"bankAccountPic" db:"bank_account_pic" gorm:"column:bank_account_pic" bson:"bankAccountPic"`
		// 联系人姓名
		ContactName string `json:"contactName" db:"contact_name" gorm:"column:contact_name" bson:"contactName"`
		// 联系人电话
		ContactPhone string `json:"contactPhone" db:"contact_phone" gorm:"column:contact_phone" bson:"contactPhone"`
	}

	// 基本资料管理器
	IProfileManager interface {
		// GetAuthenticate 获取商户认证信息
		GetAuthenticate() *Authenticate
		// SaveAuthenticate 保存商户认证信息
		SaveAuthenticate(v *Authenticate) (int, error)

		// 标记企业为审核通过
		ReviewAuthenticate(reviewed bool, message string) error

		// 修改密码
		ChangePassword(newPassword, oldPassword string) error
	}
)

func (a *Authenticate) TableName() string {
	return "mch_authenticate"
}
