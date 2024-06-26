/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2014-01-08 21:35
 * description :
 * history :
 */

package daemon

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/order"
	tool "github.com/ixre/go2o/core/infrastructure/util"
	"github.com/ixre/gof/db/orm"
)

var (
	mchDayChartKey = "cron:go2o:d:mch:day-chart-unix"
)

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
	st, et := tool.GetStartEndUnix(now)
	generateMchDayChart(st, et)
	signHandled(mchDayChartKey, unix)
}

// 生成每日报表
func generateMchDayChart(start, end int64) {
	begin := 0
	size := 20
	var mchList []int64
	var tmp int64
	dateStr := time.Unix(start, 0).Format("2006-01-02")
	// 清理数据
	appCtx.Db().ExecNonQuery(`DELETE FROM mch_day_chart WHERE date_str= $1`, dateStr)
	// 开始统计数据
	for {
		mchList = []int64{}
		appCtx.Db().Query("SELECT id FROM mch_merchant LIMIT $2 OFFSET $1", func(rs *sql.Rows) {
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
	db.QueryRow(`SELECT COUNT(1),SUM(final_amount),COUNT(distinct buyer_id)
 FROM sale_sub_order where vendor_id= $1 AND create_time BETWEEN $2 AND $3`, func(r *sql.Row) error {
		return r.Scan(&c.OrderNumber, &c.OrderAmount, &c.BuyerNumber)
	}, mchId, start, end)
	// 支付单汇总
	db.QueryRow(`SELECT COUNT(1),SUM(sale_sub_order.final_amount) FROM sale_sub_order
INNER JOIN pay_order ON pay_order.order_id = sale_sub_order.order_id
where sale_sub_order.vendor_id= $1 AND pay_order.state = 1 AND pay_order.paid_time
 BETWEEN $2 AND $3`, func(r *sql.Row) error {
		return r.Scan(&c.PaidNumber, &c.PaidAmount)
	}, mchId, start, end)
	// 今日已完成订单,应进账数量
	db.QueryRow(`SELECT COUNT(1),SUM(final_amount)
 FROM sale_sub_order where vendor_id= $1 AND state= $2
 AND update_time BETWEEN $3 AND $4`, func(r *sql.Row) error {
		return r.Scan(&c.CompleteOrders, &c.InAmount)
	}, mchId, order.StatCompleted, start, end)
	// 保存
	orm.Save(_orm, c, 0)
	//log.Println("---", c, start, end)
}
