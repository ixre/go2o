package dao

import (
	"com/ording"
	"com/ording/entity"
	"database/sql"
	"fmt"
	"ops/cf/db"
)

type memberDao struct {
	db.Connector
}

func (this *memberDao) GetIncomeLog(memberId, page, size int,
	where, orderby string) (num int, rows []map[string]interface{}) {

	d := this.Connector

	if where != "" {
		where = "WHERE " + where
	}
	if orderby != "" {
		orderby = "ORDER BY " + orderby
	}
	d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0)
			FROM mm_income_log l INNER JOIN mm_member m ON m.id=l.member_id
			WHERE member_id=? %s`, where), &num, memberId)

	d.Query(fmt.Sprintf(`SELECT l.*,
			date_format(l.record_time,%s) as record_time,
			convert(l.fee,CHAR(10)) as fee
			FROM mm_income_log l INNER JOIN mm_member m ON m.id=l.member_id
			WHERE member_id=? %s %s LIMIT ?,?`, "'%Y-%m-%d %T'",
		where, orderby),
		func(_rows *sql.Rows) {
			rows = db.ConvRowsToMapForJson(_rows)
			_rows.Close()
		}, memberId, (page-1)*size, size)

	return num, rows
}

//@deparend
//验证用户密码
func (this *memberDao) Verify(usr, pwd string) bool {
	var id int
	encPwd := ording.EncodeMemberPwd(usr, pwd)
	if err := this.Connector.ExecScalar("SELECT id FROM mm_member WHERE usr=? AND pwd=?", &id, usr, encPwd); err != nil {
		return false
	}
	return id != 0
}

/*********** 收货地址 ***********/
func (this *memberDao) GetDeliverAddrs(memberId int) []entity.DeliverAddress {
	addresses := []entity.DeliverAddress{}
	this.Connector.GetOrm().Select(&addresses, entity.DeliverAddress{}, fmt.Sprintf("member_id=%d", memberId))
	return addresses
}

//获取配送地址
func (this *memberDao) GetDeliverAddrById(memberId, deliverId int) *entity.DeliverAddress {
	addr := new(entity.DeliverAddress)
	if this.Connector.GetOrm().Get(addr, deliverId) == nil && addr.Mid == memberId {
		return addr
	}
	return nil
}

//保存配送地址
func (this *memberDao) SaveDeliverAddr(e *entity.DeliverAddress) (int, error) {
	orm := this.Connector.GetOrm()
	if e.Id <= 0 {
		//多行字符用
		_, id, err := orm.Save(nil, e)
		return int(id), err
	} else {
		_, _, err := orm.Save(e.Id, e)
		return e.Id, err
	}
}

//删除配送地址
func (this *memberDao) DeleteDeliverAddr(memberId int, deliverAddrId int) error {
	_, err := this.Connector.ExecNonQuery(
		"DELETE FROM mm_deliver_addr WHERE mid=? AND id=?",
		memberId, deliverAddrId)
	return err
}
