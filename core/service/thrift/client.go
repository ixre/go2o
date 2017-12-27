package thrift

import (
	"context"
	"crypto/tls"
	"git.apache.org/thrift.git/lib/go/thrift"
)

var (
	cliHostPort                              = "localhost:14288"
	rpcDebug                                 = false
	transportFactory                         = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory  thrift.TProtocolFactory = thrift.NewTCompactProtocolFactory()
	Context                                  = context.Background()
)

// 客户端初始化
func CliInit(hostPort string) {
	cliHostPort = hostPort
	rpcDebug = true
	if rpcDebug {
		protocolFactory = thrift.NewTDebugProtocolFactory(protocolFactory, "[ Go2o][ Rpc]:")
	}
}

func getTransportAndProtocol() (thrift.TTransport, thrift.TProtocolFactory, error) {
	var err error
	var serveTransport thrift.TTransport
	secure := false
	if secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return nil, protocolFactory, err
		}
		serveTransport, err = thrift.NewTSSLSocket(cliHostPort, cfg)
	} else {
		serveTransport, err = thrift.NewTSocket(cliHostPort)
	}
	transport, err := transportFactory.GetTransport(serveTransport)
	return transport, protocolFactory, err
}

func getClient(service string) (thrift.TTransport, thrift.TClient, error) {
	transport, protocol, err := getTransportAndProtocol()
	if err == nil {
		err = transport.Open()
		if err == nil {
			//多个服务
			proto := protocol.GetProtocol(transport)
			opProto := thrift.NewTMultiplexedProtocol(proto, service)
			return transport, thrift.NewTStandardClient(proto, opProto), nil
		}
	}
	return transport, nil, err
}
