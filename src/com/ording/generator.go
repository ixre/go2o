package ording

//todo:remove
import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const (
	minRand int = 1000000
	maxRand int = 9999999
)

//新订单号
func NewOrderNo(partnerId int) string {
	//PartnerId的首位和末尾再加7位随机数
	rand.Seed(time.Now().Unix())
	rd := minRand + rand.Intn(maxRand-minRand) //minRand - maxRand中间的随机数
	ptstr := strconv.Itoa(partnerId)
	return fmt.Sprintf("%s%s%d", ptstr[:1], ptstr[len(ptstr)-1:], rd)
}
