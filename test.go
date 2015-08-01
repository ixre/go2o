/**
 * Copyright 2015 @ S1N1 Team.
 * name : test.go
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package main

import (
	"fmt"
	"github.com/atnet/gof/crypto"
	"github.com/atnet/gof/util"
	"go2o/src/core/infrastructure/domain"
	"net/url"
	"regexp"
)

func main() {
	fmt.Println(crypto.Sha1([]byte("123456")))
	fmt.Println(domain.PartnerSha1Pwd("zs", "zs#12345"))

	fmt.Println(domain.NewApiId(105))
	fmt.Println(domain.MemberSha1Pwd("123456"))
	fmt.Println(1 << 2)
	fmt.Println(util.IsMobileAgent("Mozilla/5.0 (Linux; U; Android 2.3.7; en-us; Nexus One Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1"))

	tagRegexp := regexp.MustCompile("\\s*([^:\\|]+):([^:\\|]+)\\s*\\|*")
	matches := tagRegexp.FindAllStringSubmatch("G1:10|G2:12|G3:14", -1)
	for i := 0; i < len(matches); i++ {
		fmt.Println(i, "---", matches[i][1], matches[i][2])
	}

	fmt.Println("/main/change_device?device=3&return_url=/main/t/123" +
		url.QueryEscape("?return_url=/app"))
}
