/**
 * Copyright 2015 @ to2.net.
 * name : tcpserve.go
 * author : jarryliu
 * date : 2015-11-23 14:19
 * description :
 * history :
 */
package tcpserve

import (
	"errors"
	"github.com/ixre/gof/net/nc"
	"github.com/ixre/gof/util"
	"go2o/core/service/rsi"
	"go2o/core/service/thrift"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// 主动关闭没有活动的连接(当前减去最后活动时间)
	disconnectDuration = time.Minute * 10
	// 默认连接存活时间
	defaultReadDeadLine = time.Second * 60
	// 操作
	handlers map[string]nc.CmdFunc = map[string]nc.CmdFunc{
		"PRINT": cliPrint,
		"MGET":  cliMGet,
		"PING":  cliPing,
	}
	mux sync.Mutex
)

func NewServe(output bool) *nc.SocketServer {
	var s *nc.SocketServer
	r := func(conn net.Conn, b []byte) ([]byte, error) {
		cmd := string(b)
		id, ok := s.GetCli(conn)
		// if not join,auth first!
		if !ok {
			if err := connAuth(s, conn, cmd); err != nil {
				return nil, err
			}
			return []byte("ok"), nil
		}
		// member auth
		if strings.HasPrefix(cmd, "MAUTH:") {
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
	mux.Lock()
	defer mux.Unlock()
	handlers[cmd] = handler
}

// auth connection
func connAuth(s *nc.SocketServer, conn net.Conn, line string) error {
	if strings.HasPrefix(line, "AUTH:") {
		arr := strings.Split(line[5:], "#") // AUTH:API_ID#SECRET#VERSION
		if len(arr) == 3 {
			var af nc.AuthFunc = func() (int64, error) {
				mchId := rsi.MerchantService.GetMerchantIdByApiId(arr[0])
				apiInfo := rsi.MerchantService.GetApiInfo(mchId)
				if apiInfo != nil && apiInfo.ApiSecret == arr[1] {
					if apiInfo.Enabled == 0 {
						return int64(mchId), errors.New("api has exipres")
					}
				}
				return int64(mchId), nil
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

// member auth,command like 'MAUTH:1#3234234242342342'
func memberAuth(s *nc.SocketServer, id *nc.Client, param string) ([]byte, error) {
	var err error
	arr := strings.Split(param, "#")
	if len(arr) == 2 {
		f := func() (int64, error) {
			memberId, _ := util.I64Err(strconv.Atoi(arr[0]))
			trans, cli, err := thrift.MemberServeClient()
			if err == nil {
				defer trans.Close()
				if b, _ := cli.CheckToken(thrift.Context, memberId, arr[1]); b {
					return memberId, nil
				}
				return memberId, errors.New("auth fail")
			}
			return memberId, errors.New("connect refused")
		}

		if err = s.UAuth(id.Conn, f); err == nil {
			//验证成功
			return []byte("ok"), nil
		}
	}
	return nil, err
}

// Handle command of client sending.
func handleCommand(s *nc.SocketServer, ci *nc.Client, cmd string) ([]byte, error) {
	if time.Now().Sub(ci.LatestConnectTime) > disconnectDuration {
		//主动关闭没有活动的连接
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
