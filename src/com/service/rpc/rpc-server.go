package main

import (
	srpc "com/service/rpc"
	"com/service/server"
	"com/share/glob"
	"com/share/variable"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"strconv"
	"time"
)

func main() {
	var (
		port    *int
		debug   *bool
		httpRun *bool
		context *glob.AppContext
	)
	context = glob.NewContext()
	socketPort, _ := strconv.Atoi(context.Config().
		Get(variable.SocketPort))

	port = flag.Int("port", socketPort, "tcp server")
	debug = flag.Bool("debug", false, "")
	httpRun = flag.Bool("http", false, "running with http")

	flag.Parse()

	if *debug {
		fmt.Println("[Started]:Service (with debug) running on port [" +
			strconv.Itoa(*port) + "]:")
		context.DebugMode = true
		log.SetOutput(os.Stdout)
		//context.Db().ORM.SetTrace(true)
	} else {
		fmt.Println("[Started]:Service running on port [" +
			strconv.Itoa(*port) + "]:")
	}

	//添加RPC
	srpc.RPC_CONTEXT = context
	member := &server.Member{}
	partner := &server.Partner{}
	rpc.RegisterName("Member", member)
	rpc.RegisterName("Partner", partner)

	if *httpRun {
		hostWithHttp(port)
	} else {
		hostWithSocket(port)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("[Error]:", err.Error())
	}
}

// 以Socket方式承载
func hostWithSocket(port *int) {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", *port))
	checkErr(err)

	lis, err := net.ListenTCP("tcp", addr)
	checkErr(err)
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("[Error]:", err.Error())
			continue
		}

		// set timeout
		t := time.Now().Add(120 * time.Second)
		conn.SetReadDeadline(t)
		conn.SetWriteDeadline(t)

		//设置JsonRpc
		go jsonrpc.ServeConn(conn)
	}
}

func hostWithHttp(port *int) {
	rpc.HandleHTTP()
	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Println("[Server][Error]:", err.Error())
	}
}
