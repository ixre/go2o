package goclient

import (
	"github.com/newmin/gof/net/jsv"
)

type shareClient struct {
	conn *jsv.TCPConn
}
