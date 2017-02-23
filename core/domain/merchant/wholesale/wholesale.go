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

// 保存客户分组的批发返点率
func (w *wholesalerImpl) SaveGroupRebateRate(groupId int32, arr []*wholesaler.WsRebateRate) error {
	// 获取存在的项
	old := w.GetGroupRebateRate(groupId)
	// 分析当前数据并加入到MAP中
	delList := []int32{}
	currMap := make(map[int32]*wholesaler.WsRebateRate, len(arr))
	for _, v := range arr {
		currMap[v.RequireAmount] = v
	}
	// 筛选出要删除的项,如存在，则赋予ID
	for _, v := range old {
		new := currMap[v.RequireAmount]
		if new == nil {
			delList = append(delList, v.ID)
		} else {
			new.ID = v.ID
		}
	}
	// 删除项
	for _, v := range delList {
		w.repo.BatchDeleteWsRebateRate("id=?", v)
	}
	// 保存项
	for _, v := range arr {
		v.WsId = w.mchId
		v.BuyerGid = groupId
		i, err := util.I32Err(w.repo.SaveWsRebateRate(v))
		if err == nil {
			v.ID = i
		}
	}
	return nil
}

// 获取客户分组的批发返点率
func (w *wholesalerImpl) GetGroupRebateRate(groupId int32) []*wholesaler.WsRebateRate {
	return w.repo.SelectWsRebateRate("ws_id=? AND buyer_gid=?", w.mchId, groupId)
}

// 获取批发返点率
func (w *wholesalerImpl) GetRebateRate(groupId int32, amount int32) float64 {
	var disRate float64 = 0
	arr := w.GetGroupRebateRate(groupId)
	if len(arr) > 0 {
		var maxRequire int32
		for _, v := range arr {
			if v.RequireAmount > maxRequire && amount >= v.RequireAmount {
				maxRequire = v.RequireAmount
				disRate = v.RebateRate
			}
		}
	}
	return disRate
}
