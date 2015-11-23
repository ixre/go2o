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
	"net"
	"bufio"
	"io"
	"log"
	"github.com/jsix/gof/util"
	"time"
	"strings"
	"errors"
	"go2o/src/core/service/dps"
)

func printf(format string,args ...interface{}){
	log.Printf(format+"\n",args...)
}



type(
	TcpReceiveCaller func(conn net.Conn,read []byte)([]byte,error)
)

func listen(addr string,rc TcpReceiveCaller){
	serveAddr,err := net.ResolveTCPAddr("tcp",addr)
	if err != nil{
		panic(err)
	}
	listen,err := net.ListenTCP("tcp",serveAddr)
	for{
		if conn,err := listen.AcceptTCP();err == nil{
			printf("[ CONNECT][ NEW] - %s",conn.RemoteAddr().String())
			go receiveTcpConn(conn,rc)
		}
	}
}

func receiveTcpConn(conn *net.TCPConn,rc TcpReceiveCaller){
	for{
		buf := bufio.NewReader(conn)
		line,err := buf.ReadBytes('\n')
		if err == io.EOF{
			printf("[ DISCONNECT]- IP %s disconnect!",conn.RemoteAddr().String())
			break
		}
		if d,err := rc(conn,line);err != nil{
			conn.Write([]byte("error$"+err.Error()))
		}else{
			conn.Write(d)
		}
	}
}

var(
	clients map[string]*ClientIdentity
)

type(
	// the identity of client
	ClientIdentity struct {
		Id              int
		Addr          net.Addr
		ConnectTime     time.Time
		LastConnectTime time.Time
	}
)

// generate client id
func genClientId()string{
	for {
		id := util.RandString(5)
		if _,ok := clients[id];!ok{
			return id
		}
	}
}

func ListenTcp(addr string){
	listen(addr, myTcpReceive)
}

func myTcpReceive(conn net.Conn,b []byte)([]byte,error) {
	line := string(b)
	if !strings.HasPrefix(line, "CID:") {
		return chkClientPerm(conn, line)
	}
	splitPos := strings.Index(line,"$")
	cid := line[4:splitPos]
	return handleSocketCmd(cid,line[splitPos+1:])
}


// check client has own permission
func chkClientPerm (conn net.Conn,line string)([]byte,error) {
	arr := strings.Split(line, "&")    // API_ID&SECRET
	if len(arr) == 2 {
		partnerId := dps.PartnerService.GetPartnerIdByApiId(arr[0])
		apiInfo := dps.PartnerService.GetApiInfo(partnerId)
		if apiInfo != nil && apiInfo.ApiSecret == arr[1] {
			if apiInfo.Enabled == 0 {
				return nil, errors.New("api has exipres")
			}
			now := time.Now()
			clientId := genClientId()
			cli := &ClientIdentity{
				Id:partnerId,
				Addr:conn.RemoteAddr(),
				ConnectTime:now,
				LastConnectTime:now,
			}
			clients[clientId] = cli
			return []byte(clientId), nil    // return client id
		}
	}
	return nil, errors.New("conn reject")
}


func handleSocketCmd(cid string,cmd string)([]byte,error){
	return []byte(cmd),nil
}
