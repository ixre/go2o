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
	"go2o/core/domain/interface/member"
	"go2o/core/dto"
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
func (m *MemberQuery) GetMemberList(ids []int) []*dto.MemberSummary {
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
				m.update_time FROM mm_member m INNER JOIN mm_level lv
				ON m.level = lv.id INNER JOIN mm_account a ON
				 a.member_id = m.id AND m.id IN(%s) order by field(m.id,%s)`, inStr, inStr)
		m.Connector.GetOrm().SelectByQuery(&list, query)
	}
	return list
}

// 获取返现记录
func (m *MemberQuery) QueryBalanceLog(memberId, begin, end int,
	where, orderBy string) (num int, rows []map[string]interface{}) {

	d := m.Connector

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
	}, memberId, begin, end-begin)

	return num, rows
}

// 获取最近的余额变动信息
func (m *MemberQuery) GetLatestBalanceInfoByKind(memberId int, kind int) *member.BalanceInfo {
	var info = new(member.BalanceInfo)
	if err := m.GetOrm().GetBy(info, "member_id=? AND kind=? ORDER BY create_time DESC",
		memberId, kind); err == nil {
		return info
	}
	return nil
}

// 筛选会员根据用户或者手机
func (m *MemberQuery) FilterMemberByUsrOrPhone(key string) []*dto.SimpleMember {
	qp := "%" + key + "%"
	var list []*dto.SimpleMember = make([]*dto.SimpleMember, 0)
	var id int
	var usr, name, phone string
	m.Query(`SELECT id,usr,name,phone FROM mm_member WHERE
		usr LIKE ? OR name LIKE ? OR phone LIKE ?`, func(rows *sql.Rows) {
		for rows.Next() {
			rows.Scan(&id, &usr, &name, &phone)
			list = append(list, &dto.SimpleMember{
				Id:    id,
				User:  usr,
				Name:  name,
				Phone: phone,
			})
		}
	}, qp, qp, qp)
	return list
}

// 会员推广排名
func (m *MemberQuery) GetMemberInviRank(merchantId int, allTeam bool, levelComp string, level int,
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

	m.Query(fmt.Sprintf(`SELECT id,usr,name,invi_num,all_num,reg_time FROM ( SELECT m.*,
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
	}, startTime, endTime, startTime, endTime, startTime, endTime, merchantId, 1, num)

	return list
}

// 查询有邀请关系的会员数量
func (m *MemberQuery) GetReferNum(memberId int, layer int) int {
	total := 0
	keyword := fmt.Sprintf("''r%d'':%d", layer, memberId)
	where := "refer_str LIKE '%" + keyword + ",%' OR refer_str LIKE '%" + keyword + "}'"
	m.ExecScalar("SELECT COUNT(0) FROM mm_relation WHERE "+where, &total)
	return total
}
