/**
 * Copyright 2015 @ S1N1 Team.
 * name : alipay.go
 * author : jarryliu
 * date : 2015-07-27 18:08
 * description : https://github.com/ascoders/alipay
 * history :
 */
package alipay

import (
	_ "net/url"
	_ "regexp"
)

/* 生成支付宝即时到帐的表单参数
 * @params string 订单唯一id
 * @params int 价格
 * @params int 获得代金券的数量
 * @params string 充值账户的名称
 * @params string 充值描述
 */

/* 被动接收支付宝同步跳转的页面 */
func AlipayReturn() (int, string, string, string) {
	//列举全部传参
	type Params struct {
		Body        string `form:"body" json:"body"`                 //描述
		BuyerEmail  string `form:"buyer_email" json:"buyer_email"`   //买家账号
		BuyerId     string `form:"buyer_id" json:"buyer_id"`         //买家ID
		Exterface   string `form:"exterface" json:"exterface"`       //接口名称
		IsSuccess   string `form:"is_success" json:"is_success"`     //交易是否成功
		NotifyId    string `form:"notify_id" json:"notify_id"`       //通知校验id
		NotifyTime  string `form:"notify_time" json:"notify_time"`   //校验时间
		NotifyType  string `form:"notify_type" json:"notify_type"`   //校验类型
		OutTradeNo  string `form:"out_trade_no" json:"out_trade_no"` //在网站中唯一id
		PaymentType uint8  `form:"payment_type" json:"payment_type"` //支付类型
		SellerEmail string `form:"seller_email" json:"seller_email"` //卖家账号
		SellerId    string `form:"seller_id" json:"seller_id"`       //卖家id
		Subject     string `form:"subject" json:"subject"`           //商品名称
		TotalFee    string `form:"total_fee" json:"total_fee"`       //总价
		TradeNo     string `form:"trade_no" json:"trade_no"`         //支付宝交易号
		TradeStatus string `form:"trade_status" json:"trade_status"` //交易状态 TRADE_FINISHED或TRADE_SUCCESS表示交易成功
		Sign        string `form:"sign" json:"sign"`                 //签名
		SignType    string `form:"sign_type" json:"sign_type"`       //签名类型
	}
	/*
		//实例化参数
		param := &Params{}


		//解析表单内容，失败返回错误代码-3
		if err := contro.ParseForm(param); err != nil {
			return -3, "", "", ""
		}
		//如果最基本的网站交易号为空，返回错误代码-1
		if param.OutTradeNo == "" { //不存在交易号
			return -1, "", "", ""
		} else {
			//生成签名
			sign := alipaySign(param)

			//对比签名是否相同
			if sign == param.Sign { //只有相同才说明该订单成功了
				//判断订单是否已完成
				if param.TradeStatus == "TRADE_FINISHED" || param.TradeStatus == "TRADE_SUCCESS" { //交易成功
					return 1, param.OutTradeNo, param.BuyerEmail, param.TradeNo
				} else { //交易未完成，返回错误代码-4
					return -4, "", "", ""
				}
			} else { //签名认证失败，返回错误代码-2
				return -2, "", "", ""
			}
		}
	*/
	//位置错误类型-5
	return -5, "", "", ""
}

/* 被动接收支付宝异步通知 */
func AlipayNotify() (int, string, string, string) { /*
		//从body里读取参数，用&切割
		postArray := strings.Split(string(contro.Ctx.Input.CopyBody()), "&")

		//实例化url
		urls := &url.Values{}

		//保存传参的sign
		var paramSign string
		var sign string


		//如果字符串中包含sec_id说明是手机端的异步通知
		if strings.Index(string(contro.Ctx.Input.CopyBody()), `alipay.wap.trade.create.direct`) == -1 { //快捷支付
			for _, v := range postArray {
				detail := strings.Split(v, "=")

				//使用=切割字符串 去除sign和sign_type
				if detail[0] == "sign" || detail[0] == "sign_type" {
					if detail[0] == "sign" {
						paramSign = detail[1]
					}
					continue
				} else {
					urls.Add(detail[0], detail[1])
				}
			}

			// url解码
			urlDecode, _ := url.QueryUnescape(urls.Encode())
			sign, _ = url.QueryUnescape(urlDecode)
		} else { // 手机网页支付
			// 手机字符串加密顺序
			mobileOrder := []string{"service", "v", "sec_id", "notify_data"}
			for _, v := range mobileOrder {
				for _, value := range postArray {
					detail := strings.Split(value, "=")
					// 保存sign
					if detail[0] == "sign" {
						paramSign = detail[1]
					} else {
						// 如果满足当前v
						if detail[0] == v {
							if sign == "" {
								sign = detail[0] + "=" + detail[1]
							} else {
								sign += "&" + detail[0] + "=" + detail[1]
							}
						}
					}
				}
			}
			sign, _ = url.QueryUnescape(sign)

			//获取<trade_status></trade_status>之间的request_token
			re, _ := regexp.Compile("\\<trade_status[\\S\\s]+?\\</trade_status>")
			rt := re.FindAllString(sign, 1)
			trade_status := strings.Replace(rt[0], "<trade_status>", "", -1)
			trade_status = strings.Replace(trade_status, "</trade_status>", "", -1)
			urls.Add("trade_status", trade_status)

			//获取<out_trade_no></out_trade_no>之间的request_token
			re, _ = regexp.Compile("\\<out_trade_no[\\S\\s]+?\\</out_trade_no>")
			rt = re.FindAllString(sign, 1)
			out_trade_no := strings.Replace(rt[0], "<out_trade_no>", "", -1)
			out_trade_no = strings.Replace(out_trade_no, "</out_trade_no>", "", -1)
			urls.Add("out_trade_no", out_trade_no)

			//获取<buyer_email></buyer_email>之间的request_token
			re, _ = regexp.Compile("\\<buyer_email[\\S\\s]+?\\</buyer_email>")
			rt = re.FindAllString(sign, 1)
			buyer_email := strings.Replace(rt[0], "<buyer_email>", "", -1)
			buyer_email = strings.Replace(buyer_email, "</buyer_email>", "", -1)
			urls.Add("buyer_email", buyer_email)

			//获取<trade_no></trade_no>之间的request_token
			re, _ = regexp.Compile("\\<trade_no[\\S\\s]+?\\</trade_no>")
			rt = re.FindAllString(sign, 1)
			trade_no := strings.Replace(rt[0], "<trade_no>", "", -1)
			trade_no = strings.Replace(trade_no, "</trade_no>", "", -1)
			urls.Add("trade_no", trade_no)
		}
		//追加密钥
		sign += AlipayKey

		//md5加密
		m := md5.New()
		m.Write([]byte(sign))
		sign = hex.EncodeToString(m.Sum(nil))
		if paramSign == sign { //传进的签名等于计算出的签名，说明请求合法
			//判断订单是否已完成
			if urls.Get("trade_status") == "TRADE_FINISHED" || urls.Get("trade_status") == "TRADE_SUCCESS" { //交易成功
				contro.Ctx.WriteString("success")
				return 1, urls.Get("out_trade_no"), urls.Get("buyer_email"), urls.Get("trade_no")
			} else {
				contro.Ctx.WriteString("error")
			}
		} else {
			contro.Ctx.WriteString("error")
			//签名不符，错误代码-1
			return -1, "", "", ""
		}*/
	//未知错误-2
	return -2, "", "", ""
}
