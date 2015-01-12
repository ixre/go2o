package server

import (
	"com/domain/interface/member"
	"com/ording/dproxy"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/atnet/gof/crypto"
	"github.com/atnet/gof/net/jsv"
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
	b, e, err := dproxy.MemberService.Login(usr, pwd)
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

func (this *Member) GetMember(m *jsv.Args, r *jsv.Result) error {
	memberId, err := Verify(m)
	if err != nil {
		return err
	}

	e := dproxy.MemberService.GetMember(memberId)
	if e != nil {
		r.Data = e
		r.Result = true
	}
	return nil
}

func (this *Member) GetMemberAccount(m *jsv.Args, r *jsv.Result) error {
	memberId, err := Verify(m)
	if err != nil {
		return err
	}
	e := dproxy.MemberService.GetAccount(memberId)
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
	e := dproxy.MemberService.GetBank(memberId)
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
	err = dproxy.MemberService.SaveBankInfo(&e)

	if err != nil {
		jsv.LogErr(err)
		return err
	} else {
		r.Result = true
	}
	return nil
}

func (this *Member) GetBindPartner(m *jsv.Args, r *jsv.Result) error {
	memberId, err := Verify(m)
	if err != nil {
		return err
	}
	re := dproxy.MemberService.GetRelation(memberId)
	e := dproxy.PartnerService.GetPartner(re.Reg_PtId)

	if e != nil {
		e.Pwd = ""
	}
	//todo:
	//	if e == nil {
	//		e = dao.Partner().GetPartnerById(1000)
	//	}

	if e != nil {
		r.Data = e
		r.Result = true
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
	_, err = dproxy.MemberService.SaveMember(&e)

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
	r.Data = dproxy.MemberService.GetDeliverAddrs(memberId)
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
	r.Data = dproxy.MemberService.GetDeliverAddrById(memberId, addrId)
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

	_, err = dproxy.MemberService.SaveDeliverAddr(memberId, &e)
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

	if err = dproxy.MemberService.DeleteDeliverAddr(memberId, addrId); err == nil {
		r.Result = true
	} else {
		r.Data = err.Error()
	}
	return nil
}
