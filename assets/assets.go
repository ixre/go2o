package assets

import _ "embed"

//go:embed sensitive_dict.txt
var SensitiveDict []byte

// 下载页面

//go:embed app.html
var AppDownHtml string

// 桥接页面,用于H5跳转
//
//go:embed bridge.html
var BridgeHtml string
