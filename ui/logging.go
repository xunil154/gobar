package ui

import (
	"fmt"
)

var (
	DEBUG = true
)

func warning(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(colorize("[!]", "yellow", "black")+format, args...))
}
func info(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(colorize("[I]", "blue", "black")+format, args...))
}
func debug(format string, args ...interface{}) {
	if DEBUG {
		fmt.Println(fmt.Sprintf(colorize("[D]", "green", "black")+format, args...))
	}
}
