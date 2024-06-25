package impl

import (
	"context"

	"github.com/ixre/go2o/core/domain/interface/invoice"
	"github.com/ixre/go2o/core/infrastructure/logger"
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

// 获取发票租户
func (i *invoiceServiceImpl) GetTenant(_ context.Context, req *proto.InvoiceTenantRequest) (*proto.SInvoiceTenant, error) {
	tenant := i.repo.CreateTenant(&invoice.InvoiceTenant{
		TenantType: int(req.TenantType),
		TenantUid:  int(req.TenantUid),
	})
	if tenant == nil {
		return &proto.SInvoiceTenant{
			ErrCode: 1,
			ErrMsg:  "无法创建租户",
		}, nil
	}
	return &proto.SInvoiceTenant{
		TenantId: int64(tenant.GetAggregateRootId()),
	}, nil
}

// CreateInvoice implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) CreateInvoice(_ context.Context, req *proto.SaveRecordRequest) (*proto.SaveRecordResponse, error) {
	tenant := i.repo.GetTenant(int(req.TenantId))
	if tenant == nil {
		return &proto.SaveRecordResponse{
			ErrCode: 2,
			ErrMsg:  "无法申请发票",
		}, nil
	}
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

func (i *invoiceServiceImpl) getInvoice(tenantId, invoiceId int64) (invoice.InvoiceUserAggregateRoot, invoice.InvoiceDomain) {
	t := i.repo.GetTenant(int(tenantId))
	if t == nil {
		logger.Error("no such invoice tenant, data=%d", tenantId)
		return nil, nil
	}
	iv := t.GetInvoice(int(invoiceId))
	if iv == nil {
		logger.Error("no such invoice, data=%d", invoiceId)
		return nil, nil
	}
	return t, iv
}

// GetInvoice implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) GetInvoice(_ context.Context, req *proto.InvoiceId) (*proto.SInvoice, error) {
	it, iv := i.getInvoice(req.TenantId, req.InvoiceId)
	if it == nil || iv == nil {
		return nil, nil
	}
	v := iv.GetValue()
	ret := &proto.SInvoice{
		Id:               int64(v.Id),
		InvoiceCode:      v.InvoiceCode,
		InvoiceNo:        v.InvoiceNo,
		TenantId:         int64(v.TenantId),
		IssueTenantId:    int64(v.IssueTenantId),
		InvoiceType:      int32(v.InvoiceType),
		IssueType:        int32(v.IssueType),
		SellerName:       v.SellerName,
		SellerTaxCode:    v.SellerTaxCode,
		PurchaserName:    v.PurchaserName,
		PurchaserTaxCode: v.PurchaserTaxCode,
		InvoiceAmount:    v.InvoiceAmount,
		TaxAmount:        v.TaxAmount,
		Remark:           v.Remark,
		IssueRemark:      v.IssueRemark,
		InvoicePic:       v.InvoicePic,
		ReceiveEmail:     v.ReceiveEmail,
		InvoiceStatus:    int32(v.InvoiceStatus),
		InvoiceTime:      int64(v.InvoiceTime),
		CreateTime:       int64(v.CreateTime),
		UpdateTime:       int64(v.UpdateTime),
		Items:            []*proto.SInvoiceItem{},
	}
	for _, v := range iv.GetItems() {
		ret.Items = append(ret.Items, &proto.SInvoiceItem{
			ItemName:  v.ItemName,
			ItemSpec:  v.ItemSpec,
			Price:     v.Price,
			Quantity:  int32(v.Quantity),
			TaxRate:   v.TaxRate,
			Unit:      v.Unit,
			Amount:    v.Amount,
			TaxAmount: v.TaxAmount,
		})
	}
	return ret, nil
}

// Issue implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) Issue(_ context.Context, req *proto.InvoiceIssueRequest) (*proto.Result, error) {
	it, iv := i.getInvoice(req.TenantId, req.InvoiceId)
	if it == nil || iv == nil {
		return &proto.Result{
			ErrCode: 1,
			ErrMsg:  "no any tenant or invoice",
		}, nil
	}
	err := iv.Issue(req.InvoicePic)
	return i.error(err), nil
}

// IssueFail implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) IssueFail(_ context.Context, req *proto.InvoiceIssueFailRequest) (*proto.Result, error) {
	it, iv := i.getInvoice(req.TenantId, req.InvoiceId)
	if it == nil || iv == nil {
		return &proto.Result{
			ErrCode: 1,
			ErrMsg:  "no any tenant or invoice",
		}, nil
	}
	err := iv.IssueFail(req.Reason)
	return i.error(err), nil
}

// Revert implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) Revert(_ context.Context, req *proto.InvoiceRevertRequest) (*proto.Result, error) {
	it, iv := i.getInvoice(req.TenantId, req.InvoiceId)
	if it == nil || iv == nil {
		return &proto.Result{
			ErrCode: 1,
			ErrMsg:  "no any tenant or invoice",
		}, nil
	}
	err := iv.Revert(req.Reason)
	return i.error(err), nil
}

// SaveHeader implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) SaveHeader(_ context.Context, req *proto.SaveHeaderRequest) (*proto.SaveHeaderResponse, error) {
	t := i.repo.GetTenant(int(req.TenantId))
	if t == nil {
		logger.Error("no such invoice tenant, data=%d", req.TenantId)
		return nil, nil
	}
	v := &invoice.InvoiceHeader{
		Id:          int(req.Id),
		InvoiceType: int(req.InvoiceType),
		IssueType:   int(req.IssueType),
		HeaderName:  req.HeaderName,
		TaxCode:     req.TaxCode,
		SignAddress: req.SignAddress,
		SignTel:     req.SignTel,
		BankName:    req.BankName,
		BankAccount: req.BankAccount,
		Remarks:     req.Remarks,
		IsDefault:   int(req.GetIsDefault()),
	}
	err := t.SaveInvoiceHeader(v)
	if err != nil {
		return &proto.SaveHeaderResponse{
			ErrCode: 1,
			ErrMsg:  err.Error(),
		}, nil
	}
	return &proto.SaveHeaderResponse{
		Id: int64(v.Id),
	}, nil
}

// SendMail implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) SendMail(_ context.Context, req *proto.InvoiceSendMailRequest) (*proto.Result, error) {
	it, iv := i.getInvoice(req.TenantId, req.InvoiceId)
	if it == nil || iv == nil {
		return &proto.Result{
			ErrCode: 1,
			ErrMsg:  "no any tenant or invoice",
		}, nil
	}
	err := iv.SendMail(req.Email)
	return i.error(err), nil
}
