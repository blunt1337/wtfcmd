package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// TermType is the terminal currently runing this program.
type TermType int

const (
	// TermCmd is a cmd.exe terminal
	TermCmd TermType = iota
	// TermPowershell is a powershell.exe terminal
	TermPowershell
	// TermBash is a bash terminal
	TermBash
)

// GetLangAndCommandTemplate returns the language and the command template for this terminal.
func GetLangAndCommandTemplate(cmd *TermDependant) ([]string, string) {
	// Get the right command
	var cmdWrapper []string
	var cmdTpl string

	switch GetTerminal() {
	case TermBash:
		if cmd.Bash != "" {
			cmdTpl = cmd.Bash

			if runtime.GOOS == "windows" {
				if GetTerminalHasWSL() {
					// WSL
					cmdWrapper = []string{"wsl.exe", "/bin/sh", "-c"}
				} else {
					// Cygwin
					cmdWrapper = []string{"bash", "-c"}
				}
			} else {
				// Default bash
				cmdWrapper = []string{"/bin/sh", "-c"}
			}
		} else if GetTerminalHasWSL() && cmd.Powershell != "" {
			// WSL
			cmdWrapper = []string{"powershell.exe", "-command"}
			cmdTpl = cmd.Powershell
		} else {
			Panic("error: this command is not implemented in bash")
		}
	case TermPowershell, TermCmd:
		if cmd.Powershell != "" {
			// Powershell
			cmdWrapper = []string{"powershell.exe", "-command"}
			cmdTpl = cmd.Powershell
		} else if GetTerminalHasWSL() && cmd.Bash != "" {
			// WSL
			cmdWrapper = []string{"wsl.exe", "/bin/sh", "-c"}
			cmdTpl = cmd.Bash
		} else {
			Panic("error: this command is not implemented on windows")
		}
	}
	return cmdWrapper, cmdTpl
}

// EscapeArg escapes strings for argument in the terminal.
func EscapeArg(param interface{}) interface{} {
	if param == nil {
		param = ""
	}
	if str, ok := param.(string); ok {
		switch GetTerminal() {
		case TermBash:
			return "'" + strings.Replace(str, "'", "'\\''", -1) + "'"
		case TermCmd, TermPowershell:
			return "'" + strings.Replace(str, "'", "''", -1) + "'"
		}
	}
	return param
}

// UnescapeArg unescapes strings for argument in the terminal.
func UnescapeArg(param interface{}) interface{} {
	if str, ok := param.(string); ok {
		end := len(str) - 1

		switch GetTerminal() {
		case TermBash:
			if str[0] == '\'' && str[end] == '\'' {
				return strings.Replace(str[1:end], "'\\''", "'", -1)
			}
		case TermCmd, TermPowershell:
			if str[0] == '\'' && str[end] == '\'' {
				return strings.Replace(str[1:end], "''", "'", -1)
			}
		}
	}
	return param
}

// CmdAvailability returns a message if the command is not available on the OS.
func CmdAvailability(cmd *TermDependant) string {
	switch GetTerminal() {
	case TermBash:
		if cmd.Bash == "" {
			if GetTerminalHasWSL() && cmd.Powershell != "" {
				return orange + " (windows wsl)" + reset
			}
			return red + " (windows only)" + reset
		}
	case TermCmd, TermPowershell:
		if cmd.Powershell == "" {
			if GetTerminalHasWSL() && cmd.Bash != "" {
				return orange + " (bash wsl)" + reset
			}
			return red + " (bash only)" + reset
		}
	}
	return ""
}

// ExecBuiltin execute a builtin command.
func ExecBuiltin(args []string) {
	if len(args) <= 0 {
		Panic("No parameters")
	}

	//TODO: Loadings
	switch args[0] {
	case "Error":
		Error(strings.Join(args[1:], " "))
	case "Warn":
		Warn(strings.Join(args[1:], " "))
	case "Info":
		Info(strings.Join(args[1:], " "))
	case "Made":
		Made(strings.Join(args[1:], " "))
	case "Ask":
		if noColor {
			fmt.Print("[?] ")
		} else {
			fmt.Print("\033[38;5;99;01m[?]\033[00m ")
		}
		fmt.Println(strings.Join(args[1:], " "))
	case "AskYN":
		if AskYN(strings.Join(args[1:], " ")) {
			os.Exit(0)
		}
		os.Exit(1)
	case "Read":
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		fmt.Print(text)
	case "ReadSecure":
		fmt.Print(ReadSecure())
	case "AskList":
		values := ""
		dflt := -1

		if len(args) >= 3 {
			values = args[2]
			if len(args) >= 4 {
				if i, err := strconv.Atoi(args[3]); err == nil {
					dflt = i
				}
			}
		}
		os.Exit(AskList(strings.Split(values, ","), dflt, args[1]))
	case "Bell":
		Bell()
	}
	os.Exit(0)
}
