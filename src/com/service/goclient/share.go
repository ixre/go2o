package goclient

import (
	"com/ording/entity"
	"fmt"
	"ops/cf/net/jsv"
)

type shareClient struct {
	conn *jsv.TCPConn
}

func (this *shareClient) GetShoppingCart(cart string) (a *entity.ShoppingCart, err error) {
	var result jsv.Result
	err = this.conn.WriteAndDecode([]byte(fmt.Sprintf(
		`{"cart":"%s"}>>Share.GetShoppingCart`, cart)), &result, 1024)
	if err != nil {
		return nil, err
	}
	a = &entity.ShoppingCart{}
	err = jsv.UnmarshalMap(result.Data, &a)
	return a, err
}
