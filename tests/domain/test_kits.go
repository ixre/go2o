package domain

import (
	"testing"

	"github.com/ixre/go2o/core/msq"
	_ "github.com/ixre/go2o/tests/ti"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : test_kits.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-11-04 16:46
 * description :
 * history :
 */

func assertError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func init() {
	msq.Configure(msq.NATS, []string{"192.168.0.159:4222"})
}
