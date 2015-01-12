package jsv

import (
	"errors"
	"io"
	"net"
	"strings"
	"sync"
)

type TCPConn struct {
	net  string
	addr string
	*sync.Pool
}

func Dial(n, addr string) (*TCPConn, error) {
	return (&TCPConn{net: n, addr: addr}).init()
}

func (this *TCPConn) dial() (net.Conn, error) {
	conn, err := net.Dial(this.net, this.addr)
	if err != nil {
		Printf("[Conn][Error]: %+v\n", err)
		return nil, err
	}
	return conn, nil
}

func (this *TCPConn) init() (*TCPConn, error) {
	c, err := this.dial()
	if err != nil {
		return nil, err
	}

	this.Pool = &sync.Pool{}
	this.Pool.Put(c)
	return this, nil
}

func (this *TCPConn) getConn() net.Conn {
	var c net.Conn
	v := this.Pool.Get()
	if v != nil {
		c = v.(net.Conn)
	} else {
		c, err := this.dial()
		if err == nil {
			this.Pool.Put(c)
		}
		return c
	}
	return c
}

func (this *TCPConn) putConn(c net.Conn) {
	this.Pool.Put(c)
}

func (this *TCPConn) Close() {
	for {
		if c := this.Pool.Get(); c == nil {
			break
		} else {
			conn := c.(net.Conn)
			if conn != nil {
				conn.Close()
			}
		}
	}

}

func (this *TCPConn) Write(b []byte) (int, error) {
	var n int
	var err error
	c := this.getConn()
	n, err = c.Write(b)
	if err == nil {
		this.putConn(c)
	} else {
		//		if strings.Index(err.Error(), "broken") != -1 {
		//			this.resetConn()
		//		}
	}
	return n, err
}

func (this *TCPConn) Read(b []byte) (int, error) {
	var n int
	var err error
	c := this.getConn()
	n, err = c.Read(b)
	if err == nil {
		this.putConn(c)
	} else {
		return 0, err
	}
	return n, err
}

// write data and  read socket server response
// size is the max length for required data
// v is unmarsh struct
func (this *TCPConn) WriteAndDecode(b []byte, v interface{}, size int) error {
rewrite:
	var n int
	var err error
	c := this.getConn()
	n, err = c.Write(b)
	if err != nil {
		if strings.Index(err.Error(), "broken") != -1 {
			goto rewrite
		}
	}

	var buffer []byte
	if size == -1 {
		size = defaultBytesSize
	}
	buffer = make([]byte, size)
	n, err = c.Read(buffer)

	if err != nil {
		if err == io.EOF {
			goto rewrite
		}
		return err
	}

	this.putConn(c)

	if err = JsonCodec.Unmarshal(buffer[:n], &v); err != nil {
		return errors.New(string(buffer[:n]))
	}
	return nil
}

func (this *TCPConn) WriteAndRead(b []byte, d []byte) error {
rewrite:
	var n int
	var err error
	c := this.getConn()
	n, err = c.Write(b)
	if err != nil {
		if strings.Index(err.Error(), "broken") != -1 {
			goto rewrite
		}
	}

	n, err = c.Read(d)
	d = d[:n]

	if err != nil {
		if err == io.EOF && AutoResetConn {
			goto rewrite
		}
	}

	this.putConn(c)

	return err
}
