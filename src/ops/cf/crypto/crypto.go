package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(b []byte) string {
	cry := md5.New()
	cry.Write(b)
	return hex.EncodeToString(cry.Sum(nil))
}
