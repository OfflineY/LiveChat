package modules_test

import (
	"OnlineChat/modules"
	"testing"
)

func Benchmark_GetChat(b *testing.B) {
	modules.GetChat("ws://localhost:8080/socket")
}

func Benchmark_SendMsg(b *testing.B) {
	modules.SendMsg(nil, "test")
}
