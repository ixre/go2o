package math

import (
	"math"
)

// 四舍五入计算,通过设置pos来指定要精度位数
// 如果小于０，则四舍五入到整数为
func Round(val float64, pos int) float64 {
	if pos <= 0 {
		if val < 0 {
			return math.Ceil(val - 0.5)
		}
		return math.Floor(val + 0.5)
	}
	digit := math.Pow10(pos)
	if val < 0 {
		return math.Ceil(val*digit-0.5) / digit
	}
	return math.Floor(val*digit+0.5) / digit
}

func Round32(val float32, pos int) float32 {
	return float32(Round(float64(val), pos))
}
