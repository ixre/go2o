package impl

import (
	"context"

	"github.com/ixre/go2o/core/domain/interface/invoice"
	"github.com/ixre/go2o/core/service/proto"
)

var _ proto.InvoiceServiceServer = new(invoiceServiceImpl)

type invoiceServiceImpl struct {
	_    proto.InvoiceServiceServer
	repo invoice.IInvoiceTenantRepo
	proto.UnimplementedInvoiceServiceServer
	serviceUtil
}

func NewInvoiceService(repo invoice.IInvoiceTenantRepo) proto.InvoiceServiceServer {
	return &invoiceServiceImpl{
		repo: repo,
	}
}

// CreateRecord implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) CreateRecord(_ context.Context, req *proto.SaveRecordRequest) (*proto.SaveRecordResponse, error) {
	tenant := i.repo.CreateTenant(&invoice.InvoiceTenant{
		TenantType: int(req.TenantType),
		TenantUid:  int(req.TenantUid),
	})
	iv := tenant.CreateInvoice(&invoice.InvoiceRecord{
		IssueTenantId:    int(req.IssueTenantId),
		InvoiceType:      int(req.InvoiceType),
		IssueType:        int(req.IssueType),
		PurchaserName:    req.PurchaserName,
		PurchaserTaxCode: req.PurchaserTaxCode,
		Remark:           req.Remark,
		ReceiveEmail:     req.ReceiveEmail,
	})
	err := iv.Save()
	if err != nil {
		return &proto.SaveRecordResponse{
			ErrCode: 1,
			ErrMsg:  err.Error(),
		}, nil
	}
	return &proto.SaveRecordResponse{
		Id: int64(iv.GetDomainId()),
	}, nil
}

// GetRecord implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) GetRecord(context.Context, *proto.InvoiceRecordId) (*proto.SRecord, error) {
	panic("unimplemented")
}

// Issue implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) Issue(context.Context, *proto.InvoiceIssueRequest) (*proto.Result, error) {
	panic("unimplemented")
}

// IssueFail implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) IssueFail(context.Context, *proto.InvoiceIssueFailRequest) (*proto.Result, error) {
	panic("unimplemented")
}

// Revert implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) Revert(context.Context, *proto.InvoiceRevertRequest) (*proto.Result, error) {
	panic("unimplemented")
}

// SaveHeader implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) SaveHeader(context.Context, *proto.SaveHeaderRequest) (*proto.SaveHeaderResponse, error) {
	panic("unimplemented")
}

// SendMail implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) SendMail(context.Context, *proto.InvoiceSendMailRequest) (*proto.Result, error) {
	panic("unimplemented")
}

// mustEmbedUnimplementedInvoiceServiceServer implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) mustEmbedUnimplementedInvoiceServiceServer() {
	panic("unimplemented")
}
