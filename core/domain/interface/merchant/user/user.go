/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 16:18
 * description :
 * history :
 */
package user

type IUser interface {
	// 获取人员信息
	GetPerson() IPerson

	// 获取凭据
	GetCredential(sign string) *CredentialValue

	// 保存凭据
	SaveCredential(*CredentialValue) error

	// 保存人员信息
	//Save() (int32, error)
}
