package tencent

// 文档
//
// 	腾讯地图校验签名: https://lbs.qq.com/faq/serverFaq/webServiceKey
//
//  IP定位: https://lbs.qq.com/service/webService/webServiceGuide/webServiceIp
//
//  注: 开通后需要对应用进行调用配额

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ixre/go2o/core/domain/interface/registry"
	"github.com/ixre/go2o/core/infrastructure/lbs"
	"github.com/ixre/gof/crypto"
	api "github.com/ixre/gof/ext/api"
)

var _ lbs.LbsProvider = (*tecentLbsService)(nil)

// IPLocationResponse 是腾讯地图IP定位API的响应结构体

// "{\"status\":0,\"message\":\"Success\",\"request_id\":\"fdef8519884140398f05aecbf939fe9e\",
// \"result\":{\"ip\":\"112.93.63.206\",\"location\":{\"lat\":22.803751,\"lng\":113.293719},
// \"ad_info\":{\"nation\":\"中国\",\"province\":\"广东省\",\"city\":\"佛山市\",\"district\":\"顺德区\",
// \"adcode\":440606,\"nation_code\":156}}}\n"
type IPLocationResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Result  struct {
		IP       string `json:"ip"`
		Location struct {
			Lat float64 `json:"lat"`
			Lng float64 `json:"lng"`
		} `json:"location"`
		AdInfo AdInfo `json:"ad_info"`
	} `json:"result"`
}

type AdInfo struct {
	NationCode int `json:"nation_code"`
	AdCode     int    `json:"adcode"`
	Province   string `json:"province"`
	City       string `json:"city"`
	District   string `json:"district"`
}

type Area struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type tecentLbsService struct {
	appKey    string
	appSecret string
}

func NewLbsService(repo registry.IRegistryRepo) lbs.LbsProvider {
	initMapConfig(repo)
	appKey, _ := repo.GetValue("tencent-map.appkey")
	appSecret, _ := repo.GetValue("tencent-map.secret")
	// if len(appKey) == 0 || len(appSecret) == 0 {
	// 	logger.Warn("为配置腾讯位置定位服务", "appKey or appSecret is empty")
	// }
	return &tecentLbsService{
		appKey:    appKey,
		appSecret: appSecret,
	}
}

func initMapConfig(repo registry.IRegistryRepo) {
	repo.CreateUserKey("tencent-map.appkey", "-", "腾讯地图AppKey")
	repo.CreateUserKey("tencent-map.secret", "-", "腾讯地图密钥")
}

// 生成签名
func (t *tecentLbsService) getSn(path string, params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Add(k, v)
	}
	query := api.GetSortParams(values)
	sn := crypto.Md5([]byte(path + "?" + query + t.appSecret))
	return strings.ToLower(sn)
}

// QueryLocation implements lbs.LbsProvider.
func (t *tecentLbsService) QueryLocation(ip string) (*lbs.LocationInfo, error) {
	params := map[string]string{
		"ip":  ip,
		"key": t.appKey,
	}
	sig := t.getSn("/ws/location/v1/ip", params)
	// 腾讯地图IP定位API URL
	url := fmt.Sprintf("https://apis.map.qq.com/ws/location/v1/ip?ip=%s&key=%s&sig=%s", ip, t.appKey, sig)

	// 发送GET请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析JSON响应
	var result IPLocationResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return &lbs.LocationInfo{
		Lat:      result.Result.Location.Lat,
		Lng:      result.Result.Location.Lng,
		Location: result.Result.AdInfo.Province + result.Result.AdInfo.City + result.Result.AdInfo.District,
		City:     result.Result.AdInfo.City,
		PlaceCode: result.Result.AdInfo.AdCode,
	}, nil
}
