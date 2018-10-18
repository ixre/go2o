/**
 * Copyright 2015 @ z3q.net.
 * name : tcpserve.go
 * author : jarryliu
 * date : 2015-11-23 14:19
 * description :
 * history :
 */
package tcpserve

import (
	"errors"
	"github.com/jsix/gof"
	"github.com/jsix/gof/net/nc"
	"go2o/src/app/util"
	"go2o/src/core/service/dps"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	// 主动关闭没有活动的连接(当前减去最后活动时间)
	disconnectDuration = time.Minute * 10

	// 默认连接存活时间
	defaultReadDeadLine                       = time.Second * 60
	handlers            map[string]nc.CmdFunc = map[string]nc.CmdFunc{
		"PRINT": cliPrint,
		"MGET":  cliMGet,
		"PING":  cliPing,
	}
)

func NewServe(output bool) *nc.SocketServer {
	var s *nc.SocketServer
	r := func(conn net.Conn, b []byte) ([]byte, error) {
		cmd := string(b)
		id, ok := s.GetCli(conn)
		if !ok {
			// not join,auth first!
			if err := connAuth(s, conn, cmd); err != nil {
				return nil, err
			}
			return []byte("ok"), nil
		}
		if strings.HasPrefix(cmd, "MAUTH:") {
			//auth member
			return memberAuth(s, id, cmd[6:])
		}
		return handleCommand(s, id, cmd)
	}

	s = nc.NewSocketServer(r)
	s.ReadDeadLine = defaultReadDeadLine
	if !output {
		s.OutputOff()
	}
	return s
}

// Add socket command handler
func Handle(cmd string, handler nc.CmdFunc) {
	handlers[cmd] = handler
}

// auth connection
func connAuth(s *nc.SocketServer, conn net.Conn, line string) error {
	if strings.HasPrefix(line, "AUTH:") {
		arr := strings.Split(line[5:], "#") // AUTH:API_ID#SECRET#VERSION
		if len(arr) == 3 {
			var af nc.AuthFunc = func() (int64, error) {
				partnerId := dps.PartnerService.GetPartnerIdByApiId(arr[0])
				apiInfo := dps.PartnerService.GetApiInfo(partnerId)
				if apiInfo != nil && apiInfo.ApiSecret == arr[1] {
					if apiInfo.Enabled == 0 {
						return int64(partnerId), errors.New("api has exipres")
					}
				}
				return int64(partnerId), nil
			}

			if err := s.Auth(conn, af); err != nil {
				return err
			}

			s.Printf("[ CLIENT] - Version = %s", arr[2])
			return nil
		}
	}
	return errors.New("conn reject")
}

// member auth,command like 'MAUTH:jarrysix#3234234242342342'
func memberAuth(s *nc.SocketServer, id *nc.Client, param string) ([]byte, error) {
	var err error
	arr := strings.Split(param, "#")
	if len(arr) == 2 {

		f := func() (int64, error) {
			memberId, _ := strconv.Atoi(arr[0])
			authOk := util.CompareMemberApiToken(gof.CurrentApp.Storage(),
				memberId, arr[1])
			if !authOk {
				return int64(memberId), errors.New("auth fail")
			}
			return int64(memberId), nil
		}

		if err = s.UAuth(id.Conn, f); err == nil { //验证成功
			return []byte("ok"), nil
		}
	}
	return nil, err
}

// Handle command of client sending.
func handleCommand(s *nc.SocketServer, ci *nc.Client, cmd string) ([]byte, error) {
	if time.Now().Sub(ci.LatestConnectTime) > disconnectDuration { //主动关闭没有活动的连接
		//s.Print("--disconnect ---",ci.Addr.String())
		ci.Conn.Close()
		return nil, nil
	}
	if !strings.HasPrefix(cmd, "PING") {
		s.Printf("[ CLIENT][ MESSAGE] - send by %d ; %s", ci.Source, cmd)
		ci.LatestConnectTime = time.Now()
	}
	i := strings.Index(cmd, ":")
	if i != -1 {
		plan := cmd[i+1:]
		if v, ok := handlers[cmd[:i]]; ok {
			return v(ci, plan)
		}
	}
	return nil, errors.New("unknown command:" + cmd)
}

// print text by client sending.
func cliPrint(id *nc.Client, params string) ([]byte, error) {
	return []byte(params), nil
}

func cliPing(id *nc.Client, plan string) ([]byte, error) {
	return []byte("PONG"), nil
}
