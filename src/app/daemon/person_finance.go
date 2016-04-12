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
	"go2o/src/core"
	"github.com/garyburd/redigo/redis"
)

const batGroupSize int = 30 //跑批每组数量
var (
	settleUnixKey string = "go2o:d:pf:settled_unix"
)

func personFinanceSettle() {
	now := time.Now()
	//invokeSettle(now.Add(time.Hour * -24))
	unix := tool.GetStartDate(time.Now()).Unix()
	if todayIsSettled(unix){
		log.Println("[ PersonFinance][ Settle][ Info]:Today is settled!")
		return
	}
	invokeSettle(now)
	saveLatestSettleUnix(unix)

}

// 今日是否结算
func todayIsSettled(unix int64)bool{
	conn := core.GetRedisConn()
	defer conn.Close()
	unix2,err := redis.Int(conn.Do("GET", settleUnixKey))
	if err != nil{
		return false
	}
	return unix ==int64(unix2)
}

// 保存最新结算日期
func saveLatestSettleUnix(unix int64){
	conn := core.GetRedisConn()
	defer conn.Close()
	conn.Do("SET", settleUnixKey,unix)
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
	var err error
	settleTime := t.AddDate(0, 0, -personfinance.RiseSettleTValue) // 倒推结算日
	unixDate := tool.GetStartDate(settleTime).Unix()
	cursor := 0   // 游标,每次从db中取条数
	setupNum := 0 //步骤编号
	for {
		// 获取前1000条记录到IdArr
		idArr := []int{}
		i := 0
		err = _db.Query("SELECT id FROM pf_riselog WHERE unix_date=? AND type=? AND state=? LIMIT 0,?",
			func(rows *sql.Rows) {
				for rows.Next() {
					rows.Scan(&i)
					if i > 0 {
						idArr = append(idArr, i)
					}
				}
			}, unixDate, personfinance.RiseTypeTransferIn,
			personfinance.RiseStateDefault, 1000)
		if err != nil {
			log.Println("[ Error][ Transfer-Confirm]:", err.Error())
			break
		}
		if len(idArr) == 0 {
			break
		}

		setupNum += 1
		log.Println("[ PersonFinance][ Transfer][ Job]:Setup", setupNum,
			"; Total", len(idArr), "records! unix date =", unixDate)

		// 将IdArr按指定size切片处理
		wg := sync.WaitGroup{}
		for cursor < len(idArr) {
			var splitIdArr []int
			if cursor+batGroupSize < len(idArr) {
				splitIdArr = idArr[cursor : cursor+batGroupSize]
			} else {
				splitIdArr = idArr[cursor:]
			}
			go confirmTransferInByCursor(&wg, unixDate, splitIdArr)
			cursor += batGroupSize
			wg.Add(1)
			time.Sleep(time.Microsecond * 1000)
			//log.Println("[Output]- ", splitIdArr[0], splitIdArr[len(splitIdArr)-1],len(splitIdArr))
		}
		wg.Wait()
		cursor = 0 //重置游标
	}
}

// 分组确认转入数据
func confirmTransferInByCursor(wg *sync.WaitGroup, unixDate int64, idArr []int) {

	//log.Println(fmt.Sprintf("[SQL]: select * FROM pf_riselog WHERE id BETWEEN %d AND %d AND unix_date=%d AND type=%d AND state=%d ORDER BY id ",
	//	 idArr[0],idArr[len(idArr)-1], unixDate, personfinance.RiseTypeTransferIn,
	//	personfinance.RiseStateDefault))
	//time.Sleep(time.Second * 1)

	list := make([]*personfinance.RiseLog, 0)
	_orm.Select(&list, "id BETWEEN ? AND ? AND unix_date=? AND type=? AND state=? ORDER BY id",
		idArr[0], idArr[len(idArr)-1], unixDate, personfinance.RiseTypeTransferIn,
		personfinance.RiseStateDefault)
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
// 采用按ID分段,通过传入ID区间用多个gorouting进行处理.
func settleRiseData(settleDate time.Time) {
	var err error
	settleUnix := tool.GetStartDate(settleDate).Unix() //结算日期
	cursor := 0
	setupNum := 0 //步骤编号

	for {
		idArr := []int{}
		i := 0
		err = _db.Query("SELECT person_id FROM pf_riseinfo WHERE settlement_amount > 0 AND settled_date < ? LIMIT 0,?",
			func(rows *sql.Rows) {
				for rows.Next() {
					rows.Scan(&i)
					if i > 0 {
						idArr = append(idArr, i)
					}
				}
			}, settleUnix, 1000)
		if err != nil {
			log.Println("[ Error][ Rise-Settle]:", err.Error())
			break
		}
		if len(idArr) == 0 {
			break
		}

		setupNum += 1
		log.Println("[ PersonFinance][ RiseSettle][ Job]:Setup ", setupNum,
			" ; Total ", len(idArr), "records! unix date =", settleUnix)

		wg := sync.WaitGroup{}
		for cursor < len(idArr) {
			var splitIdArr []int
			if cursor+batGroupSize < len(idArr) {
				splitIdArr = idArr[cursor : cursor+batGroupSize]
			} else {
				splitIdArr = idArr[cursor:]
			}
			go riseGroupSettle(&wg, settleUnix, splitIdArr)
			cursor += batGroupSize
			wg.Add(1)
			time.Sleep(time.Microsecond * 1000)
			log.Println("[Output]- ", splitIdArr[0], splitIdArr[len(splitIdArr)-1], len(splitIdArr))
		}
		wg.Wait()
		cursor = 0 //重置游标
	}
}

// 分组确认转入数据
func riseGroupSettle(wg *sync.WaitGroup, settleUnix int64, personIdArr []int) {
	ds := dps.PersonFinanceService
	for _, id := range personIdArr {
		if err := ds.RiseSettleByDay(id, settleUnix, personfinance.RiseDayRatioProvider(id)); err != nil {
			log.Println("[ PersonFinance][ Settle][ Fail]: person_id=", id, "error=", err.Error())
		}
	}
	wg.Done()
}
