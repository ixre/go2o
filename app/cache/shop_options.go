/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package cache

import (
	"bytes"
	"fmt"
	"go2o/core/service/dps"
)

//func GetShopCheckboxs(mchId int, chks string) []byte {
//	shops := dps.MerchantService.GetOnlineShops(mchId)
//	buf := bytes.NewBufferString("")
//
//	if len(chks) == 0 {
//		for i, k := range shops {
//			buf.WriteString(fmt.Sprintf(
//				`<input type="checkbox" value="%d" id="shop%d" field="ApplySubs[%d]" checked="checked"/>
//			 	<label for="shop%d">%s</label>`,
//				k.Id,
//				i,
//				i,
//				i,
//				k.Name,
//			))
//		}
//	} else {
//		chks = fmt.Sprintf(",%s,", chks)
//		for i, k := range shops {
//			if strings.Index(chks, fmt.Sprintf(",%d,", k.Id)) == -1 {
//				buf.WriteString(fmt.Sprintf(
//					`<input type="checkbox" value="%d" id="shop%d" field="ApplySubs[%d]"/>
//			 	<label for="shop%d">%s</label>`,
//					k.Id,
//					i,
//					i,
//					i,
//					k.Name,
//				))
//			} else {
//				buf.WriteString(fmt.Sprintf(
//					`<input type="checkbox" value="%d" id="shop%d" field="ApplySubs[%d]" checked="checked"/>
//			 	<label for="shop%d">%s</label>`,
//					k.Id,
//					i,
//					i,
//					i,
//					k.Name,
//				))
//			}
//		}
//	}
//	return buf.Bytes()
//}

func GetShopsJson(merchantId int) []byte {
	shops := dps.MerchantService.GetShopsOfMerchant(merchantId)
	buf := bytes.NewBufferString("[")
	for i, v := range shops {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(fmt.Sprintf(`{"id":%d,"name":"%s"}`, v.Id, v.Name))
	}
	buf.WriteString("]")
	return buf.Bytes()
}

func GetShopDropList(merchantId int, selected int) []byte {
	buf := bytes.NewBuffer([]byte{})
	shops := dps.MerchantService.GetShopsOfMerchant(merchantId)
	for _, v := range shops {
		if v.Id == selected {
			buf.WriteString(fmt.Sprintf(`<option value="%d" selected="selected">%s</option>`, v.Id, v.Name))
		} else {
			buf.WriteString(fmt.Sprintf(`<option value="%d">%s</option>`, v.Id, v.Name))
		}
	}
	return buf.Bytes()
}
