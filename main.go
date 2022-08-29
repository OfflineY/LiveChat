package main

import (
	// ./modules
	"OnlineChat/modules"
)

// 设置全局 version 变量
var Version string = "1.0"

func main() {
	// 初始化引导程序
	modules.Bootstrap(Version)
}
