/**
 * Copyright 2015 @ to2.net.
 * name : person_finance
 * author : jarryliu
 * date : 2016-04-02 10:35
 * description :
 * history :
 */
package daemon

import (
	"database/sql"
	"go2o/core/domain/interface/personfinance"
	"go2o/core/infrastructure/tool"
	"go2o/core/service/rsi"
	"log"
	"math"
	"sync"
	"time"
)

var (
	settleUnixKey string = "sys:go2o:d:pf:date"
)

func personFinanceSettle() {
	now := time.Now()
	//invokeSettle(now.Add(time.Hour * -24))
	unix := tool.GetStartDate(time.Now()).Unix()
	// 今日是否结算
	if CompareLastUnix(settleUnixKey, unix) {
		log.Println("[ PersonFinance][ Settle][ Info]:Today is settled!")
		return
	}
	invokeSettle(now)
	// 保存最新结算日期
	SetLastUnix(settleUnixKey, unix)
}

// 执行结算,结算时间为当天,
// 收益计算当天前一天收益,转入转出按当天计算
func invokeSettle(t time.Time) {
	b := time.Now()
	confirmTransferIn(t)                   //今天确认T+?前的转入
	settleRiseData(t.Add(time.Hour * -24)) //今天结算昨日的收益
	log.Println("[ PersonFinance][ Settle][ Success]:Total used",
		math.Floor(time.Now().Sub(b).Minutes()*100)/100, "minutes!")
}

// 确认转入数据
// 采用按ID分段,通过传入ID区间用多个gorouting进行处理.
func confirmTransferIn(t time.Time) {
	settleTime := t.AddDate(0, 0, -personfinance.RiseSettleTValue) // 倒推结算日
	unixDate := tool.GetStartDate(settleTime).Unix()
	begin := 0
	size := 20
	for {
		idArr := []int32{}
		err := conn.Query(`SELECT id FROM pf_riselog WHERE
		unix_date<= $1 AND type= $2 AND state= $3 LIMIT $5 OFFSET $4`,
			func(rows *sql.Rows) {
				var i int32
				for rows.Next() {
					rows.Scan(&i)
					if i > 0 {
						idArr = append(idArr, i)
					}
				}
			}, unixDate, personfinance.RiseTypeTransferIn,
			personfinance.RiseStateDefault, 0, size)
		if err != nil {
			log.Println("[ Error][ Transfer-Confirm]:", err.Error())
			break
		}
		// 将IdArr按指定size切片处理
		//wg := sync.WaitGroup{}
		for _, v := range idArr {
			//wg.Add(1)
			confirmTransferInByCursor(unixDate, v)
		}
		log.Println("[ PersonFinance][ RiseSettle][ Job]:begin:", begin,
			"; size:", size, "; len:", len(idArr), "; unix date =", unixDate)
		time.Sleep(time.Second / 4)
		if l := len(idArr); l == size {
			begin += l
		} else {
			break
		}
	}
}

// 分组确认转入数据
func confirmTransferInByCursor(unixDate int64, logId int32) {
	//log.Println(fmt.Sprintf("[SQL]: select * FROM pf_riselog WHERE id BETWEEN %d AND %d AND unix_date=%d AND type=%d AND state=%d ORDER BY id ",
	//	 idArr[0],idArr[len(idArr)-1], unixDate, personfinance.RiseTypeTransferIn,
	//	personfinance.RiseStateDefault))
	//time.Sleep(time.Second * 1)
	v := personfinance.RiseLog{}
	err := _orm.GetBy(&v, "id = $1 AND unix_date<= $2 AND type= $3 AND state= $4 ORDER BY id",
		logId, unixDate, personfinance.RiseTypeTransferIn,
		personfinance.RiseStateDefault)
	if err == nil {
		err = rsi.PersonFinanceService.CommitTransfer(v.PersonId, v.Id)
		if err != nil {
			log.Println("[ PersonFinance][ Transfer][ Fail]: person_id=",
				v.PersonId, "error=", err.Error())
			v.State = -1
			_orm.Save(v.Id, v) //标记为异常
		}
	}
}

// 结算增利数据,t为结算日
// 采用按ID分段,通过传入ID区间用多个gorouting进行处理.
func settleRiseData(settleDate time.Time) {
	settleUnix := tool.GetStartDate(settleDate).Unix() //结算日期
	begin := 0
	size := 20
	for {
		idArr := []int64{}
		err := conn.Query(`SELECT person_id FROM pf_riseinfo WHERE
            settlement_amount > 0 AND settled_date < $1 LIMIT $3 OFFSET $2`,
			func(rows *sql.Rows) {
				var i int64
				for rows.Next() {
					rows.Scan(&i)
					if i > 0 {
						idArr = append(idArr, i)
					}
				}
			}, settleUnix, 0, size)
		if err != nil {
			log.Println("[ Error][ Rise-Settle]:", err.Error())
			break
		}
		wg := sync.WaitGroup{}
		for _, personId := range idArr {
			wg.Add(1)
			go riseGroupSettle(&wg, settleUnix, personId)
		}
		wg.Wait()
		log.Println("[ PersonFinance][ RiseSettle][ Job]:begin:", begin,
			"; size:", size, "; len:", len(idArr), "; unix date =", settleUnix)
		time.Sleep(time.Second / 4)
		if l := len(idArr); l == size {
			begin += l
		} else {
			break
		}
	}
}

// 结算每日数据
func riseGroupSettle(wg *sync.WaitGroup, settleUnix int64, personId int64) {
	err := rsi.PersonFinanceService.RiseSettleByDay(personId, settleUnix,
		personfinance.RiseDayRatioProvider(personId))
	if err != nil {
		log.Println("[ PersonFinance][ Settle][ Fail]: person_id=",
			personId, "error=", err.Error())
	}
	wg.Done()
}
