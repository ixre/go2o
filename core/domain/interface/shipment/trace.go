package shipment

type (
	// 发货单追踪
	ShipOrderTrace struct {
		// 物流单号
		LogisticCode string
		// 承运商代码
		ShipperCode string
		// 发货状态
		ShipState int
		// 更新时间
		UpdateTime int64
		// 包含发货单流
		Flows []*ShipFlow
	}
	// 发货流
	ShipFlow struct {
		// 记录标题
		Subject string
		// 记录时间
		CreateTime int64
		// 备注
		Remark string
	}
)

/*

 {
        "EBusinessID": "1109259",
        "OrderCode": "",
        "ShipperCode": "SF",
        "LogisticCode": "118461988807",
        "Success": true,
        "State": 3,
        "Reason": null,
        "Traces": [
        {
        "AcceptTime": "2014/06/25 08:05:37",
        "AcceptStation": "正在派件..(派件人:邓裕富,电话:18718866310)[深圳 市]",
        "Remark": null
        },
        {
        "AcceptTime": "2014/06/25 04:01:28",
        "AcceptStation": "快件在 深圳集散中心 ,准备送往下一站 深圳 [深圳市]",
        "Remark": null
        },
        {
        "AcceptTime": "2014/06/25 01:41:06",
        "AcceptStation": "快件在 深圳集散中心 [深圳市]",
        "Remark": null
        },
        {
        "AcceptTime": "2014/06/24 20:18:58",
        "AcceptStation": "已收件[深圳市]",
        "Remark": null
        },
        {
        "AcceptTime": "2014/06/24 20:55:28",
        "AcceptStation": "快件在 深圳 ,准备送往下一站 深圳集散中心 [深圳市]",
        "Remark": null
        },
        {
        "AcceptTime": "2014/06/25 10:23:03",
        "AcceptStation": "派件已签收[深圳市]",
        "Remark": null
        },
        {
        "AcceptTime": "2014/06/25 10:23:03",
        "AcceptStation": "签收人是：已签收[深圳市]",
        "Remark": null
        }
        ]
        }

 */
