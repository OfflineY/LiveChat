package util

import (
	"fmt"
	"log"
	"os"

	"github.com/go-ini/ini"
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
	openFile, e := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if e != nil {
		fmt.Println(e)
	}
	str := `
[MongoDB]
applyUrl = "mongodb://127.0.0.1:27017"
maxPoolSize = 20
databaseName = "LiveChat"
	
[web]
port = ":5215"
address = 127.0.0.1:5215
	
[websocket]
port = ":8080"
address = 127.0.0.1:8080`

	_, err := openFile.WriteString(str)
	if err != nil {
		log.Println(err)
	}

	err = openFile.Close()
	if err != nil {
		log.Println(err)
	}
}
