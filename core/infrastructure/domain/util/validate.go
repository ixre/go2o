/**
 * Copyright 2015 @ to2.net.
 * name : validate
 * author : jarryliu
 * date : 2016-07-23 15:42
 * description :
 * history :
 */
package util

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// JS 验证

/*var aCity={11:"北京",12:"天津",13:"河北",14:"山西",15:"内蒙古",21:"辽宁",22:"吉林",23:"黑龙江",31:"上海",32:"江苏",33:"浙江",34:"安徽",
35:"福建",36:"江西",37:"山东",41:"河南",42:"湖北",43:"湖南",44:"广东",45:"广西",46:"海南",50:"重庆",51:"四川",52:"贵州",53:"云南",54:"西藏",61:"陕西",62:"甘肃",63:"青海",64:"宁夏",65:"新疆",71:"台湾",81:"香港",82:"澳门",91:"国外"}
*/
//function isCardID(sId){
// var iSum=0 ;
// var info="" ;
// if(!/^\d{17}(\d|x)$/i.test(sId)) return "&nbsp;";
// sId=sId.replace(/x$/i,"a");
// if(aCity[parseInt(sId.substr(0,2))]==null) return "&nbsp;";
// sBirthDay=sId.substr(6,4)+"-"+Number(sId.substr(10,2))+"-"+Number(sId.substr(12,2));
// var d=new Date(sBirthDay.replace(/-/g,"/")) ;
// if(sBirthDay!=(d.getFullYear()+"-"+ (d.getMonth()+1) + "-" + d.getDate()))return "&nbsp;";
// for(var i = 17;i>=0;i --) iSum += (Math.pow(2,i) % 11) * parseInt(sId.charAt(17 - i),11) ;
// if(iSum%11!=1) return "&nbsp;";
// return true;//aCity[parseInt(sId.substr(0,2))]+","+sBirthDay+","+ (sId.substr(16,1)%2?"男":"女")
//}

var (
	cityCodeMap = map[int]string{
		11: "北京", 12: "天津", 13: "河北", 14: "山西", 15: "内蒙古",
		21: "辽宁", 22: "吉林", 23: "黑龙江", 31: "上海", 32: "江苏",
		33: "浙江", 34: "安徽", 35: "福建", 36: "江西", 37: "山东",
		41: "河南", 42: "湖北", 43: "湖南", 44: "广东", 45: "广西",
		46: "海南", 50: "重庆", 51: "四川", 52: "贵州", 53: "云南",
		54: "西藏", 61: "陕西", 62: "甘肃", 63: "青海", 64: "宁夏",
		65: "新疆", 71: "台湾", 81: "香港", 82: "澳门", 91: "国外"}
	cardLenRegexp = regexp.MustCompile(`^\d{17}(\d|X)$`)
)

// 检查中国公民身份证号码是否正确
func CheckChineseCardID(sId string) error {
	if !cardLenRegexp.MatchString(sId) {
		return errors.New("身份证长度不正确")
	}
	// 以X结尾的身份证号
	sId = strings.Replace(sId, "X", "a", 1)
	cityCode, _ := strconv.Atoi(sId[:2])
	if _, exists := cityCodeMap[cityCode]; !exists {
		return errors.New("未知的身份证信息地区")
	}
	birthDate := strings.Join([]string{sId[6:10],
		sId[10:12], sId[12:14]}, "-")
	if _, err := time.ParseInLocation("2006-01-02", birthDate, time.Local); err != nil {
		return errors.New("身份证出生日期错误")
	}
	//for (var i = 17; i>=0; i --) iSum += (Math.pow(2, i) % 11) * parseInt(sId.charAt(17 - i), 11);
	////aCity[parseInt(sId.substr(0,2))]+","+sBirthDay+","+ (sId.substr(16,1)%2?"男":"女")
	var iSum int64 = 0
	for i := 17; i >= 0; i-- {
		b, err := strconv.ParseInt(sId[17-i:17-i+1], 11, 32)
		if err != nil {
			b = 0
		}
		iSum += int64(math.Pow(2, float64(i))) % 11 * b
	}
	if iSum%11 != 1 {
		return errors.New("身份证号码错误")
	}
	return nil
}
