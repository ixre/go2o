package query

import "github.com/ixre/go2o/core/infrastructure/fw"

type InvoiceQuery struct {
	orm fw.ORM
}

func NewInvoiceQuery(o fw.ORM) *InvoiceQuery {
	return &InvoiceQuery{
		orm: o,
	}
}

func (i *InvoiceQuery) QueryMerchantIssueInvoices(p *fw.PagingParams) (*fw.PagingResult, error) {
	tables := ""
	fields := ""
	return fw.UnifinedQueryPaging(i.orm, p, tables, fields)
}
