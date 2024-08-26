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

// 查询工单列表
func (q *WorkQuery) QueryPagingWorkorders(p *fw.PagingParams) (*fw.PagingResult, error) {
	tabels := `workorder w
			   INNER JOIN mm_member m ON w.member_id = m.id
			   LEFT JOIN rbac_user u ON w.allocate_aid= u.id`
	fields := "w.*,m.nickname as nickname,m.profile_photo,u.nickname as allocate_agent_name"
	ret, err := fw.UnifinedQueryPaging(q.ORM, p, tabels, fields)
	if err != nil {
		return nil, err
	}
	for _, v := range ret.Rows {
		r := fw.ParsePagingRow(v)
		r.Excludes("content")
	}
	return ret, nil
}

// 查询工单最新评论
func (q *WorkQuery) QueryLatestWorkorderComments(workorderId int, p *fw.PagingParams) []*workorder.WorkorderComment {
	p.Equal("order_id", workorderId)
	p.OrderBy("id asc")
	return q.commentRepo.FindList(&fw.QueryOption{
		Limit: p.Size,
		Order: p.Order,
	}, p.Arguments[0].(string), p.Arguments[1:])
}
