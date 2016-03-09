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
	"github.com/garyburd/redigo/redis"
	"github.com/jsix/gof"
	"github.com/jsix/gof/net/nc"
	"go2o/src/app/util"
	"go2o/src/core"
	"go2o/src/core/service/dps"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	s                  *nc.SocketServer
	disconnectDuration                       = time.Second * 60
	readDeadLine                             = time.Second * 40
	handlers           map[string]nc.CmdFunc = map[string]nc.CmdFunc{
		"PRINT": cliPrint,
		"MGET":  cliMGet,
		"PING":  cliPing,
	}
)

func ListenTcp(addr string) {
	s = nc.NewSocketServer()
	s.ReadDeadLine = readDeadLine
	serveLoop(s) // server loop,send some message to client
	//s.OutputOff()
	s.Listen(addr, func(conn net.Conn, b []byte) ([]byte, error) {
		cmd := string(b)
		id, ok := s.GetCli(conn)
		if !ok { // not join,auth first!
			if err := connAuth(s, conn, cmd); err != nil {
				return nil, err
			}
			return []byte("ok"), nil
		}
		if strings.HasPrefix(cmd, "MAUTH:") { //auth member
			return memberAuth(s, id, cmd[6:])
		}
		return handleCommand(id, cmd)
	})
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
			var af nc.AuthFunc = func() (int, error) {
				partnerId := dps.PartnerService.GetPartnerIdByApiId(arr[0])
				apiInfo := dps.PartnerService.GetApiInfo(partnerId)
				if apiInfo != nil && apiInfo.ApiSecret == arr[1] {
					if apiInfo.Enabled == 0 {
						return partnerId, errors.New("api has exipres")
					}
				}
				return partnerId, nil
			}

			if err := s.Auth(conn, af); err != nil {
				return err
			}

			s.Print("[ CLIENT] - Version = %s", arr[2])
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

		f := func() (int, error) {
			memberId, _ := strconv.Atoi(arr[0])
			authOk := util.CompareMemberApiToken(gof.CurrentApp.Storage(),
				memberId, arr[1])
			if !authOk {
				return memberId, errors.New("auth fail")
			}
			return memberId, nil
		}

		if err = s.UAuth(id.Conn, f); err == nil { //验证成功
			return []byte("ok"), nil
		}
	}
	return nil, err
}

// Handle command of client sending.
func handleCommand(ci *nc.Client, cmd string) ([]byte, error) {
	if time.Now().Sub(ci.LatestConnectTime) > disconnectDuration { //主动关闭没有活动的连接
		//s.Print("--disconnect ---",ci.Addr.String())
		ci.Conn.Close()
		return nil, nil
	}
	if !strings.HasPrefix(cmd, "PING") {
		s.Print("[ CLIENT][ MESSAGE] - send by %d ; %s", ci.Source, cmd)
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

func serveLoop(s *nc.SocketServer) {
	conn := core.GetRedisConn()
	go notifyMup(s, conn)
}

func notifyMup(s *nc.SocketServer, conn redis.Conn) {
	time.Sleep(time.Second * 10) // 等待监听服务
	for {
		err := mmSummaryNotify(s, conn)
		err1 := mmAccountNotify(s, conn)
		if err != nil || err1 != nil {
			time.Sleep(time.Second * 1) //阻塞,避免轮询占用CPU
		}
	}
}

// print text by client sending.
func cliPrint(id *nc.Client, params string) ([]byte, error) {
	return []byte(params), nil
}

func cliPing(id *nc.Client, plan string) ([]byte, error) {
	return []byte("PONG"), nil
}
