package modules

import (
	"embed"
)

//go:embed assets/* static/*
var f embed.FS
