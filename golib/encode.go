package golib

import "crypto/md5"

func MD5Sum(source string) (hash string) {
	h := md5.New()
	return string(h.Sum([]byte(source)))
}
