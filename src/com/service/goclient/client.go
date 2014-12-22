package goclient

import (
	"ops/cf/app"
	"ops/cf/net/jsv"
)

var (
	_conn    *jsv.TCPConn
	Member   *memberClient
	Partner  *partnerClient
	Redirect *redirectClient
	Share    *shareClient
)

func Configure(net, addr string, c app.Context) {
	_conn = jsv.Dial(net, addr)
	jsv.Configure(c)
	Member = &memberClient{conn: _conn}
	Partner = &partnerClient{conn: _conn}
	Redirect = &redirectClient{conn: _conn}
	Share = &shareClient{conn: _conn}
}
