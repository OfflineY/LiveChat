package util_test

import "testing"

// test file

import (
	"LiveChat/util"
	"fmt"
	"log"
)

func TestIsIniExist(t *testing.T) {
	if util.IsIniExist("../conf.ini") {
		log.Println("Conf exists, skip initialization.")
	} else {
		log.Println("Conf does not exist, start initialize.")
	}
	cfg := util.IniLoad("../conf.ini")
	fmt.Println(cfg.Section("web").Key("address").Value())
}

func TestIniInitial(t *testing.T) {
	util.IniInitial("./conf.ini")
}

func TestPrint(t *testing.T) {
	util.Paint("v3.0 BETA")
}
