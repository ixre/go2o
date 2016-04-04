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
    "github.com/jsix/gof/net/nc"
    "github.com/garyburd/redigo/redis"
    "go2o/src/core"
    "strconv"
    "go2o/src/core/variable"
    "net"
    "encoding/json"
)

func AccountNotifyJob(s *nc.SocketServer){
    conn := core.GetRedisConn()
    defer conn.Close()
    for {
        values, err := redis.Values(conn.Do("BLPOP",
            variable.KvAccountUpdateTcpNotifyQueue))
        if err == nil {
            id,err := strconv.Atoi(string(values[1].([]byte)))
            if err == nil{
                connList := s.GetConnections(id)
                if len(connList) > 0 {
                    go pushMemberAccount(s,connList,id)
                }
            }
        }
    }
}


// push member summary to tcp client
func pushMemberAccount(s *nc.SocketServer,connList []net.Conn, memberId int) {
    s.Print("[ TCP][ NOTIFY] - notify account update - %d", memberId)
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


func MemberSummaryNotifyJob(s *nc.SocketServer) error {
    conn := core.GetRedisConn()
    defer conn.Close()
    for {
        values, err := redis.Values(conn.Do("BLPOP",
            variable.KvMemberUpdateTcpNotifyQueue))
        if err == nil {
            id,err := strconv.Atoi(string(values[1].([]byte)))
            if err == nil{
                connList := s.GetConnections(id)
                if len(connList) > 0 {
                    go pushMemberSummary(s,connList,id)
                }
            }
        }
    }
}


// push member summary to tcp client
func pushMemberSummary(s *nc.SocketServer,connList []net.Conn, memberId int) {
    s.Print("[ TCP][ NOTIFY] - notify member update - %d", memberId)
    sm := GetMemberSummary(memberId, 0)
    if d, err := json.Marshal(sm); err == nil {
        d = append([]byte("MSUM:"), d...)
        for _, conn := range connList {
            conn.Write(append(d, '\n'))
        }
    }
}


