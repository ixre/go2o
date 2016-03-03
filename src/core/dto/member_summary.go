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
<<<<<<< HEAD
	Avatar         string `db:"-"`
=======
>>>>>>> 2616cf765706f843f62d942c38b85a9a18214d6d
	Exp            int    `db:"exp"`
	Level          int    `db:"level"`
	LevelName      string
	Integral       int
	Balance        float32
	PresentBalance float32
	UpdateTime     int64 `db:"update_time"`
}
