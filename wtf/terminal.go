package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

// TermType is the terminal currently runing this program.
type TermType int

const (
	TermCmd TermType = iota
	TermPowershell
	TermBash
)

// Currently parent terminal.
var term TermType
var termIsFilled = false

// GetTerminal returns the type of the terminal running our program.
func GetTerminal() TermType {
	if !termIsFilled {
		if runtime.GOOS == "windows" {
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
				} else if regexp.MustCompile("cygwin\\\\bin\\\\zsh.exe").Match(out) {
					term = TermBash
				}
			}
		} else {
			term = TermBash
		}
	}

	termIsFilled = true
	return term
}

// GetLangAndCommandTemplate returns the language and the command template for this terminal.
func GetLangAndCommandTemplate(cmd *TermDependant) ([]string, string) {
	// Get the right command
	var cmdWrapper []string
	var cmdTpl string

	switch term {
	case TermBash:
		if cmd.Bash != "" {
			if runtime.GOOS == "windows" {
				cmdWrapper = []string{"bash", "-c"}
			} else {
				cmdWrapper = []string{"/bin/bash", "-c"}
			}
			cmdTpl = cmd.Bash
		} else {
			Panic("error: this command is not implemented in bash")
		}
	case TermPowershell, TermCmd:
		if cmd.Powershell != "" {
			cmdWrapper = []string{"powershell.exe", "-command"}
			cmdTpl = cmd.Powershell
		} else {
			Panic("error: this command is not implemented on windows")
		}
	}
	/*case TermPowershell:
		if cmd.Powershell != "" {
			cmdWrapper = []string{"powershell.exe", "-command"}
			cmdTpl = cmd.Powershell
		} else if cmd.Cmd != "" {
			cmdWrapper = []string{"cmd.exe", "/C"}
			cmdTpl = cmd.Cmd
		} else {
			Panic("error: this command is not implemented on windows")
		}
	case TermCmd:
		if cmd.Cmd != "" {
			cmdWrapper = []string{"cmd.exe", "/C"}
			cmdTpl = cmd.Cmd
		} else if cmd.Powershell != "" {
			cmdWrapper = []string{"powershell.exe", "-command"}
			cmdTpl = cmd.Powershell
		} else {
			Panic("error: this command is not implemented on windows")
		}
	}*/
	return cmdWrapper, cmdTpl
}

// EscapeArg escapes strings for argument in the terminal.
func EscapeArg(param interface{}) interface{} {
	if param == nil {
		param = ""
	}
	if str, ok := param.(string); ok {
		switch term {
		case TermBash:
			return "'" + strings.Replace(str, "'", "'\\''", -1) + "'"
		case TermPowershell:
			return "'" + strings.Replace(str, "'", "''", -1) + "'"
		}
	}
	return param
}

// UnescapeArg unescapes strings for argument in the terminal.
func UnescapeArg(param interface{}) interface{} {
	if str, ok := param.(string); ok {
		end := len(str) - 1

		switch term {
		case TermBash:
			if str[0] == '\'' && str[end] == '\'' {
				return strings.Replace(str[1:end], "'\\''", "'", -1)
			}
		case TermPowershell:
			if str[0] == '\'' && str[end] == '\'' {
				return strings.Replace(str[1:end], "''", "'", -1)
			}
		}
	}
	return param
}

// CmdAvailability returns a message if the command is not available on the OS.
func CmdAvailability(cmd *TermDependant) string {
	switch term {
	case TermBash:
		if cmd.Bash == "" {
			return red + " (windows only)" + reset
		}
	case TermCmd, TermPowershell:
		if cmd.Powershell == "" {
			return red + " (bash only)" + reset
		}
	}
	/*case TermPowershell:
		if cmd.Powershell == "" && cmd.Cmd == "" {
			return red + " (bash only)" + reset
		}
	case TermCmd:
		if cmd.Cmd == "" {
			if cmd.Powershell == "" {
				return red + " (bash only)" + reset
			}
			return red + " (powershell preferably)" + reset
		}
	}*/
	return ""
}

// ExecBuiltin execute a builtin command.
func ExecBuiltin(args []string) {
	if len(args) <= 1 {
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
		casted := make([]interface{}, len(args))
		for i, v := range args {
			casted[i] = v
		}

		if AskYN(casted[1:]...) {
			os.Exit(0)
		}
		os.Exit(1)
	/*case "AskList":
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

	        os.Exit(AskList(strings.Split(values, ","), dflt, args[1]))*/
	case "Bell":
		Bell()
	}
	os.Exit(0)
}
