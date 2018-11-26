// +build windows

package main

import (
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"syscall"
)

// Currently parent terminal.
var term TermType
var termIsFilled = false

// GetTerminal returns the type of the terminal running our program.
func GetTerminal() TermType {
	if !termIsFilled {
		// Default
		term = TermCmd

		// Parent pid
		ppid := os.Getppid()

		// Executable path
		out, err := exec.Command("wmic", "process", "where", "processid="+strconv.Itoa(ppid), "get", "ExecutablePath").Output()
		if err == nil {
			if regexp.MustCompile("\\\\WindowsPowerShell\\\\v[^\\\\]+\\\\powershell.exe").Match(out) {
				term = TermPowershell

				// Enable colors
				fd := os.Stdout.Fd()
				var mode uint32
				if err = syscall.GetConsoleMode(syscall.Handle(fd), &mode); err == nil {
					mode |= 0x0004 // ENABLE_VIRTUAL_TERMINAL_PROCESSING
					syscall.NewLazyDLL("kernel32.dll").NewProc("SetConsoleMode").Call(fd, uintptr(mode), 0)
				}
			} else if regexp.MustCompile("cygwin\\\\bin\\\\").Match(out) {
				term = TermBash
			}
		}
	}

	termIsFilled = true
	return term
}
