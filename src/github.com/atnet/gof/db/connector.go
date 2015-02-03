/**
 * Copyright 2014 @ Ops Inc.
 * name :
 * author : newmin
 * date : 2013-12-03 20:19
 * description :
 * history :
 */

package db

import (
	"database/sql"
	"github.com/atnet/gof/db/orm"
)

type Connector interface {
	GetDb() *sql.DB

	GetOrm() orm.Orm

	Query(sql string, f func(*sql.Rows), arg ...interface{}) error

	// 查询Rows
	QueryRow(sql string, f func(*sql.Row), arg ...interface{}) error

	ExecScalar(s string, result interface{}, arg ...interface{}) error

	// 执行
	Exec(sql string, args ...interface{}) (rows int, lastInsertId int, err error)

	ExecNonQuery(sql string, args ...interface{}) (int, error)
}
