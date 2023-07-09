package util

import (
	"github.com/go-ini/ini"
	"log"
)

// IsIniExist ini配置文件是否存在
func IsIniExist(path string) bool {
	_, err := ini.Load(path)
	if err != nil {
		return false
	}

	return true
}

// IniLoad 加载ini文件
func IniLoad(path string) *ini.File {
	cfg, err := ini.Load(path)
	if err != nil {
		log.Panicln(err)
	}

	return cfg
}
