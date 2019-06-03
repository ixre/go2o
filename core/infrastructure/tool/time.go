/**
 * Copyright 2015 @ to2.net.
 * name : time.go
 * author : jarryliu
 * date : 2015-12-28 03:38
 * description :
 * history :
 */
package tool

import "time"

// 获取起始日期
func GetStartDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

// 获取昨天的起始时间
func GetYesterdayStartDate() time.Time {
	return GetStartDate(time.Now().Add(time.Hour * -24))
}

// 获取一天开始和结束的时间戳
func GetStartEndUnix(t time.Time) (int64, int64) {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local).Unix(),
		time.Date(y, m, d, 23, 59, 59, 999, time.Local).Unix()
}

var monthDays []int = []int{31, 0, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// 获取一个月的最后一天
func LastDay(year int, m time.Month) int {
	if m == 2 {
		isLeapYear := year%4 == 0 &&
			(year%100 != 0 || year%400 == 0)
		if isLeapYear {
			return 29
		}
		return 28
	}
	return monthDays[m-1]
}
