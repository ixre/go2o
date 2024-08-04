/**
 * Copyright 2014 @ 56x.net.
 * name :
 * author : jarryliu
 * date : 2013-12-03 23:20
 * description :
 * history :
 */
package query

import (
	"regexp"

	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/go2o/core/initial/provide"
	"github.com/ixre/go2o/core/variable"
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/storage"
)

type MerchantQuery struct {
	db.Connector
	Storage storage.Interface
	fw.BaseRepository[merchant.Merchant]
	AuthRepo fw.BaseRepository[merchant.Authenticate]
}

func NewMerchantQuery(c gof.App, fo fw.ORM) *MerchantQuery {
	q := &MerchantQuery{
		Connector: c.Db(),
		Storage:   c.Storage(),
	}
	q.ORM = fo
	q.AuthRepo.ORM = fo
	return q
}

var (
	commHostRegexp *regexp.Regexp
)

func getHostRegexp() *regexp.Regexp {
	if commHostRegexp == nil {
		cfg := provide.GetApp().Config()
		commHostRegexp = regexp.MustCompile("([^\\.]+)." +
			cfg.GetString(variable.ServerDomain))
	}
	return commHostRegexp
}

// 根据主机查询商户编号
func (m *MerchantQuery) QueryMerchantIdByHost(host string) int64 {
	var mchId int64
	var err error

	reg := getHostRegexp()
	if reg.MatchString(host) {
		matches := reg.FindAllStringSubmatch(host, 1)
		user := matches[0][1]
		err = m.Connector.ExecScalar(`SELECT id FROM mch_merchant WHERE login_user = $1`, &mchId, user)
	} else {
		err = m.Connector.ExecScalar(
			`SELECT id FROM mch_merchant INNER JOIN pt_siteconf
                     ON pt_siteconf.merchant_id = mch_merchant.id
                     WHERE host= $1`, &mchId, host)
	}
	if err != nil {
		//gof.CurrentApp.Log().Error(err)
	}
	return mchId
}

// 验证用户密码并返回编号
func (m *MerchantQuery) Verify(user, pwd string) int {
	var id int
	m.Connector.ExecScalar("SELECT id FROM mch_merchant WHERE login_user = $1 AND login_pwd= $2", &id, user, pwd)
	return id
}

// QueryPagingMerchantList 查询分页的商户列表
func (m *MerchantQuery) QueryPagingMerchantList(p *fw.PagingParams) (_ *fw.PagingResult, err error) {
	tables := `mch_merchant p
         LEFT JOIN mch_authenticate a ON a.mch_id=p.id AND a.version = 1`

	ret, err := fw.UnifinedPagingQuery(m.ORM, p, tables, `
			 p.*,p.mch_name,person_name,p.tel,
        (SELECT COUNT(1) FROM mch_online_shop s WHERE s.vendor_id=p.id) as online_shops`)
	// (SELECT COUNT(1) FROM mch_offline_shop s WHERE s.mch_id=p.id AND shop_type=2) as ofs_num`)
	for _, v := range ret.Rows {
		r := fw.ParsePagingRow(v)
		r.Excludes("password")
	}
	return ret, err
}
