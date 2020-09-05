/**
 * Copyright 2015 @ to2.net.
 * name : profilemanager.go
 * author : jarryliu
 * date : 2016-05-26 21:19
 * description :
 * history :
 */
package merchant

type (
	// 企业信息
	EnterpriseInfo struct {
		// 编号
		ID int32 `db:"id" pk:"yes" auto:"yes"`
		// 商户编号
		MchId int32 `db:"mch_id"`
		// 公司名称
		CompanyName string `db:"company_name"`
		// 公司营业执照编号
		CompanyNo string `db:"company_no"`
		// 法人
		PersonName string `db:"person_name"`
		// 法人身份证编号
		PersonIdNo string `db:"person_id"`
		// 身份证验证图片(人捧身份证照相)
		PersonImage string `db:"person_image"`
		// 公司电话
		Tel string `db:"tel"`
		// 省
		Province int32 `db:"province"`
		// 市
		City int32 `db:"city"`
		// 区
		District int32 `db:"district"`
		// 省+市+区字符串表示
		Location string `db:"location"`
		// 公司地址
		Address string `db:"address"`
		// 营业执照图片
		CompanyImage string `db:"company_image"`
		// 授权书
		AuthDoc string `db:"auth_doc"`
		//是否已审核
		Reviewed int32 `db:"review_state"`
		// 审核时间
		ReviewTime int64 `db:"review_time"`
		// 审核备注
		ReviewRemark string `db:"review_remark"`
		//更新时间
		UpdateTime int64 `db:"update_time"`
	}

	// 基本资料管理器
	IProfileManager interface {
		// 获取企业信息
		GetEnterpriseInfo() *EnterpriseInfo

		// 保存企业信息
		SaveEnterpriseInfo(v *EnterpriseInfo) (int32, error)

		// 标记企业为审核通过
		ReviewEnterpriseInfo(reviewed bool, message string) error

		// 修改密码
		ModifyPassword(newPwd, oldPwd string) error
	}
)
