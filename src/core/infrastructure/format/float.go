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
	"strconv"
	"strings"
)

func FormatFloat(f float32) string {
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
