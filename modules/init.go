package modules

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

// 引导
func Bootstrap(Version string) {
	PrintLogo(Version)

	// 检查更新
	go CheckNew(Version)

	setPath := "set.ini"
	s, Seterr := ini.Load(setPath)

	if Seterr != nil {

		// 初始化创建文件
		log.Print("未找到配置文件开始初始化\n")
		_, err := os.Create(setPath)
		if err != nil {
			log.Printf("初始化失败(%s)\n", err)
			GetOut()
		}

		f, err := os.OpenFile(setPath, os.O_RDWR, 0600)
		if err != nil {
			log.Printf("初始化失败(%s)\n", err)
			GetOut()
		}

		f.WriteString(`# mode 为设置服务模式, 0 为服务端, 1 为客户端
mode = 0

# 开启服务端模式则必须设置此内容，客户端情况下删去保留没有影响
[Server]
# Url可以当作密码或者标签
Url = "/socket"
Port = "8080"
# 待补充配置
# 客户端模式必须设置否则闪退，举例ws://【ServerIp】:【Port】/【Url】
[User]
Url = "ws://localhost:8080/socket"
# 待补充配置`)

		defer f.Close()

		log.Print("初始化结束，请更改配置文件，并重启应用\n")

		// 注意事项
		fmt.Print(`
		


【注意事项】
		
  + 请在 cmd 或 vscode 终端或一些常驻终端使用，否则可能会出现出现错误闪退的情况。

  + 目前仅处于开发中且项目未完全完工，如有错误请提出。

  + 虽已经实现效果但安全性有待测试，仅供学习，请勿用于其他用途。



`)
		GetOut()
	} else {

		// 跳过初始化
		log.Printf("找到配置文件，跳过初始化过程\n")

		mode := s.Section("").Key("mode").String()
		switch {
		case mode == "0":
			// 服务端
			go ServeCommand()
			WorkPrimaryServer(s)
		case mode == "1":
			// 网络服务端
			go WebServer()
			WorkPrimaryServer(s)
		case mode == "2":
			// 客户端
			log.Print("已开启【客户端】模式，此模式是为使用者所准备。\n")
			GetChat(s.Section("User").Key("Url").String())
		}
	}
}

// 打印
func PrintLogo(version string) {
	fmt.Println(`
    _____     _ _         _____ _       _   
   |     |___| |_|___ ___|     | |_ ___| |_     OnlineChat
   |  |  |   | | |   | -_|   --|   | .'|  _|    Version ` + version + `
   |_____|_|_|_|_|_|_|___|_____|_|_|__,|_|      Powered by BillyYuan

========================================================================
   `)
}
