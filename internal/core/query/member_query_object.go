/**
 * Copyright 2015 @ 56x.net.
 * name : member
 * author : jarryliu
 * date : 2015-10-29 15:06
 * description :
 * history :
 */
package query

type (

	// 会员排名信息
	RankMember struct {
		Id       int64
		Name     string
		Usr      string
		RankNum  int
		InviNum  int // 邀请数量
		TotalNum int // 总数
		RegTime  int
	}
)
