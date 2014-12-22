package dao

import (
	"com/ording/entity"
	"database/sql"
	"errors"
	"fmt"
	"ops/cf/db"
	"time"
)

type shopDao struct {
	db.Connector
}

//获取门店信息
func (this *shopDao) GetShopById(id int) (shop *entity.Shop) {
	this.Connector.QueryRow("SELECT id,pt_id,name,address,phone,order_index,state,create_time FROM pt_shop WHERE id=?",
		func(row *sql.Row) {
			shop = &entity.Shop{}
			//result := db.ConvSqlRowToMap(rows)
			//shop.Address = string(result["Address"])
			//fmt.Println(result["CreateTime"])
			//shop.CreateTime,_ = time.p

			var createTime string
			row.Scan(&shop.Id, &shop.PartnerId, &shop.Name,
				&shop.Address, &shop.Phone, &shop.OrderIndex,
				&shop.State, &createTime)
			shop.CreateTime, _ = time.Parse("2006-01-02 15:04:05", createTime)

		}, id)
	return shop
}

func (this *shopDao) SaveShop(shop *entity.Shop) (int, error) {
	if shop.Id <= 0 {
		//多行字符用``
		_, id, err := this.Connector.Exec(`INSERT INTO pt_shop(pt_id,name,address,phone,order_index,state,create_time)
				VALUES(?,?,?,?,?,?,?)`, shop.PartnerId, shop.Name, shop.Address,
			shop.Phone, shop.OrderIndex, shop.State, shop.CreateTime)
		return id, err
	} else {
		_, _, err := this.Connector.Exec(`UPDATE pt_shop
                            SET
                            name=?,
                            address=?,
                            phone=?,
                            order_index=?,
                            state=?
                            WHERE id=? AND pt_id=?`, shop.Name, shop.Address,
			shop.Phone, shop.OrderIndex, shop.State, shop.Id, shop.PartnerId)
		return shop.Id, err
	}
}

func (this *shopDao) GetShopsOfPartner(partnerId int) []entity.Shop {
	shops := []entity.Shop{}

	this.Connector.Query("SELECT id,pt_id,name,address,phone,order_index,state,create_time FROM pt_shop WHERE pt_id=?",
		func(rows *sql.Rows) {
			for rows.Next() {
				shop := entity.Shop{}
				var createTime string
				rows.Scan(&shop.Id, &shop.PartnerId, &shop.Name,
					&shop.Address, &shop.Phone, &shop.OrderIndex,
					&shop.State, &createTime)
				shop.CreateTime, _ = time.Parse("2006-01-02 15:04:05", createTime)
				shops = append(shops, shop)
			}
		}, partnerId)

	return shops
}

func (this *shopDao) DeleteShop(partnerId, shopId int) error {
	var row int
	this.Connector.ExecScalar(
		`SELECT COUNT(0) FROM pt_order where pt_id=? AND shop_id=?`,
		&row, partnerId, shopId)
	if row != 0 {
		return errors.New(fmt.Sprintf("该门店有%d条相关订单，无法删除！", row))
	}

	this.Connector.GetOrm().Delete(entity.Shop{},
		fmt.Sprintf("pt_id=%d AND id=%d", partnerId, shopId))
	return nil
}
