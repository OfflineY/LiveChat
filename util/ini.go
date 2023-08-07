package util

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

// IsIniExist ini配置文件是否存在
func IsIniExist(path string) bool {
	_, err := ini.Load(path)
	return err == nil
}

// IniLoad 加载ini文件
func IniLoad(path string) *ini.File {
	cfg, err := ini.Load(path)
	if err != nil {
		log.Panicln(err)
	}
	return cfg
}

// IniInitial 初始化ini文件
func IniInitial(path string) {
	_, e := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if e != nil {
		fmt.Println(e)
	}
	//cfg := IniLoad(path)
	//_, err := cfg.NewSection("new section")
	//if err != nil {
	//	log.Println(err)
	//}
	//_, err = cfg.Section("").NewKey("name", "value")
	//if err != nil {
	//	log.Println(err)
	//}
}
