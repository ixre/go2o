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
	"math"
	"sync"
	"time"
)

func personFinanceSettle() {
	b := time.Now()
	confirmTransferIn(time.Now()) //今天确认T+?前的转入
	settleRiseData() //今天结算昨日的收益
	log.Println("[ PersonFinance][ Settle][ Success]:Total used",
		math.Floor(time.Now().Sub(b).Minutes()*100)/100, "minutes!")
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
		time.Sleep(time.Second)
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
			log.Println("[ PersonFinance][ Transfer][ Fail]: person_id=", v.PersonId, "error=", err.Error())
			v.State = -1
			_orm.Save(v.Id, v) //标记为异常
		}
	}
	wg.Done()
}

// 结算增利数据,t为结算日
func settleRiseData() {
	var err error
	dt := time.Now().Add(time.Hour * -24)
	settleDate := tool.GetStartDate(dt).Unix() //结算日期
	total := 0  //总数
	cursor := 0 // 游标,每次从db中取条数
	const size int = 50

	err = _db.ExecScalar("SELECT COUNT(0) FROM pf_riseinfo WHERE balance > 0 AND settled_date < ? ",
		&total, settleDate)
	if err != nil {
		log.Println("[ Error][ Rise-Settle]:", err.Error())
		return
	}
	log.Println("[ PersonFinance][ RiseSettle][ Job]:Total ", total, "records! unix date =", settleDate)
	wg := &sync.WaitGroup{}
	for cursor < total {
		go riseGroupSettle(wg, settleDate, cursor, size)
		cursor += size
		wg.Add(1)
		time.Sleep(time.Second)
	}
	wg.Wait()
}

// 分组确认转入数据
func riseGroupSettle(wg *sync.WaitGroup, settleDate int64, cursor, size int) {
	list := []int{}
	var id int
	_db.Query("SELECT person_id FROM pf_riseinfo WHERE balance > 0 AND settled_date < ? LIMIT ?,?",
		func(rows *sql.Rows) {
			for rows.Next() {
				rows.Scan(&id)
				list = append(list, id)
			}
		}, settleDate,cursor, size)
	ds := dps.PersonFinanceService
	for _, id := range list {
		if err := ds.RiseSettleByDay(id, personfinance.RiseDayRatioProvider(id)); err != nil {
			log.Println("[ PersonFinance][ Settle][ Fail]: person_id=", id, "error=", err.Error())
		}
	}
	wg.Done()
}
