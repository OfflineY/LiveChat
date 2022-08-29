package modules

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

// 开启时引导程序
// 包括：初始化、判断配置文件、启动服务
func Bootstrap(Version string) {
	PrintLogo(Version)

	// 检查版本的更新
	// 启动一个新的 goroutine 检查，防止github连接失败一直卡在连接步骤
	go CheckNew(Version)

	// ini 的位置和名称
	setPath := "set.ini"

	// 检测配置文件的存在性
	s, Seterr := ini.Load(setPath)

	// 无 初始化，创建ini
	// 有 跳过直接开始
	if Seterr != nil {

		// 初始化创建文件
		log.Print("未找到配置文件开始初始化\n")
		_, err := os.Create(setPath)
		if err != nil {
			log.Printf("初始化失败(%s)\n", err)
			GetOut()
		}

		// 打开文件写入
		f, err := os.OpenFile(setPath, os.O_RDWR, 0600)
		if err != nil {
			log.Printf("初始化失败(%s)\n", err)
			GetOut()
		}

		// 初始化写入ini
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

		// 写入完成后之间关闭文件
		defer f.Close()

		// 省略（感觉log写的太多了）
		// log.Print("创建配置文件完成...\n")

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

		// 获取 ini 的 mode 数据，0为服务端，1为客户端
		// 外部内容 Section 空白即可
		mode := s.Section("").Key("mode").String()
		switch {
		case mode == "0":
			// mode 0 服务端 普通模式下的 server
			// 服务端指令控制
			go ServeCommand()
			WorkPrimaryServer(s)
		case mode == "1":
			// mode 1 服务端 web模式下的 server
			go WebServer()
			WorkPrimaryServer(s)
		case mode == "2":
			// mode 2 客户端
			// 文字在客户端处理完成后发回服务端返回给其他用户
			// 可以尽可能减少服务器负担
			log.Print("已开启【客户端】模式，此模式是为使用者所准备。\n")

			GetChat(s.Section("User").Key("Url").String())
		}
	}
}

// 打印 logo 和版本的信息
func PrintLogo(version string) {
	fmt.Println(`
    _____     _ _         _____ _       _   
   |     |___| |_|___ ___|     | |_ ___| |_     OnlineChat
   |  |  |   | | |   | -_|   --|   | .'|  _|    Version ` + version + `
   |_____|_|_|_|_|_|_|___|_____|_|_|__,|_|      Powered by BillyYuan

========================================================================
   `)
}
