package modules

import (
	"fmt"
	"os"
)

func GetOut() {
	fmt.Printf("\n\n\n\n按回车键退出...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
