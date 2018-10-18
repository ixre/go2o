/**
 * Copyright 2015 @ z3q.net.
 * name : tcp_observer
 * author : jarryliu
 * date : 2016-04-04 11:17
 * description :
 * history :
 */
package tcpserve

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/jsix/gof/net/nc"
	"go2o/src/core"
	"go2o/src/core/variable"
	"net"
	"strconv"
)

//note: 业务逻辑上可能会出现通知多次的情况
func AccountNotifyJob(s *nc.SocketServer) {
	conn := core.GetRedisConn()
	defer conn.Close()
	for {
		values, err := redis.Values(conn.Do("BLPOP",
			variable.KvAccountUpdateTcpNotifyQueue, 0))
		if err == nil {
			id, err := strconv.Atoi(string(values[1].([]byte)))
			if err == nil {
				connList := s.GetConnections(int64(id))
				if len(connList) > 0 {
					go pushMemberAccount(s, connList, id)
				}
			}
		}
	}
}

// push member summary to tcp client
func pushMemberAccount(s *nc.SocketServer, connList []net.Conn, memberId int) {
	s.Printf("[ TCP][ NOTIFY] - notify account update - %d", memberId)
	sm := getMemberAccount(memberId, 0)
	if sm != nil {
		if d, err := json.Marshal(sm); err == nil {
			d = append([]byte("MACC:"), d...)
			for _, conn := range connList {
				conn.Write(append(d, '\n'))
			}
		}
	}
}

func MemberSummaryNotifyJob(s *nc.SocketServer) {
	conn := core.GetRedisConn()
	defer conn.Close()
	for {
		values, err := redis.Values(conn.Do("BLPOP",
			variable.KvMemberUpdateTcpNotifyQueue, 0))
		if err == nil {
			id, err := strconv.Atoi(string(values[1].([]byte)))
			if err == nil {
				connList := s.GetConnections(int64(id))
				if len(connList) > 0 {
					go pushMemberSummary(s, connList, id)
				}
			}
		}
	}
}

// push member summary to tcp client
func pushMemberSummary(s *nc.SocketServer, connList []net.Conn, memberId int) {
	s.Printf("[ TCP][ NOTIFY] - notify member update - %d", memberId)
	sm := GetMemberSummary(memberId, 0)
	if d, err := json.Marshal(sm); err == nil {
		d = append([]byte("MSUM:"), d...)
		for _, conn := range connList {
			conn.Write(append(d, '\n'))
		}
	}
}
