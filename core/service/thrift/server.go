/**
 * Copyright 2015 @ at3.net.
 * name : server.go
 * author : jarryliu
 * date : 2016-11-12 18:52
 * description :
 * history :
 */
package thrift

import (
	"crypto/tls"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"go2o/core/service/dps"
	"go2o/core/service/thrift/idl/gen-go/define"
)

func ListenAndServe(addr string, secure bool) error {
	var err error
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTCompactProtocolFactory()
	var transport thrift.TServerTransport
	if secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		transport, err = thrift.NewTSSLServerSocket(addr, cfg)
	} else {
		transport, err = thrift.NewTServerSocket(addr)
	}

	transport, err = thrift.NewTServerSocket(addr)
	if err == nil {
		processor := thrift.NewTMultiplexedProcessor()
		processor.RegisterProcessor("member", define.NewMemberServiceProcessor(dps.MemberService))
		processor.RegisterProcessor("foundation", define.NewFoundationServiceProcessor(dps.FoundationService))
		server := thrift.NewTSimpleServer4(processor, transport,
			transportFactory, protocolFactory)
		fmt.Println("Starting the thrift server... on ", addr)
		err = server.Serve()
	}
	return err
}
