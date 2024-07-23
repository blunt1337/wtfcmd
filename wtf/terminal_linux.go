//go:build linux
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

		b := make([]byte, len(utsname.Release))
		for i, v := range utsname.Release {
			b[i] = byte(v)
		}

		release := string(b)
		release = strings.ToLower(release)

		hasWSL = strings.Contains(release, "microsoft")
	}

	return hasWSL
}
