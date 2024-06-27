package repos

import (
	"github.com/ixre/go2o/core/domain/interface/invoice"
	impl "github.com/ixre/go2o/core/domain/invoice"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/go2o/core/infrastructure/logger"
)

var _ invoice.IInvoiceTenantRepo = new(invoiceTenantRepoImpl)

type invoiceTenantRepoImpl struct {
	fw.BaseRepository[invoice.InvoiceTenant]
	headerRepo invoice.IInvoiceTitleRepo
	itemRepo   invoice.IInvoiceItemRepo
	recordRepo invoice.IInvoiceRecordRepo
}

// NewInvoiceTenantRepo 创建发票租户仓储
func NewInvoiceTenantRepo(o fw.ORM) invoice.IInvoiceTenantRepo {
	r := &invoiceTenantRepoImpl{}
	r.ORM = o
	return r
}

// Header implements invoice.IInvoiceTenantRepo.
func (i *invoiceTenantRepoImpl) Header() invoice.IInvoiceTitleRepo {
	if i.headerRepo == nil {
		i.headerRepo = NewInvoiceTitleRepo(i.ORM)
	}
	return i.headerRepo
}

// Items implements invoice.IInvoiceTenantRepo.
func (i *invoiceTenantRepoImpl) Items() invoice.IInvoiceItemRepo {
	if i.itemRepo == nil {
		i.itemRepo = NewInvoiceItemRepo(i.ORM)
	}
	return i.itemRepo
}

// Records implements invoice.IInvoiceTenantRepo.
func (i *invoiceTenantRepoImpl) Records() invoice.IInvoiceRecordRepo {
	if i.recordRepo == nil {
		i.recordRepo = NewInvoiceRecordRepo(i.ORM)
	}
	return i.recordRepo
}

// CreateTenant implements invoice.IInvoiceTenantRepo.
func (i *invoiceTenantRepoImpl) CreateTenant(v *invoice.InvoiceTenant) invoice.InvoiceUserAggregateRoot {
	e := i.FindBy("tenant_type=? AND tenant_uid=?", v.TenantType, v.TenantUid)
	if e != nil {
		// 已经存在租户
		return impl.NewInvoiceTenant(e, i)
	}
	t := impl.NewInvoiceTenant(v, i)
	err := t.Create()
	if err != nil {
		logger.Error("创建租户失败: %+v", v)
		return nil
	}
	return t
}

// GetTenant implements invoice.IInvoiceTenantRepo.
func (i *invoiceTenantRepoImpl) GetTenant(id int) invoice.InvoiceUserAggregateRoot {
	v := i.Get(id)
	if v != nil {
		return impl.NewInvoiceTenant(v, i)
	}
	return nil
}

var _ invoice.IInvoiceTitleRepo = new(invoiceHeaderRepoImpl)

type invoiceHeaderRepoImpl struct {
	fw.BaseRepository[invoice.InvoiceTitle]
}

// NewInvoiceTitleRepo 创建发票抬头仓储
func NewInvoiceTitleRepo(o fw.ORM) invoice.IInvoiceTitleRepo {
	r := &invoiceHeaderRepoImpl{}
	r.ORM = o
	return r
}

var _ invoice.IInvoiceRecordRepo = new(invoiceRecordRepoImpl)

type invoiceRecordRepoImpl struct {
	fw.BaseRepository[invoice.InvoiceRecord]
}

// NewInvoiceRecordRepo 创建发票仓储
func NewInvoiceRecordRepo(o fw.ORM) invoice.IInvoiceRecordRepo {
	r := &invoiceRecordRepoImpl{}
	r.ORM = o
	return r
}

var _ invoice.IInvoiceItemRepo = new(invoiceItemRepoImpl)

type invoiceItemRepoImpl struct {
	fw.BaseRepository[invoice.InvoiceItem]
}

// NewInvoiceItemRepo 创建发票项目仓储
func NewInvoiceItemRepo(o fw.ORM) invoice.IInvoiceItemRepo {
	r := &invoiceItemRepoImpl{}
	r.ORM = o
	return r
}
