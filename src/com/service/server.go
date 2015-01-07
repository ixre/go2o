package service

import (
	"com/service/server"
	"io"
	"net"
	"github.com/newmin/gof/app"
	"github.com/newmin/gof/net/jsv"
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
	serve := jsv.NewServer()
	serve.RegisterName("Member", &server.Member{})
	serve.RegisterName("Partner", &server.Partner{})
	serve.RegisterName("Share", &server.Share{})

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
		t := time.Now().Add(5 * time.Minute)
		conn.SetDeadline(t)

		go receiveConn(conn, serve)
	}
}

// Receive client request
// command defined example :
// 	[Member.Test]param1,param2\n
func receiveConn(conn net.Conn, serve *jsv.Server) {
	var buffer []byte
	//var err error

	if context.Debug() {
		context.Log().Println("[CLIENT]: Client connecting...", conn.RemoteAddr().String())
	}

	buffer = make([]byte, maxLength)
	for {
		n, err := conn.Read(buffer)
		// if client connection closed will happen io eof
		// other reason example network timeout will happen
		// timeout error
		if err != nil {
			if err == io.EOF {
				context.Log().Println("[Client]: Client disconnect :", conn.RemoteAddr().String())
			}
			conn.Close()
			break
		}

		if context.Debug() {
			context.Log().Println("[Client][Send]:", string(buffer[:n]))
		}

		serve.HandleRequest(conn, buffer[:n])
		buffer = make([]byte, maxLength)
	}
}

func checkErr(err error) {
	if err != nil {
		context.Log().Fatalf("[Error]:", err.Error())
	}
}
