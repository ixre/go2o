/**
 * Copyright 2015 @ z3q.net.
 * name : person_finance
 * author : jarryliu
 * date : 2016-04-02 10:35
 * description :
 * history :
 */
package daemon

import (
	"database/sql"
	"go2o/src/core/domain/interface/personfinance"
	"go2o/src/core/infrastructure/tool"
	"go2o/src/core/service/dps"
	"log"
	"sync"
	"time"
)

func personFinanceSettle() {
	confirmTransferIn(time.Now())
	settleRiseData(time.Now().Add(time.Hour * -24)) //结算昨日的收益
}

// 确认转入数据
func confirmTransferIn(t time.Time) {
	var err error
	settleTime := t.AddDate(0, 0, -personfinance.RiseSettleTValue) // 倒推结算日
	unixDate := tool.GetStartDate(settleTime).Unix()
	total := 0  //总数
	cursor := 0 // 游标,每次从db中取条数
	const size int = 50

	err = _db.ExecScalar("SELECT COUNT(0) FROM pf_riselog WHERE unix_date=? AND type=? AND state=?",
		&total, unixDate, personfinance.RiseTypeTransferIn, personfinance.RiseStateDefault)
	if err != nil {
		log.Println("[ Error][ Transfer-Confirm]:", err.Error())
		return
	}
	log.Println("[ PersonFinance][ Transfer][ Job]:Total ", total, "records! unix date =", unixDate)
	wg := &sync.WaitGroup{}
	for cursor < total {
		go confirmTransferInByCursor(wg, unixDate, cursor, size)
		cursor += size
		wg.Add(1)
	}
	wg.Wait()
}

// 分组确认转入数据
func confirmTransferInByCursor(wg *sync.WaitGroup, unixDate int64, cursor, size int) {
	list := make([]*personfinance.RiseLog, 0)
	_orm.Select(&list, "unix_date=? AND type=? AND state=? ORDER BY id LIMIT ?,?",
		unixDate, personfinance.RiseTypeTransferIn,
		personfinance.RiseStateDefault, cursor, size)
	ds := dps.PersonFinanceService
	for _, v := range list {
		if err := ds.CommitTransfer(v.PersonId, v.Id); err != nil {
			log.Println("[ PersonFinance][ Transfer][ Fail]:", err.Error())
			v.State = -1
			_orm.Save(v.Id, v) //标记为异常
		}
	}
	wg.Done()
}

// 结算增利数据
func settleRiseData(t time.Time) {
	var err error
	unixDate := tool.GetStartDate(t).Unix()
	total := 0  //总数
	cursor := 0 // 游标,每次从db中取条数
	const size int = 50

	err = _db.ExecScalar("SELECT COUNT(0) FROM pf_riseinfo WHERE balance > 0", &total)
	if err != nil {
		log.Println("[ Error][ Rise-Settle]:", err.Error())
		return
	}
	log.Println("[ PersonFinance][ RiseSettle][ Job]:Total ", total, "records! unix date =", unixDate)
	wg := &sync.WaitGroup{}
	for cursor < total {
		go riseGroupSettle(wg, unixDate, cursor, size)
		cursor += size
		wg.Add(1)
	}
	wg.Wait()
}

// 分组确认转入数据
func riseGroupSettle(wg *sync.WaitGroup, unixDate int64, cursor, size int) {
	list := []int{}
	var id int
	_db.Query("SELECT person_id FROM pf_riseinfo WHERE balance > 0 LIMIT ?,?",
		func(rows *sql.Rows) {
			for rows.Next() {
				rows.Scan(&id)
				list = append(list, id)
			}
		}, cursor, size)
	ds := dps.PersonFinanceService
	for _, id := range list {
		if err := ds.RiseSettleByDay(id, personfinance.RiseDayRatioProvider(id)); err != nil {
			log.Println("[ PersonFinance][ Settle][ Fail]:", err.Error())
		}
	}
	wg.Done()
}
