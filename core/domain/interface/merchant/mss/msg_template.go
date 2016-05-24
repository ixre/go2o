/**
 * Copyright 2015 @ z3q.net.
 * name : msg_template
 * author : jarryliu
 * date : 2015-07-26 21:57
 * description :
 * history :
 */
package mss

type IMsgTemplate interface {
	// 应用数据
	ApplyData(MsgData)
	// 加入到发送对列
	JoinQueen(to []string) error
}
