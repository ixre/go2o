package domain

import (
	"testing"
)

func TestAndPayMethod(t *testing.T) {
	flag := MathPaymentMethodFlag([]int{1, 2, 5})
	println("--value=", flag)
	println("and MCash:", AndPayMethod(flag, 1))
	println("and MCash:", AndPayMethod(flag, 5))
	flag2 := MathPaymentMethodFlag([]int{2, 9})
	println("and MCash:", flag&flag2 == flag2)
}
