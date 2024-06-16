/**
 * Copyright 2014 @ 56x.net.
 * name : delivery
 * author : jarryliu
 * date : 2014-10-06 14:21 :
 * description :
 * history :
 */
package domain

import (
	"errors"
	"regexp"
)

var (
	areaRegexp  = regexp.MustCompile("(市)((.+)(区|县))")
	errNotMatch = errors.New("未识别的地址")
	cityRegexp  = regexp.MustCompile("(省|自治区|行政区)((.+)市)")
)

// 获取地区名称
func GetAreaName(addr string) (string, error) {
	var matches = areaRegexp.FindAllStringSubmatch(addr, -1)
	if len(matches) == 0 {
		return "", errNotMatch
	}
	return matches[0][2], nil
}

// 获取城市名称
func GetCityName(addr string) (string, error) {
	var matches = cityRegexp.FindAllStringSubmatch(addr, -1)
	if len(matches) == 0 {
		return "", errNotMatch
	}
	return matches[0][2], nil
}

// 去除省市区中的直辖区和县
func TrimAreaNames(arr []string) []string {
	if len(arr) >= 3 {
		if arr[1] == "市辖区" || arr[1] == "市辖县" || arr[1] == "县" {
			return []string{arr[0], arr[2]}
		}
	}
	return arr
}
