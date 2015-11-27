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
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jsix/gof"
	"go2o/src/cache"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/dto"
	"go2o/src/core/service/dps"
	"go2o/src/core/variable"
	"net"
	"strconv"
	"strings"
)

// get summary of member,if dbGet will get summary from database.
func GetMemberSummary(memberId int, updateTime int) *dto.MemberSummary {
	sto := gof.CurrentApp.Storage()
	var kvMut int
	mutKey := fmt.Sprintf("%s%d", variable.KvMemberUpdateTime, memberId)
	kvMut, _ = sto.GetInt(mutKey)
	//get from redis
	var v *dto.MemberSummary = new(dto.MemberSummary)
	var key = fmt.Sprintf("cac:mm:summary:%d", memberId)
	if kvMut != 0 && kvMut == updateTime {
		if cache.GetKVS().Get(key, v) == nil {
			return v
		}
	}
	v = dps.MemberService.GetMemberSummary(memberId)
	sto.SetExpire(key, v, 3600*360) // cache 15 hours
	sto.SetExpire(mutKey, v.UpdateTime, 3600*400)
	return v
}

func getMemberAccount(memberId int, updateTime int) *member.AccountValue {
	sto := gof.CurrentApp.Storage()
	var kvAut int
	autKey := fmt.Sprintf("%s%d", variable.KvAccountUpdateTime, memberId)
	kvAut, _ = sto.GetInt(autKey)
	//get from redis
	var v *member.AccountValue = new(member.AccountValue)
	var key = fmt.Sprintf("cac:mm:acc:%d", memberId)
	if kvAut != 0 && kvAut == updateTime {
		if cache.GetKVS().Get(key, v) == nil {
			return v
		}
	}
	v = dps.MemberService.GetAccount(memberId)
	sto.SetExpire(key, v, 3600*360) // cache 15 hours
	sto.SetExpire(autKey, v.UpdateTime, 3600*400)

	return v

}

// push member summary to tcp client
func pushMemberSummary(connList []net.Conn, memberId int) {
	printf(false, "[ TCP][ NOTIFY] - notify member update - %d", memberId)
	sm := GetMemberSummary(memberId, 0)
	if d, err := json.Marshal(sm); err == nil {
		d = append([]byte("MSUM:"), d...)
		for _, conn := range connList {
			go conn.Write(append(d, '\n'))
		}
	}
}

// push member summary to tcp client
func pushMemberAccount(connList []net.Conn, memberId int) {
	printf(false, "[ TCP][ NOTIFY] - notify account update - %d", memberId)
	sm := getMemberAccount(memberId, 0)
	if d, err := json.Marshal(sm); err == nil {
		d = append([]byte("MACC:"), d...)
		for _, conn := range connList {
			go conn.Write(append(d, '\n'))
		}
	}
}

// get profile of member
func cliMGet(ci *ClientIdentity, plan string) ([]byte, error) {
	var obj interface{} = nil
	var d []byte = []byte{}

	i := strings.Index(plan, ":")
	ut, _ := strconv.Atoi(plan[i+1:])

	switch plan[0:i] {
	case "SUMMARY":
		obj = GetMemberSummary(ci.UserId, ut)
		d = []byte("MSUM:")
	case "ACCOUNT":
		obj = getMemberAccount(ci.UserId, ut)
		d = []byte("MACC:")
	}
	if obj != nil {
		d1, err := json.Marshal(obj)
		return append(d, d1...), err

	}
	return nil, errors.New("unknown type:" + plan)
}

func mmSummaryNotify(conn redis.Conn) error {
	mid, err := redis.Int(conn.Do("LPOP", variable.KvMemberUpdateTcpNotifyQueue))
	if err == nil {
		arr := strings.Split(users[mid], "$")
		var connList []net.Conn = make([]net.Conn, 0)
		for _, v := range arr {
			if ide, ok := clients[v]; ok && ide.Conn != nil {
				connList = append(connList, ide.Conn)
			}
		}
		if len(connList) > 0 {
			pushMemberSummary(connList, mid)
		}
	}
	return err
}

func mmAccountNotify(conn redis.Conn) error {
	mid, err := redis.Int(conn.Do("LPOP", variable.KvAccountUpdateTcpNotifyQueue))
	if err == nil {
		arr := strings.Split(users[mid], "$")
		var connList []net.Conn = make([]net.Conn, 0)
		for _, v := range arr {
			if ide, ok := clients[v]; ok && ide.Conn != nil {
				connList = append(connList, ide.Conn)
			}
		}
		if len(connList) > 0 {
			pushMemberAccount(connList, mid)
		}
	}
	return err
}
