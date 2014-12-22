package jsv

import (
	"errors"
	"io"
	"net"
	"strings"
)

type TCPConn struct {
	net  string
	addr string
	conn net.Conn
}

func Dial(n, addr string) *TCPConn {
	return &TCPConn{net: n, addr: addr}
}

func (this *TCPConn) dial() (c net.Conn, err error) {

	// 如果非长链接，则请求完关掉
	// 常链接可以发送心跳到服务端保持链接
	// 服务端可以设置连接活动时间，到达指定时间
	// 自动关闭
	if this.conn == nil {
		// tcpAddr, err := net.ResolveTCPAddr(this.network, this.addr)
		// if err != nil {
		//	 return nil, err
		// }
		// this.conn, err = net.Dial(this.network, nil, tcpAddr)
		this.conn, err = net.Dial(this.net, this.addr)
		if err != nil {
			return nil, err
		}

		//this.conn.SetKeepAlive(true)
		println("[Socket]:new connection...")
	}
	return this.conn, nil
}

func (this *TCPConn) resetConn() error {
	Println("[Connect]:Reset the connection for socket")
	if this.conn != nil {
		err := this.conn.Close()
		this.conn = nil
		return err
	}
	return nil
}

func (this *TCPConn) Write(b []byte) (int, error) {
	var n int
	c, err := this.dial()
	if err != nil {
		return 0, err
	}
	n, err = c.Write(b)
	if err != nil {
		if strings.Index(err.Error(), "broken") != -1 {
			this.resetConn()
		}
	}
	return n, err
}

func (this *TCPConn) Read(b []byte) (int, error) {
	c, _ := this.dial()
	return c.Read(b)
}

//func (this *TCPConn) ReadResult(v interface{}, size int) error {
//	var buffer []byte
//	if size == -1 {
//		size = defaultBytesSize
//	}
//	buffer = make([]byte, size)
//	n, err := this.conn.Read(buffer)
//
//	if err != nil {
//		if err == io.EOF {
//			this.resetConn()
//		}
//		return err
//	}
//	return JsonCodec.Unmarshal(buffer[:n], &v)
//}

// write data and  read socket server response
// size is the max length for required data
// v is unmarsh struct
func (this *TCPConn) WriteAndDecode(b []byte, v interface{}, size int) error {
rewrite:
	var n int
	c, err := this.dial()
	if err != nil {
		return err
	}
	n, err = c.Write(b)
	if err != nil {
		if strings.Index(err.Error(), "broken") != -1 {
			this.resetConn()
		}
	}

	var buffer []byte
	if size == -1 {
		size = defaultBytesSize
	}
	buffer = make([]byte, size)
	n, err = this.conn.Read(buffer)

	if err != nil {
		if err == io.EOF {
			this.resetConn()
			goto rewrite
		}
		return err
	}

	if err = JsonCodec.Unmarshal(buffer[:n], &v); err != nil {
		return errors.New(string(buffer[:n]))
	}
	return nil
}

func (this *TCPConn) WriteAndRead(b []byte, d []byte) error {
rewrite:
	var n int
	c, err := this.dial()
	if err != nil {
		return err
	}
	n, err = c.Write(b)
	if err != nil {
		if strings.Index(err.Error(), "broken") != -1 {
			this.resetConn()
		}
	}

	n, err = this.conn.Read(d)
	d = d[:n]

	if err != nil {
		if err == io.EOF && AutoResetConn {
			this.resetConn()
			goto rewrite
		}
		return err
	}
	return err
}
