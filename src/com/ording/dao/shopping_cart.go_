package dao

import (
	"bytes"
	"com/ording/entity"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	//data example : 16*1|12*2|80
	cartRegex = regexp.MustCompile("(\\d+)\\s*\\*\\s*(\\d+)")
)

type shoppingCart struct {
	//购物车商品
	Items []entity.FoodItem
	//购物车商品数量
	Quantities map[int]int
	//客户端计算的金额
	ClientFee float32

	summary string
}

// 将字符表示转换为购物车
func ParseShoppingCart(s string) (cart *shoppingCart, err error) {
	if !cartRegex.MatchString(s) {
		return nil, errors.New("cart string is error.example:16*1|12*2")
	}
	cart = new(shoppingCart)
	matches := cartRegex.FindAllStringSubmatch(s, -1)

	length := len(matches)
	var ids []string = make([]string, length) //ID数组
	cart.Quantities = make(map[int]int, length)
	//cart.Items = []entity.FoodItem{}

	var id int
	var qua int

	for i, v := range matches {
		id, err = strconv.Atoi(v[1])
		if err != nil {
			continue
		}
		qua, err = strconv.Atoi(v[2])
		if err != nil {
			continue
		}

		ids[i] = strconv.Itoa(id)
		cart.Quantities[id] = qua
	}

	//todo:改成database/sql方式，不使用orm
	err = context.Db().GetOrm().SelectByQuery(&cart.Items, entity.FoodItem{},
		`SELECT * FROM it_item WHERE id IN(`+strings.Join(ids, ",")+`)`)

	if err != nil {
		return nil, err
	}

	return cart, nil
}

// 获取总览信息
func (this *shoppingCart) GetSummary() string {
	if len(this.summary) != 0 {
		return this.summary
	}
	buf := bytes.NewBufferString("")
	length := len(this.Items)
	for i, v := range this.Items {
		buf.WriteString(v.Name)
		if len(v.Note) != 0 {
			buf.WriteString("(" + v.Note + ")")
		}
		buf.WriteString("*" + strconv.Itoa(this.Quantities[v.Id]))
		if i < length-1 {
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

// 获取订单金额,返回totalFee为总额，
// orderFee为实际订单的金额(扣去促销优惠等后的金额)
func (this *shoppingCart) GetFee() (totalFee float32, orderFee float32) {
	var qua float32
	for _, v := range this.Items {
		qua = float32(this.Quantities[v.Id])
		totalFee = totalFee + v.Price*qua
		orderFee = orderFee + v.SalePrice*qua
	}
	return totalFee, orderFee
}
