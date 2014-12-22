package dao

import (
	"com/share/glob"
	"ops/cf/app"
)

var (
	context      app.Context
	member_dao   *memberDao
	item_dao     *itemDao
	partner_dao  *partnerDao
	shop_dao     *shopDao
	category_dao *categoryDao
	order_dao    *orderDao
	comm_dao     *commDao
)

func init() {
	context = glob.CurrContext()
}

func Member() *memberDao {
	if member_dao == nil {
		member_dao = &memberDao{context.Db()}
	}
	return member_dao
}

func Item() *itemDao {
	if item_dao == nil {
		item_dao = &itemDao{context.Db()}
	}
	return item_dao
}

func Category() *categoryDao {
	if category_dao == nil {
		category_dao = &categoryDao{context.Db()}
	}
	return category_dao
}

func Partner() *partnerDao {
	if partner_dao == nil {
		partner_dao = &partnerDao{Context: context, Connector: context.Db()}
	}
	return partner_dao
}

func Shop() *shopDao {
	if shop_dao == nil {
		shop_dao = &shopDao{context.Db()}
	}
	return shop_dao
}

func Order() *orderDao {
	if order_dao == nil {
		order_dao = &orderDao{context.Db()}
	}
	return order_dao
}

func Common() *commDao {
	if comm_dao == nil {
		comm_dao = &commDao{context.Db()}
	}
	return comm_dao
}
