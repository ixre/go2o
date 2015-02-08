/**
 * Copyright 2014 @ ops Inc.
 * name :
 * author : newmin
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package goclient

import (
	"fmt"
	"github.com/atnet/gof/net/jsv"
	"github.com/atnet/gof/web/pager"
	"go2o/core/domain/interface/member"
	"go2o/core/ording/entity"
	"go2o/core/service/goclient"
	"go2o/core/share/glob"
	"go2o/core/share/variable"
	"strconv"
	"strings"
	"time"
)

func main() {
	p := pager.NewUrlPager(10, 1, nil)
	fmt.Println(p.PagerString())
	fmt.Println("--------------------------")

	p = pager.NewUrlPager(0, 0, nil)
	fmt.Println(p.PagerString())
	fmt.Println("--------------------------")
	p = pager.NewUrlPager(10, 10, nil)
	fmt.Println(p.PagerString())

	fmt.Println("--------------------------")
	p = pager.NewUrlPager(10, 9, nil)
	fmt.Println(p.PagerString())

	return

	context := glob.NewContext()
	goclient.Configure("tcp", ":"+context.Config().GetString(variable.SocketPort), context)
	context.DebugMode = true
	jsv.Configure(context)

	//fmt.Println(ording.NewSecret(666888))

	go testRegister()
	go testPartner()
	go testSubmitOrder()

	testMemberLogin()
	for {
		time.Sleep(10 * time.Second)
		go testMemberLogin()
	}
}

func testRegister() {
	m := member.ValueMember{
		Usr:      "test",
		Pwd:      "test",
		Name:     "测试员",
		Sex:      1,
		Avatar:   "",
		Birthday: "1988-11-09",
		Phone:    "18616999822",
		Address:  "",
		Qq:       "",
		Email:    "",
		RegIp:    "127.0.0.1",
	}
	b, err := goclient.Member.Register(&m, 666888, 0, "")
	if err != nil {
		jsv.LogErr(err)

	} else {
		jsv.Println("注册成功")
	}
	b, _, _ = goclient.Member.Login(m.Usr, m.Pwd)
	if b {
		jsv.Println("登录成功")
	} else {

		jsv.Printf("登录失败：Usr:%s,Pwd:%s\n", m.Usr, m.Pwd)
	}
}

func testPartner() {
	p, err := goclient.Partner.GetPartner(666888, "d435a520e50e960b")
	if err != nil {
		jsv.LogErr(err)
	} else {
		jsv.Println(p)
	}
}

func testMemberLogin() {
	b, t, err := goclient.Member.Login("newmin", "123000")
	if b {
		jsv.Println("[Login]:Sucessfull.", t)
		arr := strings.Split(t, "$")
		id, _ := strconv.Atoi(arr[0])
		token := arr[1]
		r := goclient.Member.Verify(id, token)
		if r.Result {
			jsv.Println("成功\n")
		} else {
			jsv.Println("验证失败:", r.Message)
		}

		acc, _ := goclient.Member.GetMemberRelation(id, token)
		if acc != nil {
			jsv.Println(*acc)
		}
	} else {
		jsv.Println("[Login]:Failed.", err)
	}

}

func testSubmitOrder() {
	items := "2*1|3*2|4*1"
	orderNo, err := goclient.Partner.SubmitOrder(666888, "d435a520e50e960b",
		1, 0, 1, items, "")
	if err != nil {
		jsv.LogErr(err)
	} else {
		jsv.Println("[OrderNo]", orderNo)
	}
}
