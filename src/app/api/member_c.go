/**
 * Copyright 2015 @ S1N1 Team.
 * name : partner_c.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package api
import (
    "github.com/atnet/gof/web"
    "go2o/src/core/service/dps"
    "go2o/src/core/dto"
    "github.com/atnet/gof/crypto"
    "github.com/atnet/gof"
    "errors"
    "fmt"
    "github.com/atnet/gof/web/mvc"
    "strconv"
)

var _ mvc.Filter = new(memberC)
type memberC struct{
    *baseC
}

// 会员登陆后才能调用接口
func (this *memberC) Requesting(ctx *web.Context)bool {
    if this.baseC != nil && this.baseC.Requesting(ctx) {
        r := ctx.Request
        memberId, _ := strconv.Atoi(r.FormValue("member_id"))
        token := r.FormValue("token")
        if chkMemberToken(ctx.App.Storage(), memberId, token) {
            return true
        }
        this.errorOutput(ctx,"invalid request!")
    }
    return false
}

const offset string = "%$^&@#"
func chkStorage(sto gof.Storage){
    if sto == nil{
        panic(errors.New("[ Api] - api token storage is null !"))
    }
}
func getMemberTokenKey(memberId int)string{
    return fmt.Sprintf("api:member:token:%d",memberId)
}

// 设置令牌，并返回
func setMemberToken(sto gof.Storage,memberId int,pwd string)string{
    chkStorage(sto)
    cyp := crypto.NewUnixCrypto(pwd+offset,offset)
    var token string = string(cyp.Encode())
    var key string = getMemberTokenKey(memberId)

    sto.Set(key,token)      // 存储令牌
    sto.Set(key+"base",pwd) // 存储令牌凭据

    return token
}

// 校验令牌
func chkMemberToken(sto gof.Storage,memberId int,token string)bool{
    chkStorage(sto)

    if len(token) == 0{
        return false
    }

    var key = getMemberTokenKey(memberId)
    var srcToken,tokenBase string

    sto.Get(key,&srcToken)
    sto.Get(key+"base",&tokenBase)

    if len(srcToken) ==0 || len(tokenBase) == 0{
        return false
    }

    cyp := crypto.NewUnixCrypto(tokenBase+offset,offset)
    b,_,_ := cyp.Compare(token)
    return b
}

// 处理请求
func (this *memberC) handle(ctx *web.Context){
    mvc.Handle(this,ctx,false)
}

// 登陆
func (this *memberC) login(ctx *web.Context){
    if this.baseC.Requesting(ctx) {

        r := ctx.Request
        var usr, pwd string = r.FormValue("usr"), r.FormValue("pwd")
        var result dto.MemberLoginResult;

        if len(usr) == 0 || len(pwd) == 0 {
            result.Message = "会员不存在"
        }else {
            b, e, err := dps.MemberService.Login(usr, pwd)
            result.Result = b

            if b {
                // 生成令牌
                e.DynamicToken = setMemberToken(ctx.App.Storage(), e.Id, e.Pwd)
                result.Member = e
            }
            if err != nil {
                result.Message = err.Error()
            }
        }

        this.jsonOutput(ctx, result)
    }
}

// 注册
func (this *memberC) register(ctx *web.Context){
    if this.baseC.Requesting(ctx) {
        var result = dto.MessageResult{}

        //        r, w := ctx.Request, ctx.ResponseWriter
        //        p := this.GetPartner(ctx)
        //        r.ParseForm()
        //        var member member.ValueMember
        //        web.ParseFormToEntity(r.Form, &member)
        //        if i := strings.Index(r.RemoteAddr, ":"); i != -1 {
        //            member.RegIp = r.RemoteAddr[:i]
        //        }
        //        b, err := goclient.Partner.RegisterMember(&member, p.Id, 0, "")
        //        if b {
        //            w.Write([]byte(`{"result":true}`))
        //        } else {
        //            if err != nil {
        //                w.Write([]byte(fmt.Sprintf(`{"result":false,"message":"%s"}`, err.Error())))
        //            } else {
        //                w.Write([]byte(`{"result":false}`))
        //            }
        //
        //        }
        //
        //
        this.jsonOutput(ctx, result)
    }
}