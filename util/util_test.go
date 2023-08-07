package util_test

import (
	"LiveChat/util"
	"testing"
)

func TestIniInitial(t *testing.T) {
	util.IniInitial("./conf01.ini")
}
