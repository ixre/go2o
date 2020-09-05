/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 17:37
 * description :
 * history :
 */
package user

type IUserRepo interface {
	// 保存角色
	SaveRole(*RoleValue) (int32, error)

	// 保存人员
	SavePerson(*PersonValue) (int32, error)

	// 保存凭据
	SaveCredential(*CredentialValue) (int32, error)

	// 获取人员
	GetPersonValue(id int32) *PersonValue

	// 获取配送人员
	GetDeliveryStaffPersons(mchId int32) []*PersonValue
}
