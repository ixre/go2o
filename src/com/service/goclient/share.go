package goclient

import (
	"github.com/atnet/gof/net/jsv"
)

type shareClient struct {
	conn *jsv.TCPConn
}
