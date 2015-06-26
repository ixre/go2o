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
	Rep ad.IAdvertisementRep
	Value *ad.ValueAdvertisement
}

// 获取领域对象编号
func (this *Advertisement) GetDomainId() int{
	if this.Value != nil {
		return this.Value.Id
	}
	return 0
}

// 是否为系统内置的广告
func (this *Advertisement) System()bool{
	return this.Value.IsInternal == 1
}

// 广告类型
func (this *Advertisement) Type()int{
	return this.Value.Type
}

// 广告名称
func (this *Advertisement) Name()string{
	return this.Value.Name
}

// 设置值
func (this *Advertisement) SetValue(v *ad.ValueAdvertisement)error {
	// 如果为系统内置广告，不能修改名称
	if this.System() {
		v.Name = this.Value.Name
		v.Enabled = 1
		v.IsInternal = 1
	}
	this.Value = v

	return nil
}

// 获取值
func (this *Advertisement)GetValue()*ad.ValueAdvertisement{
	return this.Value
}

// 保存广告
func (this *Advertisement) Save()(int,error){
	id := this.Rep.GetIdByName(this.Value.PartnerId,this.Value.Name)
	if id > 0 && id != this.GetDomainId(){
		return this.GetDomainId(),ad.ErrNameExists
	}
	return this.Rep.SaveAdvertisementValue(this.Value)
}
