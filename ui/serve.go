package ui

import (
	"embed"
	"io/fs"
)

//nolint:golint
//go:embed build
var uiDir embed.FS
var Build fs.FS

func init() {
	var err error
	Build, err = fs.Sub(uiDir, "build")
	if err != nil {
		panic(err)
	}
}
