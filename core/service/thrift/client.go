/**
 * Copyright 2015 @ at3.net.
 * name : client.go
 * author : jarryliu
 * date : 2016-11-13 12:32
 * description :
 * history :
 */
package thrift

import (
	"crypto/tls"
	"git.apache.org/thrift.git/lib/go/thrift"
	"go2o/core/service/thrift/idl/gen-go/define"
)

var (
	cliHostPort      = "localhost:14288"
	transportFactory = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory  = thrift.NewTCompactProtocolFactory()
)

// 客户端初始化
func CliInit(hostPort string) {
	cliHostPort = hostPort
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
	transport := transportFactory.GetTransport(serveTransport)
	return transport, protocolFactory, err
}

func MemberClient() (*define.MemberServiceClient, error) {
	transport, protocol, err := getTransportAndProtocol()
	if err == nil {
		err = transport.Open()
		if err == nil {
			//多个服务
			proto := protocol.GetProtocol(transport)
			opProto := thrift.NewTMultiplexedProtocol(proto, "member")
			return define.NewMemberServiceClientProtocol(transport, proto, opProto), err
			//单个服务
			//return define.NewMemberServiceClientFactory(transport, protocol), err
		}
	}
	return nil, err
}
