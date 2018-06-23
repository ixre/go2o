package thrift

import (
	"context"
	"crypto/tls"
	"git.apache.org/thrift.git/lib/go/thrift"
)

var Context context.Context = nil

type ClientFactory struct {
	thriftServer     string // "localhost:7019"
	secureTransport  bool   // false
	tlsCertFile      string // "./cert/server.crt"
	tlsKeyFile       string // "./cert/server.key"
	transportFactory thrift.TTransportFactory
	protocolFactory  thrift.TProtocolFactory
}

func NewClientFactory(server string, secure bool, tslKeyFile string,
	tslCertFile string) *ClientFactory {
	return &ClientFactory{
		thriftServer:     server,
		secureTransport:  secure,
		tlsKeyFile:       tslKeyFile,
		tlsCertFile:      tslCertFile,
		transportFactory: thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory()),
		protocolFactory:  thrift.NewTBinaryProtocolFactoryDefault(),
	}
}

// enable debug mode
func (c *ClientFactory) Debug(prefix string) {
	c.protocolFactory = thrift.NewTDebugProtocolFactory(c.protocolFactory, prefix)
}

func (c *ClientFactory) prepare() (thrift.TTransport, thrift.TProtocolFactory, error) {
	var err error
	var transport thrift.TTransport
	if c.secureTransport {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair(c.tlsCertFile, c.tlsKeyFile); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return nil, c.protocolFactory, err
		}
		transport, err = thrift.NewTSSLSocket(c.thriftServer, cfg)
	} else {
		transport, err = thrift.NewTSocket(c.thriftServer)
	}
	if err == nil {
		transport, err = c.transportFactory.GetTransport(transport)
	}
	return transport, c.protocolFactory, err
}

// get thrift client
func (c *ClientFactory) GetClient(service string) (thrift.TTransport, thrift.TClient, error) {
	transport, protocol, err := c.prepare()
	if err == nil {
		err = transport.Open()
		if err == nil {
			proto := protocol.GetProtocol(transport)
			opProto := thrift.NewTMultiplexedProtocol(proto, service)
			return transport, thrift.NewTStandardClient(proto, opProto), nil
		}
	}
	return transport, nil, err
}
