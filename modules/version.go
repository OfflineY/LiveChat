package modules

// 检测版本的更新

import (
	"io"
	"log"
	"net/http"
)

// 检查最新的版本和是否需要更新
func CheckNew(version string) {
	log.Print("正在连接 Github 检查更新...")

	// 用 gist 这个来记录更新
	url := "https://raw.githubusercontent.com/OfflineY/OfflineY/main/online-chat-version"

	// 将要读取的网站放入到 get 方法中
	resp, err := http.Get(url)
	if err != nil {
		log.Println("检查更新失败，无法连接 Github")
	} else {
		// 读取 body 里面的内容
		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("检查更新失败，无法解析返回数据")
		}
		// 关闭链接
		defer resp.Body.Close()
		data := string(bytes)
		// 判断新的版本提示更新or不更新
		if data != version {
			log.Printf("检查更新完成，已经是最新版本：%s", data)
		} else {
			log.Printf("检查更新完成，有新的可更新版本：%s", data)
		}
	}
}
