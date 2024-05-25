package domain

import (
	"errors"
	"math"
)

// 计算支付方式标志
func MathPaymentMethodFlag(methods []int) int {
	f := 0
	for _, v := range methods {
		f |= 1 << uint(v-1)
	}
	return f
}

// 支付方式且运算
func AndPayMethod(payFlag int, method int) bool {
	f := 1 << uint(method-1)
	return payFlag&f == f
}

// GrantFlag 添加标志
func GrantFlag(field, flag int) (int, error) {
	f := int(math.Abs(float64(flag)))
	if f&(f-1) != 0 {
		return flag, errors.New("not right flag value")
	}
	if flag > 0 {
		// 添加标志
		if field&f != f {
			field |= flag
		}
	} else {
		// 去除标志
		if field&f == f {
			field ^= f
		}
	}
	return field, nil
}

// TestFlag 测试是否包含标志
func TestFlag(field, flag int) bool {
	f := int(math.Abs(float64(flag)))
	return field&f == f
}
