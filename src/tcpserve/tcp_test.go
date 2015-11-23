/**
 * Copyright 2015 @ z3q.net.
 * name : tcp_test.go
 * author : jarryliu
 * date : 2015-11-23 16:15
 * description :
 * history :
 */
package tcpserve

import (
	"testing"
	"net"
	"time"
	"fmt"
	"log"
)

func TestConn(t *testing.T){
	fmt.Println("---beigin test ---")
	raddr,err := net.ResolveTCPAddr("tcp",":1005")
	if err !=nil{
		t.Error(err)
		t.Fail()
	}

	cli,err := net.DialTCP("tcp",nil,raddr)
	if err != nil{
		t.Error(err)
		t.Fail()
	}
	var buffer []byte= make([]byte,6048)
	for i:=0;i<10;i++{
		//b,_ := encodeContent(time.Now().Format("2006年01月02日 15时04分05秒"))
		cli.Write([]byte(time.Now().Format("2006年01月02日 15时04分05秒\n")))
		n,_ := cli.Read(buffer)
		log.Println("<",string(buffer[:n]),">",n)
		time.Sleep(time.Second * 1)
	}

	//cli.Close()
}
