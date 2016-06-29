/**
 * Copyright 2015 @ z3q.net.
 * name : snapshot
 * author : jarryliu
 * date : 2016-06-28 23:52
 * description :
 * history :
 */
package goods

import "go2o/core/domain/interface/sale/goods"

//var _ goods.ISnapshotManager = new(snapshotManagerImpl)
type snapshotManagerImpl struct {
    _rep     goods.IGoodsRep
    _skuId int
}

func NewSnapshotManagerImpl(skuId int,rep goods.IGoodsRep)goods.ISnapshotManager{
    return &snapshotManagerImpl{
        _rep:rep,
        _skuId:skuId,
    }
}

//
//// 获取最新的快照
//func (this *snapshotManagerImpl) GetLatestSnapshot() goods.Snapshot{
//
//}
//
//// 更新快照
//func (this *snapshotManagerImpl) UpdateSnapshot() (int, error){
//
//}
//
//// 生成交易快照
//func (this *snapshotManagerImpl) GenerateSaleSnapshot() (int, error){
//
//}
//
//// 根据KEY获取已销售商品的快照
//func (this *snapshotManagerImpl) GetSaleSnapshotByKey(key string) goods.GoodsSnapshot{
//
//}
//
//// 根据ID获取已销售商品的快照
//func (this *snapshotManagerImpl) GetSaleSnapshot(id int) goods.GoodsSnapshot{
//
//}