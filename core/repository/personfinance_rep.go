/**
 * Copyright 2015 @ z3q.net.
 * name : personfinance_rep
 * author : jarryliu
 * date : 2016-04-01 09:30
 * description :
 * history :
 */
package repository

import (
	"github.com/jsix/gof/db"
	"github.com/jsix/gof/db/orm"
	"go2o/core/domain/interface/member"
	"go2o/core/domain/interface/personfinance"
	pf "go2o/core/domain/personfinance"
)

var _ personfinance.IPersonFinanceRepository = new(personFinanceRepository)

type personFinanceRepository struct {
	_db     db.Connector
	_orm    orm.Orm
	_accRep member.IMemberRep
}

func NewPersonFinanceRepository(conn db.Connector, mRep member.IMemberRep) personfinance.IPersonFinanceRepository {
	return &personFinanceRepository{
		_db:     conn,
		_orm:    conn.GetOrm(),
		_accRep: mRep,
	}
}

func (p *personFinanceRepository) GetPersonFinance(personId int64) personfinance.IPersonFinance {
	return pf.NewPersonFinance(personId, p, p._accRep)
}

func (p *personFinanceRepository) GetRiseByTime(personId int64, begin,
	end int64) []*personfinance.RiseDayInfo {
	list := []*personfinance.RiseDayInfo{}
	p._orm.Select(&list, "person_id=? AND unix_date BETWEEN ? AND ?", personId, begin, end)
	return list
}

func (p *personFinanceRepository) GetRiseValueByPersonId(id int64) (
	*personfinance.RiseInfoValue, error) {
	e := &personfinance.RiseInfoValue{}
	err := p._orm.Get(id, e)
	return e, err
}

func (p *personFinanceRepository) SaveRiseInfo(v *personfinance.RiseInfoValue) (
	id int, err error) {
	if _, err = p.GetRiseValueByPersonId(v.PersonId); err == nil {
		_, _, err = p._orm.Save(v.PersonId, v)
	} else {
		_, _, err = p._orm.Save(nil, v)
	}
	return v.PersonId, err
}

// 获取日志
func (p *personFinanceRepository) GetRiseLog(personId, logId int64) *personfinance.RiseLog {
	e := &personfinance.RiseLog{}
	if p._orm.GetBy(e, "person_id=? AND id=?", personId, logId) == nil {
		return e
	}
	return nil
}

// 保存日志
func (p *personFinanceRepository) SaveRiseLog(v *personfinance.RiseLog) (int64, error) {
	if v.Id > 0 {
		_, _, err = p._orm.Save(v.Id, v)
	} else {
		_, _, err = p._orm.Save(nil, v)
		p._db.ExecScalar("SELECT MAX(id) FROM pf_riselog", &v.Id)
	}
	return v.Id, err
}

// 获取日志
func (p *personFinanceRepository) GetRiseLogs(personId int64, date int64, riseType int) []*personfinance.RiseLog {
	list := []*personfinance.RiseLog{}
	p._orm.Select(&list, "person_id=? AND unix_date=? AND type=?", personId, date, riseType)
	return list
}

// 保存每日收益
func (p *personFinanceRepository) SaveRiseDayInfo(v *personfinance.RiseDayInfo) (int64, error) {
	if v.Id > 0 {
		_, _, err = p._orm.Save(v.Id, v)
	} else {
		_, _, err = p._orm.Save(nil, v)
		p._db.ExecScalar("SELECT MAX(id) FROM pf_riseday", &v.Id)
	}
	return v.Id, err
}
