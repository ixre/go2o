/**
 * Copyright 2015 @ z3q.net.
 * name : lang_test
 * author : jarryliu
 * date : 2016-07-18 12:26
 * description :
 * history :
 */
package tests

import (
	"github.com/ixre/gof/log"
	"strconv"
	"testing"
	"time"
)

type (
	A struct{}
	B struct {
		*A
	}
)

func (a *A) Call() {
	a.Hello()
}

func (a *A) Hello() {
	log.Println("--hello from A")
}

func (b *B) Hello() {
	log.Println("--hello from B")
}

func (a *B) Call2() {
	a.Call()
}

func TestOverride(t *testing.T) {
	a := &A{}
	b := &B{a}
	b.Call2()
}

func TestI64ToStr(t *testing.T) {
	s := strconv.Itoa(int(time.Now().UnixNano()))
	t.Log(s)
}

// 求幂
func TestPow(t *testing.T) {
	i := 7
	t.Log(i & 2)
}
