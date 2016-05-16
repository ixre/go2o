/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:20
 * description :
 * history :
 */
package query

import (
	"database/sql"
	"fmt"
	"github.com/jsix/gof/db"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/dto"
	"strconv"
	"strings"
)

type MemberQuery struct {
	db.Connector
}

func NewMemberQuery(c db.Connector) *MemberQuery {
	return &MemberQuery{c}
}

// 获取会员列表
func (this *MemberQuery) GetMemberList(partnerId int, ids []int) []*dto.MemberSummary {
	list := []*dto.MemberSummary{}
	strIds := make([]string, len(ids))
	for i, v := range ids {
		strIds[i] = strconv.Itoa(v)
	}
	if len(ids) > 0 {
		inStr := strings.Join(strIds, ",") // order by field(field,val1,val2,val3)按IN的顺序排列
		query := fmt.Sprintf(`SELECT m.id,m.usr,m.name,m.avatar,m.exp,m.level,
				lv.name as level_name,a.integral,a.balance,a.present_balance,
				a.grow_balance,a.grow_amount,a.grow_earnings,a.grow_total_earnings,
				m.update_time FROM mm_member m INNER JOIN pt_member_level lv
				ON m.level = lv.value INNER JOIN mm_account a ON
				 a.member_id = m.id WHERE lv.merchant_id=? AND m.id IN(%s) order by field(m.id,%s)`, inStr, inStr)
		this.Connector.GetOrm().SelectByQuery(&list, query, partnerId)
	}
	return list
}

// 获取返现记录
func (this *MemberQuery) QueryBalanceLog(memberId, page, size int,
	where, orderBy string) (num int, rows []map[string]interface{}) {

	d := this.Connector

	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	}
	d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM mm_balance_info bi
	 	INNER JOIN mm_member m ON m.id=bi.member_id
			WHERE bi.member_id=? %s`, where), &num, memberId)

	sqlLine := fmt.Sprintf(`SELECT bi.* FROM mm_balance_info bi
			INNER JOIN mm_member m ON m.id=bi.member_id
			WHERE member_id=? %s %s LIMIT ?,?`,
		where, orderBy)

	d.Query(sqlLine, func(_rows *sql.Rows) {
		rows = db.RowsToMarshalMap(_rows)
		_rows.Close()
	}, memberId, (page-1)*size, size)

	return num, rows
}

// 查询分页订单
func (this *MemberQuery) QueryPagerOrder(memberId, page, size int,
	where, orderBy string) (num int, rows []map[string]interface{}) {

	d := this.Connector

	if where != "" {
		where = "AND " + where
	}
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy
	} else {
		orderBy = " ORDER BY update_time DESC,create_time desc "
	}

	d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM pt_order WHERE
	 		member_id=? %s`, where), &num, memberId)

	d.Query(fmt.Sprintf(` SELECT id,
			order_no,
			member_id,
			merchant_id,
			shop_id,
			replace(items_info,'\n','<br />') as items_info,
			total_fee,
			fee,
			pay_fee,
			payment_opt,
			is_paid,
			note,
			status,
			paid_time,
            create_time,
            deliver_time,
            update_time
            FROM pt_order WHERE member_id=? %s %s LIMIT ?,?`,
		where, orderBy),
		func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
			_rows.Close()
		}, memberId, (page-1)*size, size)

	return num, rows
}

// 获取最近的余额变动信息
func (this *MemberQuery) GetLatestBalanceInfoByKind(memberId int, kind int) *member.BalanceInfoValue {
	var info = new(member.BalanceInfoValue)
	if err := this.GetOrm().GetBy(info, "member_id=? AND kind=? ORDER BY create_time DESC",
		memberId, kind); err == nil {
		return info
	}
	return nil
}

// 筛选会员根据用户或者手机
func (this *MemberQuery) FilterMemberByUsrOrPhone(partnerId int, key string) []*dto.SimpleMember {
	qp := "%" + key + "%"
	var list []*dto.SimpleMember = make([]*dto.SimpleMember, 0)
	var id int
	var usr, name, phone string
	this.Query(`SELECT id,usr,name,phone FROM mm_member INNER JOIN mm_relation ON
		mm_relation.member_id = mm_member.id WHERE mm_relation.reg_merchant_id=?
		AND usr LIKE ? OR name LIKE ?`, func(rows *sql.Rows) {
		for rows.Next() {
			rows.Scan(&id, &usr, &name, &phone)
			list = append(list, &dto.SimpleMember{
				Id:    id,
				User:  usr,
				Name:  name,
				Phone: phone,
			})
		}
	}, partnerId, qp, qp)
	return list
}

// 会员推广排名
func (this *MemberQuery) GetMemberInviRank(partnerId int, allTeam bool, levelComp string, level int,
	startTime int64, endTime int64, num int) []*dto.RankMember {
	var list []*dto.RankMember = make([]*dto.RankMember, 0)
	var id int
	var usr, name string
	var inviNum, totalNum, regTime int
	var rank int = 0

	var sortField string = "t.all_num DESC"
	if !allTeam {
		sortField = "t.invi_num DESC"
	}

	var levelCompStr string = fmt.Sprintf("%s%d", levelComp, level)
	//{level_comp}{level_value}

	this.Query(fmt.Sprintf(`SELECT id,usr,name,invi_num,all_num,reg_time FROM ( SELECT m.*,
 (SELECT COUNT(0) FROM mm_relation r INNER JOIN mm_member m1 ON m1.id = r.member_id WHERE
  (m1.level%s) AND r.invi_member_id = m.id
	AND r.reg_merchant_id=rl.reg_merchant_id  AND m1.reg_time BETWEEN
  ? AND ? ) as invi_num,
	((SELECT COUNT(0) FROM mm_relation r INNER JOIN mm_member m1 ON m1.id = r.member_id WHERE
  (m1.level%s) AND r.invi_member_id = m.id
	AND r.reg_merchant_id=rl.reg_merchant_id AND m1.reg_time BETWEEN
  ? AND ? )+
 (SELECT COUNT(0) FROM mm_relation r INNER JOIN mm_member m1
  ON m1.id = r.member_id WHERE (m1.level%s) AND invi_member_id IN
	(SELECT member_id FROM mm_relation r INNER JOIN mm_member m1 ON m1.id = r.member_id WHERE
  (m1.level%s) AND r.invi_member_id =
    m.id AND r.reg_merchant_id=rl.reg_merchant_id AND m1.reg_time BETWEEN
  ? AND ? ))) as all_num
 FROM mm_member m INNER JOIN mm_relation rl ON m.id= rl.member_id
 WHERE rl.reg_merchant_id = ? AND state= ?) t ORDER BY %s,t.reg_time asc
 LIMIT 0,?`, levelCompStr, levelCompStr, levelCompStr, levelCompStr, sortField), func(rows *sql.Rows) {
		for rows.Next() {
			rows.Scan(&id, &usr, &name, &inviNum, &totalNum, &regTime)
			rank++
			list = append(list, &dto.RankMember{
				Id:       id,
				Usr:      usr,
				Name:     name,
				RankNum:  rank,
				InviNum:  inviNum,
				TotalNum: totalNum,
				RegTime:  regTime,
			})
		}
	}, startTime, endTime, startTime, endTime, startTime, endTime, partnerId, 1, num)

	return list
}
