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
	"go2o/src/core/service/dps"
	"io"
	"log"
	"net"
	"strings"
	"time"
	"go2o/src/app/util"
	"github.com/jsix/gof"
	"strconv"
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
			// remove client
			addr :=  conn.RemoteAddr().String()
			uid := clients[addr].UserId
			delete(clients,addr)
			delete(users,uid)
			printf("[ CLIENT][ DISCONN] - IP : %s disconnect!active clients : %d",
				conn.RemoteAddr().String(), len(clients))
			break
		}

		if d, err := rc(conn, line[:len(line)-1]); err != nil {  // remove '\n'
			conn.Write([]byte("error$" + err.Error()))
		} else if d != nil {
			conn.Write(d)
		}
		conn.SetDeadline(time.Now().Add(time.Second * 300))  // dead after 5m
	}
}

var (
	clients map[string]*ClientIdentity = make(map[string]*ClientIdentity)
	users map[int]string = make(map[int]string)
)

type (
	// the identity of client
	ClientIdentity struct {
		Id              int		  // client id
		UserId          int       // user id
		Addr            net.Addr
		ConnectTime     time.Time
		LastConnectTime time.Time
	}
)

func ListenTcp(addr string) {
	go serveLoop()
	listen(addr, func(conn net.Conn, b []byte) ([]byte, error) {
		cmd := string(b)
		id, ok := clients[conn.RemoteAddr().String()]
		if !ok { // auth
			if err := createConnection(conn, cmd); err != nil {
				return nil, err
			}
			return []byte("ok"), nil
		}
		printf("message send by %d , content:%s", id.Id,cmd)
		return handleSocketCmd(id,cmd)
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

func handleSocketCmd(id *ClientIdentity, cmd string) ([]byte, error) {
	i := strings.Index(cmd,":")
	if i == -1{
		return nil,errors.New("unknown command!")
	}
	plan := cmd[i+1:]
	switch cmd[:i] {
	case "MAUTH":
		return cliMAuth(id,plan)
	case "PRINT":
		return cliPrint(id,plan)

	}
	return []byte(cmd), nil
}

func serveLoop(){
	//redis := gof.CurrentApp.(*core.MainApp).Redis()
}

// member auth,command like 'MAUTH:jarrysix#3234234242342342'
func cliMAuth(id *ClientIdentity,param string)([]byte,error){
	arr := strings.Split(param,"#")
	if len(arr) == 2{
		memberId,_ := strconv.Atoi(arr[0])
		b := util.CompareMemberApiToken(gof.CurrentApp.Storage(),
			memberId,arr[1])
		b = true
		if b{ // auth success
			id.UserId = memberId
			users[id.UserId] = id.Addr.String()
			return []byte("ok"),nil
		}
	}
	return nil,errors.New("auth fail")
}

// print text by client sending.
func cliPrint(id *ClientIdentity,params string)([]byte,error){
	return []byte(params),nil
}