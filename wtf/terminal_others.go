//go:build !linux && !windows
// +build !linux,!windows

package main

// GetTerminal returns the type of the terminal running our program.
func GetTerminal() TermType {
	return TermBash
}

// GetTerminalHasWSL returns true if wsl is available.
func GetTerminalHasWSL() bool {
	return false
}
