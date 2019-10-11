// +build linux

package main

import (
	"strings"
	"syscall"
)

// GetTerminal returns the type of the terminal running our program.
func GetTerminal() TermType {
	return TermBash
}

var hasWSL bool
var hasWSLIsFilled = false

// GetTerminalHasWSL returns true if wsl is available.
func GetTerminalHasWSL() bool {
	if !hasWSLIsFilled {
		hasWSLIsFilled = true

		// Uname
		utsname := syscall.Utsname{}
		syscall.Uname(&utsname)

		release := int8ToString(utsname.Release)
		release = strings.ToLower(release)

		hasWSL = strings.Contains(release, "microsoft")
	}

	return hasWSL
}

// int8ToString simply converts int8 array to string
func int8ToString(array [65]int8) string {
	b := make([]byte, len(array))
	for i, v := range array {
		b[i] = byte(v)
	}
	return string(b)
}
