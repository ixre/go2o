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
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/storage"
	"go2o/core/infrastructure"
	"go2o/core/variable"
	"regexp"
)

type MerchantQuery struct {
	db.Connector
	Storage storage.Interface
}

func NewMerchantQuery(c gof.App) *MerchantQuery {
	return &MerchantQuery{
		Connector: c.Db(),
		Storage:   c.Storage(),
	}
}

var (
	commHostRegexp *regexp.Regexp
)

func getHostRegexp() *regexp.Regexp {
	if commHostRegexp == nil {
		commHostRegexp = regexp.MustCompile("([^\\.]+)." +
			infrastructure.GetApp().Config().GetString(variable.ServerDomain))
	}
	return commHostRegexp
}

// 根据主机查询商户编号
func (m *MerchantQuery) QueryMerchantIdByHost(host string) int32 {
	var mchId int32
	var err error

	reg := getHostRegexp()
	if reg.MatchString(host) {
		matches := reg.FindAllStringSubmatch(host, 1)
		usr := matches[0][1]
		err = m.Connector.ExecScalar(`SELECT id FROM mch_merchant WHERE usr= $1`, &mchId, usr)
	} else {
		err = m.Connector.ExecScalar(
			`SELECT id FROM mch_merchant INNER JOIN pt_siteconf
                     ON pt_siteconf.merchant_id = mch_merchant.id
                     WHERE host= $1`, &mchId, host)
	}
	if err != nil {
		gof.CurrentApp.Log().Error(err)
	}
	return mchId
}

// 验证用户密码并返回编号
func (m *MerchantQuery) Verify(usr, pwd string) int32 {
	var id int32
	m.Connector.ExecScalar("SELECT id FROM mch_merchant WHERE usr= $1 AND pwd= $2", &id, usr, pwd)
	return id
}
