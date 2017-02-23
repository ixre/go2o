package wholesaler

import (
	"errors"
	"github.com/jsix/gof/util"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/merchant/wholesaler"
)

var _ wholesaler.IWholesaler = new(wholesalerImpl)

type wholesalerImpl struct {
	mchId int32
	value *wholesaler.WsWholesaler
	repo  wholesaler.IWholesaleRepo
}

func NewWholesaler(mchId int32, v *wholesaler.WsWholesaler,
	repo wholesaler.IWholesaleRepo) wholesaler.IWholesaler {
	return &wholesalerImpl{
		mchId: mchId,
		value: v,
		repo:  repo,
	}
}

// 获取领域编号
func (w *wholesalerImpl) GetDomainId() int32 {
	return w.mchId
}

// 获取值
func (w *wholesalerImpl) Value() *wholesaler.WsWholesaler {
	return w.value
}

// 审核批发商
func (w *wholesalerImpl) Review(pass bool, reason string) error {
	if w.value.ReviewState == enum.ReviewAwaiting {
		if pass {
			w.value.ReviewState = enum.ReviewPass
		} else {
			w.value.ReviewState = enum.ReviewReject
		}
		_, err := w.Save()
		return err
	}
	return errors.New("review state not awaiting review!")
}

// 停止批发权限
func (w *wholesalerImpl) Abort() error {
	w.value.ReviewState = enum.ReviewAbort
	_, err := w.Save()
	return err
}

// 保存
func (w *wholesalerImpl) Save() (int32, error) {
	return util.I32Err(w.repo.SaveWsWholesaler(w.value, false))
}

// 保存批发返点率
func (w *wholesalerImpl) SaveRebateRate(v *wholesaler.WsRebateRate) (int32, error) {
	if v.WsId != w.mchId {
		return 0, errors.New("wsid not match")
	}
	list := w.repo.SelectWsRebateRate(
		"ws_id=? AND buyer_gid=? AND require_amount=?",
		w.mchId, v.BuyerGid, v.RequireAmount)
	if len(list) > 0 {
		if rId := list[0].ID; v.ID != rId {
			if v.ID <= 0 {
				v.ID = rId
			} else {
				return 0, errors.New("require amount exists")
			}
		}
	}
	return util.I32Err(w.repo.SaveWsRebateRate(v))
}

// 获取批发返点率
func (w *wholesalerImpl) RebateRates() []*wholesaler.WsRebateRate {
	return w.repo.SelectWsRebateRate("ws_id=?", w.mchId)
}
