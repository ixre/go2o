package assets

import _ "embed"

//go:embed sensitive_dict.txt
var SensitiveDict []byte

//go:embed app.html
var AppDownHtml string
