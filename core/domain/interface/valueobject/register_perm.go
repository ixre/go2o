package valueobject

import "go2o/core/domain/interface/registry"

/**
 * Copyright 2009-2019 @ to2.net
 * name : register_perm.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-07-01 10:25
 * description :
 * history :
 */


// 注册权限设置
type RegisterPerm1 struct {
	// 注册模式,等于member.RegisterMode
	RegisterMode int `key:"MemberRegisterAllowAnonymous"`
	// 是否允许匿名注册
	AnonymousRegistered bool `key:"MemberRegisterAllowAnonymous"`
	// 手机号码作为用户名
	PhoneAsUser bool `key:"MemberRegisterPhoneAsUser"`
	// 是否需要填写手机
	NeedPhone bool `key:"MemberRegisterNeedPhone"`
	// 必须绑定手机
	MustBindPhone bool `key:"MemberRegisterMustBindPhone"`
	// 是否需要填写即时通讯
	NeedIm bool `key:"MemberRegisterNeedIm"`
	// 注册提示
	Notice string `key:"MemberRegisterNotice"`
	// 用户条款内容
	Licence string `key:"MemberRegisterReturnUrl"`
	// 注册回调页
	CallBackUrl string `key:"MemberRegisterPresentIntegral"`
	keys []string
}

func (r RegisterPerm1) Keys(){
	if r.keys == nil {
		r.keys = []string{
			//registry.MemberReg
		}
	}
}

func LoadRegisterPerm(repo registry.IRegistryRepo){

}