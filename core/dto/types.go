/**
 * Copyright 2015 @ to2.net.
 * name : types
 * author : jarryliu
 * date : 2015-10-29 15:33
 * description :
 * history :
 */
package dto

type (
	TextObject struct {
		Text  string `json:"text"`
		Value int    `json:"value"`
		Title string `json:"title"`
	}

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
