package wholesaler

import (
	"errors"
	"github.com/ixre/gof/util"
	"go2o/core/domain/interface/enum"
	"go2o/core/domain/interface/item"
	"go2o/core/domain/interface/merchant/wholesaler"
	"log"
)

var _ wholesaler.IWholesaler = new(wholesalerImpl)

type wholesalerImpl struct {
	mchId    int32
	value    *wholesaler.WsWholesaler
	repo     wholesaler.IWholesaleRepo
	itemRepo item.IGoodsItemRepo
}

func NewWholesaler(mchId int32, v *wholesaler.WsWholesaler,
	repo wholesaler.IWholesaleRepo, itemRepo item.IGoodsItemRepo) wholesaler.IWholesaler {
	return &wholesalerImpl{
		mchId:    mchId,
		value:    v,
		repo:     repo,
		itemRepo: itemRepo,
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

// 同步商品
func (w *wholesalerImpl) SyncItems(syncPrice bool) map[string]int32 {
	add := []int{}
	failed := []int{}
	newest := w.repo.GetAwaitSyncItems(w.mchId)
	for _, itemId := range newest {
		err := w.syncSingleItem(int64(itemId), syncPrice)
		if err == nil {
			add = append(add, itemId)
		} else {
			failed = append(failed, itemId)
		}
	}
	_, del := w.repo.SyncItems(w.mchId, item.ShelvesInWarehouse, enum.ReviewPass)
	return map[string]int32{
		"add":    int32(len(add)),
		"del":    int32(del),
		"failed": int32(len(failed)),
	}
}

func (w *wholesalerImpl) syncSingleItem(itemId int64, syncPrice bool) error {
	it := w.itemRepo.GetItem(itemId)
	if it != nil {
		ws := it.Wholesale()
		_, err := ws.Save()
		if err != nil {
			return err
		}
		for _, v := range it.SkuArray() {
			if v.ID <= 0 {
				continue
			}
			err = ws.SaveSkuPrice(v.ID, []*item.WsSkuPrice{
				{
					ItemId:          itemId,
					SkuId:           v.ID,
					RequireQuantity: 1,
					WholesalePrice:  float64(v.Price),
				},
			})
			if err != nil {
				log.Println("[ Go2o][ Wholesale][ Sync]:", err.Error(),
					"ID:", itemId, "; SkuId:", v.ID)
			}
		}
	}
	return nil
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
		w.repo.BatchDeleteWsRebateRate("id= $1", v)
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
	return w.repo.SelectWsRebateRate("ws_id= $1 AND buyer_gid= $2", w.mchId, groupId)
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
