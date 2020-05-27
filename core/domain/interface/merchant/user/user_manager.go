/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-02-14 16:55
 * description :
 * history :
 */
package user

type IUserManager interface {
	// 获取单个用户
	GetUser(id int32) IUser

	// 获取所有配送员
	GetDeliveryStaff() []IDeliveryStaff
}
