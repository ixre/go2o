/**
 * Copyright 2015 @ z3q.net.
 * name : lang_test
 * author : jarryliu
 * date : 2016-07-18 12:26
 * description :
 * history :
 */
package testing

import (
	"github.com/jsix/gof/log"
	"testing"
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
