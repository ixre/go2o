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

	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/go2o/core/infrastructure/fw/collections"
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
	AuthRepo      fw.BaseRepository[merchant.Authenticate]
	_staffRepo    staff.IStaffRepo
	_approvalRepo approval.IApprovalRepository
}

func NewMerchantQuery(c gof.App, fo fw.ORM, staffRepo staff.IStaffRepo,
	approvalRepo approval.IApprovalRepository,
) *MerchantQuery {
	q := &MerchantQuery{
		Connector:     c.Db(),
		Storage:       c.Storage(),
		_staffRepo:    staffRepo,
		_approvalRepo: approvalRepo,
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

	ret, err := fw.UnifinedQueryPaging(m.ORM, p, tables, `
			 p.*,p.mch_name,person_name,p.tel,
        (SELECT COUNT(1) FROM mch_online_shop s WHERE s.vendor_id=p.id) as online_shops`)
	// (SELECT COUNT(1) FROM mch_offline_shop s WHERE s.mch_id=p.id AND shop_type=2) as ofs_num`)
	for _, v := range ret.Rows {
		r := fw.ParsePagingRow(v)
		r.Excludes("password")
	}
	return ret, err
}

// 查询分页商户待审核记录
func (m *MerchantQuery) QueryPagingAuthenticates(p *fw.PagingParams) (_ *fw.PagingResult, err error) {
	tables := `mch_authenticate a
         LEFT JOIN mch_merchant p ON p.id=a.mch_id`
	ret, err := fw.UnifinedQueryPaging(m.ORM, p, tables, `
			 a.*,p.mch_name`)
	for _, v := range ret.Rows {
		r := fw.ParsePagingRow(v)
		r.Excludes("password")
	}
	return ret, err
}

// 查询商户的认证信息
func (m *MerchantQuery) QueryMerchantAuthenticates(mchId int) []*merchant.Authenticate {
	var ret []*merchant.Authenticate
	m.ORM.Find(&ret, "mch_id = ?", mchId)
	return ret
}

// GetStaffTransferInfo 获取转商户信息
func (m *MerchantQuery) GetStaffTransferInfo(staffId int64) (ret struct {
	TxId            int                     `json:"txId"`
	MchName         string                  `json:"mchName"`
	TransferMchName string                  `json:"transferMchName"`
	ApprovalLogs    []*approval.ApprovalLog `json:"approvalLogs"`
}) {
	sf := m._staffRepo.TransferRepo().FindBy("staff_id = ? ORDER BY id DESC", staffId)
	if sf == nil {
		return ret
	}
	ret.TxId = sf.Id
	mch := m.Get(sf.OriginMchId)
	tarMch := m.Get(sf.TransferMchId)
	ret.MchName = mch.MchName
	ret.TransferMchName = tarMch.MchName
	ret.ApprovalLogs = m._approvalRepo.GetLogRepo().FindList(&fw.QueryOption{
		Skip:  0,
		Limit: 10,
		Order: "id ASC",
	}, " approval_id = ?", sf.ApprovalId)
	return ret
}

// QueryMerchantByName 根据商户名称查询商户信息
func (m *MerchantQuery) QueryMerchantByName(name string) []map[string]interface{} {
	list := m.FindList(nil, "mch_name = ?", name)
	return collections.MapList(list, func(m *merchant.Merchant) map[string]interface{} {
		return map[string]interface{}{
			"id":      m.Id,
			"mchName": m.MchName,
			"address": m.Address,
		}
	})
}
