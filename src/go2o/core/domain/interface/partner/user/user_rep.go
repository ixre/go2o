/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-14 17:37
 * description :
 * history :
 */
package user

type IUserRep interface {
	// 保存角色
	SaveRole(*RoleValue) (int, error)

	// 保存人员
	SavePerson(*RoleValue) (int, error)

	// 保存凭据
	SaveCredential(*CredentialValue) (int, error)
}
