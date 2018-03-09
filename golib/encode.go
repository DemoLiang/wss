package golib

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5Sum(source string) (hash string) {
	h := md5.New()
	h.Write([]byte(source))
	return hex.EncodeToString(h.Sum(nil))
}
