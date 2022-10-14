package modules_test

import (
	"OnlineChat/modules"
	"testing"
)

func Benchmark_Bootstrap(b *testing.B) {
	modules.Bootstrap("test")
}

func Benchmark_Print(b *testing.B) {
	modules.PrintLogo("test")
}
