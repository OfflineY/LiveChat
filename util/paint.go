package util

import (
	"fmt"
	// "io"
	// "net/http"
	// v "github.com/hashicorp/go-version"
)

func Paint(version string) {
	fmt.Print("\n   __ _             ___ _           _   \n")
	fmt.Println("  / /(_)_   _____  / __\\ |__   __ _| |_ ")
	fmt.Println(" / / | \\ \\ / / _ \\/ /  | '_ \\ / _` | __|")
	fmt.Println("/ /__| |\\ V /  __/ /___| | | | (_| | |_ ")
	fmt.Println("\\____/_| \\_/ \\___\\____/|_| |_|\\__,_|\\__|")
	fmt.Printf("LiveChat %s\n\n", version)

	// TODO：新版本检查更新

	// url := "https://raw.githubusercontent.com/OfflineY/OfflineY/main/online-chat-version"

	// resp, err := http.Get(url)
	// if err != nil {
	// 	fmt.Println("Cannot find new version")
	// } else {
	// 	bytes, err := io.ReadAll(resp.Body)
	// 	if err != nil {
	// 		fmt.Println("Cannot find new version")
	// 	}

	// 	defer resp.Body.Close()
	// 	data := string(bytes)

	// 	v1, err := v.NewVersion("1.2")
	// 	v2, err := v.NewVersion("1.5+metadata")

	// 	// Comparison example. There is also GreaterThan, Equal, and just
	// 	// a simple Compare that returns an int allowing easy >=, <=, etc.
	// 	if v1.LessThan(v2) {
	// 		fmt.Printf("%s is less than %s", v1, v2)
	// 	}

	// 	if data != version {
	// 		fmt.Printf("Has new version: %s", data)
	// 	} else {
	// 		return
	// 	}
	// }
}
