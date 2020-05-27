/**
 * Copyright 2015 @ to2.net.
 * name : model
 * author : jarryliu
 * date : 2015-10-31 00:35
 * description :
 * history :
 */
package restapi

type (
	AsyncResult struct {
		MemberId       int64 // 会员编号
		MemberUpdated  bool  //会员已经更新
		AccountUpdated bool  //会员账户已经更新
	}
)
