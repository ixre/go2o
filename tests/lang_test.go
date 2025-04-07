/**
 * Copyright 2015 @ 56x.net.
 * name : lang_test
 * author : jarryliu
 * date : 2016-07-18 12:26
 * description :
 * history :
 */
package tests

import (
	"log"
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

func TestModifySlice(t *testing.T) {
	s := []int{1, 1, 2, 3, 4, 8, 9}
	for i, v := range s {
		if v%2 == 0 {
			s = append(s[:i], s[i+1:]...)
		}
		t.Log(i, v, s)
	}
}


func TestFormatFloat64(t *testing.T){
	var v float64 = 3.34
	s := strconv.FormatFloat(v, 'g', 3, 64)
	log.Printf("float number :%g => :%s",v,s)
}