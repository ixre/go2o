// 合作商户的接口
package server

import (
	"bytes"
	"com/domain/interface/member"
	"com/ording"
	"com/ording/dao"
	"com/ording/dproxy"
	"encoding/json"
	"ops/cf/net/jsv"
	"strconv"
	"time"
)

type Partner struct{}

//var partnerId int
//	var partner *entity.Partner
//	var err error
//	_, err, partner = VerifyPartner((*m)["partner_id"], (*m)["secret"])
//	if err != nil {
//		r.Result = false
//		r.Code = jsv.C_PERMISSION_DENIED
//		r.Message = err.Error()
//		return r
//	}
func (this *Partner) GetPartner(m *jsv.Args, r *jsv.Result) error {
	_, err, e := VerifyPartner(m)
	if err != nil {
		return err
	}
	r.Result = true
	r.Data = e
	return nil
}

func (this *Partner) GetSiteConf(m *jsv.Args, r *jsv.Result) error {
	partnerId, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}

	siteConf :=dproxy.PartnerService.GetSiteConf(partnerId)
	r.Result = true
	r.Data = siteConf
	return nil
}

func (this *Partner) GetHost(m *jsv.Args, r *jsv.Result) error {
	partnerId, err := strconv.Atoi((*m)["partner_id"].(string))
	if err != nil {
		return err
	}

	host := dao.Partner().GetHostById(partnerId)
	r.Data = host
	r.Result = true
	return nil
}

func (this *Partner) GetShops(m *jsv.Args, r *jsv.Result) error {
	partnerId, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}

	shops := dao.Shop().GetShopsOfPartner(partnerId)
	r.Result = true
	r.Data = shops
	return nil
}

func (this *Partner) Category(m *jsv.Args, r *jsv.Result) error {
	partnerId, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}

	categories := dao.Category().GetCategoriesOfPartner(partnerId)
	r.Result = true
	r.Data = categories
	return nil
}

func (this *Partner) GetItems(m *jsv.Args, r *jsv.Result) error {
	partnerId, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}
	cid, _ := strconv.Atoi((*m)["cid"].(string))
	num, _ := strconv.Atoi((*m)["num"].(string))

	items := dao.Item().GetItemsByCid(partnerId, cid, num)
	r.Result = true
	r.Data = items

	return nil
}

func (this *Partner) RegistMember(m *jsv.Args, r *jsv.Result) error {

	var err error

	e := member.ValueMember{}

	if err = jsv.UnmarshalMap((*m)["json"], &e); err != nil {
		return err
	}
	var cardId string
	var tgid int
	var partnerId int

	cardId = (*m)["card_id"].(string)
	tgid, _ = strconv.Atoi((*m)["tg_id"].(string))
	partnerId, _ = strconv.Atoi((*m)["partner_id"].(string))

	//如果卡片ID为空时，自动生成
	if cardId == "" {
		cardId = time.Now().Format("200601021504")
	}

	e.Pwd = ording.EncodeMemberPwd(e.Usr, e.Pwd)
	id, err := dproxy.MemberService.SaveMember(&e)

	if err == nil {
		dproxy.MemberService.SaveRelation(id, cardId, tgid, partnerId)
		r.Result = true
		return nil
	}

	return err
}

func (this *Partner) BuildOrder(m *jsv.Args, r *jsv.Result) error {
	partnerId, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}

	memberId, err := strconv.Atoi((*m)["member_id"].(string))
	cartData := (*m)["cart"].(string)
	couponCode := (*m)["coupon_code"].(string)
	if err != nil {
		return err
	}

	order, err := dproxy.SpService.BuildOrder(partnerId,
		memberId, cartData, couponCode)
	if err != nil {
		return err
	}

	v := order.GetValue()
	buf := bytes.NewBufferString("")

	for _, v := range order.GetCoupons() {
		buf.WriteString(v.GetDescribe())
		buf.WriteString("\n")
	}

	var data map[string]interface{}
	data = make(map[string]interface{})
	if couponCode != "" {
		if v.CouponFee == 0 {
			data["result"] = v.CouponFee != 0
			data["message"] = "优惠券无效"
		} else {
			// 成功应用优惠券
			data["totalFee"] = v.TotalFee
			data["fee"] = v.Fee
			data["payFee"] = v.PayFee
			data["discountFee"] = v.DiscountFee
			data["couponFee"] = v.CouponFee
			data["couponDescribe"] = buf.String()
		}
	} else {
		//　取消优惠券
		data["totalFee"] = v.TotalFee
		data["fee"] = v.Fee
		data["payFee"] = v.PayFee
		data["discountFee"] = v.DiscountFee
	}

	js, _ := json.Marshal(data)

	r.Result = true
	r.Data = string(js)
	return nil
}

// 需要传递配送地址
func (this *Partner) SubmitOrder(m *jsv.Args, r *jsv.Result) error {
	partnerId, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}
	memberId, _ := strconv.Atoi((*m)["member_id"].(string))
	shopId, _ := strconv.Atoi((*m)["shop_id"].(string))
	pay_method, _ := strconv.Atoi((*m)["pay_method"].(string))
	deliverAddrId, _ := strconv.Atoi((*m)["addr_id"].(string))
	cart := (*m)["cart"].(string)
	couponCode := (*m)["coupon_code"].(string)
	note := (*m)["note"].(string)

	orderNo, err := dproxy.SpService.SubmitOrder(
		partnerId, memberId, shopId, pay_method,
		deliverAddrId, cart, couponCode, note)
	if err != nil {
		return err
	} else {
		r.Result = true
		r.Data = orderNo
	}
	return nil
}

func (this *Partner) GetOrderByNo(m *jsv.Args, r *jsv.Result) error {
	partnerId, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}
	order := dproxy.SpService.GetOrderByNo(partnerId,
		(*m)["order_no"].(string))
	if order != nil {
		r.Result = true
		r.Data = *order
	}
	return nil
}

func (this *Partner) CheckUsrExist(m *jsv.Args, r *jsv.Result) error {
	_, err, _ := VerifyPartner(m)
	if err != nil {
		return err
	}
	r.Result = true
	r.Data = dproxy.MemberService.CheckUsrExist((*m)["usr"].(string))
	return nil
}
