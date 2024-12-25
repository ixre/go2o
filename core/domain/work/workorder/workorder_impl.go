package workorder

import (
	"errors"
	"time"

	"github.com/ixre/go2o/core/domain/interface/work/workorder"
	"github.com/ixre/go2o/core/infrastructure/domain"
	"github.com/ixre/go2o/core/infrastructure/fw/types"
)

var _ workorder.IWorkorderAggregateRoot = new(workorderAggregateRootImpl)

// 假设这是实现了IWorkorderAggregateRoot接口的结构体
type workorderAggregateRootImpl struct {
	repo  workorder.IWorkorderRepo
	value *workorder.Workorder
}

func NewWorkorder(value *workorder.Workorder, repo workorder.IWorkorderRepo) workorder.IWorkorderAggregateRoot {
	return &workorderAggregateRootImpl{
		repo:  repo,
		value: value,
	}
}

// Submit implements workorder.IWorkorderAggregateRoot.
func (w *workorderAggregateRootImpl) Submit() error {
	if w.GetAggregateRootId() > 0 {
		return errors.New("workorder has been submitted")
	}
	if len(w.value.Subject) == 0 {
		return errors.New("workorder subject is empty")
	}
	if len(w.value.Content) == 0 {
		return errors.New("workorder content is empty")
	}
	if w.value.MemberId == 0 {
		return errors.New("workorder member is empty")
	}
	if w.value.ClassId == workorder.ClassSuggest {
		// 建议不开放评论
		w.value.IsOpened = 0
	} else if w.value.ClassId == workorder.ClassAppeal {
		// 投诉允许评论
		w.value.IsOpened = 1
	}

	w.value.OrderNo = domain.NewTradeNo(7, w.value.MemberId)
	w.value.Status = workorder.StatusPending
	w.value.CreateTime = int(time.Now().Unix())
	w.value.UpdateTime = int(time.Now().Unix())
	_, err := w.repo.Save(w.value)
	return err
}

// AllocateAgentId implements workorder.IWorkorderAggregateRoot.
func (w *workorderAggregateRootImpl) AllocateAgentId(userId int) error {
	if w.value.AllocateAid > 0 {
		return errors.New("agent has been allocated")
	}
	w.value.AllocateAid = userId
	if w.value.Status == workorder.StatusPending {
		w.value.Status = workorder.StatusProcessing
	}
	_, err := w.repo.Save(w.value)
	return err
}

// Apprise implements workorder.IWorkorderAggregateRoot.
func (w *workorderAggregateRootImpl) Apprise(isUsefully bool, rank int, apprise string) error {
	if w.value.ServiceRank > 0 {
		return errors.New("workorder has been apprised")
	}
	if w.value.Status != workorder.StatusFinished {
		return errors.New("workorder is not finished")
	}
	if (w.value.Flag & workorder.FlagUserClosed) == workorder.FlagUserClosed {
		// 用户关闭后, 不能评价
		return errors.New("workorder closed by user, can not apprise")
	}
	w.value.ServiceRank = rank
	w.value.IsUsefully = types.Ternary(isUsefully, 1, 0)
	w.value.ServiceApprise = apprise
	w.value.UpdateTime = int(time.Now().Unix())
	_, err := w.repo.Save(w.value)
	return err
}

// Close implements workorder.IWorkorderAggregateRoot.
func (w *workorderAggregateRootImpl) Close() error {
	if w.value.Status == workorder.StatusFinished {
		return errors.New("workorder has been closed")
	}
	// 用户关闭
	w.value.Status = workorder.StatusFinished
	w.value.Flag |= workorder.FlagUserClosed
	w.value.IsUsefully = 1
	w.value.UpdateTime = int(time.Now().Unix())
	_, err := w.repo.Save(w.value)
	return err
}

// Finish implements workorder.IWorkorderAggregateRoot.
func (w *workorderAggregateRootImpl) Finish() error {
	if w.value.Status == workorder.StatusFinished {
		return errors.New("workorder has been finished")
	}
	w.value.Status = workorder.StatusFinished
	w.value.UpdateTime = int(time.Now().Unix())
	_, err := w.repo.Save(w.value)
	if err == nil {
		err = w.SubmitComment("本次服务已经结束,感谢您对我们工作的支持,祝您生活愉快!", true, 0)
	}
	return err
}

// GetAggregateRootId implements workorder.IWorkorderAggregateRoot.
func (w *workorderAggregateRootImpl) GetAggregateRootId() int {
	return w.value.Id
}

// SubmitComment implements workorder.IWorkorderAggregateRoot.
func (w *workorderAggregateRootImpl) SubmitComment(content string, isReplay bool, refCommentId int) error {
	if w.value.ClassId == workorder.ClassSuggest && !isReplay {
		// 建议不允许用户提交评论
		return errors.New("suggest workorder only replay can be submitted")
	}
	if w.value.IsOpened == 0 && !isReplay {
		// 未开放评论
		return errors.New("工单未开放评论")
	}
	comment := &workorder.WorkorderComment{
		Id:         refCommentId,
		OrderId:    w.GetAggregateRootId(),
		IsReplay:   types.Ternary(isReplay, 1, 0),
		Content:    content,
		IsRevert:   0,
		RefCid:     refCommentId,
		CreateTime: int(time.Now().Unix()),
	}
	_, err := w.repo.CommentRepo().Save(comment)
	return err
}

// Value implements workorder.IWorkorderAggregateRoot.
func (w *workorderAggregateRootImpl) Value() *workorder.Workorder {
	return types.DeepClone(w.value)
}
