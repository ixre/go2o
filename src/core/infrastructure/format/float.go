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
	m "github.com/jsix/gof/math"
	"strconv"
	"strings"
)

func FormatFloat(f float32) string {
	//regexp : ([^\.]+)(\.|(\.[1-9]))0*$  =>  $1$3
	s := fmt.Sprintf("%.2f", f)
	if strings.HasSuffix(s, ".00") {
		return s[:len(s)-3]
	} else if strings.HasSuffix(s, "0") {
		return s[:len(s)-1]
	}
	return s
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

func FixedDecimal(amount float64) float64 {
	return m.FixedFloat(amount, 2)
}
