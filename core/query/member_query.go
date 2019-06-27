/**
 * Copyright 2014 @ to2.net.
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
	"github.com/ixre/gof/db"
	"github.com/labstack/gommon/log"
	"go2o/core/domain/interface/member"
	"go2o/core/dto"
	"go2o/core/infrastructure/format"
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
func (m *MemberQuery) GetMemberList(ids []int64) []*dto.MemberSummary {
	list := []*dto.MemberSummary{}
	strIds := make([]string, len(ids))
	for i, v := range ids {
		strIds[i] = strconv.Itoa(int(v))
	}
	if len(ids) > 0 {
		inStr := strings.Join(strIds, ",") // order by field(field,val1,val2,val3)按IN的顺序排列
		query := fmt.Sprintf(`SELECT m.id,m.user,m.name,m.avatar,m.exp,m.level,
				lv.name as level_name,a.integral,a.balance,a.wallet_balance,
				a.grow_balance,a.grow_amount,a.grow_earnings,a.grow_total_earnings,
				m.update_time FROM mm_member m INNER JOIN mm_level lv
				ON m.level = lv.id INNER JOIN mm_account a ON
				 a.member_id = m.id AND m.id IN(%s) order by field(m.id,%s)`, inStr, inStr)
		m.Connector.GetOrm().SelectByQuery(&list, query)
	}
	return list
}

// 获取账户余额分页记录
func (m *MemberQuery) PagedBalanceAccountLog(memberId int64, begin, end int,
	where, orderBy string) (num int, rows []map[string]interface{}) {
	d := m.Connector
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy + ",bi.id DESC"
	}
	d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM mm_balance_log bi
	 	INNER JOIN mm_member m ON m.id=bi.member_id
			WHERE bi.member_id= $1 %s`, where), &num, memberId)

	sqlLine := fmt.Sprintf(`SELECT bi.* FROM mm_balance_log bi
			INNER JOIN mm_member m ON m.id=bi.member_id
			WHERE member_id= $1 %s %s LIMIT $3 OFFSET $2`,
		where, orderBy)

	d.Query(sqlLine, func(_rows *sql.Rows) {
		rows = db.RowsToMarshalMap(_rows)
	}, memberId, begin, end-begin)

	return num, rows
}

// 获取账户余额分页记录
func (m *MemberQuery) PagedWalletAccountLog(memberId int64, begin, end int,
	where, orderBy string) (num int, rows []map[string]interface{}) {
	d := m.Connector
	if orderBy != "" {
		orderBy = "ORDER BY " + orderBy + ",bi.id DESC"
	}
	d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM mm_wallet_log bi
	 	INNER JOIN mm_member m ON m.id=bi.member_id
			WHERE bi.member_id= $1 %s`, where), &num, memberId)

	if num > 0 {
		sqlLine := fmt.Sprintf(`SELECT bi.* FROM mm_wallet_log bi
			INNER JOIN mm_member m ON m.id=bi.member_id
			WHERE member_id= $1 %s %s LIMIT $3 OFFSET $2`,
			where, orderBy)
		d.Query(sqlLine, func(_rows *sql.Rows) {
			rows = db.RowsToMarshalMap(_rows)
		}, memberId, begin, end-begin)
	} else {
		rows = []map[string]interface{}{}
	}

	return num, rows
}

// 获取最近的余额变动信息
func (m *MemberQuery) GetLatestWalletLogByKind(memberId int64, kind int) *member.WalletAccountLog {
	var info = new(member.WalletAccountLog)
	if err := m.GetOrm().GetBy(info, "member_id= $1 AND kind= $2 ORDER BY create_time DESC",
		memberId, kind); err == nil {
		return info
	}
	return nil
}

// 筛选会员根据用户或者手机
func (m *MemberQuery) FilterMemberByUsrOrPhone(key string) []*dto.SimpleMember {
	qp := "%" + key + "%"
	list := make([]*dto.SimpleMember, 0)
	var id int
	var user, name, phone, avatar string
	m.Query(`SELECT id,user,mm_profile.name,mm_profile.phone,
        mm_profile.avatar FROM mm_member
        INNER JOIN mm_profile ON mm_profile.member_id=mm_member.id
        WHERE user LIKE $1 OR mm_profile.name LIKE $2 OR
        mm_profile.phone LIKE $3`, func(rows *sql.Rows) {
		for rows.Next() {
			rows.Scan(&id, &user, &name, &phone, &avatar)
			list = append(list, &dto.SimpleMember{
				Id:     id,
				User:   user,
				Name:   name,
				Phone:  phone,
				Avatar: format.GetResUrl(avatar),
			})
		}
	}, qp, qp, qp)
	return list
}

func (m *MemberQuery) GetMemberByUserOrPhone(key string) *dto.SimpleMember {
	e := dto.SimpleMember{}
	err := m.QueryRow(`SELECT id,user,mm_profile.name,mm_profile.phone,
        mm_profile.avatar FROM mm_member
        INNER JOIN mm_profile ON mm_profile.member_id=mm_member.id
        WHERE user = $1 OR mm_profile.phone = $2`, func(rows *sql.Row) error {
		er := rows.Scan(&e.Id, &e.User, &e.Name, &e.Phone, &e.Avatar)
		e.Avatar = format.GetResUrl(e.Avatar)
		return er
	}, key, key)
	if err == nil {
		return &e
	}
	return nil
}

// 根据手机获取会员编号
func (m *MemberQuery) GetMemberIdByPhone(phone string) int64 {
	var id int64
	m.ExecScalar(`SELECT id FROM mm_member
        INNER JOIN mm_profile ON mm_profile.member_id=mm_member.id
        WHERE mm_profile.phone = $1 LIMIT 1`, &id, phone)
	return id
}

// 会员推广排名
func (m *MemberQuery) GetMemberInviRank(mchId int32, allTeam bool, levelComp string, level int,
	startTime int64, endTime int64, num int) []*dto.RankMember {
	list := make([]*dto.RankMember, 0)
	var id int64
	var user, name string
	var inviNum, totalNum, regTime int
	var rank int = 0

	var sortField string = "t.all_num DESC"
	if !allTeam {
		sortField = "t.invi_num DESC"
	}

	var levelCompStr string = fmt.Sprintf("%s%d", levelComp, level)
	//{level_comp}{level_value}

	m.Query(fmt.Sprintf(`SELECT id,user,name,invi_num,all_num,reg_time FROM ( SELECT m.*,
 (SELECT COUNT(0) FROM mm_relation r INNER JOIN mm_member m1 ON m1.id = r.member_id WHERE
  (m1.level%s) AND r.inviter_id = m.id
	AND r.reg_mchid=rl.reg_mchid  AND m1.reg_time BETWEEN
  $1 AND $2 ) as invi_num,
	((SELECT COUNT(0) FROM mm_relation r INNER JOIN mm_member m1 ON m1.id = r.member_id WHERE
  (m1.level%s) AND r.inviter_id = m.id
	AND r.reg_mchid=rl.reg_mchid AND m1.reg_time BETWEEN
  $3 AND $4 )+
 (SELECT COUNT(0) FROM mm_relation r INNER JOIN mm_member m1
  ON m1.id = r.member_id WHERE (m1.level%s) AND inviter_id IN
	(SELECT member_id FROM mm_relation r INNER JOIN mm_member m1 ON m1.id = r.member_id WHERE
  (m1.level%s) AND r.inviter_id =
    m.id AND r.reg_mchid=rl.reg_mchid AND m1.reg_time BETWEEN
  $5 AND $6 ))) as all_num
 FROM mm_member m INNER JOIN mm_relation rl ON m.id= rl.member_id
 WHERE rl.reg_mchid = $7 AND state= $8) t ORDER BY %s,t.reg_time asc
 LIMIT $9`, levelCompStr, levelCompStr, levelCompStr, levelCompStr, sortField), func(rows *sql.Rows) {
		for rows.Next() {
			rows.Scan(&id, &user, &name, &inviNum, &totalNum, &regTime)
			rank++
			list = append(list, &dto.RankMember{
				Id:       id,
				Usr:      user,
				Name:     name,
				RankNum:  rank,
				InviNum:  inviNum,
				TotalNum: totalNum,
				RegTime:  regTime,
			})
		}
	}, startTime, endTime, startTime, endTime, startTime, endTime, mchId, 1, num)

	return list
}

// 获取分页商铺收藏
func (m *MemberQuery) PagedShopFav(memberId int64, begin, end int,
	where string) (num int, rows []*dto.PagedShopFav) {
	d := m.Connector
	if len(where) > 0 {
		where = " AND " + where
	}
	d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM mm_favorite f
	INNER JOIN  mch_shop s ON f.refer_id =s.id
    INNER JOIN mch_online_shop o ON s.id = o.shop_id
    INNER JOIN mch_merchant mch ON mch.id = s.vendor_id
    WHERE f.member_id= $1 AND f.fav_type= $2 %s`, where), &num,
		memberId, member.FavTypeShop)

	if num > 0 {
		sqlLine := fmt.Sprintf(`SELECT f.id,s.id as shop_id,mch.id as mch_id,
    s.name as shop_name,o.logo,f.update_time FROM mm_favorite f
    INNER JOIN  mch_shop s ON f.refer_id =s.id
    INNER JOIN mch_online_shop o ON s.id = o.shop_id
    INNER JOIN mch_merchant mch ON mch.id = s.vendor_id
    WHERE f.member_id= $1 AND f.fav_type= $2 %s ORDER BY f.update_time DESC LIMIT $4 OFFSET $3`,
			where)
		d.Query(sqlLine, func(rs *sql.Rows) {
			for rs.Next() {
				e := dto.PagedShopFav{}
				rs.Scan(&e.Id, &e.ShopId, &e.MchId, &e.ShopName,
					&e.Logo, &e.UpdateTime)
				e.Logo = format.GetResUrl(e.Logo)
				rows = append(rows, &e)
			}
		}, memberId, member.FavTypeShop, begin, end-begin)
	} else {
		rows = make([]*dto.PagedShopFav, 0)
	}
	return num, rows
}

// 获取分页商铺收藏
func (m *MemberQuery) PagedGoodsFav(memberId int64, begin, end int,
	where string) (num int, rows []*dto.PagedGoodsFav) {
	d := m.Connector
	if len(where) > 0 {
		where = " AND " + where
	}
	d.ExecScalar(fmt.Sprintf(`SELECT COUNT(0) FROM mm_favorite f
    INNER JOIN item_info gs ON gs.id = f.refer_id
    INNER JOIN pro_product product ON gs.product_id=product.id
    WHERE f.member_id= $1 AND f.fav_type= $2 %s`, where), &num,
		memberId, member.FavTypeGoods)

	if num > 0 {
		sqlLine := fmt.Sprintf(`SELECT f.id,gs.id as goods_id,product.name as goods_name,
            img,sale_price,gs.stock_num,product.update_time
            FROM mm_favorite f INNER JOIN item_info gs ON gs.id = f.refer_id
            INNER JOIN pro_product product ON gs.product_id=product.id
            WHERE f.member_id= $1 AND f.fav_type= $2 %s ORDER BY f.update_time DESC
            LIMIT $4 OFFSET $3`,
			where)
		d.Query(sqlLine, func(rs *sql.Rows) {
			for rs.Next() {
				e := dto.PagedGoodsFav{}
				rs.Scan(&e.Id, &e.SkuId, &e.GoodsName, &e.Image, &e.SalePrice,
					&e.StockNum, &e.UpdateTime)
				e.Image = format.GetResUrl(e.Image)
				rows = append(rows, &e)
			}
		}, memberId, member.FavTypeGoods, begin, end-begin)

	} else {
		rows = make([]*dto.PagedGoodsFav, 0)
	}
	return num, rows
}

// 获取从指定时间到现在推荐指定等级会员的数量
func (m *MemberQuery) GetInviteQuantity(memberId int64, where string) int32 {
	var total int32
	m.Connector.ExecScalar(`SELECT COUNT(0) FROM mm_relation
        INNER JOIN mm_member ON mm_member.id = mm_relation.member_id
        LEFT JOIN mm_trusted_info mt ON mt.member_id=mm_member.id
        WHERE inviter_id = $1 `+where, &total, memberId)
	return total
}

// 获取从指定时间到现在推荐指定等级会员的数量
func (m *MemberQuery) GetInviteArray(memberId int64, where string) []int64 {
	arr := []int64{}
	m.Connector.Query(`SELECT mm_relation.member_id FROM mm_relation
        INNER JOIN mm_member ON mm_member.id = mm_relation.member_id
        LEFT JOIN mm_trusted_info mt ON mt.member_id=mm_member.id
        WHERE inviter_id = $1 `+where, func(rows *sql.Rows) {
		var i int64
		for rows.Next() {
			if err := rows.Scan(&i); err == nil && i > 0 {
				arr = append(arr, i)
			}
		}
	}, memberId)
	return arr
}

// 获取邀请会员数量
func (m *MemberQuery) InviteMembersQuantity(memberId int64, where string) int {
	total := 0
	err := m.Connector.ExecScalar(`SELECT COUNT(0) FROM mm_relation WHERE `+where, &total)
	if err != nil {
		log.Printf("[ Query][ Error]: invite member quantity error:%s", err.Error())
	}
	return total
}
