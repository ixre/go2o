/**
 * Copyright 2015 @ z3q.net.
 * name : tcp_member
 * author : jarryliu
 * date : 2015-11-24 11:49
 * description :
 * history :
 */
package tcpserve

import (
	"encoding/json"
	"fmt"
	"go2o/src/cache"
	"go2o/src/core/dto"
	"go2o/src/core/service/dps"
	"net"
)

// get summary of member,if dbGet will get summary from database.
func GetMemberSummary(memberId int, dbGet bool) *dto.MemberSummary {
	var v *dto.MemberSummary = new(dto.MemberSummary)
	var key = fmt.Sprintf("cache:member:summary:%d", memberId)
	if dbGet || cache.GetKVS().Get(key, &v) != nil {
		v = dps.MemberService.GetMemberSummary(memberId)
		cache.GetKVS().SetExpire(key, v, 3600*48) // cache 48 hours
	}
	return v
}

// push member summary to tcp client
func pushMemberSummary(connList []net.Conn, memberId int) {
	printf(false, "[ TCP][ NOTIFY] - notify member update - %d", memberId)
	sm := GetMemberSummary(memberId, true)
	if d, err := json.Marshal(sm); err == nil {
		d = append([]byte("MUP:"), d...)
		for _,conn := range connList{
			conn.Write(d)
			conn.Write([]byte("\n"))
		}
	}
}
