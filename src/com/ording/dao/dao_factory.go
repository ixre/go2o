package dao

import (
	"com/share/glob"
	"github.com/atnet/gof/app"
)

var (
	context      app.Context
	member_dao   *memberDao
	partner_dao  *partnerDao
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
