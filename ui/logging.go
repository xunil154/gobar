package ui

import (
	"fmt"
)

var (
	DEBUG = true
)

func warning(s string) {
	fmt.Println(colorize("[!]", "yellow", "black"), s)
}
func info(s string) {
	fmt.Println(colorize("[I]", "blue", "black"), s)
}
func debug(s string) {
	if DEBUG {
		fmt.Println(colorize("[D]", "green", "black"), s)
	}
}
