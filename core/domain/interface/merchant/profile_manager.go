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
		Id int `db:"id" pk:"yes" auto:"yes" json:"id" bson:"id"`
		// 商户编号
		MchId int `db:"mch_id" json:"mchId" bson:"mchId"`
		// 公司名称
		OrgName string `db:"org_name" json:"orgName" bson:"orgName"`
		// 营业执照编号
		OrgNo string `db:"org_no" json:"orgNo" bson:"orgNo"`
		// 营业执照照片
		OrgPic string `db:"org_pic" json:"orgPic" bson:"orgPic"`
		// 办公地
		WorkCity int `db:"work_city" json:"workCity" bson:"workCity"`
		// 资质图片
		QualificationPic string `db:"qualification_pic" json:"qualificationPic" bson:"qualificationPic"`
		// 法人身份证号
		PersonId string `db:"person_id" json:"personId" bson:"personId"`
		// 法人姓名
		PersonName string `db:"person_name" json:"personName" bson:"personName"`
		// 法人身份证照片
		PersonPic string `db:"person_pic" json:"personPic" bson:"personPic"`
		// 联系人手机
		PersonPhone string `db:"person_phone" json:"personPhone" bson:"personPhone"`
		// 授权书
		AuthorityPic string `db:"authority_pic" json:"authorityPic" bson:"authorityPic"`
		// 开户银行
		BankName string `db:"bank_name" json:"bankName" bson:"bankName"`
		// 开户户名
		BankAccount string `db:"bank_account" json:"bankAccount" bson:"bankAccount"`
		// 开户账号
		BankNo string `db:"bank_no" json:"bankNo" bson:"bankNo"`
		// 扩展数据
		ExtraData string `db:"extra_data" json:"extraData" bson:"extraData"`
		// 审核状态
		ReviewStatus int `db:"review_status" json:"reviewStatus" bson:"reviewStatus"`
		// 审核备注
		ReviewRemark string `db:"review_remark" json:"reviewRemark" bson:"reviewRemark"`
		// 审核时间
		ReviewTime int `db:"review_time" json:"reviewTime" bson:"reviewTime"`
		// 版本号: 0: 待审核 1: 已审核
		Version int `db:"version" json:"version" bson:"version"`
		// 更新时间
		UpdateTime int `db:"update_time" json:"updateTime" bson:"updateTime"`
	}

	// 基本资料管理器
	IProfileManager interface {
		// SaveAuthenticate 保存商户认证信息
		SaveAuthenticate(v *Authenticate) (int, error)

		// 标记企业为审核通过
		ReviewAuthenticate(reviewed bool, message string) error

		// 修改密码
		ChangePassword(newPassword, oldPwd string) error
	}
)
