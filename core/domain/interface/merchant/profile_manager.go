/**
 * Copyright 2015 @ z3q.net.
 * name : profilemanager.go
 * author : jarryliu
 * date : 2016-05-26 21:19
 * description :
 * history :
 */
package merchant

type(
    // 企业信息
    EnterpriseInfo struct {
        // 编号
        Id              int    `db:"id"`

        // 商户编号
        MerchantId      int   `db:"merchant_id"`

        // 公司名称
        Name            string `db:"name"`

        // 公司营业执照编号
        CompanyNo       string  `db:"company_no"`

        // 法人
        PersonName      string `db:"person_name"`

        // 公司电话
        Tel             string `db:"tel"`

        // 公司地址
        Address         string `db:"address"`

        // 身份证验证图片(人捧身份证照相)
        PersonImageUrl  string `db:"person_imageurl"`

        // 营业执照图片
        CompanyImageUrl string  `db:"company_imageurl"`

        //是否已审核
        Reviewed        int `db:"reviewed"`

        // 审核时间
        ReviewTime      int64 `db:"review_time"`

        // 审核备注
        Remark          string `db:"remark"`

        //更新时间
        updateTime      int64 `db:"update_time"`
    }
)
