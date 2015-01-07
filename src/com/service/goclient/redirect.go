package goclient

import (
	"github.com/newmin/gof/net/jsv"
)

type redirectClient struct {
	conn *jsv.TCPConn
}

func (this *redirectClient) Post(text []byte, len int) []byte {
	this.conn.Write(text)

	var buffer []byte
	if len < 1 {
		len = 512
	}
	buffer = make([]byte, len)
	n, _ := this.conn.Read(buffer)
	return buffer[:n]
}
