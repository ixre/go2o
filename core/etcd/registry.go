package etcd

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : etcd_registry
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-02 15:38
 * description :
 * history :
 */

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"hash/crc32"
	"log"
	"net"
	"strings"
	"time"
)

var prefix = "/registry/server/"

type Registry interface {
	// 创建租期/注册节点,返回租期ID和错误, 如果IP为空,则默认为第一个网卡首个IP
	Register(ip string, port int) (int64, error)
	// 撤销租期/注销节点
	Revoke(LeaseID int64) error
	UnRegister()
}

var _ Registry = new(registryServer)

type registryServer struct {
	cli        *clientv3.Client
	stop       chan bool
	isRegistry bool
	leaseID    clientv3.LeaseID // 租约ID
	service    string
	ttl        int64 // 租约时间
}

// 创建服务注册, ttl租约时间
func NewRegistry(service string, ttl int64, config clientv3.Config) (Registry, error) {
	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	return &registryServer{
		service:    service,
		ttl:        ttl,
		stop:       make(chan bool),
		isRegistry: false,
		cli:        cli,
	}, nil
}
func (s *registryServer) resolveIp() string {
	addrList, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalln(err.Error())
	}
	for _, address := range addrList {
		// 检查ip地址判断是否回环地址
		if i, ok := address.(*net.IPNet); ok && !i.IP.IsLoopback() {
			if i.IP.To4() != nil {
				return i.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
func (s *registryServer) Register(ip string, port int) (leaseId int64, err error) {
	if s.isRegistry {
		panic("only one nodes can be registered")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.ttl)*time.Second)
	defer cancel()
	if len(strings.TrimSpace(ip)) == 0 {
		ip = s.resolveIp()
	}
	addr := fmt.Sprintf("%s:%d", ip, port)
	// 创建租约
	grant, err := s.cli.Grant(context.Background(), s.ttl)
	if err != nil {
		return -1, err
	}
	var node = Node{
		Id:   s.HashKey(addr),
		Addr: addr,
	}
	nodeVal, err := s.GetVal(node)
	if err != nil {
		return -1, err
	}
	//　存储键值,注册服务
	_, err = s.cli.Put(ctx, s.GetKey(node), nodeVal, clientv3.WithLease(grant.ID))
	if err != nil {
		return -1, err
	}
	s.leaseID = grant.ID
	s.isRegistry = true
	go s.KeepAlive()
	return int64(s.leaseID), nil
}
func (s *registryServer) UnRegister() {
	if s.isRegistry {
		s.stop <- true
	}
}

// 注销服务
func (s *registryServer) Revoke(leaseId int64) error {
	return s.revoke(clientv3.LeaseID(leaseId))
}

// 注销服务
func (s *registryServer) revoke(leaseID clientv3.LeaseID) error {
	// 撤销租期
	_, err := s.cli.Revoke(context.TODO(), leaseID)
	if err != nil {
		log.Printf("[Revoke] err : %s", err.Error())
	}
	s.isRegistry = false
	return err
}

// 续租/监听服务
func (s *registryServer) KeepAlive() error {
	//　续租
	keepAliveCh, err := s.cli.KeepAlive(context.TODO(), s.leaseID)
	if err != nil {
		log.Printf("[KeepAlive] err : %s", err.Error())
		return err
	}
	for {
		select {
		case <-s.stop:
			_ = s.revoke(s.leaseID)
			return nil
		case _, ok := <-keepAliveCh:
			if !ok {
				_ = s.revoke(s.leaseID)
				return nil
			}
		}
	}
}
func (s *registryServer) GetKey(node Node) string {
	return fmt.Sprintf("%s%s/%d", prefix, s.service, s.HashKey(node.Addr))
}
func (s *registryServer) GetVal(node Node) (string, error) {
	data, err := json.Marshal(&node)
	return string(data), err
}
func (s *registryServer) HashKey(addr string) uint32 {
	return crc32.ChecksumIEEE([]byte(addr))
}
