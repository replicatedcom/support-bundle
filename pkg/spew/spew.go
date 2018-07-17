package spew

import "github.com/davecgh/go-spew/spew"

var Config = &spew.ConfigState{
	Indent:                  "  ",
	DisablePointerAddresses: true,
}

func Sdump(a ...interface{}) string {
	return Config.Sdump(a...)
}
