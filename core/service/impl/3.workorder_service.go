package impl

import (
	"context"

	"github.com/ixre/go2o/core/domain/interface/work/workorder"
	"github.com/ixre/go2o/core/service/proto"
)

var _ proto.WorkorderServiceServer = new(workorderServiceImpl)

type workorderServiceImpl struct {
	repo workorder.IWorkorderRepo
	proto.UnimplementedWorkorderServiceServer
	serviceUtil
}

func NewWorkorderService(repo workorder.IWorkorderRepo) proto.WorkorderServiceServer {
	return &workorderServiceImpl{
		repo: repo,
	}
}

// AllocateAgentId implements proto.WorkorderServiceServer.
func (w *workorderServiceImpl) AllocateAgentId(_ context.Context, req *proto.AllocateWorkorderAgentRequest) (*proto.TxResult, error) {
	iw := w.repo.GetWorkorder(int(req.WorkorderId))
	if iw == nil {
		return &proto.TxResult{
			Code:    1,
			Message: "工单不存在",
		}, nil
	}
	err := iw.AllocateAgentId(int(req.UserId))
	return w.ret(err)
}

// Apprise implements proto.WorkorderServiceServer.
func (w *workorderServiceImpl) Apprise(_ context.Context, req *proto.WorkorderAppriseRequest) (*proto.TxResult, error) {
	iw := w.repo.GetWorkorder(int(req.WorkorderId))
	if iw == nil {
		return &proto.TxResult{
			Code:    1,
			Message: "工单不存在",
		}, nil
	}
	err := iw.Apprise(req.IsUsefully, int(req.ServiceRank), req.ServiceApprise)
	return w.ret(err)
}

// Close implements proto.WorkorderServiceServer.
func (w *workorderServiceImpl) Close(_ context.Context, req *proto.WorkorderId) (*proto.TxResult, error) {
	iw := w.repo.GetWorkorder(int(req.WorkorderId))
	if iw == nil {
		return &proto.TxResult{
			Code:    1,
			Message: "工单不存在",
		}, nil
	}
	err := iw.Close()
	return w.ret(err)
}
func (w *workorderServiceImpl) ret(err error) (*proto.TxResult, error) {
	if err != nil {
		return &proto.TxResult{
			Code:    1,
			Message: err.Error(),
		}, nil
	}
	return &proto.TxResult{}, nil
}

// DeleteWorkorder implements proto.WorkorderServiceServer.
func (w *workorderServiceImpl) DeleteWorkorder(_ context.Context, req *proto.WorkorderId) (*proto.TxResult, error) {
	// iw := w.repo.GetWorkorder(int(req.WorkorderId))
	// if iw == nil {
	// 	return &proto.TxResult{
	// 		Code: 1,
	// 		Msg:  "工单不存在",
	// 	}, nil
	// }
	// err := iw.AllocateAgentId(int(req.UserId))
	// return w.ret(err)
	panic("unimplemented")
}

// Finish implements proto.WorkorderServiceServer.
func (w *workorderServiceImpl) Finish(_ context.Context, req *proto.WorkorderId) (*proto.TxResult, error) {
	iw := w.repo.GetWorkorder(int(req.WorkorderId))
	if iw == nil {
		return &proto.TxResult{
			Code:    1,
			Message: "工单不存在",
		}, nil
	}
	err := iw.Finish()
	return w.ret(err)
}

// GetWorkorder implements proto.WorkorderServiceServer.
func (w *workorderServiceImpl) GetWorkorder(_ context.Context, req *proto.WorkorderId) (*proto.SWorkorder, error) {
	iw := w.repo.GetWorkorder(int(req.WorkorderId))
	if iw == nil {
		return nil, nil
	}
	iv := iw.Value()
	return &proto.SWorkorder{
		Id:             int64(iv.Id),
		OrderNo:        iv.OrderNo,
		MemberId:       int64(iv.MemberId),
		ClassId:        int32(iv.ClassId),
		MchId:          int64(iv.MchId),
		Flag:           int32(iv.Flag),
		Wip:            iv.Wip,
		Subject:        iv.Subject,
		Content:        iv.Content,
		IsOpened:       int32(iv.IsOpened),
		HopeDesc:       iv.HopeDesc,
		FirstPhoto:     iv.FirstPhoto,
		PhotoList:      iv.PhotoList,
		ContactWay:     iv.ContactWay,
		Status:         int32(iv.Status),
		AllocateAid:    int64(iv.AllocateAid),
		ServiceRank:    int32(iv.ServiceRank),
		ServiceApprise: iv.ServiceApprise,
		IsUsefully:     int32(iv.IsUsefully),
		CreateTime:     int64(iv.CreateTime),
		UpdateTime:     int64(iv.UpdateTime),
	}, nil
}

// SubmitComment implements proto.WorkorderServiceServer.
func (w *workorderServiceImpl) SubmitComment(_ context.Context, req *proto.SubmitWorkorderCommentRequest) (*proto.TxResult, error) {
	iw := w.repo.GetWorkorder(int(req.WorkorderId))
	if iw == nil {
		return &proto.TxResult{
			Code:    1,
			Message: "工单不存在",
		}, nil
	}
	err := iw.SubmitComment(req.Content, req.IsReplay, int(req.RefCommentId))
	return w.ret(err)
}

// SubmitWorkorder implements proto.WorkorderServiceServer.
func (w *workorderServiceImpl) SubmitWorkorder(_ context.Context, req *proto.SubmitWorkorderRequest) (*proto.SubmitWorkorderResponse, error) {
	dst := &workorder.Workorder{
		MemberId:   int(req.MemberId),
		ClassId:    int(req.ClassId),
		MchId:      int(req.MchId),
		Wip:        req.Wip,
		Subject:    req.Subject,
		Content:    req.Content,
		IsOpened:   int(req.IsOpened),
		HopeDesc:   req.HopeDesc,
		PhotoList:  req.PhotoList,
		ContactWay: req.ContactWay,
	}
	wo := w.repo.CreateWorkorder(dst)
	err := wo.Submit()

	ret := &proto.SubmitWorkorderResponse{
		WorkorderId: int64(wo.GetAggregateRootId()),
	}
	if err != nil {
		ret.Code = 1
		ret.Message = err.Error()
	}
	return ret, nil
}
