/**
 * Copyright 2015 @ z3q.net.
 * name : member_summary
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package dto

// 会员概览信息
type MemberSummary struct {
	Id             int    `db:"id" auto:"yes" pk:"yes"`
	Usr            string `db:"usr"`
	Name           string `db:"name"`
	Avatar		   string  `db:"-"`
	Exp            int    `db:"exp"`
	Level          int    `db:"level"`
	LevelName      string
	Integral       int
	Balance        float32
	PresentBalance float32
	UpdateTime     int64 `db:"update_time"`
}
