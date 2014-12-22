/**
 * Copyright 2014 @ Ops Inc.
 * name :
 * author : newmin
 * date : 2013-11-21 17:28
 * description :
 * history :
 */

package jsv

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	shellPrefix = "shell > "
)

//运行shell
func InShell() {
	var toServ *string
	toServ = flag.String("s", "", "object socket server")
	flag.Parse()

	const slogan string = " [ Json Socket client v 1.0 ] "
	fmt.Println(strings.Repeat("*", len(slogan)))
	fmt.Println(slogan)
	fmt.Println(strings.Repeat("*", len(slogan)))

	var buffer []byte = make([]byte, 100)
	var readerBuffer []byte = make([]byte, 204800)
	var conn *TCPConn

	AutoResetConn = false //disable auto reset connection

connect:
	if len(*toServ) == 0 {
		fmt.Print("Input Server: ")
		n, _ := os.Stdin.Read(buffer)
		*toServ = string(buffer[:n-1])
	}
	conn = Dial("tcp", *toServ)
	_, err := conn.Write([]byte(""))
	if err != nil {
		fmt.Println("Connect failed.details:", err)
		*toServ = ""
		goto connect
	} else {
		fmt.Print(shellPrefix)
	}

	for {
		n, _ := os.Stdin.Read(buffer)
		//		var line string
		//		fmt.Scanln(&line)

		if n > 1 {
			conn.Write(buffer[:n-1]) //not contain "\n"
			n, err = conn.Read(readerBuffer)
			if err != nil {
				fmt.Println("[Error]:", err)
			} else {
				fmt.Println("[Return]:", string(readerBuffer[:n]))
			}
			fmt.Print("\n" + shellPrefix)
		} else {
			fmt.Print(shellPrefix)
		}

	}
}

//func handle(conn *jsserv.TCPConn,line string){
//	conn.Write([]byte(line + "\n"))
//	n, err := conn.Read(readerBuffer)
//	if err != nil {
//		fmt.Println("[Error]:", err)
//	} else {
//		fmt.Println("[Return]:", string(readerBuffer[:n]))
//	}
//
//	fmt.Print("\n" + shellPrefix)
//}
//
//func initHandle(conn *jsserv.TCPConn){
//	handle(conn,`{"baby_id":"125","lng":"102.1","lat":"24.01","time":"1416564146"}>>Point.Post`)
//	handle(conn,`{"baby_id":"125","time_span":"1805"}>>Point.Gets`)
//}
