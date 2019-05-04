/**
 * Copyright 2015 @ z3q.net.
 * name : personfinance_rep
 * author : jarryliu
 * date : 2016-04-01 09:30
 * description :
 * history :
 */
package repos

import (
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/personfinance"
	pf "go2o/core/domain/personfinance"
)

var _ personfinance.IPersonFinanceRepository = new(personFinanceRepository)

type personFinanceRepository struct {
	_db      db.Connector
	_orm     orm.Orm
	_accRepo member.IMemberRepo
}

func NewPersonFinanceRepository(conn db.Connector, mRepo member.IMemberRepo) personfinance.IPersonFinanceRepository {
	return &personFinanceRepository{
		_db:      conn,
		_orm:     conn.GetOrm(),
		_accRepo: mRepo,
	}
}

func (p *personFinanceRepository) GetPersonFinance(personId int64) personfinance.IPersonFinance {
	return pf.NewPersonFinance(personId, p, p._accRepo)
}

func (p *personFinanceRepository) GetRiseByTime(personId int64, begin,
	end int64) []*personfinance.RiseDayInfo {
	var list []*personfinance.RiseDayInfo
	p._orm.Select(&list, "person_id= $1 AND unix_date BETWEEN $2 AND $3", personId, begin, end)
	return list
}

func (p *personFinanceRepository) GetRiseValueByPersonId(id int64) (
	*personfinance.RiseInfoValue, error) {
	e := &personfinance.RiseInfoValue{}
	err := p._orm.Get(id, e)
	return e, err
}

func (p *personFinanceRepository) SaveRiseInfo(v *personfinance.RiseInfoValue) (int, error) {
	var err error
	if _, err = p.GetRiseValueByPersonId(v.PersonId); err == nil {
		_, _, err = p._orm.Save(v.PersonId, v)
	} else {
		_, _, err = p._orm.Save(nil, v)
	}
	return int(v.PersonId), err
}

// 获取日志
func (p *personFinanceRepository) GetRiseLog(personId int64, logId int32) *personfinance.RiseLog {
	e := &personfinance.RiseLog{}
	if p._orm.GetBy(e, "person_id= $1 AND id= $2", personId, logId) == nil {
		return e
	}
	return nil
}

// 保存日志
func (p *personFinanceRepository) SaveRiseLog(v *personfinance.RiseLog) (int32, error) {
	return orm.I32(orm.Save(p._db.GetOrm(), v, int(v.Id)))
}

// 获取日志
func (p *personFinanceRepository) GetRiseLogs(personId int64, date int64, riseType int) []*personfinance.RiseLog {
	list := []*personfinance.RiseLog{}
	p._orm.Select(&list, "person_id= $1 AND unix_date= $2 AND type= $3", personId, date, riseType)
	return list
}

// 保存每日收益
func (p *personFinanceRepository) SaveRiseDayInfo(v *personfinance.RiseDayInfo) (int32, error) {
	return orm.I32(orm.Save(p._db.GetOrm(), v, int(v.Id)))
}
