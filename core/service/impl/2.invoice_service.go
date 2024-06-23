package impl

import (
	"context"

	"github.com/ixre/go2o/core/service/proto"
)

var _ proto.InvoiceServiceServer = new(invoiceServiceImpl)

type invoiceServiceImpl struct {
	_ proto.InvoiceServiceServer
	proto.UnimplementedInvoiceServiceServer
}

func NewInvoiceService() proto.InvoiceServiceServer {
	return new(invoiceServiceImpl)
}

// DeleteTenant implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) DeleteTenant(context.Context, *proto.InvoiceTenantId) (*proto.Result, error) {
	panic("unimplemented")
}

// GetTenant implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) GetTenant(context.Context, *proto.InvoiceTenantId) (*proto.STenant, error) {
	panic("unimplemented")
}

// PagingTenant implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) PagingTenant(context.Context, *proto.TenantPagingRequest) (*proto.TenantPagingResponse, error) {
	panic("unimplemented")
}

// QueryTenantList implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) QueryTenantList(context.Context, *proto.QueryTenantRequest) (*proto.QueryTenantResponse, error) {
	panic("unimplemented")
}

// SaveTenant implements proto.InvoiceServiceServer.
func (i *invoiceServiceImpl) SaveTenant(context.Context, *proto.SaveTenantRequest) (*proto.SaveTenantResponse, error) {
	panic("unimplemented")
}
