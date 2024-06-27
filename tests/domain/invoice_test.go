package domain

import (
	"testing"

	"github.com/ixre/go2o/core/domain/interface/invoice"
	"github.com/ixre/go2o/core/inject"
)

// 添加发票抬头
func TestInvoiceTitle(t *testing.T) {
	r := inject.GetInvoiceTenantRepo()
	tn := r.CreateTenant(&invoice.InvoiceTenant{
		TenantType: int(invoice.TenantUser),
		TenantUid:  1,
	})
	err := tn.CreateInvoiceTitle(&invoice.InvoiceTitle{
		InvoiceType: 1,
		IssueType:   2,
		HeaderName:  "上海丁丁网络科技有限公司",
		TaxCode:     "64443223446656622",
		SignAddress: "",
		SignTel:     "",
		BankName:    "",
		BankAccount: "",
		Remarks:     "",
		IsDefault:   1,
	})
	if err != nil {
		t.Fatal(err)
	}
}
