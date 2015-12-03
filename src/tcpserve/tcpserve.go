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
	"github.com/garyburd/redigo/redis"
	"github.com/jsix/gof"
	"go2o/src/app/util"
	"go2o/src/core"
	"go2o/src/core/service/dps"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type (
	TcpReceiveCaller func(conn net.Conn, read []byte) ([]byte, error)
	SocketCmdHandler func(ci *ClientIdentity, plan string) ([]byte, error)
	// the identity of client
	ClientIdentity struct {
		Id              int // client id
		UserId          int // user id
		Addr            net.Addr
		Conn            net.Conn
		ConnectTime     time.Time
		LastConnectTime time.Time
	}
)

var (
	DebugOn      bool                        = false
	ReadDeadLine time.Duration               = time.Second * 300
	clients      map[string]*ClientIdentity  = make(map[string]*ClientIdentity)
	users        map[int]string              = make(map[int]string)
	handlers     map[string]SocketCmdHandler = map[string]SocketCmdHandler{
		"MAUTH": cliMAuth,
		"PRINT": cliPrint,
		"MGET":  cliMGet,
		"PING":  cliPing,
	}
)

func printf(force bool, format string, args ...interface{}) {
	if DebugOn {
		log.Printf(format+"\n", args...)
	}
}

func listen(addr string, rc TcpReceiveCaller) {
	serveAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err)
	}
	listen, err := net.ListenTCP("tcp", serveAddr)
	for {
		if conn, err := listen.AcceptTCP(); err == nil {
			printf(true, "[ CLIENT][ CONNECT] - new client connection IP: %s ; active clients : %d",
				conn.RemoteAddr().String(), len(clients)+1)
			go receiveTcpConn(conn, rc)
		}
	}
}

func receiveTcpConn(conn *net.TCPConn, rc TcpReceiveCaller) {
	for {
		buf := bufio.NewReader(conn)
		line, err := buf.ReadBytes('\n')
		if err != nil {
			// remove client
			addr := conn.RemoteAddr().String()
			if v, ok := clients[addr]; ok {
				uid := v.UserId
				delete(clients, addr)
				addr2 := users[uid]
				if strings.Index(addr2, "$") == -1 {
					delete(users, uid)
				} else {
					users[uid] = strings.Replace(strings.Replace(addr2, addr, "", 1), "$$", "$", -1)
				}
			}
			printf(true, "[ CLIENT][ DISCONN] - IP : %s disconnect!active clients : %d",
				conn.RemoteAddr().String(), len(clients))
			break

		}
		if d, err := rc(conn, line[:len(line)-1]); err != nil { // remove '\n'
			conn.Write([]byte("error$" + err.Error()))
		} else if d != nil {
			conn.Write(d)
		}
		conn.Write([]byte("\n"))
		conn.SetReadDeadline(time.Now().Add(ReadDeadLine)) // discount after 5m
	}
}

func ListenTcp(addr string) {
	serveLoop() // server loop,send some message to client
	listen(addr, func(conn net.Conn, b []byte) ([]byte, error) {
		cmd := string(b)
		id, ok := clients[conn.RemoteAddr().String()]
		if !ok { // auth
			if err := createConnection(conn, cmd); err != nil {
				return nil, err
			}
			return []byte("ok"), nil
		}
		if !strings.HasPrefix(cmd, "PING") {
			printf(false, "[ CLIENT][ MESSAGE] - send by %d ; %s", id.Id, cmd)
		}
		return handleSocketCmd(id, cmd)
	})
}

// register socket command handler
func AddHandler(cmd string, handler SocketCmdHandler) {
	handlers[cmd] = handler
}

// create partner connection
func createConnection(conn net.Conn, line string) error {
	if strings.HasPrefix(line, "AUTH:") {
		arr := strings.Split(line[5:], "#") // AUTH:API_ID#SECRET#VERSION
		if len(arr) == 3 {
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
					Conn:            conn,
					ConnectTime:     now,
					LastConnectTime: now,
				}
				clients[conn.RemoteAddr().String()] = cli
				printf(true, "[ CLIENT][ AUTH] - auth success! client id = %d ; version = %s", partnerId, arr[2])
				return nil
			}
		}
	}
	return errors.New("conn reject")
}

func handleSocketCmd(ci *ClientIdentity, cmd string) ([]byte, error) {
	i := strings.Index(cmd, ":")
	if i != -1 {
		plan := cmd[i+1:]
		if v, ok := handlers[cmd[:i]]; ok {
			return v(ci, plan)
		}
	}
	return nil, errors.New("unknown command:" + cmd)
}

func serveLoop() {
	conn := core.GetRedisConn()
	go notifyMup(conn)
}

func notifyMup(conn redis.Conn) {
	for {
		err := mmSummaryNotify(conn)
		err1 := mmAccountNotify(conn)
		if err != nil || err1 != nil {
			time.Sleep(time.Second * 1) //阻塞,避免轮询占用CPU
		}
	}
}

// member auth,command like 'MAUTH:jarrysix#3234234242342342'
func cliMAuth(id *ClientIdentity, param string) ([]byte, error) {
	arr := strings.Split(param, "#")
	if len(arr) == 2 {
		memberId, _ := strconv.Atoi(arr[0])
		b := util.CompareMemberApiToken(gof.CurrentApp.Storage(),
			memberId, arr[1])
		b = true
		if b { // auth success
			id.UserId = memberId
			// bind user activated clients
			if v, ok := users[id.UserId]; ok {
				users[id.UserId] = v + "$" + id.Addr.String()
			} else {
				users[id.UserId] = id.Addr.String()
			}
			return []byte("ok"), nil
		}
	}
	return nil, errors.New("auth fail")
}

// print text by client sending.
func cliPrint(id *ClientIdentity, params string) ([]byte, error) {
	return []byte(params), nil
}

func cliPing(id *ClientIdentity, plan string) ([]byte, error) {
	return []byte("PONG"), nil
}
