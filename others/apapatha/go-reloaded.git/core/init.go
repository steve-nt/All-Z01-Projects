package core

import (
	"regexp"
	"strings"
)

const (
	alphaOffset = 'a' - 'A'
	expArgc     = 3
	punc        = "[.,!?:;']+"
)

var (
	re     = regexp.MustCompile(punc)
	cmdMap map[string]func(string) string
)

func init() {
	cmdMap = make(map[string]func(string) string)
	cmdMap["(bin"] = bin
	cmdMap["(hex"] = hex
	cmdMap["(up"] = strings.ToUpper
	cmdMap["(low"] = strings.ToLower
	cmdMap["(cap"] = capitalize
}
