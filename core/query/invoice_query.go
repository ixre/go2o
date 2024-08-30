package query

import (
	"github.com/ixre/go2o/core/domain/interface/invoice"
	"github.com/ixre/go2o/core/infrastructure/fw"
)

type InvoiceQuery struct {
	orm        fw.ORM
	repo       invoice.IInvoiceRecordRepo
	titleRepo  invoice.IInvoiceTitleRepo
	tenantRepo fw.Repository[invoice.InvoiceTenant]
}

func NewInvoiceQuery(o fw.ORM) *InvoiceQuery {
	return &InvoiceQuery{
		orm: o,
		repo: &fw.BaseRepository[invoice.InvoiceRecord]{
			ORM: o,
		},
		titleRepo: &fw.BaseRepository[invoice.InvoiceTitle]{
			ORM: o,
		},
		tenantRepo: &fw.BaseRepository[invoice.InvoiceTenant]{
			ORM: o,
		},
	}
}

func (i *InvoiceQuery) QueryMerchantIssueInvoices(p *fw.PagingParams) (*fw.PagingResult, error) {
	tables := ""
	fields := ""
	return fw.UnifinedQueryPaging(i.orm, p, tables, fields)
}

// QueryPagingInvoices 查询发票分页
func (i *InvoiceQuery) QueryPagingInvoices(p *fw.PagingParams) (*fw.PagingResult, error) {
	return i.repo.QueryPaging(p)
}

// QueryPagingInvoiceTitles 查询发票抬头分页
func (i *InvoiceQuery) QueryPagingInvoiceTitles(p *fw.PagingParams) (*fw.PagingResult, error) {
	return i.titleRepo.QueryPaging(p)
}

// GetInvoiceTitle 获取发票抬头
func (i *InvoiceQuery) GetInvoiceTitle(id int) *invoice.InvoiceTitle {
	return i.titleRepo.Get(id)
}

// GetMemberTenantId 获取会员的租户ID
func (i *InvoiceQuery) GetMemberTenantId(memberId int) int {
	tenant := i.tenantRepo.FindBy("tenant_uid = ? AND tenant_type = ?", memberId, invoice.TenantUser)
	if tenant == nil {
		return 0
	}
	return tenant.Id
}

// GetMerchantTenantId 获取商户的租户ID
func (i *InvoiceQuery) GetMerchantTenantId(mchId int) int {
	tenant := i.tenantRepo.FindBy("tenant_uid = ? AND tenant_type = ?", mchId, invoice.TenantMerchant)
	if tenant == nil {
		return 0
	}
	return tenant.Id
}

// GetTitles 获取发票抬头
func (i *InvoiceQuery) GetTitles(tenantId int) []*invoice.InvoiceTitle {
	return i.titleRepo.FindList(nil, "tenant_id = ?", tenantId)
}
