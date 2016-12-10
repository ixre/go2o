package promodel

import (
    "go2o/core/domain/interface/pro_model"
    "database/sql"
)

var _ promodel.IAttrService = new(attrServiceImpl)

type attrServiceImpl struct {
    rep promodel.IProModelRepo
}

func NewAttrService(rep promodel.IProModelRepo) *attrServiceImpl {
    return &attrServiceImpl{
        rep: rep,
    }
}

// 获取属性
func (a *attrServiceImpl)GetAttr(attrId int32) *promodel.Attr {
    return a.rep.GetAttr(attrId)
}

// 保存属性
func (a *attrServiceImpl)SaveAttr(v *promodel.Attr) (int32, error) {
    id, err := a.rep.SaveAttr(v)
    return int32(id), err
}
// 保存属性项
func (a *attrServiceImpl)SaveItem(v *promodel.AttrItem) (int32, error) {
    id, err := a.rep.SaveAttrItem(v)
    return int32(id), err
}
// 删除属性
func (a *attrServiceImpl)DeleteAttr(attrId int32) error {
    _, err := a.rep.BatchDeleteAttrItem("attr_id=?", attrId)
    if err == nil || err == sql.ErrNoRows {
        err = a.rep.DeleteAttr(attrId)
    }
    return err
}
// 删除属性项
func (a *attrServiceImpl)DeleteItem(itemId int32) error {
    return a.rep.DeleteAttrItem(itemId)
}
// 获取属性的属性项
func (a *attrServiceImpl)GetItems(attrId int32) []*promodel.AttrItem {
    return a.rep.SelectAttrItem("attr_id=?", attrId)
}
// 获取产品模型的属性
func (a *attrServiceImpl)GetModelAttrs(proModel int32) []*promodel.Attr {
    return a.rep.SelectAttr("pro_model=?", proModel)
}
// 获取产品的属性
//func (a *attrServiceImpl)GetGoodsAttrs(proId int32) []*ProAttr



