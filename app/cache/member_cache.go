/**
 * Copyright 2015 @ to2.net.
 * name : member_cache
 * author : jarryliu
 * date : 2016-07-23 12:13
 * description :
 * history :
 */
package cache

import (
	"encoding/json"
	"go2o/core/domain/interface/member"
	"go2o/core/service/rsi"
	"strconv"
)

// 获取最高等级
func GetHighestLevel() *member.Level {
	key := "go2o:repo:level:glob:max"
	sto := GetKVS()
	lv := member.Level{}
	if sto.Get(key, &lv) != nil {
		lv = rsi.MemberService.GetHighestLevel()
		if lv.ID > 0 {
			sto.SetExpire(key, lv, DefaultMaxSeconds)
		}
	}
	return &lv
}

// 获取等级JSON
func GetLevelMapJson() string {
	key := "go2o:repo:level:mp-json"
	sto := GetKVS()
	str, err := sto.GetString(key)
	if err != nil {
		list := rsi.MemberService.GetMemberLevels()
		mp := make(map[string]string, 0)
		for _, v := range list {
			if v.Enabled == 1 {
				mp[strconv.Itoa(int(v.ID))] = v.Name
			}
		}
		data, _ := json.Marshal(mp)
		str = string(data)
		sto.SetExpire(key, str, DefaultMaxSeconds)
	}
	return str
}
