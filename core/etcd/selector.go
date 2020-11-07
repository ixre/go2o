package etcd

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : selector
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-09-02 17:18
 * description :
 * history :
 */

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type SelectAlgorithm int

const AlgRandom SelectAlgorithm = 0
const AlgRoundRobin SelectAlgorithm = 1

type Selector interface {
	Next() (Node, error)
}

type serverSelector struct {
	cli   *clientv3.Client
	nodes []Node
	alg   SelectAlgorithm
	last  int
	name  string
}

func NewSelector(name string, config clientv3.Config, alg SelectAlgorithm) (Selector, error) {
	cli, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	var s = &serverSelector{
		cli:  cli,
		name: name,
		alg:  alg,
		last: -1,
	}
	s.loadNodes()
	go s.Watch()
	return s, nil
}

// 使用轮询算法
func (s *serverSelector) RoundRobin() {
	s.alg = AlgRoundRobin
}

func (s *serverSelector) Next() (Node, error) {
	l := len(s.nodes)
	if l == 0 {
		return Node{}, fmt.Errorf("no nodes found on %s", s.name)
	}
	if s.alg == AlgRoundRobin {
		if s.last += 1; s.last >= l {
			s.last = 0
		}
		return s.nodes[s.last], nil
	}
	// 随机算法
	//i := rand.Int() % len(s.nodes)
	i := rand.Intn(l)
	return s.nodes[i], nil
}

func (s *serverSelector) Watch() {
	//　监听变化,并动态处理节点
	ch := s.cli.Watch(context.TODO(), prefix, clientv3.WithPrefix())
	for {
		select {
		case c := <-ch:
			for _, e := range c.Events {
				switch e.Type {
				case clientv3.EventTypePut:
					node, err := s.GetVal(e.Kv.Value)
					if err != nil {
						log.Printf("[EventTypePut] err : %s", err.Error())
						continue
					}
					s.AddNode(node)
				case clientv3.EventTypeDelete:
					keyArray := strings.Split(string(e.Kv.Key), "/")
					if len(keyArray) == 0 {
						println("[ Etcd][ Event]: delete node key is empty")
						return
					}
					nodeId, err := strconv.Atoi(keyArray[len(keyArray)-1])
					if err != nil {
						log.Printf("[EventTypePut] key Atoi : %s", err.Error())
						continue
					}
					s.DelNode(uint32(nodeId))
				}
			}
		}
	}
}

func (s *serverSelector) DelNode(id uint32) {
	var node []Node
	for _, v := range s.nodes {
		if v.Id != id {
			node = append(node, v)
		}
	}
	s.nodes = node
}

func (s *serverSelector) AddNode(node Node) {
	var exist bool
	for _, v := range s.nodes {
		if v.Id == node.Id {
			exist = true
		}
	}
	if !exist {
		s.nodes = append(s.nodes, node)
	}
}

func (s *serverSelector) GetKey() string {
	return fmt.Sprintf("%s%s", prefix, s.name)
}

func (s *serverSelector) GetVal(val []byte) (Node, error) {
	var node Node
	err := json.Unmarshal(val, &node)
	if err != nil {
		return node, err
	}
	return node, nil
}

func (s *serverSelector) loadNodes() {
	res, err := s.cli.Get(context.TODO(), s.GetKey(), clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		log.Printf("[Watch] err : %s", err.Error())
		return
	}
	for _, kv := range res.Kvs {
		node, err := s.GetVal(kv.Value)
		if err != nil {
			log.Println("[ Etcd][ Error]: load nodes failed! error : " + err.Error())
			continue
		}
		s.nodes = append(s.nodes, node)
	}
}
