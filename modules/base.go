package modules

import (
	"fmt"
	"log"
	"os"
)

func GetOut() {
	fmt.Printf("\n\n\n\n按回车键退出...")
	b := make([]byte, 1)
	_, err := os.Stdin.Read(b)
	if err != nil {
		log.Print("os:", err)
	}
}
