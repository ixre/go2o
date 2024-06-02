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

	"github.com/ixre/go2o/core/domain/interface/order"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	tool "github.com/ixre/go2o/core/infrastructure/util"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"github.com/ixre/gof/storage"
)

type SummaryStatistics struct {
	// 汇总总数
	TotalMembers int64 `db:"totalMembers"`
	// 今日注册会员数
	TodayJoinMembers int64 `db:"todayJoinMembers"`
	// 今日登录数
	TodayLoginMembers int64 `db:"todayLoginMembers"`
	// 今日新增订单数
	TodayCreateOrders int64 `db:"todayCreateOrders"`
	// 待发货订单数
	AwaitShipmentOrders int64 `db:"awaitShipmentOrders"`
	// 待审核提现申请数量
	AwaitReviewWithdrawRequests int64 `db:"awaitReviewWithdrawRequests"`
}

type StatisticsQuery struct {
	db.Connector
	o              orm.Orm
	Storage        storage.Interface
	commHostRegexp *regexp.Regexp
}

func NewStatisticsQuery(o orm.Orm, s storage.Interface) *StatisticsQuery {
	return &StatisticsQuery{
		Connector: o.Connector(),
		o:         o,
		Storage:   s,
	}
}

// QuerySummary 查询汇总信息
func (s *StatisticsQuery) QuerySummary() *SummaryStatistics {
	var ss SummaryStatistics
	var todayBeginTime = tool.GetStartDate(time.Now()).Unix()
	s.Connector.QueryRow(`
		(SELECT COUNT(1) FROM mm_member) as totalMembers,
		(SELECT COUNT(1) FROM mm_member WHERE > $1) as totalMembers,
		(SELECT COUNT(1) FROM mm_member WHERE last_login_time > $1) as todayLoginMembers,
		(SELECT COUNT(1) FROM sale_sub_order WHERE create_time > $1) as todayCreateOrders,
		(SELECT COUNT(1) FROM sale_sub_order WHERE status = $2) as awaitShipmentOrders,
		(SELECT COUNT(1) FROM wal_wallet_log WHERE review_status = $3 AND kind IN (22,23)) as awaitReviewWithdrawRequests,	
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
		wallet.ReviewAwaiting)
	return &ss
}
