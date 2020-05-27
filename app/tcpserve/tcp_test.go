/**
 * Copyright 2015 @ to2.net.
 * name : tcp_test.go
 * author : jarryliu
 * date : 2015-11-23 16:15
 * description :
 * history :
 */
package tcpserve

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"testing"
)

func TestConn(t *testing.T) {
	var ch chan bool = make(chan bool)
	fmt.Println("---beigin test ---")
	rAddr, err := net.ResolveTCPAddr("tcp", ":14197")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	cli, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	var buffer []byte = make([]byte, 6048)

	cli.Write([]byte("AUTH:6000037440#0befdb52f387cc93#1.0\n"))
	n, _ := cli.Read(buffer)
	line := string(buffer[:n])
	if line != "ok\n" {
		log.Println(line)
		return
	}
	log.Println("merchant auth success")
	cli.Write([]byte("MAUTH:1#25245e2640237ea0681ed8ce1542756543111b1e750e238eafa926\n"))
	n, _ = cli.Read(buffer)
	line = string(buffer[:n])
	if line != "ok\n" {
		log.Println(line)
		return
	}
	log.Println("member auth success")
	go listenTcp(cli)
	<-ch
}

func listenTcp(conn net.Conn) {
	for {
		buf := bufio.NewReader(conn)
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		log.Println(line)
	}
}
