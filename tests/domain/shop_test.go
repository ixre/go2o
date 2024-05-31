package domain

import (
	"testing"

	"github.com/ixre/go2o/core/inject"
)

/**
 * Copyright (C) 2007-2020 56X.NET,All rights reserved.
 *
 * name : shop_test.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2020-10-08 10:24
 * description :
 * history :
 */

func TestGetShop(t *testing.T) {
	repo := inject.GetShopRepo()
	isp := repo.GetShop(1)
	if isp == nil {
		t.FailNow()
	}
}
