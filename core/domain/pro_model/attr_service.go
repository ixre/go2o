package promodel

import (
	"bytes"
	"database/sql"
	"fmt"
	"go2o/core/domain/interface/pro_model"
	"strconv"
)

var _ promodel.IAttrService = new(attrServiceImpl)

type attrServiceImpl struct {
	repo    promodel.IProModelRepo
	builder *attrHtmlBuilder
}

func NewAttrService(repo promodel.IProModelRepo) *attrServiceImpl {
	return &attrServiceImpl{
		repo:    repo,
		builder: &attrHtmlBuilder{},
	}
}

// 获取属性
func (a *attrServiceImpl) GetAttr(attrId int32) *promodel.Attr {
	return a.repo.GetAttr(attrId)
}

// 保存属性
func (a *attrServiceImpl) SaveAttr(v *promodel.Attr) (id int32, err error) {
	var i int
	// 如不存在，则新增
	if v.ID <= 0 {
		i, err = a.repo.SaveAttr(v)
		v.ID = int32(i)
		if v == nil {
			return v.ID, err
		}
	}
	// 保存项
	if v.Items != nil {
		v.ItemValues = ""
		for i, iv := range v.Items {
			iv.ProModel = v.ProModel
			iv.AttrId = v.ID
			if i > 0 {
				v.ItemValues += ","
			}
			v.ItemValues += iv.Value
		}
		err = a.saveAttrItems(v.ID, v.Items)
	}
	// 再次保存
	if err == nil {
		_, err = a.repo.SaveAttr(v)
	}
	return v.ID, err
}

// 保存属性项
func (a *attrServiceImpl) saveAttrItems(attrId int32, items []*promodel.AttrItem) (err error) {
	var i int
	pk := attrId
	// 获取存在的项
	old := a.repo.SelectAttrItem("attr_id = $1", pk)
	// 分析当前项目并加入到MAP中
	delList := []int32{}
	currMap := make(map[int32]*promodel.AttrItem, len(items))
	for _, v := range items {
		currMap[v.ID] = v
	}
	// 筛选出要删除的项
	for _, v := range old {
		if currMap[v.ID] == nil {
			delList = append(delList, v.ID)
		}
	}
	// 删除项
	for _, v := range delList {
		a.repo.DeleteAttrItem(v)
	}
	// 保存项
	for _, v := range items {
		if v.AttrId == 0 {
			v.AttrId = pk
		}
		if v.AttrId == pk {
			if i, err = a.repo.SaveAttrItem(v); err == nil {
				v.ID = int32(i)
			}
		}
	}
	return err
}

// 保存属性项
func (a *attrServiceImpl) SaveItem(v *promodel.AttrItem) (int32, error) {
	id, err := a.repo.SaveAttrItem(v)
	return int32(id), err
}

// 删除属性
func (a *attrServiceImpl) DeleteAttr(attrId int32) error {
	_, err := a.repo.BatchDeleteAttrItem("attr_id= $1", attrId)
	if err == nil || err == sql.ErrNoRows {
		err = a.repo.DeleteAttr(attrId)
	}
	return err
}

// 删除属性项
func (a *attrServiceImpl) DeleteItem(itemId int32) error {
	return a.repo.DeleteAttrItem(itemId)
}

// 获取属性的属性项
func (a *attrServiceImpl) GetItems(attrId int32) []*promodel.AttrItem {
	return a.repo.SelectAttrItem("attr_id= $1", attrId)
}

// 获取产品模型的属性
func (a *attrServiceImpl) GetModelAttrs(proModel int32) []*promodel.Attr {
	arr := a.repo.SelectAttr("pro_model= $1", proModel)
	for _, v := range arr {
		v.Items = a.GetItems(v.ID)
	}
	return arr
}

// 获取属性的HTML表示
func (a *attrServiceImpl) AttrsHtml(arr []*promodel.Attr) string {
	buf := bytes.NewBuffer(nil)
	if len(arr) == 0 {
		buf.WriteString("<div class=\"no-attr\">该分类下未包含属性</div>")
	} else {
		for _, v := range arr {
			a.builder.Append(buf, v)
		}
	}
	return buf.String()
}

type attrHtmlBuilder struct {
}

func (a *attrHtmlBuilder) Append(buf *bytes.Buffer, attr *promodel.Attr) {
	buf.WriteString("<div class=\"attr-item attr\" attr-id=\"")
	buf.WriteString(strconv.Itoa(int(attr.ID)))
	buf.WriteString("\">")
	a.buildLabel(buf, attr.Name)
	buf.WriteString("<div class=\"attr-list attr\">")
	if attr.MultiChk == 1 {
		a.buildCheckBox(buf, attr)
	} else {
		a.buildDropDown(buf, attr)
	}
	buf.WriteString("</div></div>\n")
}
func (a *attrHtmlBuilder) buildDropDown(buf *bytes.Buffer,
	attr *promodel.Attr) {
	buf.WriteString("<select class=\"attr-val\" _field=\"_AttrData\">")
	for _, v := range attr.Items {
		buf.WriteString("<option value=\"")
		buf.WriteString(strconv.Itoa(int(v.ID)))
		buf.WriteString("\">")
		buf.WriteString(v.Value)
		buf.WriteString("</option>")
	}
	buf.WriteString("</select>")
}
func (a *attrHtmlBuilder) buildCheckBox(buf *bytes.Buffer,
	attr *promodel.Attr) {
	for i, v := range attr.Items {
		str := fmt.Sprintf("%d-%d", v.AttrId, i)
		buf.WriteString("<input type=\"checkbox\" class=\"attr-val\" _field=\"_AttrData[")
		buf.WriteString(str)
		buf.WriteString("]\" value=\"")
		buf.WriteString(strconv.Itoa(int(v.ID)))
		buf.WriteString("\" id=\"ck_attr_")
		buf.WriteString(str)
		buf.WriteString("\"/><label class=\"ck_label\" for=\"ck_attr_")
		buf.WriteString(str)
		buf.WriteString("\">")
		buf.WriteString(v.Value)
		buf.WriteString("</label>")
	}
}
func (a *attrHtmlBuilder) buildLabel(buf *bytes.Buffer, label string) {
	buf.WriteString("<div class=\"attr-label\">")
	buf.WriteString(label)
	buf.WriteString(": </div>")
}
