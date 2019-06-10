package rsi

import "github.com/ixre/gof/util"

/**
 * Copyright 2009-2019 @ to2.net
 * name : service_init.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-06-10 13:17
 * description :
 * history :
 */

func (s *memberService) generateMemberCode() string {
	var code string
	for {
		code = util.RandString(6)
		if memberId := s.repo.GetMemberIdByCode(code); memberId == 0 {
			break
		}
	}
	return code
}
