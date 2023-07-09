package util

import (
	"fmt"
	"log"
	"testing"
)

func TestIsIniExist(t *testing.T) {
	if IsIniExist("../conf.ini") {
		log.Println("Conf exists, skip initialization.")
	} else {
		//log.Println("Conf does not exist, start initialize.")
	}

	cfg := IniLoad("../conf.ini")
	fmt.Println(cfg.Section("web").Key("address").Value())
}
