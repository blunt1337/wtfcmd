//go:build windows
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

// GetTerminal returns the type of the terminal running our program
func GetTerminal() TermType {
	if !termIsFilled {
		termIsFilled = true

		// Default
		term = TermCmd

		// Powershell windows 10/11 with wmic deprecated
		if _, ok := os.LookupEnv("PSModulePath"); ok {
			term = TermPowershell

			// Enable colors
			fd := os.Stdout.Fd()
			var mode uint32
			if err := syscall.GetConsoleMode(syscall.Handle(fd), &mode); err == nil {
				mode |= 0x0004 // ENABLE_VIRTUAL_TERMINAL_PROCESSING
				syscall.NewLazyDLL("kernel32.dll").NewProc("SetConsoleMode").Call(fd, uintptr(mode), 0)
			}
			return term
		}

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
			} else if regexp.MustCompile("\\\\wsl\\.exe").Match(out) {
				hasWSL = true
				hasWSLIsFilled = true
				term = TermBash
			}
		}
	}

	return term
}

var hasWSL bool
var hasWSLIsFilled = false

// GetTerminalHasWSL returns true if wsl is available.
func GetTerminalHasWSL() bool {
	if !hasWSLIsFilled {
		hasWSLIsFilled = true

		// WSL binary exists
		_, err := exec.LookPath("wsl.exe")
		hasWSL = err == nil
	}

	return hasWSL
}
