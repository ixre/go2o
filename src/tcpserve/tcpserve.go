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
	"bufio"
	"errors"
	"fmt"
	"go2o/src/core/service/dps"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

func printf(format string, args ...interface{}) {
	log.Printf(format+"\n", args...)
}

type (
	TcpReceiveCaller func(conn net.Conn, read []byte) ([]byte, error)
)

func listen(addr string, rc TcpReceiveCaller) {
	serveAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	listen, err := net.ListenTCP("tcp", serveAddr)
	for {
		if conn, err := listen.AcceptTCP(); err == nil {
			printf("[ CLIENT][ CONNECT] - new client connection IP: %s ; active clients : %d",
				conn.RemoteAddr().String(), len(clients)+1)
			go receiveTcpConn(conn, rc)
		}
	}
}

func receiveTcpConn(conn *net.TCPConn, rc TcpReceiveCaller) {
	for {
		buf := bufio.NewReader(conn)
		line, err := buf.ReadBytes('\n')
		if err == io.EOF {
			delete(clients, conn.RemoteAddr().String())
			printf("[ CLIENT][ DISCONN] - IP : %s disconnect!active clients : %d",
				conn.RemoteAddr().String(), len(clients))
			break
		}

		if d, err := rc(conn, line[:len(line)-1]); err != nil {
			conn.Write([]byte("error$" + err.Error()))
		} else if d != nil {
			conn.Write(d)
		}
	}
}

var (
	clients map[string]*ClientIdentity = make(map[string]*ClientIdentity)
)

type (
	// the identity of client
	ClientIdentity struct {
		Id              int
		Addr            net.Addr
		ConnectTime     time.Time
		LastConnectTime time.Time
	}
)

func ListenTcp(addr string) {
	listen(addr, func(conn net.Conn, b []byte) ([]byte, error) {
		id, ok := clients[conn.RemoteAddr().String()]
		// auth

		if !ok {
			if err := createConnection(conn, string(b)); err != nil {
				return nil, err
			}
			return nil, nil
		}
		//
		return []byte(fmt.Sprintf("message send by %d", id.Id)), nil
	})
}

// create partner connection
func createConnection(conn net.Conn, line string) error {
	if strings.HasPrefix(line, "AUTH:") {
		arr := strings.Split(line[5:], "#") // AUTH:API_ID#SECRET
		if len(arr) == 2 {
			partnerId := dps.PartnerService.GetPartnerIdByApiId(arr[0])
			apiInfo := dps.PartnerService.GetApiInfo(partnerId)

			if apiInfo != nil && apiInfo.ApiSecret == arr[1] {
				if apiInfo.Enabled == 0 {
					return errors.New("api has exipres")
				}
				now := time.Now()
				cli := &ClientIdentity{
					Id:              partnerId,
					Addr:            conn.RemoteAddr(),
					ConnectTime:     now,
					LastConnectTime: now,
				}
				clients[conn.RemoteAddr().String()] = cli
				return nil
			}
		}
	}
	return errors.New("conn reject")
}

func handleSocketCmd(cid string, cmd string) ([]byte, error) {
	return []byte(cmd), nil
}
