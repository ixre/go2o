package invoice

import (
	"errors"
	"fmt"
	"time"

	"github.com/ixre/go2o/core/domain/interface/invoice"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
	"github.com/ixre/gof/domain/eventbus"
)

var _ invoice.InvoiceUserAggregateRoot = new(invoiceTenantAggregateRootImpl)

type invoiceTenantAggregateRootImpl struct {
	value *invoice.InvoiceTenant
	repo  invoice.IInvoiceTenantRepo
}

func NewInvoiceTenant(v *invoice.InvoiceTenant, repo invoice.IInvoiceTenantRepo) invoice.InvoiceUserAggregateRoot {
	return &invoiceTenantAggregateRootImpl{
		value: v,
		repo:  repo,
	}
}

func (i *invoiceTenantAggregateRootImpl) GetAggregateRootId() int {
	return i.value.Id
}

// TeantType 获取租户类型
func (i *invoiceTenantAggregateRootImpl) TenantType() invoice.TenantType {
	return invoice.TenantType(i.value.TenantType)
}

// TenantUserId 获取租户用户编号
func (i *invoiceTenantAggregateRootImpl) TenantUserId() int {
	return i.value.TenantUid
}

// Create 创建租户
func (i *invoiceTenantAggregateRootImpl) Create() error {
	if i.GetAggregateRootId() > 0 {
		return errors.New("invoice tenant has been created")
	}
	_, err := i.repo.Save(i.value)
	return err
}

// GetInvoiceTitle 获取发票抬头
func (i *invoiceTenantAggregateRootImpl) GetInvoiceTitle(id int) *invoice.InvoiceTitle {
	return i.repo.Header().Get(id)
}

// CreateInvoiceTitle 新增发票抬头
func (i *invoiceTenantAggregateRootImpl) CreateInvoiceTitle(header *invoice.InvoiceTitle) error {
	if header.Id > 0 {
		return errors.New("invoice header has been created")
	}
	if header.TenantId > 0 && header.TenantId != i.GetAggregateRootId() {
		return errors.New("invoice tenant error")
	}
	header.CreateTime = int(time.Now().Unix())
	header.TenantId = i.GetAggregateRootId()
	_, err := i.repo.Header().Save(header)
	return err
}

// CreateInvoice 创建发票
func (i *invoiceTenantAggregateRootImpl) RequestInvoice(v *invoice.InvoiceRequestData) (invoice.InvoiceDomain, error) {
	r := &invoice.InvoiceRecord{
		InvoiceNo:     "T" + domain.NewTradeNo(11, i.GetAggregateRootId()),
		TenantId:      i.GetAggregateRootId(),
		Remark:        v.Remark,
		IssueRemark:   "",
		InvoicePic:    "",
		ReceiveEmail:  v.ReceiveEmail,
		InvoiceStatus: invoice.IssueAwaiting,
	}
	// 申请人信息
	h := i.repo.Header().Get(v.HeaderId)
	if h == nil || h.Id <= 0 || h.TenantId != i.GetAggregateRootId() {
		return nil, errors.New("invoice header is error")
	}
	r.PurchaserName = h.HeaderName
	r.PurchaserTaxCode = h.TaxCode
	r.IssueType = h.IssueType
	r.InvoiceType = h.InvoiceType
	// 开具人信息
	if v.IssueTenantId == 0 {
		r.SellerName = "系统"
		r.SellerTaxCode = ""
	} else {
		tn := i.repo.GetTenant(v.IssueTenantId)
		if tn == nil {
			return nil, errors.New("issue tenant not exists")
		}
		r.SellerName = "商户"
		r.SellerTaxCode = ""
	}
	// 开票项目
	if len(v.Items) == 0 {
		return nil, errors.New("no such invoice items")
	}
	for _, v := range v.Items {
		amount := v.Price * float64(v.Quantity)
		r.InvoiceAmount += amount
		r.TaxAmount += amount * v.TaxRate
	}
	if r.InvoiceAmount <= 0 {
		return nil, errors.New("invoice amount is zero")
	}
	r.CreateTime = int(time.Now().Unix())
	r.IssueTenantId = v.IssueTenantId
	return i.createInvoice(r, v.Items), nil
}

func (i *invoiceTenantAggregateRootImpl) createInvoice(v *invoice.InvoiceRecord, items []*invoice.InvoiceItem) invoice.InvoiceDomain {
	return newInvoiceRecord(v, items, i.repo.Records(), i.repo.Items())
}

// GetInvoice 获取发票
func (i *invoiceTenantAggregateRootImpl) GetInvoice(id int) invoice.InvoiceDomain {
	v := i.repo.Records().Get(id)
	if v != nil {
		items := i.repo.Items().FindList(nil, "invoice_id=?", id)
		return i.createInvoice(v, items)
	}
	return nil
}

var _ invoice.InvoiceDomain = new(invoiceRecordDomainImpl)

type invoiceRecordDomainImpl struct {
	value    *invoice.InvoiceRecord
	repo     invoice.IInvoiceRecordRepo
	itemRepo invoice.IInvoiceItemRepo
	_items   []*invoice.InvoiceItem
}

func newInvoiceRecord(v *invoice.InvoiceRecord,
	items []*invoice.InvoiceItem,
	repo invoice.IInvoiceRecordRepo, itemRepo invoice.IInvoiceItemRepo) invoice.InvoiceDomain {
	return &invoiceRecordDomainImpl{
		value:    v,
		repo:     repo,
		itemRepo: itemRepo,
		_items:   items,
	}
}

// GetDomainId implements invoice.InvoiceDomain.
func (i *invoiceRecordDomainImpl) GetDomainId() int {
	return i.value.Id
}

// GetValue implements invoice.InvoiceDomain.
func (i *invoiceRecordDomainImpl) GetValue() *invoice.InvoiceRecord {
	return types.DeepClone(i.value)
}

// GetItems implements invoice.InvoiceDomain.
func (i *invoiceRecordDomainImpl) GetItems() []*invoice.InvoiceItem {
	if i._items == nil {
		i._items = i.itemRepo.FindList(nil, "invoice_id=?", i.GetDomainId())
	}
	return i._items
}

// Issue implements invoice.InvoiceDomain.
func (i *invoiceRecordDomainImpl) Issue(picture string) error {
	if i.value.InvoiceStatus != invoice.IssueAwaiting {
		return errors.New("invoice status error")
	}
	i.value.InvoiceStatus = invoice.IssueSuccess
	i.value.InvoicePic = picture
	i.value.IssueRemark = ""
	i.value.UpdateTime = int(time.Now().Unix())
	return i.Save()
}

// IssueFail implements invoice.InvoiceDomain.
func (i *invoiceRecordDomainImpl) IssueFail(reason string) error {
	i.value.InvoiceStatus = invoice.IssueFail
	i.value.IssueRemark = reason
	i.value.UpdateTime = int(time.Now().Unix())
	return i.Save()
}

// Revert implements invoice.InvoiceDomain.
func (i *invoiceRecordDomainImpl) Revert(reason string) error {
	if i.value.InvoiceStatus == invoice.IssueRevert {
		return errors.New("invoice status error")
	}
	return errors.New("not implemented")
}

// Save implements invoice.InvoiceDomain.
func (i *invoiceRecordDomainImpl) Save() error {
	_, err := i.repo.Save(i.value)
	return err
}

// SendMail implements invoice.InvoiceDomain.
func (i *invoiceRecordDomainImpl) SendMail(mail string) error {
	if len(i.value.InvoicePic) == 0 {
		return errors.New("invoice picture is empty")
	}
	eventbus.Publish(&events.SendEmailEvent{
		Subject: "请查收发票",
		To:      mail,
		Body:    fmt.Sprintf(`<img src="%s" alt="发票图片"/>`, i.value.InvoicePic),
	})
	return nil
}
