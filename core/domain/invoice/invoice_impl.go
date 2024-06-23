package invoice

import (
	"errors"
	"fmt"
	"time"

	"github.com/ixre/go2o/core/domain/interface/invoice"
	"github.com/ixre/go2o/core/event/events"
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

// GetInvoiceHeader 获取发票抬头
func (i *invoiceTenantAggregateRootImpl) GetInvoiceHeader(id int) *invoice.InvoiceHeader {
	return i.repo.Header().Get(id)
}

// SaveInvoiceHeader 保存发票抬头
func (i *invoiceTenantAggregateRootImpl) SaveInvoiceHeader(header *invoice.InvoiceHeader) error {
	_, err := i.repo.Header().Save(header)
	return err
}

// CreateInvoice 创建发票
func (i *invoiceTenantAggregateRootImpl) CreateInvoice(record *invoice.InvoiceRecord) invoice.InvoiceDomain {
	return newInvoiceRecord(record, i.repo.Records(), i.repo.Items())
}

// GetInvoice 获取发票
func (i *invoiceTenantAggregateRootImpl) GetInvoice(id int) invoice.InvoiceDomain {
	v := i.repo.Records().Get(id)
	if v != nil {
		return i.CreateInvoice(v)
	}
	return nil
}

var _ invoice.InvoiceDomain = new(invoiceRecordDomainImpl)

type invoiceRecordDomainImpl struct {
	value    *invoice.InvoiceRecord
	repo     invoice.IInvoiceRecordRepo
	itemRepo invoice.IInvoiceItemRepo
}

func newInvoiceRecord(v *invoice.InvoiceRecord, repo invoice.IInvoiceRecordRepo, itemRepo invoice.IInvoiceItemRepo) invoice.InvoiceDomain {
	return &invoiceRecordDomainImpl{
		value:    v,
		repo:     repo,
		itemRepo: itemRepo,
	}
}

// GetDomainId implements invoice.InvoiceDomain.
func (i *invoiceRecordDomainImpl) GetDomainId() int {
	return i.value.Id
}

// GetValue implements invoice.InvoiceDomain.
func (i *invoiceRecordDomainImpl) GetValue(id int) *invoice.InvoiceRecord {
	return types.DeepClone(i.value)
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
