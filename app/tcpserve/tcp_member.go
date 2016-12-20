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
	"github.com/jsix/gof/net/nc"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift/idl/gen-go/define"
	"strconv"
	"strings"
)

// get summary of member,if dbGet will get summary from database.
func GetMemberSummary(memberId int, updateTime int) *define.MemberSummary {
	v, _ := rsi.MemberService.Summary(int32(memberId))
	if v != nil {
		return v
	}
	return nil
}

func getMemberAccount(memberId int, updateTime int) *define.Account {
	v, _ := rsi.MemberService.GetAccount(int32(memberId))
	return v
}

// get profile of member
func cliMGet(ci *nc.Client, plan string) ([]byte, error) {
	var obj interface{} = nil
	var d []byte = []byte{}

	i := strings.Index(plan, ":")
	ut, _ := strconv.Atoi(plan[i+1:])

	switch plan[0:i] {
	case "SUMMARY":
		obj = GetMemberSummary(ci.User, ut)
		d = []byte("MSUM:")
	case "ACCOUNT":
		obj = getMemberAccount(ci.User, ut)
		d = []byte("MACC:")
	}
	if obj != nil {
		d1, err := json.Marshal(obj)
		return append(d, d1...), err

	}
	return nil, errors.New("unknown type:" + plan)
}
