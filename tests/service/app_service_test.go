/**
 * Copyright (C) 2007-2024 fze.NET, All rights reserved.
 *
 * name: app_service_test.go
 * author: jarrysix (jarrysix@gmail.com)
 * date: 2024-11-11 23:37:18
 * description:
 * history:
 */

package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/inject"
	"github.com/ixre/go2o/core/service/proto"
)

func TestCheckAppVersion(t *testing.T) {
	s := inject.GetAppService()
	ret, _ := s.CheckAppVersion(context.TODO(), &proto.CheckAppVersionRequest{
		AppName:         "app",
		TerminalOS:      "android",
		TerminalChannel: "beta",
		Version:         "1.0.0",
	})
	t.Log(ret)
}
