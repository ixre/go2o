package member

import (
	"github.com/ixre/go2o/app/api/util"
	api "github.com/ixre/gof/jwt-api"
)

/**
 * Copyright (C) 2007-2021 56X.NET,All rights reserved.
 *
 * name : complain.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2021-02-16 14:02
 * description :
 * history :
 */

var _ api.Handler = new(ComplainApi)

type ComplainApi struct {
	util.Utils
}

func (c ComplainApi) Group() string {
	return "complain"
}

func (c ComplainApi) Process(fn string, ctx api.Context) *api.Response {
	return c.Success("投诉成功！我们将会尽快处理并对违规账号进行处罚.")
}
