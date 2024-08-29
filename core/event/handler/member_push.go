package handler

import (
	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/event/events"
	"github.com/ixre/go2o/core/event/msq"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/typeconv"
)

func (h *EventHandler) HandleMemberPushEvent(data interface{}) {
	v := data.(*events.MemberPushEvent)
	if v == nil {
		return
	}
	m := v.Member
	ev := &proto.EVMemberPushEventData{
		MemberId:      int64(v.Member.Id),
		IsNewMember:   v.IsCreate,
		UserCode:      m.UserCode,
		UserFlag:      int32(m.UserFlag),
		Role:          int32(m.RoleFlag),
		Username:      m.Username,
		Exp:           int64(m.Exp),
		Level:         int32(m.Level),
		Nickname:      m.Nickname,
		ProfilePhoto:  m.ProfilePhoto,
		Phone:         m.Phone,
		Email:         m.Email,
		RegFrom:       m.RegFrom,
		InviterId:     int64(v.InviterId),
		RealName:      m.RealName,
		LastLoginTime: int64(v.Member.LastLoginTime),
	}

	msq.Push(msq.MemberUpdated, typeconv.MustJson(ev))
}

func (h *EventHandler) HandleMemberAccountPushEvent(data interface{}) {
	v := data.(*events.MemberAccountPushEvent)
	if v == nil {
		return
	}
	pushEnabled := h.registryRepo.Get(registry.MemberAccountPushEnabled).BoolValue()
	if pushEnabled {
		ev := &proto.EVMemberAccountEventData{
			//MemberId:v.MemberId,
			Integral:      int64(v.Integral),
			Balance:       int64(v.Balance),
			WalletCode:    v.WalletCode,
			WalletBalance: int64(v.WalletBalance),
			FlowBalance:   int64(v.FlowBalance),
			GrowBalance:   int64(v.GrowBalance),
			TotalExpense:  int64(v.TotalExpense),
			TotalCharge:   int64(v.TotalCharge),
			UpdateTime:    int64(v.UpdateTime),
		}
		msq.PushDelay(msq.MemberAccountUpdated, typeconv.MustJson(ev), 500)
	}
}

func (h EventHandler) HandleWithdrawalPushEvent(data interface{}) {
	v := data.(*events.WithdrawalPushEvent)
	if v == nil {
		return
	}

	isPush := h.registryRepo.Get(registry.MemberWithdrawalPushEnabled).BoolValue()
	if isPush {
		ev := &proto.EVMemberWithdrawalPushEventData{
			MemberId:       int64(v.MemberId),
			RequestId:      int64(v.RequestId),
			Amount:         int64(v.Amount),
			TransactionFee: int64(v.TransactionFee),
			ReviewResult:   v.ReviewResult,
			IsReviewEvent:  v.IsReviewEvent,
		}
		msq.PushDelay(msq.MembertWithdrawalTopic, typeconv.MustJson(ev), 500)
	}
}

func (h EventHandler) HandleMemberAccountLogPushEvent(data interface{}) {
	v := data.(*events.AccountLogPushEvent)
	if v == nil {
		return
	}
	isPush := h.registryRepo.Get(registry.MemberAccountLogPushEnabled).BoolValue()
	if isPush {
		ev := &proto.EVAccountLogPushEventData{
			Account:        int32(v.Account),
			IsUpdateEvent:  v.IsUpdateEvent,
			UserId:         int64(v.MemberId),
			LogId:          int64(v.LogId),
			LogKind:        int32(v.LogKind),
			Subject:        v.Subject,
			OuterNo:        v.OuterNo,
			ChangeValue:    int64(v.ChangeValue),
			Balance:        int64(v.Balance),
			TransactionFee: int64(v.ProcedureFee),
			ReviewStatus:   int32(v.ReviewStatus),
			CreateTime:     int64(v.CreateTime),
		}
		msq.PushDelay(msq.MemberAccountLogTopic, typeconv.MustJson(ev), 500)
	}
}
