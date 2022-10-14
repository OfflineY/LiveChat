package modules

import (
	"fmt"
	"log"
)

func GetUserList() {
	fmt.Print("UserList指令存在但还在开发中\n")
}

// 向在线用户中添加用户
func AddUser(name string) {
	log.Print(name, " 加入了")
}
