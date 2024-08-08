package query

import (
	"github.com/ixre/go2o/core/domain/interface/work/workorder"
	"github.com/ixre/go2o/core/infrastructure/fw"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
)

type WorkQuery struct {
	fw.BaseRepository[workorder.Workorder]
	commentRepo fw.BaseRepository[workorder.WorkorderComment]
}

func NewWorkQuery(orm fw.ORM) *WorkQuery {
	q := &WorkQuery{}
	q.ORM = orm
	q.commentRepo.ORM = orm
	return q
}

// 查询会员的提交的工单
func (q *WorkQuery) QueryMemberPagingWorkOrder(memberId int, p *fw.PagingParams) (*fw.PagingResult, error) {
	p.Equal("member_id", memberId)
	return q.QueryPaging(p)
	// tables := ""
	// fields := ""
	// return fw.UnifinedQueryPaging(q.ORM,p,tables,fields)
}

// 查询会员工单评论
func (q *WorkQuery) QueryPagingWorkorderComments(workorderId int, p *fw.PagingParams) (*fw.PagingResult, error) {
	p.Equal("order_id", workorderId)
	p.OrderBy("create_time desc")
	ret, err := q.commentRepo.QueryPaging(p)
	for i, v := range ret.Rows {
		d, _ := types.ParseJSONObject(v)
		r := fw.ParsePagingRow(d)
		r.Excludes("orderId")
		ret.Rows[i] = d
	}
	return ret, err
}
