package tool

import (
	"github.com/ixre/tto"
	"go2o/core/domain/interface/member"
	"testing"
)

/**
 * Copyright 2009-2019 @ to2.net
 * name : thrift_gen_test.go.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-06-24 18:55
 * description :
 * history :
 */

func TestGenerateThriftStruct(t *testing.T) {
	bytes, _ := tto.ThriftStruct(member.ComplexMember{})
	t.Log(string(bytes))
}
