/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 17:37
 * description :
 * history :
 */
package user

type IUserRep interface {
	// 保存角色
	SaveRole(*RoleValue) (int, error)

	// 保存人员
	SavePerson(*PersonValue) (int, error)

	// 保存凭据
	SaveCredential(*CredentialValue) (int, error)

	// 获取人员
	GetPersonValue(int) *PersonValue

	// 获取配送人员
	GetDeliveryStaffPersons(partnerId int) []*PersonValue
}
