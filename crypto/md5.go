package crypto

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encode(str string) (md string, err error) {
	h := md5.New()
	_, err = h.Write([]byte(str))
	md = hex.EncodeToString(h.Sum(nil))
	return
}
