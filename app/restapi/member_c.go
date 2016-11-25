/**
 * Copyright 2015 @ z3q.net.
 * name : partner_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package restapi

import (
    "fmt"
    "github.com/jsix/goex/echox"
    "github.com/jsix/gof"
    "github.com/labstack/echo"
    "go2o/app/cache"
    "go2o/core/domain/interface/member"
    "go2o/core/dto"
    "go2o/core/infrastructure/domain"
    "go2o/core/service/rsi"
    "go2o/core/service/thrift/idl/gen-go/define"
    "go2o/core/variable"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "go2o/core/service/thrift"
    "time"
    "github.com/jsix/gof/util"
    "errors"
)

// 会员登录后才能调用接口
type MemberC struct {
}

// 登录
func (mc *MemberC) Login(c echo.Context) error {
    var result dto.MemberLoginResult
    r := c.Request()
    usr := strings.TrimSpace(r.FormValue("usr"))
    pwd := strings.TrimSpace(r.FormValue("pwd"))
    if len(usr) == 0 || len(pwd) == 0 {
        result.Message = "会员不存在"
    } else {
        encodePwd := domain.MemberSha1Pwd(pwd)
        mp, _ := rsi.MemberService.Login(usr, encodePwd, true)
        id := mp["Id"]
        rst := mp["Result"]
        if id > 0 {
            cli, err := thrift.MemberClient()
            if err == nil {
                defer cli.Transport.Close()
                token, _ := cli.ResetToken(id)
                result.Member = &dto.LoginMember{
                    Id:         int(id),
                    Token:      token,
                    UpdateTime: time.Now().Unix(),
                }
                result.Result = true
            } else {
                result.Member = "网络连接失败"
            }
        } else {
            switch rst {
            case -1:
                result.Message = member.ErrNoSuchMember.Error()
            case -2:
                result.Message = member.ErrCredential.Error()
            case -3:
                result.Message = member.ErrDisabled.Error()
            default:
                result.Message = "登陆失败"
            }
        }
    }
    return c.JSON(http.StatusOK, result)
}

// 注册
func (mc *MemberC) Register(c echo.Context) error {
    r := c.Request()
    result := gof.Message{}
    mchId := getMerchantId(c)
    usr := r.FormValue("usr")
    pwd := r.FormValue("pwd")
    phone := r.FormValue("phone")
    registerFrom := r.FormValue("reg_from")          // 注册来源
    invitationCode := r.FormValue("invitation_code") // 邀请码
    var regIp string
    if i := strings.Index(r.RemoteAddr, ":"); i != -1 {
        regIp = r.RemoteAddr[:i]
    }
    m := &define.Member{}
    pro := &define.Profile{}
    m.Usr = usr
    m.Pwd = domain.MemberSha1Pwd(pwd)
    m.RegIp = regIp
    m.RegFrom = registerFrom
    pro.Phone = phone
    pro.Name = m.Usr
    _, err := rsi.MemberService.RegisterMember(mchId,
        m, pro, "", invitationCode)
    return c.JSON(http.StatusOK, result.Error(err))
}

func (mc *MemberC) Ping(c echo.Context) error {
    //log.Println("---", ctx.Request.FormValue("member_id"), ctx.Request.FormValue("member_token"))
    return c.String(http.StatusOK, "PONG")
}

// 同步
func (mc *MemberC) Async(c echo.Context) error {
    var rlt AsyncResult
    var form = url.Values(c.Request().Form)
    var mut, aut, kvMut, kvAut int
    memberId := int32(GetMemberId(c))
    mut, _ = strconv.Atoi(form.Get("member_update_time"))
    aut, _ = strconv.Atoi(form.Get("account_update_time"))
    mutKey := fmt.Sprintf("%s%d", variable.KvMemberUpdateTime, memberId)
    sto.Get(mutKey, &kvMut)
    autKey := fmt.Sprintf("%s%d", variable.KvAccountUpdateTime, memberId)
    sto.Get(autKey, &kvAut)
    if kvMut == 0 {
        m, _ := rsi.MemberService.GetMember(memberId)
        kvMut = int(m.UpdateTime)
        sto.Set(mutKey, kvMut)
    }
    //kvAut = 0
    if kvAut == 0 {
        acc := rsi.MemberService.GetAccount(memberId)
        kvAut = int(acc.UpdateTime)
        sto.Set(autKey, kvAut)
    }
    rlt.MemberId = memberId
    rlt.MemberUpdated = kvMut != mut
    rlt.AccountUpdated = kvAut != aut
    return c.JSON(http.StatusOK, rlt)
}

// 获取最新的会员信息
func (mc *MemberC) Get(c echo.Context) error {
    memberId := GetMemberId(c)
    m, _ := rsi.MemberService.GetMember(memberId)
    cli, err := thrift.MemberClient()
    if err == nil {
        defer cli.Transport.Close()
        m.DynamicToken, _ = cli.GetToken(memberId)
    }
    return c.JSON(http.StatusOK, m)
}

// 汇总信息
func (mc *MemberC) Summary(c echo.Context) error {
    memberId := GetMemberId(c)
    var updateTime int64 = rsi.MemberService.GetMemberLatestUpdateTime(memberId)
    var v *dto.MemberSummary = new(dto.MemberSummary)
    var key = fmt.Sprintf("cac:mm:summary:%d", memberId)
    if cache.GetKVS().Get(key, &v) != nil || v.UpdateTime < updateTime {
        v = rsi.MemberService.GetMemberSummary(memberId)
        cache.GetKVS().SetExpire(key, v, 3600 * 48) // cache 48 hours
    }
    return c.JSON(http.StatusOK, v)
}

// 获取最新的会员账户信息
func (mc *MemberC) Account(c echo.Context) error {
    memberId := GetMemberId(c)
    m := rsi.MemberService.GetAccount(memberId)
    return c.JSON(http.StatusOK, m)
}

// 断开
func (mc *MemberC) Disconnect(c *echox.Context) error {
     result := gof.Message{}
    mStr := c.QueryParam("member_id")
    memberId, err := util.I32Err(strconv.Atoi(mStr))
    token := c.QueryParam("token")
    cli, err := thrift.MemberClient()
    if err == nil {
        defer cli.Transport.Close()
        if b,_ := cli.CheckToken(memberId,token);b{
            cli.RemoveToken(memberId)
        }else{
            err = errors.New("error credential")
        }
    }
    return c.JSON(http.StatusOK, result.Error(err))
}
