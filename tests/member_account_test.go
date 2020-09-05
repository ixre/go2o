package tests

import (
	"go2o/core/domain/interface/member"
	"go2o/tests/ti"
	"testing"
)

/**
 * Copyright 2009-2019 @ to2.net
 * name : member_account_test.go.go
 * author : jarrysix (jarrysix#gmail.com)
 * date : 2019-06-25 07:25
 * description :
 * history :
 */

func TestFlowAccount(t *testing.T) {
	repo := ti.Factory.GetMemberRepo()
	m := repo.GetMember(1)
	acc := m.GetAccount()
	balance := acc.GetValue().FlowBalance
	err := acc.Adjust(member.AccountFlow, "系统赠送", 10000, "系统", 1)
	if err == nil {
		err = acc.Charge(member.AccountFlow, "用户充值50元", 5000, "-", "")
		if err == nil {
			err = acc.Consume(member.AccountFlow, "消费150", 15000, "-", "")
		}
	}
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	nowBalance := repo.GetMember(m.GetAggregateRootId()).GetAccount().GetValue().FlowBalance
	if nowBalance != balance {
		t.Logf("before:%.2f  now:%.2f", balance, nowBalance)
		t.FailNow()
	}
}
