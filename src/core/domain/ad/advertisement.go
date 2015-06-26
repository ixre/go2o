/**
 * Copyright 2015 @ S1N1 Team.
 * name : advertisement
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package ad
import "go2o/src/core/domain/interface/ad"

var _ ad.IAdvertisement = new(Advertisement)

type Advertisement struct {
	_rep ad.IAdvertisementRep
	_value *ad.ValueAdvertisement
}

// 获取领域对象编号
func (this *Advertisement) GetDomainId() int{
	if this._value != nil {
		return this._value.Id
	}
	return 0
}

// 是否为系统内置的广告
func (this *Advertisement) System()bool{
	return this._value.IsInternal
}

// 广告类型
func (this *Advertisement) Type()int{
	return this._value.Type
}

// 广告名称
func (this *Advertisement) Name()string{
	return this._value.Name
}

// 设置值
func (this *Advertisement) SetValue(v *ad.ValueAdvertisement)error{
	if v != nil {
		this._value = v
	}
	return nil
}

// 获取值
func (this *Advertisement)GetValue()*ad.ValueAdvertisement{
	return this._value
}

// 保存广告
func (this *Advertisement) Save()(int,error){
	return this._rep.SaveAdvertisementValue(this._value)
}
