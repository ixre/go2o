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
	"fmt"
	"regexp"
	"time"

	"github.com/ixre/go2o/core/domain/interface/approval"
	"github.com/ixre/go2o/core/domain/interface/domain/enum"
	"github.com/ixre/go2o/core/domain/interface/merchant"
	"github.com/ixre/go2o/core/domain/interface/merchant/staff"
	"github.com/ixre/go2o/core/domain/interface/wallet"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/go2o/core/infrastructure/fw/collections"
	"github.com/ixre/go2o/core/infrastructure/util"
	"github.com/ixre/go2o/core/initial/provide"
	"github.com/ixre/go2o/core/variable"
	"github.com/ixre/gof"
	"github.com/ixre/gof/db"
	"github.com/ixre/gof/storage"
)

type MerchantQuery struct {
	db.Connector
	fw.ORM
	Storage        storage.Interface
	_repo          merchant.IMerchantRepo
	AuthRepo       fw.BaseRepository[merchant.Authenticate]
	_staffRepo     staff.IStaffRepo
	_approvalRepo  approval.IApprovalRepository
	_walletRepo    wallet.IWalletRepo
	_walletLogRepo fw.BaseRepository[wallet.WalletLog]
}

func NewMerchantQuery(c gof.App, fo fw.ORM, staffRepo staff.IStaffRepo,
	approvalRepo approval.IApprovalRepository,
	mchRepo merchant.IMerchantRepo,
	walletRepo wallet.IWalletRepo,
) *MerchantQuery {
	q := &MerchantQuery{
		Connector:     c.Db(),
		Storage:       c.Storage(),
		_staffRepo:    staffRepo,
		_approvalRepo: approvalRepo,
		_repo:         mchRepo,
		_walletRepo:   walletRepo,
	}
	q.ORM = fo
	q.AuthRepo.ORM = fo
	q._walletLogRepo.ORM = fo
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

// QueryMerchantList 查询商户列表
func (m *MerchantQuery) QueryMerchantList(begin, size int) []*merchant.Merchant {
	return m._repo.FindList(&fw.QueryOption{
		Skip:  begin,
		Limit: size,
		Order: "id ASC",
	}, "")
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

// QueryPagingMerchants 查询分页的商户列表(基本信息)
func (m *MerchantQuery) QueryPagingMerchants(p *fw.PagingParams) (_ *fw.PagingResult, err error) {
	return m._repo.QueryPaging(p)
}

// QueryPagingMerchantList 查询分页的商户列表(完整信息)
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

// GetMerchantAuthenticate 查询商户的认证信息
func (m *MerchantQuery) GetMerchantAuthenticate(mchId int) *merchant.Authenticate {
	return m.AuthRepo.FindBy("mch_id = ?", mchId)
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
	StaffName       string                  `json:"staffName"`
	MchName         string                  `json:"mchName"`
	TransferMchName string                  `json:"transferMchName"`
	CreateTime      int                     `json:"createTime"`
	ApprovalLogs    []*approval.ApprovalLog `json:"approvalLogs"`
}) {
	sf := m._staffRepo.TransferRepo().FindBy("staff_id = ? ORDER BY id DESC", staffId)
	if sf == nil {
		return ret
	}
	if sf.ReviewStatus == int(enum.ReviewApproved) {
		// 如果已通过，也不返回数据
		return ret
	}
	staff := m._staffRepo.Get(staffId)
	ret.TxId = sf.Id
	mch := m._repo.Get(sf.OriginMchId)
	tarMch := m._repo.Get(sf.TransferMchId)
	ret.MchName = mch.MchName
	ret.StaffName = staff.CertifiedName
	ret.TransferMchName = tarMch.MchName
	ret.CreateTime = sf.CreateTime
	ret.ApprovalLogs = m._approvalRepo.GetLogRepo().FindList(&fw.QueryOption{
		Skip:  0,
		Limit: 10,
		Order: "id ASC",
	}, " approval_id = ?", sf.ApprovalId)
	return ret
}

// QueryMerchantByName 根据商户名称查询商户信息
func (m *MerchantQuery) QueryMerchantByName(name string) []map[string]interface{} {
	list := m._repo.FindList(nil, "mch_name = ?", name)
	return collections.MapList(list, func(m *merchant.Merchant) map[string]interface{} {
		return map[string]interface{}{
			"mchId":   m.Id,
			"mchName": m.MchName,
			"address": m.Address,
		}
	})
}

// QueryPagingStaffs 查询商户待认证员工列表(商户端)
func (m *MerchantQuery) QueryMerchantPendingStaffs(p *fw.PagingParams) (*fw.PagingResult, error) {
	tables := fmt.Sprintf(`mm_member m
		INNER JOIN mch_staff s ON s.member_id=m.id
		LEFT JOIN mm_profile pro ON pro.member_id = m.id
		LEFT JOIN mm_cert_info c ON c.member_id = m.id AND version=0 AND review_status <> %d`, enum.ReviewApproved)
	fields := `
	distinct(s.id),m.nickname,c.real_name,m.username,m.exp,m.profile_photo,pro.gender,
	m.phone,m.level,m.user_flag,
	m.reg_from,m.reg_time,m.login_time,
	s.certified_name,s.is_certified,c.review_status,c.remark,c.manual_review
	`
	return fw.UnifinedQueryPaging(m.ORM, p, tables, fields)
}

// 查询商户员工转商户信息
func (m *MerchantQuery) QueryTransferStaffs(mchId int, transferType int, p *fw.PagingParams) (*fw.PagingResult, error) {
	tables := `mch_staff_transfer t
		INNER JOIN mch_staff s ON s.id = t.staff_id
		INNER JOIN approval a ON a.id = t.approval_id
		`
	fields := `t.*,
	s.certified_name,
	s.gender,
	a.assign_uid,
	a.assign_name,
	(SELECT mch_name FROM mch_merchant WHERE id = t.origin_mch_id) as origin_mch_name,
	(SELECT mch_name FROM mch_merchant WHERE id = t.transfer_mch_id) as transfer_mch_name
	`
	if transferType == 1 {
		p.And("transfer_mch_id = ?", mchId)
	} else {
		p.And("origin_mch_id = ?", mchId)
	}
	rows, err := fw.UnifinedQueryPaging(m.ORM, p, tables, fields)
	for _, row := range rows.Rows {
		r := fw.ParsePagingRow(row)
		isApproval := r.Get("assignUid").(int64) == int64(mchId) && r.Get("reviewStatus").(int64) == 1
		r.Put("isApproval", isApproval)
		r.Put("isTransferIn", r.Get("transferMchId").(int64) == int64(mchId))
	}
	return rows, err
}

// 查询商户月度账单
func (m *MerchantQuery) QueryPagingBills(p *fw.PagingParams) (*fw.PagingResult, error) {
	tables := `mch_bill b
		INNER JOIN mch_merchant m ON m.id = b.mch_id`
	fields := `b.*,m.mch_name`
	return fw.UnifinedQueryPaging(m.ORM, p, tables, fields)
}

// 获取商户月度账单
func (m *MerchantQuery) GetBill(id int) *merchant.MerchantBill {
	return m._repo.BillRepo().Get(id)
}

// 查询商户月度账单明细
func (m *MerchantQuery) QueryPagingBillItems(billId int, p *fw.PagingParams) (*fw.PagingResult, error) {
	return m._walletLogRepo.QueryPaging(p)
}

// 查询商户员工列表
func (m *MerchantQuery) QueryPagingStaffs(p *fw.PagingParams) (*fw.PagingResult, error) {
	return m._staffRepo.QueryPaging(p)
}

// QueryWaitGenerateDailyBills 查询待生成日度账单的商户
func (m *MerchantQuery) QueryWaitGenerateDailyBills(size int, lastId int) []*merchant.MerchantBill {
	// 以当天开始时间，作为账单生成的结束时间
	endTime, _ := util.GetStartEndUnix(time.Now())
	return m._repo.BillRepo().FindList(&fw.QueryOption{
		Limit: size,
	}, "end_time <? AND bill_type = ? AND status = ? AND id > ?",
		endTime,
		merchant.BillTypeDaily,
		merchant.BillStatusPending,
		lastId)
}

// 查询离线员工列表
func (m *MerchantQuery) QueryOfflineStaffList(beginTime, overTime int, begin, size int) []*staff.Staff {
	return m._staffRepo.FindList(&fw.QueryOption{
		Skip:  begin,
		Limit: size,
	}, "work_status = ? AND last_online_time BETWEEN ? AND ?", staff.WorkStatusOffline, beginTime, overTime)
}
