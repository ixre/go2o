package member

import (
	"strings"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/domain/interface/valueobject"
)

var _ member.IDeliverAddress = new(addressImpl)

type addressImpl struct {
	_value      *member.ConsigneeAddress
	_memberRepo member.IMemberRepo
	_valRepo    valueobject.IValueRepo
}

func newDeliverAddress(v *member.ConsigneeAddress, memberRepo member.IMemberRepo,
	valRepo valueobject.IValueRepo) member.IDeliverAddress {
	d := &addressImpl{
		_value:      v,
		_memberRepo: memberRepo,
		_valRepo:    valRepo,
	}
	return d
}

func (p *addressImpl) GetDomainId() int64 {
	return p._value.Id
}

func (p *addressImpl) GetValue() member.ConsigneeAddress {
	return *p._value
}

func (p *addressImpl) SetValue(v *member.ConsigneeAddress) error {
	if p._value.MemberId == v.MemberId {
		if err := p.checkValue(v); err != nil {
			return err
		}
		p._value = v
	}
	return nil
}

// 设置地区中文名
func (p *addressImpl) renewAreaName(v *member.ConsigneeAddress) string {
	return p._valRepo.GetAreaString(v.Province, v.City, v.District)
}

func (p *addressImpl) checkValue(v *member.ConsigneeAddress) error {
	v.DetailAddress = strings.TrimSpace(v.DetailAddress)
	v.ConsigneeName = strings.TrimSpace(v.ConsigneeName)
	v.ConsigneePhone = strings.TrimSpace(v.ConsigneePhone)

	if len([]rune(v.ConsigneeName)) < 2 {
		return member.ErrDeliverContactConsigneeName
	}

	if v.Province <= 0 || v.City <= 0 || v.District <= 0 {
		return member.ErrNotSetArea
	}

	if !phoneRegex.MatchString(v.ConsigneePhone) {
		return member.ErrDeliverContactPhone
	}
	// 判断字符长度
	if len([]rune(v.DetailAddress)) == 0 {
		return member.ErrEmptyDeliverAddress
	}
	return nil
}

// Save 保存收货地址
func (p *addressImpl) Save() error{
	if err := p.checkValue(p._value); err != nil {
		return err
	}
	p._value.Area = p.renewAreaName(p._value)
	id,err := p._memberRepo.SaveDeliverAddress(p._value)
	if p.GetDomainId() == 0{
		p._value.Id = id
	}
	return err
}

// resetDefaultAddress 重置默认收货地址
// func (p *addressImpl) resetDefaultAddress() {
// 	if p._value.IsDefault != 1 {
// 		return
// 	}
// 	list := p._memberRepo.GetDeliverAddress(p._value.MemberId)
// 	for _, v := range list {
// 		if v.IsDefault == 1 {
// 			v.IsDefault = 0
// 			p._memberRepo.SaveDeliverAddress(v)
// 		}
// 	}
// }
