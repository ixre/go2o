package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/atnet/gof/db/orm"
	"github.com/atnet/gof/log"
	_ "github.com/go-sql-driver/mysql"
)

var _ Connector = new(CommonConnector)

//数据库连接器
type CommonConnector struct {
	driverName   string  //驱动名称
	driverSource string  //驱动连接地址
	_db          *sql.DB //golang db只需要open一次即可
	_orm         orm.Orm
	logger       log.ILogger
}

//create a new connector
func NewCommonConnector(driverName, driverSource string,
	l log.ILogger, maxConn int) Connector {
	db, err := sql.Open(driverName, driverSource)

	if err == nil {
		err = db.Ping()
	}

	if err != nil {
		defer db.Close()
		//如果异常，则显示并退出
		log.Fatalln("[" + driverName + "] " + err.Error())
		return nil
	}

	// 设置最大连接数
	if maxConn > 0 {
		db.SetMaxOpenConns(maxConn)
	}

	return &CommonConnector{
		_db:          db,
		_orm:         orm.NewOrm(db),
		driverName:   driverName,
		driverSource: driverName,
		logger:       l,
	}
}

func (this *CommonConnector) println(v ...interface{}) {
	if this.logger != nil {
		this.logger.Println(v...)
	}
}

func (this *CommonConnector) GetDb() *sql.DB {
	return this._db
}

func (this *CommonConnector) GetOrm() orm.Orm {
	return this._orm
}

func (this *CommonConnector) Query(sql string, f func(*sql.Rows), arg ...interface{}) error {
	stmt, err := this.GetDb().Prepare(sql)
	if err != nil {
		err = errors.New(fmt.Sprintf("[SQL][Error]:", err.Error(), " [SQL]:", sql))
		this.println(err.Error())
		return err
	}
	rows, err := stmt.Query(arg...)
	if err != nil {
		this.println(err.Error())
		return err
	}
	defer stmt.Close()
	if f != nil {
		f(rows)
	}
	return nil
}

//查询Rows
func (this *CommonConnector) QueryRow(sql string, f func(*sql.Row), arg ...interface{}) error {
	stmt, err := this.GetDb().Prepare(sql)
	if err != nil {
		err = errors.New(fmt.Sprintf("[SQL][Error]:", err.Error(), " [SQL]:", sql))
		this.println(err.Error())
		return err
	} else {
		defer stmt.Close()
		row := stmt.QueryRow(arg...)
		if f != nil && row != nil {
			f(row)
		}
	}
	return nil
}

func (this *CommonConnector) ExecScalar(s string, result interface{}, arg ...interface{}) (err error) {
	if result == nil {
		return errors.New("Result is null")
	}

	this.QueryRow(s, func(row *sql.Row) {
		err = row.Scan(result)
	}, arg...)

	if err != nil {
		err = errors.New(fmt.Sprintf("[SQL][Error]:", err.Error(), " [SQL]:", s))
		this.println(err.Error())
		return err
	}

	return err
}

//执行
func (this *CommonConnector) Exec(sql string, args ...interface{}) (rows int, lastInsertId int, err error) {
	stmt, err := this.GetDb().Prepare(sql)
	if err != nil {
		return 0, -1, err
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		err = errors.New(fmt.Sprintf("[SQL][Error]:", err.Error(), " [SQL]:", sql))
		this.println(err.Error())
		return 0, -1, err
	}
	defer stmt.Close()

	lastId, _ := result.LastInsertId()
	affect, _ := result.RowsAffected()

	return int(affect), int(lastId), nil
}

func (this *CommonConnector) ExecNonQuery(sql string, args ...interface{}) (int, error) {
	n, _, err := this.Exec(sql, args...)
	return n, err
}
