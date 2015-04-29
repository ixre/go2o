/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package server

import (
	"fmt"
	"github.com/atnet/gof/crypto"
	"github.com/atnet/gof/net/jsv"
	"github.com/garyburd/redigo/redis"
	"go2o/src/core/domain/interface/member"
	"go2o/src/core/domain/interface/partner"
	"go2o/src/core/service/dps"
	"strconv"
	"strings"
	"time"
)

type Member struct {
	Redis *redis.Pool
}

//登录验证
func (this *Member) Login(m *jsv.Args, r *jsv.Result) error {
	usr, pwd := (*m)["usr"].(string), (*m)["pwd"].(string)
	b, e, err := dps.MemberService.Login(usr, pwd)
	r.Result = b
	if !b {
		r.Message = err.Error()
	} else {
		md5 := strings.ToLower(crypto.Md5([]byte(time.Now().String())))
		rds := Redis().Get()
		rds.Do("SETEX", fmt.Sprintf("member$%d_session_key", e.Id), 3600*300, md5)
		r.Data = fmt.Sprintf("%d$%s", e.Id, md5)
		if jsv.Context.Debug() {
			jsv.Printf("[Member][Login]%d -- %s", e.Id, md5)
		}
		rds.Close()
	}
	return nil
}

func (this *Member) Verify(m *jsv.Args, r *jsv.Result) error {
	_, err := Verify(m)
	if err != nil {
		r.Message = err.Error()
	} else {
		r.Result = true
	}
	return nil
}

func (this *Member) GetMember(m *jsv.Args, r *member.ValueMember) error {
	memberId, err := Verify(m)
	if err != nil {
		return err
	}

	e := dps.MemberService.GetMember(memberId)
	if e != nil {
		*r = *e
	}
	return nil
}

func (this *Member) GetMemberAccount(m *jsv.Args, r *jsv.Result) error {
	memberId, err := Verify(m)
	if err != nil {
		return err
	}
	e := dps.MemberService.GetAccount(memberId)
	if e != nil {
		r.Data = e
		r.Result = true
	}
	return nil
}

func (this *Member) GetBankInfo(m *jsv.Args, r *jsv.Result) error {
	memberId, err := Verify(m)
	if err != nil {
		return err
	}
	e := dps.MemberService.GetBank(memberId)
	if e != nil {
		r.Data = e
		r.Result = true
	}
	return nil
}

func (this *Member) SaveBankInfo(m *jsv.Args, r *jsv.Result) error {
	memberId, err := Verify(m)
	if err != nil {
		return err
	}

	var e member.BankInfo
	err = jsv.UnmarshalMap((*m)["json"], &e)
	if err != nil {
		return err
	}
	e.MemberId = memberId
	err = dps.MemberService.SaveBankInfo(&e)

	if err != nil {
		jsv.LogErr(err)
		return err
	} else {
		r.Result = true
	}
	return nil
}

func (this *Member) GetBindPartner(m *jsv.Args, r *partner.ValuePartner) error {
	memberId, err := Verify(m)
	if err != nil {
		return err
	}
	re := dps.MemberService.GetRelation(memberId)
	e, err := dps.PartnerService.GetPartner(re.Reg_PtId)
	if err != nil {
		return err
	}

	if e != nil {
		e.Pwd = ""
	}
	//todo:
	//	if e == nil {
	//		e = dao.Partner().GetPartnerById(1000)
	//	}

	if e != nil {
		*r = *e
	}
	return nil
}

func (this *Member) SaveMember(m *jsv.Args, r *jsv.Result) error {
	memberId, err := Verify(m)
	if err != nil {
		return err
	}

	var e member.ValueMember
	err = jsv.UnmarshalMap((*m)["json"], &e)
	if err != nil {
		return err
	}
	e.Id = memberId
	_, err = dps.MemberService.SaveMember(&e)

	if err != nil {
		jsv.LogErr(err)
		r.Message = err.Error()
	} else {
		r.Result = true
	}
	return nil
}

func (this *Member) GetDeliverAddrs(m *jsv.Args, r *jsv.Result) error {
	memberId, err := Verify(m)
	if err != nil {
		return err
	}
	r.Result = true
	r.Data = dps.MemberService.GetDeliverAddrs(memberId)
	return nil
}

func (this *Member) GetDeliverAddrById(m *jsv.Args, r *jsv.Result) error {
	memberId, err := Verify(m)
	if err != nil {
		return err
	}
	addrId, err := strconv.Atoi((*m)["addr_id"].(string))
	if err != nil {
		return err
	}
	r.Result = true
	r.Data = dps.MemberService.GetDeliverAddrById(memberId, addrId)
	return nil
}

func (this *Member) SaveDeliverAddr(m *jsv.Args, r *jsv.Result) error {
	memberId, err := Verify(m)
	if err != nil {
		return err
	}

	var e member.DeliverAddress
	err = jsv.UnmarshalMap((*m)["json"], &e)
	if err != nil {
		return err
	}
	e.MemberId = memberId

	_, err = dps.MemberService.SaveDeliverAddr(memberId, &e)
	if err != nil {
		jsv.LogErr(err)
		r.Message = err.Error()
	} else {
		r.Result = true
	}
	return nil
}

func (this *Member) DeleteDeliverAddr(m *jsv.Args, r *jsv.Result) error {
	memberId, err := Verify(m)
	if err != nil {
		return err
	}
	addrId, err := strconv.Atoi((*m)["addr_id"].(string))
	if err != nil {
		return err
	}

	if err = dps.MemberService.DeleteDeliverAddr(memberId, addrId); err == nil {
		r.Result = true
	} else {
		r.Data = err.Error()
	}
	return nil
}
