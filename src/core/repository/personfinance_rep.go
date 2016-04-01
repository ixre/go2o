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
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/personfinance"
	pf "go2o/src/core/domain/personfinance"
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

func (this *personFinanceRepository) GetPersonFinance(personId int) personfinance.IPersonFinance {
	return pf.NewPersonFinance(personId, this, this._accRep)
}

func (this *personFinanceRepository) GetRiseByTime(personId int, begin,
	end int64) []*personfinance.RiseDayInfo {
	list := []*personfinance.RiseDayInfo{}
	this._orm.Select(&list, "person_id=? AND unix_date BETWEEN ? AND ?", personId, begin, end)
	return list
}

func (this *personFinanceRepository) GetRiseValueByPersonId(id int) (
	*personfinance.RiseInfoValue, error) {
	e := &personfinance.RiseInfoValue{}
	err := this._orm.Get(id, e)
	return e, err
}

func (this *personFinanceRepository) SaveRiseInfo(v *personfinance.RiseInfoValue) (
	id int, err error) {
	if _, err = this.GetRiseValueByPersonId(v.PersonId); err == nil {
		_, _, err = this._orm.Save(v.PersonId, v)
	} else {
		_, _, err = this._orm.Save(nil, v)
	}
	return v.PersonId, err
}

// 获取日志
func (this *personFinanceRepository) GetRiseLog(personId, logId int) *personfinance.RiseLog {
	e := &personfinance.RiseLog{}
	if this._orm.GetBy(e, "person_id=? AND id=?", personId, logId) == nil {
		return e
	}
	return nil
}

// 保存日志
func (this *personFinanceRepository) SaveRiseLog(v *personfinance.RiseLog) (id int, err error) {
	if v.Id > 0 {
		_, _, err = this._orm.Save(v.Id, v)
	} else {
		_, _, err = this._orm.Save(nil, v)
		this._db.ExecScalar("SELECT MAX(id) FROM pf_riselog", &v.Id)
	}
	return v.Id, err
}

// 获取日志
func (this *personFinanceRepository) GetRiseLogs(personId int, date int64, riseType int) []*personfinance.RiseLog {
	list := []*personfinance.RiseLog{}
	this._orm.Select(&list, "person_id=? AND unix_date=? AND type=?", personId, date, riseType)
	return list
}
