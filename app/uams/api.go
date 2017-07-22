package uams

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"net/url"
	"sort"
)

// 接口响应
type Response struct {
	Result  int    `json:"result"`
	Data    string `json:"data"`
	Message string `json:"message"`
}

// 参数排序后，转换为字节,排除sign和sign_type
func paramsToBytes(r url.Values, token string) []byte {
	i := 0
	buf := bytes.NewBuffer(nil)
	// 键排序
	keys := []string{}
	for k, _ := range r {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// 拼接参数和值
	for _, k := range keys {
		if k == "sign" || k == "sign_type" {
			continue
		}
		if i > 0 {
			buf.WriteString("&")
		}
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(r[k][0])
		i++
	}
	buf.WriteString(token)
	return buf.Bytes()
}

// 签名
func Sign(signType string, r url.Values, token string) string {
	data := paramsToBytes(r, token)
	switch signType {
	case "md5":
		return md5Encode(data)
	case "sha1":
		return sha1Encode(data)
	}
	return ""
}

// MD5加密
func md5Encode(data []byte) string {
	m := md5.New()
	m.Write(data)
	dec := m.Sum(nil)
	return hex.EncodeToString(dec)
}

// SHA1加密
func sha1Encode(data []byte) string {
	s := sha1.New()
	s.Write(data)
	d := s.Sum(nil)
	return hex.EncodeToString(d)
}
