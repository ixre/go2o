/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:20
 * description :
 * history :
 */
package query

import (
	"database/sql"
	"regexp"
	"time"

	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/infrastructure/fw"
	tool "github.com/ixre/go2o/core/infrastructure/util"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/log"
	"github.com/ixre/gof/storage"
	"github.com/ixre/gof/typeconv"
)

type SummaryStatistics struct {
	// 汇总总数
	TotalMembers int64 `db:"totalMembers" json:"totalMembers"`
	// 今日注册会员数
	TodayJoinMembers int64 `db:"todayJoinMembers" json:"todayJoinMembers"`
	// 今日登录数
	TodayLoginMembers int64 `db:"todayLoginMembers" json:"todayLoginMembers"`
	// 今日新增订单数
	TodayCreateOrders int64 `db:"todayCreateOrders" json:"todayCreateOrders"`
	// 待发货订单数
	AwaitShipmentOrders int64 `db:"awaitShipmentOrders" json:"awaitShipmentOrders"`
	// 待审核提现申请数量
	AwaitReviewWithdrawRequests int64 `db:"awaitReviewWithdrawRequests" json:"awaitReviewWithdrawRequests"`
	// 商户总数
	TotalMerchants int `db:"totalMerchants" json:"totalMerchants"`
	// 待审核商户
	AwaitApproveMerchants int `db:"awaitApproveMerchants" json:"awaitApproveMerchants"`
	// 新增商户员工
	TotalMerchantStaffs int `db:"totalMerchantStaffs" json:"totalMerchantStaffs"`
	// 待审核商户员工
	AwaitApproveStaffs int `db:"awaitApproveStaffs" json:"awaitApproveStaffs"`
	// 今日商户服务数
	TodayMchServiceCount int `db:"todayMchServiceCount" json:"todayMchServiceCount"`
	// 今日商户服务金额
	TodayMchServiceAmount int `db:"todayMchServiceAmount" json:"todayMchServiceAmount"`
}

type StatisticsQuery struct {
	db.Connector
	o              orm.Orm
	on             fw.ORM
	Storage        storage.Interface
	commHostRegexp *regexp.Regexp
}

func NewStatisticsQuery(o orm.Orm, on fw.ORM, s storage.Interface) *StatisticsQuery {
	return &StatisticsQuery{
		Connector: o.Connector(),
		o:         o,
		Storage:   s,
		on:        on,
	}
}

// QuerySummary 查询汇总信息
func (s *StatisticsQuery) QuerySummary() *SummaryStatistics {
	var ss SummaryStatistics
	var todayBeginTime = tool.GetStartDate(time.Now()).Unix()
	err := s.Connector.QueryRow(`
		SELECT (SELECT COUNT(1) FROM mm_member) as totalMembers,
		(SELECT COUNT(1) FROM mm_member WHERE reg_time > $1) as totalMembers,
		(SELECT COUNT(1) FROM mm_member WHERE login_time > $1) as todayLoginMembers,
		(SELECT COUNT(1) FROM sale_sub_order WHERE create_time > $1) as todayCreateOrders,
		(SELECT COUNT(1) FROM sale_sub_order WHERE status = $2) as awaitShipmentOrders,
		(SELECT COUNT(1) FROM wal_wallet_log WHERE review_status = $3 
			AND kind IN (22,23)) as awaitReviewWithdrawRequests
		`,
		func(row *sql.Row) error {
			return row.Scan(
				&ss.TotalMembers,
				&ss.TodayJoinMembers,
				&ss.TodayLoginMembers,
				&ss.TodayCreateOrders,
				&ss.AwaitShipmentOrders,
				&ss.AwaitReviewWithdrawRequests,
			)
		},
		todayBeginTime,
		order.StatAwaitingShipment,
		wallet.ReviewPending)
	if err != nil {
		log.Error(err)
	}
	mp := make(map[string]interface{}, 0)
	s.on.Raw(`
	SELECT * FROM (SELECT (SELECT COUNT(*) FROM mch_merchant) as total_merchants,
		(SELECT COUNT(*) FROM mch_authenticate WHERE review_status = ?) as await_approve_merchants,
		(SELECT COUNT(*) FROM mch_staff) as total_merchant_staffs,
		(SELECT COUNT(c.*) FROM mm_cert_info c INNER JOIN mch_staff s ON s.member_id=c.member_id WHERE review_status = ?) as await_approve_staffs) t
		,(SELECT COUNT(*) as today_mch_service_orders,SUM(final_fee) as today_mch_service_amount FROM mch_service_order where create_time > ?) t2
`, enum.ReviewPending, enum.ReviewPending, todayBeginTime).Scan(&mp)

	ss.TotalMerchants = typeconv.Int(mp["total_merchants"])
	ss.AwaitApproveMerchants = typeconv.Int(mp["await_approve_merchants"])
	ss.TotalMerchantStaffs = typeconv.Int(mp["total_merchant_staffs"])
	ss.AwaitApproveStaffs = typeconv.Int(mp["await_approve_staffs"])
	ss.TodayMchServiceCount = typeconv.Int(mp["today_mch_service_orders"])
	ss.TodayMchServiceAmount = typeconv.Int(mp["today_mch_service_amount"])
	return &ss
}

type ServiceStatistics struct {
	Date   string  `json:"date"`
	Count  int     `json:"count"`
	Amount float64 `json:"amount"`
	Users  int     `json:"users"`
}

// QueryServiceStatistics 查询服务统计信息
func (s *StatisticsQuery) QueryServiceStatistics(beginTime, endTime int) []*ServiceStatistics {
	rows := make([]*ServiceStatistics, 0)
	s.on.Raw(`select b::date as date,count(id) as count,
			count(distinct member_id) as users,
			coalesce(SUM(final_fee),0.00) as amount
			FROM generate_series(
				to_timestamp(?),
				to_timestamp(?),'1 days'
			) as b
			LEFT JOIN mch_service_order
			ON b::date = date_trunc('day', to_timestamp(create_time)::date)

			group by b,final_fee
			ORDER BY b`, beginTime, endTime).Scan(&rows)
	return rows
}

// QueryCityMerchantStaffs 查询城市商户员工统计信息(前十)
func (s *StatisticsQuery) QueryCityMerchantStaffs() (rows []*struct {
	StationId int    `json:"stationId"`
	Label     string `json:"label"`
	Value     int    `json:"value"`
}) {
	s.on.Raw(`
		select station_id, coalesce(name,'其他') as label,count as value FROM sys_sub_station s 
		INNER JOIN sys_district d ON d.code = s.city_code
		RIGHT JOIN (select station_id,count(0) as count FROM mch_staff group by station_id) t
		ON t.station_id = s.id
		ORDER BY count DESC LIMIT 10`).Scan(&rows)
	return rows
}
