/**
 * Copyright 2015 @ z3q.net.
 * name : member_cache
 * author : jarryliu
 * date : 2016-07-23 12:13
 * description :
 * history :
 */
package cache

import (
	"go2o/core/domain/interface/member"
	"go2o/core/service/dps"
)

func GetHighestLevel() *member.Level {
	key := "go2o:cache:max-level"
	sto := GetKVS()
	lv := member.Level{}
	if sto.Get(key, &lv) != nil {
		lv = dps.MemberService.GetHighestLevel()
		if lv.Id > 0 {
			sto.SetExpire(key, lv, DefaultMaxSeconds)
		}
	}
	return &lv
}
