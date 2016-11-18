/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-01-08 21:35
 * description :
 * history :
 */

package daemon

import (
	"database/sql"
	"github.com/jsix/gof"
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/order"
	"go2o/core/infrastructure/tool"
	"go2o/core/service/dps"
	"log"
	"sync"
	"time"
)

var (
	mchIds []int
)

func getMerchants() []int {
	if mchIds == nil {
		mchIds = dps.MerchantService.GetMerchantsId()
	}
	return mchIds
}

/***** OLD CODE *****/
// todo: 等待重构

func orderDaemon(app gof.App) {
	defer recoverDaemon()
	ids := getMerchants()
	for _, v := range ids {
		autoSetOrder(v)
	}
}

func autoSetOrder(mchId int64) {
	f := func(err error) {
		appCtx.Log().Error(err)
	}
	dps.ShoppingService.OrderAutoSetup(mchId, f)
}

var (
	mchDayChartKey string = "cron:go2o:d:mch:day-chart-unix"
)

func testGenerateMchDayChart() {
	dt := time.Now().Add(time.Hour * -24 * 15)
	for i := 0; i < 15; i++ {
		st, et := tool.GetTodayStartEndUnix(dt.Add(time.Hour * 24 * time.Duration(i)))
		generateMchDayChart(st, et)
	}
}

// 商户每日报表
func mchDayChart() {
	//test
	//generateMchDayChart(0, time.Now().Unix())
	//testGenerateMchDayChart()

	unix := tool.GetStartDate(time.Now()).Unix()
	if CompareLastUnix(mchDayChartKey, unix) {
		log.Println("[ Mch][ Day][ Chart]: today chart is generated!")
		return
	}
	now := time.Now().Add(time.Hour * -24)
	st, et := tool.GetTodayStartEndUnix(now)
	generateMchDayChart(st, et)
	signHandled(mchDayChartKey, unix)
}

// 生成每日报表
func generateMchDayChart(start, end int64) {
	begin := 0
	size := 20
	var mchList []int
	tmp := 0
	dateStr := time.Unix(start, 0).Format("2006-01-02")
	// 清理数据
	appCtx.Db().ExecNonQuery(`DELETE FROM mch_day_chart WHERE date_str=?`, dateStr)
	// 开始统计数据
	for {
		mchList = []int{}
		appCtx.Db().Query("SELECT id FROM mch_merchant LIMIT ?,?", func(rs *sql.Rows) {
			for rs.Next() {
				rs.Scan(&tmp)
				mchList = append(mchList, tmp)
			}
		}, begin, size)
		wg := &sync.WaitGroup{}
		for _, v := range mchList {
			wg.Add(1)
			genDayChartForMch(wg, v, dateStr, start, end)
		}
		wg.Wait()
		if l := len(mchList); l == size {
			begin += l
		} else {
			break
		}
	}
}

func genDayChartForMch(wg *sync.WaitGroup, mchId int64, dateStr string, start int64, end int64) {
	defer wg.Done()
	c := &merchant.MchDayChart{
		MchId:   mchId,
		DateStr: dateStr,
		Date:    start,
	}
	//13826408857
	db := appCtx.Db()
	// 统计订单
	db.QueryRow(`SELECT COUNT(0),SUM(final_amount),COUNT(distinct buyer_id)
 FROM sale_sub_order where vendor_id=? AND create_time BETWEEN ? AND ?`, func(r *sql.Row) {
		r.Scan(&c.OrderNumber, &c.OrderAmount, &c.BuyerNumber)
	}, mchId, start, end)
	// 支付单汇总
	db.QueryRow(`SELECT COUNT(0),SUM(sale_sub_order.final_amount) FROM sale_sub_order
INNER JOIN pay_order ON pay_order.order_id = sale_sub_order.parent_order
where sale_sub_order.vendor_id=? AND pay_order.state = 1 AND pay_order.paid_time
 BETWEEN ? AND ?`, func(r *sql.Row) {
		r.Scan(&c.PaidNumber, &c.PaidAmount)
	}, mchId, start, end)
	// 今日已完成订单,应进账数量
	db.QueryRow(`SELECT COUNT(0),SUM(final_amount)
 FROM sale_sub_order where vendor_id=? AND state=?
 AND update_time BETWEEN ? AND ?`, func(r *sql.Row) {
		r.Scan(&c.CompleteOrders, &c.InAmount)
	}, mchId, order.StatCompleted, start, end)
	// 保存
	orm.Save(db.GetOrm(), c, 0)
	//log.Println("---", c, start, end)
}
