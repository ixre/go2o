package thrift

import (
	"context"
	"crypto/tls"
	"git.apache.org/thrift.git/lib/go/thrift"
)

var (
	thriftServer                             = "localhost:14288"
	secureTransport                          = false
	rpcDebug                                 = false
	tlsCertFile                              = "./cert/server.crt"
	tlsKeyFile                               = "./cert/server.key"
	transportFactory                         = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory  thrift.TProtocolFactory = thrift.NewTCompactProtocolFactory()
	Context                                  = context.Background()
)

// 客户端初始化
func CliInit(hostPort string) {
	thriftServer = hostPort
	if rpcDebug {
		protocolFactory = thrift.NewTDebugProtocolFactory(protocolFactory, "[ Go2o][ Rpc]:")
	}
}

func getTransportAndProtocol() (thrift.TTransport, thrift.TProtocolFactory, error) {
	var err error
	var transport thrift.TTransport
	if secureTransport {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair(tlsCertFile, tlsKeyFile); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return nil, protocolFactory, err
		}
		transport, err = thrift.NewTSSLSocket(thriftServer, cfg)
	} else {
		transport, err = thrift.NewTSocket(thriftServer)
	}
	if err == nil {
		transport, err = transportFactory.GetTransport(transport)
	}
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
