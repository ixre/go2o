/**
 * Copyright 2014 @ to2.net.
 * name :
 * author : jarryliu
 * date : 2014-01-12 21:02
 * description :
 * history :
 */

package lbs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

const (
	MAP_API        = "http://api.map.baidu.com/geocoder/v2/"
	MAP_API_SECRET = "ElGNZ2ihulmww5cnHfga4HET"
)

func reqApi(url string) ([]byte, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return []byte{}, err
	} else {
		return ioutil.ReadAll(rsp.Body)
	}
}

func GetLocation(address string) (lng, lat float64, err error) {
	req := fmt.Sprintf("%s?ak=%s&output=json&address=%s&city=",
		MAP_API, MAP_API_SECRET, address)
	d, err := reqApi(req)
	if err != nil {
		return 0, 0, err
	}

	var m map[string]interface{} = make(map[string]interface{})
	err = json.Unmarshal(d, &m)
	result := m["result"]
	if result == nil {
		return 0, 0, errors.New("unknown location")
	}
	m = result.(map[string]interface{})["location"].(map[string]interface{})
	return m["lng"].(float64), m["lat"].(float64), err
}

func rad(f float64) float64 {
	return f * math.Pi / 180.0
}

func GetLocDistance(lng1, lat1, lng2, lat2 float64) float64 {
	const EARTH_RADIUS = 6378.137
	var radLat1 float64 = rad(lat1)
	var radLat2 float64 = rad(lat2)
	a := radLat1 - radLat2
	b := rad(lng1) - rad(lng2)
	s := 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+
		math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(b/2), 2)))
	s = s * EARTH_RADIUS
	s = math.Floor(s*10000*1000+0.5) / 10000 //米
	return s
}

//func main() {
//	lng, lat, _ := getLocation("软件园二期望海路10号")
//	lng2, lat2, _ := getLocation("湖里万达广场")
//
//
//	fmt.Println(lng, lat)
//	fmt.Println(lng2, lat2)
//
//	fmt.Println("距离为:",GetLocDistance(lng,lat,lng2,lat2))
//}
