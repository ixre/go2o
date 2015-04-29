/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package service

import (
	"github.com/atnet/gof"
	"github.com/atnet/gof/net/jsv"
	"go2o/src/core/service/server"
	"io"
	"net"
	"time"
)

var (
	context gof.App
)

const (
	maxLength int = 1024
)

//正在服务器上拿数据
func ServerListen(n, host string, c gof.App) {
	context = c
	jsv.Configure(c)
	serve := jsv.NewServer()
	serve.RegisterName("Member", &server.Member{})
	serve.RegisterName("Partner", &server.Partner{})

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
		if buffer[0] != byte('{') {
			conn.Write([]byte("Invalid Request"))
			conn.Close()
			break
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
