package service

import (
	"com/service/server"
	"io"
	"net"
	"ops/cf/app"
	"ops/cf/net/jsv"
	"time"
)

var (
	context app.Context
)

const (
	maxLength int = 1024
)

//正在服务器上拿数据
func ServerListen(n, host string, c app.Context) {
	context = c
	jsv.Configure(c)
	serv := jsv.NewServer()
	serv.RegisterName("Member", &server.Member{})
	serv.RegisterName("Partner", &server.Partner{})
	serv.RegisterName("Share", &server.Share{})

	addr, err := net.ResolveTCPAddr(n, host)
	checkErr(err)

	lis, err := net.ListenTCP(n, addr)
	checkErr(err)

	for {
		conn, err := lis.Accept()
		if err != nil {
			context.Log().Println("[Error]:", err.Error())
			continue
		}
		// set timeout
		t := time.Now().Add(120 * time.Second)
		conn.SetReadDeadline(t)
		conn.SetWriteDeadline(t)
		go receiveConn(conn, serv)
	}
}

// Receivce client request
// command defined example :
// 	[Member.Test]param1,param2\n
func receiveConn(conn net.Conn, serv *jsv.Server) {
	var buffer []byte
	//var err error

	if context.Debug() {
		context.Log().Println("[Client]: new client connecting...", conn.RemoteAddr().String())
	}

	buffer = make([]byte, maxLength)
	for {
		n, err := conn.Read(buffer)
		// if client connection closed will happen io eof
		// other reason example network timeout will happen
		// timeout error
		if err != nil {
			if err == io.EOF {
				context.Log().Println("[Disconnect]: client address :", conn.RemoteAddr().String())
			}
			conn.Close()
			break
		}

		if context.Debug() {
			context.Log().Println("[Client][Send]:", string(buffer[:n]))
		}
		serv.HandleRequest(conn, buffer[:n])
		buffer = make([]byte, maxLength)
	}
}

func checkErr(err error) {
	if err != nil {
		context.Log().Fatalf("[Error]:", err.Error())
	}
}
