/**
 * Copyright 2014 @ to2.net.
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
	"go2o/core/service/rsi"
)

//func GetShopCheckboxs(mchId  int32, chks string) []byte {
//	shops := dps.ShopService.GetOnlineShops(mchId)
//	buf := bytes.NewBufferString("")
//
//	if len(chks) == 0 {
//		for i, k := range shops {
//			buf.WriteString(fmt.Sprintf(
//				`<input type="checkbox" value="%d" id="shop%d" field="ApplySubs[%d]" checked="checked"/>
//			 	<label for="shop%d">%s</label>`,
//				k.ID,
//				i,
//				i,
//				i,
//				k.Name,
//			))
//		}
//	} else {
//		chks = fmt.Sprintf(",%s,", chks)
//		for i, k := range shops {
//			if strings.Index(chks, fmt.Sprintf(",%d,", k.ID)) == -1 {
//				buf.WriteString(fmt.Sprintf(
//					`<input type="checkbox" value="%d" id="shop%d" field="ApplySubs[%d]"/>
//			 	<label for="shop%d">%s</label>`,
//					k.ID,
//					i,
//					i,
//					i,
//					k.Name,
//				))
//			} else {
//				buf.WriteString(fmt.Sprintf(
//					`<input type="checkbox" value="%d" id="shop%d" field="ApplySubs[%d]" checked="checked"/>
//			 	<label for="shop%d">%s</label>`,
//					k.ID,
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

func GetShopsJson(mchId int32) []byte {
	shops := rsi.MerchantService.GetShopsOfMerchant(mchId)
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

func GetShopDropList(mchId int32, selected int32) []byte {
	buf := bytes.NewBuffer([]byte{})
	shops := rsi.MerchantService.GetShopsOfMerchant(mchId)
	for _, v := range shops {
		if v.Id == selected {
			buf.WriteString(fmt.Sprintf(`<option value="%d" selected="selected">%s</option>`, v.Id, v.Name))
		} else {
			buf.WriteString(fmt.Sprintf(`<option value="%d">%s</option>`, v.Id, v.Name))
		}
	}
	return buf.Bytes()
}
