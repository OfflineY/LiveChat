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
	s, err := ini.Load(setPath)

	if err != nil {
		// 初始化创建文件
		log.Print("未找到配置文件开始初始化\n")
		_, err := os.Create(setPath)
		if err != nil {
			log.Print("os:", err)
			GetOut()
		}

		f, err := os.OpenFile(setPath, os.O_RDWR, 0600)
		if err != nil {
			log.Print("os:", err)
			GetOut()
		}

		setFile := "# 设置服务模式\n"
		setFile += "mode = 1\n"
		setFile += "\n"
		setFile += "# 开启服务端模式则必须设置此内容\n"
		setFile += "# Server -> Url -> none or \"/\"+sth.\n"
		setFile += "# Server -> Server -> [ip]:port\n"
		setFile += "[Server]\n"
		setFile += "Url = \"/socket\"\n"
		setFile += "Server = \"192.168.1.2:52158\"\n"
		setFile += "WebPort = 8080\n"
		setFile += "\n"
		setFile += "# 客户端模式必须设置否则闪退\n"
		setFile += "[User]\n"
		setFile += "Url = \"ws://localhost:8080/socket\"\n"

		_, err = f.WriteString(setFile)
		if err != nil {
			log.Print("os:", err)
		}

		defer f.Close()

		log.Print("初始化结束，请更改配置文件，并重启应用\n")

		// 回车
		GetOut()

	} else {

		// 跳过初始化
		log.Printf("找到配置文件，跳过初始化过程\n")

		// 获取服务模式
		mode := s.Section("").Key("mode").String()

		switch {
		case mode == "0":
			// 服务端
			go ServeCommand()
			WorkPrimaryServer(s)
		case mode == "1":
			// 网络服务端
			go WebServer(s)
			WorkPrimaryServer(s)
		case mode == "2":
			// 客户端
			log.Print("已开启【客户端】模式，此模式是为使用者所准备\n")
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
