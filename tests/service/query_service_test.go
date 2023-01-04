package service

import (
	"context"
	"testing"

	"github.com/ixre/go2o/core/domain/interface/member"
	"github.com/ixre/go2o/core/service/impl"
	"github.com/ixre/go2o/core/service/proto"
	"github.com/ixre/gof/types/typeconv"
)

func TestPagingMemberAccountLog(t *testing.T) {
	s := impl.QueryService
	ret, _ := s.PagingMemberAccountLog(context.TODO(), &proto.PagingAccountLogRequest{
		MemberId:    702,
		ValueFilter: 2,
		AccountType: int32(member.AccountWallet),
		Params: &proto.SPagingParams{
			Begin: 0,
			End:   10,
		},
	})
	t.Log(typeconv.MustJson(ret))
}

func TestPagingIntegralLog(t *testing.T) {
	params := &proto.SPagingParams{
		Parameters: nil,
		SortBy:     "",
		Begin:      0,
		End:        10,
	}
	r, _ := impl.QueryService.PagingMemberAccountLog(context.TODO(),
		&proto.PagingAccountLogRequest{
			MemberId:    1,
			AccountType: int32(member.AccountIntegral),
			Params:      params,
		})
	t.Log(typeconv.MustJson(r))
}

func TestPagingBalanceLog(t *testing.T) {
	memberId := 702
	params := &proto.SPagingParams{
		Parameters: nil,
		SortBy:     "",
		Begin:      0,
		End:        10,
	}
	r, _ := impl.QueryService.PagingMemberAccountLog(context.TODO(),
		&proto.PagingAccountLogRequest{
			MemberId:    int64(memberId),
			AccountType: int32(member.AccountBalance),
			Params:      params,
		})
	t.Log(typeconv.MustJson(r))
}
