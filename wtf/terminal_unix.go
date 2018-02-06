// +build !windows

package main

// Currently parent terminal.
var term TermType

// GetTerminal returns the type of the terminal running our program.
func GetTerminal() TermType {
	term = TermBash
	return term
}
