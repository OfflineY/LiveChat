package util

import (
	"crypto/md5"
	"fmt"
)

// MD5 对str进行md5加密
func MD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}
