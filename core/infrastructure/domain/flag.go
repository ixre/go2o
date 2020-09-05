package domain

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
