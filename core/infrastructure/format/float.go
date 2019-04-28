/**
 * Copyright 2014 @ z3q.net.
 * name : float.go
 * author : jarryliu
 * date : 2013-12-02 21:34
 * description :
 * history :
 */

package format

import (
	"fmt"
	m "github.com/ixre/gof/math"
	"log"
	"strconv"
	"strings"
)

func FormatFloat64(f float64) string {
	//regexp : ([^\.]+)(\.|(\.[1-9]))0*$  =>  $1$3
	s := fmt.Sprintf("%.2f", f)
	if s == "NaN" {
		log.Println("----[float][Nan] ", f)
	}
	if strings.HasSuffix(s, ".00") {
		return s[:len(s)-3]
	} else if strings.HasSuffix(s, "0") {
		return s[:len(s)-1]
	}
	return s
}

func FormatFloat(f float32) string {
	return FormatFloat64(float64(f))
}

func IntToFloatAmount(i int) string {
	return DecimalToString(float64(i) / 100)
}

func DecimalToString(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func ToDiscountStr(discount int) string {
	if discount == 0 || discount == 100 {
		return ""
	}
	s := strconv.Itoa(discount)
	if s[:1] == "0" {
		return s[:1]
	}
	return s[:1] + "." + s[1:]
}

func RoundAmount(amount float32) float32 {
	return m.Round32(amount, 2)
}

// 普通近似值计算, 不四舍五入,n为小数点精度
func FixedDecimalN(amount float64, n int) float64 {
	return m.FixedFloat(amount, n)
}

// 精确到2位浮点数
func FixedDecimal(amount float64) float64 {
	return m.FixedFloat(amount, 2)
}
