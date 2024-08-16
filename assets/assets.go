package assets

import _ "embed"

//go:embed sensitive_dict.txt
var SensitiveDict []byte
