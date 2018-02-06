package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// TermType is the terminal currently runing this program.
type TermType int

const (
	TermCmd TermType = iota
	TermPowershell
	TermBash
)

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
