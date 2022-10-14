package modules

import (
	"io"
	"log"
	"net/http"
)

// 检查最新的版本和是否需要更新

// 等待修改，即将废弃

func CheckNew(version string) {
	log.Print("正在连接 Github 检查更新...")

	// url
	url := "https://raw.githubusercontent.com/OfflineY/OfflineY/main/online-chat-version"

	resp, err := http.Get(url)
	if err != nil {
		log.Println("检查更新失败，无法连接 Github")
	} else {

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("检查更新失败，无法解析返回数据")
		}

		defer resp.Body.Close()
		data := string(bytes)

		if data != version {
			log.Printf("检查更新完成，有新的可更新版本：%s", data)
		} else {
			log.Printf("检查更新完成，已经是最新版本：%s", data)
		}
	}
}
