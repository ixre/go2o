package goclient

import (
	"fmt"
	"github.com/atnet/gof/app"
	"github.com/atnet/gof/net/jsv"
	"os"
)

var (
	_conn    *jsv.TCPConn
	Member   *memberClient
	Partner  *partnerClient
	Redirect *redirectClient
	Share    *shareClient
)

func Configure(net, addr string, c app.Context) {
	var err error
	_conn, err = jsv.Dial(net, addr)

	if err != nil {
		fmt.Println("[TCP]: Connect Refused,", addr)
		os.Exit(1)
	}

	jsv.Configure(c)
	Member = &memberClient{conn: _conn}
	Partner = &partnerClient{conn: _conn}
	Redirect = &redirectClient{conn: _conn}
	Share = &shareClient{conn: _conn}
}
