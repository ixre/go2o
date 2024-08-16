package approval

import "github.com/ixre/go2o/core/domain/interface/approval"

var _ approval.IApprovalAggregateRoot = new(staffTransferApprovalImpl)

type staffTransferApprovalImpl struct {
	*ApprovalImpl
}

func NewStaffTransferApproval(value *approval.Approval, repo approval.IApprovalRepository) approval.IApprovalAggregateRoot {
	s := &staffTransferApprovalImpl{}
	s._value = value
	s._repo = repo
	return s
}

func (s *staffTransferApprovalImpl) Flow() approval.ApprovalFlow {

}
