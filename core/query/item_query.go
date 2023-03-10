/**
 * Copyright 2015 @ 56x.net.
 * name : goods_query.go
 * author : jarryliu
 * date : 2016-04-23 12:37
 * description :
 * history :
 */
package query

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/item"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
	"github.com/ixre/go2o/core/dto"
	"github.com/ixre/go2o/core/infrastructure/format"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/db/orm"
)

type ItemQuery struct {
	db.Connector
	o orm.Orm
}

func NewItemQuery(o orm.Orm) *ItemQuery {
	return &ItemQuery{
		Connector: o.Connector(),
		o:         o,
	}
}

// GetPagingOnShelvesItem 根据关键词搜索上架的商品
func (i ItemQuery) GetPagingOnShelvesItem(catId int32,
	start, end int32, where, orderBy string) (int32, []*item.GoodsItem) {
	var sql string
	var total int32

	if len(orderBy) != 0 {
		orderBy += ","
	}
	if catId > 0 {
		where += fmt.Sprintf(" AND item_snapshot.cat_id= %d", catId)
	}
	// i.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(1) FROM item_snapshot
	//      INNER JOIN product ON product.id = item_info.product_id
	i.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(1) FROM item_snapshot
		 WHERE item_snapshot.review_state= $1
		 AND item_snapshot.shelve_state= $2 %s`, where), &total,
		enum.ReviewPass, item.ShelvesOn)
	var list []*item.GoodsItem
	if total > 0 {
		sql = fmt.Sprintf(`SELECT * FROM item_snapshot
		 WHERE item_snapshot.review_state= $1
		 AND item_snapshot.shelve_state= $2 %s
		 ORDER BY %s item_snapshot.update_time DESC LIMIT $4 OFFSET $3`,
			where, orderBy)
		i.o.SelectByQuery(&list, sql,
			enum.ReviewPass, item.ShelvesOn, start, end-start)
	}
	return total, list
}

// SearchOnShelvesItem 根据关键词搜索上架的商品
func (i ItemQuery) SearchOnShelvesItem(word string, start, end int32,
	where, orderBy string) (int32, []*item.GoodsItem) {
	var sql string
	var total int32

	if len(orderBy) != 0 {
		orderBy += ","
	}

	buf := bytes.NewBuffer(nil)
	if word != "" {
		buf.WriteString(" AND (item_snapshot.title LIKE '%")
		buf.WriteString(word)
		buf.WriteString("%' OR item_snapshot.short_title LIKE '%")
		buf.WriteString(word)
		buf.WriteString("%')")
	}
	buf.WriteString(where)
	where = buf.String()

	i.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(1) FROM item_snapshot
		 WHERE item_snapshot.review_state= $1
		 AND item_snapshot.shelve_state= $2 %s`, where), &total,
		enum.ReviewPass, item.ShelvesOn)
	var list []*item.GoodsItem
	if total > 0 {
		sql = fmt.Sprintf(`SELECT * FROM item_snapshot
		 WHERE item_snapshot.review_state= $1
		 AND item_snapshot.shelve_state= $2 %s
		 ORDER BY %s item_snapshot.update_time DESC LIMIT $4 OFFSET $3`,
			where, orderBy)
		i.o.SelectByQuery(&list, sql,
			enum.ReviewPass, item.ShelvesOn, start, end-start)
	}
	return total, list
}

// 根据分类获取上架的商品
func (i ItemQuery) GetOnShelvesItem(catIdArr []int, begin, end int,
	where string) []*item.GoodsItem {
	var list []*item.GoodsItem
	if len(catIdArr) > 0 {
		catIdStr := format.IntArrStrJoin(catIdArr)
		sql := fmt.Sprintf(`SELECT * FROM item_snapshot
		 WHERE item_snapshot.cat_id IN(%s) %s ORDER BY item_snapshot.update_time DESC
		   LIMIT $3 OFFSET $2`, catIdStr, where)
		i.o.SelectByQuery(&list, sql,
			 begin, end-begin)
	}
	return list
}

// 搜索随机的商品列表
func (i ItemQuery) GetRandomItem(catIdArr []int, begin, end int, where string) []*item.GoodsItem {
	/*
	       随机查询： 要减去获取的条数，以确保至少有2条数据
	   SELECT * FROM item_info

	   JOIN (SELECT ROUND(RAND() * (
	     SELECT MAX(id)-2 FROM item_info
	        )) AS id) AS r2

	        WHERE item_info.id > r2.id LIMIT 2

	*/

	s := []string{where}
	if catIdArr != nil && len(catIdArr) > 0 {
		catIdStr := format.IntArrStrJoin(catIdArr)
		if where != "" {
			s = append(s, " AND")
		}
		s = append(s, fmt.Sprintf("item_info.cat_id IN (%s)", catIdStr))
	}
	search := strings.Join(s, "")

	var list []*item.GoodsItem
	sql := fmt.Sprintf(`SELECT * FROM item_info
    JOIN (SELECT ROUND(RAND() * (
      SELECT MAX(id)-? FROM item_snapshot WHERE  item_info.review_state= $1
         AND item_info.shelve_state= ? %s
         )) AS id) AS r2
		 WHERE item_info.Id > r2.id
		  AND item_info.review_state= ?
		 AND item_info.shelve_state= ? %s LIMIT ? OFFSET $3`,
		search, search)
	i.o.SelectByQuery(&list, sql,
		enum.ReviewPass, item.ShelvesOn,
		enum.ReviewPass, item.ShelvesOn, begin, end-begin)
	return list
}

// GetPagingOnShelvesGoods 获取已上架的商品
func (i ItemQuery) GetPagingOnShelvesGoods(shopId int64,
	catIds []int, flag int, start, end int,
	where, orderBy string) (int, []*valueobject.Goods) {
	total := 0
	if len(catIds) > 0 {
		where += fmt.Sprintf(" AND item_snapshot.cat_id IN (%s)",
			format.IntArrStrJoin(catIds))
	}
	if flag > 0 {
		where += fmt.Sprintf(" AND (item_snapshot.item_flag & %d = %d)", flag, flag)
	}
	var list = make([]*valueobject.Goods, 0)
	s := fmt.Sprintf(`SELECT item_snapshot."item_id",
		item_snapshot."cat_id",item_snapshot."title",
		item_snapshot."image",item_snapshot."origin_price",item_snapshot."price",
		item_snapshot."price_range",
		item_snapshot."item_flag",
		item_info."stock_num",item_info."shop_id"
		 FROM item_snapshot
		 LEFT JOIN item_info ON item_snapshot.item_id = item_info.id
		 WHERE ($1 <= 0 OR item_snapshot.shop_id = $2)
		  %s ORDER BY %s LIMIT $4 OFFSET $3`, where, orderBy)
	err := i.Query(s, func(_rows *sql.Rows) {
		for _rows.Next() {
			e := valueobject.Goods{}
			_rows.Scan(&e.ItemId, &e.CategoryId, &e.Title, &e.Image,
				&e.OriginPrice, &e.Price, &e.PriceRange, &e.ItemFlag,
				&e.StockNum, &e.ShopId)
			list = append(list, &e)
		}
	}, shopId, shopId, start, end-start)
	if err != nil {
		log.Println("[ GO2O][ Repo][ Error]:", err.Error(), s)
	}
	log.Println("sql=", s)
	return total, list
}

// QueryItemSalesHistory 查询商品销售记录
func (i *ItemQuery) QueryItemSalesHistory(itemId int64, size int, random bool) (rows []*dto.ItemSalesHistoryDto) {
	s := `SELECT m.user_code,m.nickname,m.portrait,ord.create_time,
		ord.status FROM sale_sub_item it 
		INNER JOIN sale_normal_order ord ON ord.id = it.order_id
		INNER JOIN mm_member m ON m.member_id = ord.buyer_id
		WHERE it.item_id = $1 LIMIT $2 `
	err := i.Query(s, func(_rows *sql.Rows) {
		for _rows.Next() {
			e := dto.ItemSalesHistoryDto{}
			_rows.Scan(&e.BuyerUserCode, &e.BuyerName, &e.BuyerPortrait, &e.BuyTime, &e.OrderState)
			rows = append(rows, &e)
		}
	}, itemId, size)
	if err != nil {
		log.Println("[ GO2O][ ERROR]", err.Error(), s)
	}
	return rows
}

func (i *ItemQuery) SearchItem(shopId int, keyword string, size int) (rows []*dto.SearchItemResultDto) {
	where := "item_snapshot.title LIKE '%" + keyword + "%'"
	if shopId > 0 {
		where += fmt.Sprintf(" AND item_snapshot.shop_id = %d", shopId)
	}
	cmd := fmt.Sprintf(`SELECT item_id,
			item_snapshot.item_flag, 
			item_snapshot.code,
			item_snapshot.title,
			item_snapshot.image,
			item_snapshot.vendor_id,
			item_snapshot.price_range,
			item_info.stock_num
			FROM item_snapshot
			LEFT JOIN item_info ON item_info.id = item_snapshot.item_id
			 WHERE %s LIMIT $1 `, where)
	err := i.Query(cmd, func(_rows *sql.Rows) {
		for _rows.Next() {
			e := dto.SearchItemResultDto{}
			_rows.Scan(&e.ItemId, &e.ItemFlag, &e.Code, &e.Title, &e.Image,
				&e.SellerId, &e.PriceRange, &e.StockNum)
			rows = append(rows, &e)
		}
	}, size)
	if err != nil {
		log.Println("[ GO2O][ ERROR]", err.Error(), cmd)
	}
	return rows
}

// GetPagingOnShelvesItemForWholesale 根据关键词搜索上架的商品
func (i ItemQuery) GetPagingOnShelvesItemForWholesale(catId int32,
	start, end int32, where, orderBy string) (int32, []*item.GoodsItem) {
	var sql string
	var total int32

	if len(orderBy) != 0 {
		orderBy += ","
	}
	if catId > 0 {
		where += fmt.Sprintf(" AND item_info.cat_id= %d", catId)
	}
	i.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(1) FROM ws_item
         INNER JOIN item_info ON item_info.id=ws_item.item_id
         INNER JOIN product ON product.id = item_info.product_id
		 WHERE ws_item.review_state= $1
		 AND ws_item.shelve_state= $2 %s`, where), &total,
		enum.ReviewPass, item.ShelvesOn)
	var list []*item.GoodsItem
	if total > 0 {
		sql = fmt.Sprintf(`SELECT * FROM  ws_item
         INNER JOIN item_info ON item_info.id=ws_item.item_id
         INNER JOIN product ON product.id = item_info.product_id
		 WHERE ws_item.review_state= $1
		 AND ws_item.shelve_state= $2 %s
		 ORDER BY %s item_info.update_time DESC LIMIT $4 OFFSET $3`,
			where, orderBy)
		i.o.SelectByQuery(&list, sql,
			enum.ReviewPass, item.ShelvesOn, start, end-start)
	}
	return total, list
}

// 根据关键词搜索上架的商品
func (i ItemQuery) SearchOnShelvesItemForWholesale(word string, start, end int32,
	where, orderBy string) (int32, []*item.GoodsItem) {
	var sql string
	var total int32

	if len(orderBy) != 0 {
		orderBy += ","
	}

	buf := bytes.NewBuffer(nil)
	if word != "" {
		buf.WriteString(" AND (item_info.title LIKE '%")
		buf.WriteString(word)
		buf.WriteString("%' OR item_info.short_title LIKE '%")
		buf.WriteString(word)
		buf.WriteString("%')")
	}
	buf.WriteString(where)
	where = buf.String()

	i.Connector.ExecScalar(fmt.Sprintf(`SELECT COUNT(1) FROM ws_item
         INNER JOIN item_info ON item_info.id=ws_item.item_id
         INNER JOIN product ON product.id = item_info.product_id
		 WHERE ws_item.review_state= $1
		 AND ws_item.shelve_state= $2  %s`, where), &total,
		enum.ReviewPass, item.ShelvesOn)
	var list []*item.GoodsItem

	if total > 0 {
		sql = fmt.Sprintf(`SELECT item_info.id,item_info.product_id,item_info.item_flag,
		item_info.cat_id,item_info.vendor_id,item_info.brand_id,item_info.shop_id,
		item_info.shop_cat_id,item_info.express_tid,item_info.title,
		item_info.short_title,item_info.code,item_info.image,
		ws_item.price_range,item_info.stock_num,
		item_info.sale_num,item_info.sku_num,item_info.sku_id,item_info.cost,
		ws_item.price,item_info.origin_price,item_info.weight,item_info.bulk,
		item_info.shelve_state,item_info.review_state,item_info.review_remark,
		item_info.sort_num,item_info.create_time,item_info.update_time
		 FROM ws_item INNER JOIN item_info ON item_info.id=ws_item.item_id
         INNER JOIN product ON product.id = item_info.product_id
		 WHERE ws_item.review_state= $1
		 AND ws_item.shelve_state= $2 %s
		 ORDER BY %s item_info.update_time DESC LIMIT $4 OFFSET $3`,
			where, orderBy)
		i.o.SelectByQuery(&list, sql,
			enum.ReviewPass, item.ShelvesOn, start, end-start)
	}
	return total, list
}
