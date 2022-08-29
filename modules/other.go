package modules

import (
	"fmt"
	"os"
)

// 实现按任意键退出，但好像只有回车管用哈哈哈哈
func GetOut() {
	fmt.Printf("\n\n\n\n按任意键退出...")
	b := make([]byte, 1)
	os.Stdin.Read(b)
}
