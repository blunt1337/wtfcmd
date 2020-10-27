package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"text/template"
)

// ExecCmd run the command as a child process.
// Use command.Cmd template and params to build the command.
// Execute it in bash / powershell depending on the os.
func ExecCmd(group *Group, command *Command, params map[string]interface{}, debug bool, stdout io.Writer) {
	config := command.Config

	// Get the right command
	cmdWrapper, cmdTpl := GetLangAndCommandTemplate(config.Cmd)

	// Generate the command from the template
	tmpl, err := template.New("cmd").Funcs(getTplFuncs(config)).Parse(cmdTpl)
	if err != nil {
		Panic(fmt.Sprintf("Error in %s: The template for the command %s %s cannot be compiled: %s", config.File, group.Name, command.Name, err.Error()))
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, params)
	if err != nil {
		Panic(fmt.Sprintf("Error in %s: The template for the command %s %s failed to execute: %s", config.File, group.Name, command.Name, err.Error()))
	}
	cmd := buffer.String()

	// Debug the command
	if debug {
		Made("The command to execute is:\n")
		fmt.Println(cmd + "\n")

		// Ask to execute
		if !AskYN("Execute?", false) {
			os.Exit(0)
		}
	}

	// Create the process
	cmdWrapper0, cmdWrapperN := cmdWrapper[0], cmdWrapper[1:]
	cmdWrapperN = append(cmdWrapperN, cmd)
	process := exec.Command(cmdWrapper0, cmdWrapperN...)
	process.Dir = ResolveCwd(command.Config)

	// Pipes
	if stdout != nil || len(command.Config.Expect) == 0 {
		if stdout != nil {
			process.Stdout = stdout
		} else {
			process.Stdout = os.Stdout
		}
		process.Stderr = os.Stderr
		process.Stdin = os.Stdin
	} else {
		// Run every expects
		for i, exp := range command.Config.Expect {
			if exp.Cmd != nil {
				exp.Send += ExecExpectCmd(i, group, command, params, debug)
			}
		}

		// "Expect" watch "output" to send input code
		process.Stdout, process.Stderr, process.Stdin, err = expectPipes(process, command.Config.Expect)
		if err != nil {
			Panic(err.Error())
		}
	}

	// Start
	if err := process.Start(); err != nil {
		Panic(err.Error())
	}

	// Normal mode
	if stdout == nil {
		// Wait for the end
		if err := process.Wait(); err != nil {
			// The program has exited with an exit code != 0
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					os.Exit(status.ExitStatus())
				}
			} else {
				Panic(err.Error())
			}
		}
		os.Exit(0)
	} else {
		process.Wait()
	}
}

// ExecExpectCmd executes the "expect" command and return the value.
func ExecExpectCmd(index int, group *Group, command *Command, params map[string]interface{}, debug bool) string {
	config := command.Config

	// Clone arguments
	config2 := Config{
		File:  config.File,
		Cmd:   command.Config.Expect[index].Cmd,
		Args:  config.Args,
		Flags: config.Flags,
		Cwd:   config.Cwd,
	}
	group2 := Group{Name: group.Name + " " + command.Name}
	command2 := Command{Name: fmt.Sprintf("expect[%d]", index), Config: &config2}

	// Execute
	var stdout bytes.Buffer
	ExecCmd(&group2, &command2, params, debug, bufio.NewWriter(&stdout))

	return stdout.String()
}

// ResolveCwd finds where the current working dir will be
func ResolveCwd(cfg *Config) string {
	var toResolve string
	if cfg.Cwd == nil {
		toResolve = ""
	} else {
		switch GetTerminal() {
		case TermBash:
			toResolve = cfg.Cwd.Bash
		case TermCmd, TermPowershell:
			toResolve = cfg.Cwd.Powershell
		}
	}
	lg := len(toResolve)

	// Current working dir
	var cwd string
	var err error
	if cwd, err = os.Getwd(); err != nil {
		cwd = "/"
	}

	// Empty
	if lg == 0 {
		return cwd
	}

	// Starting with dot: config dir + toResolve
	if toResolve[0] == '.' {
		return path.Join(filepath.Dir(cfg.File), toResolve)
	}

	// Starting with / or x:/ absolute path
	switch GetTerminal() {
	case TermBash:
		if toResolve[0] == '/' {
			return toResolve
		}
	case TermCmd, TermPowershell:
		if lg >= 3 && toResolve[1] == ':' && (toResolve[1] == '/' || toResolve[1] == '\\') {
			return toResolve
		}
	}

	// Default: cwd + toResolve
	return path.Join(cwd, toResolve)
}

// getTplFuncs creates a big funcMap with all strings functions and more for the template.
func getTplFuncs(config *Config) template.FuncMap {
	return template.FuncMap{
		// All strings
		"compare":        strings.Compare,
		"contains":       strings.Contains,
		"containsAny":    strings.ContainsAny,
		"containsRune":   strings.ContainsRune,
		"count":          strings.Count,
		"equalFold":      strings.EqualFold,
		"fields":         strings.Fields,
		"fieldsFunc":     strings.FieldsFunc,
		"hasPrefix":      strings.HasPrefix,
		"hasSuffix":      strings.HasSuffix,
		"index":          strings.Index,
		"indexAny":       strings.IndexAny,
		"indexFunc":      strings.IndexFunc,
		"indexRune":      strings.IndexRune,
		"join":           strings.Join,
		"lastIndex":      strings.LastIndex,
		"lastIndexAny":   strings.LastIndexAny,
		"lastIndexByte":  strings.LastIndexByte,
		"lastIndexFunc":  strings.LastIndexFunc,
		"map":            strings.Map,
		"repeat":         strings.Repeat,
		"replace":        strings.Replace,
		"split":          strings.Split,
		"splitAfter":     strings.SplitAfter,
		"splitAfterN":    strings.SplitAfterN,
		"splitN":         strings.SplitN,
		"title":          strings.Title,
		"toLower":        strings.ToLower,
		"toLowerSpecial": strings.ToLowerSpecial,
		"toTitle":        strings.ToTitle,
		"toTitleSpecial": strings.ToTitleSpecial,
		"toUpper":        strings.ToUpper,
		"toUpperSpecial": strings.ToUpperSpecial,
		"trim":           strings.Trim,
		"trimFunc":       strings.TrimFunc,
		"trimLeft":       strings.TrimLeft,
		"trimLeftFunc":   strings.TrimLeftFunc,
		"trimPrefix":     strings.TrimPrefix,
		"trimRight":      strings.TrimRight,
		"trimRightFunc":  strings.TrimRightFunc,
		"trimSpace":      strings.TrimSpace,
		"trimSuffix":     strings.TrimSuffix,
		// Escape a string for bash/powershell.
		"esc":    EscapeArg,
		"escape": EscapeArg,
		// Unescape a string from bash/powershell.
		"raw":      UnescapeArg,
		"unescape": UnescapeArg,
		// Convert first argument to json.
		// Second argument is pretty print, default false.
		// Return "false" on error.
		"json": func(arg_and_pretty ...interface{}) string {
			if len(arg_and_pretty) == 0 {
				return "null"
			}

			pretty := false
			if len(arg_and_pretty) == 2 {
				if argPretty, ok := arg_and_pretty[1].(bool); ok && argPretty {
					pretty = argPretty
				}
			}

			var b []byte
			var err error

			if pretty {
				b, err = json.MarshalIndent(arg_and_pretty[0], "", "\t")
			} else {
				b, err = json.Marshal(arg_and_pretty[0])
			}

			if err != nil {
				return "false"
			}
			return string(b)
		},
		// Convert a json string into an interface{}.
		// Return false on error.
		"jsonParse": func(arg string) interface{} {
			var res interface{}
			err := json.Unmarshal([]byte(arg), &res)
			if err != nil {
				return false
			}
			return res
		},
		// Directory of the configuration file of the command running.
		"configdir": func() string {
			return filepath.Dir(config.File)
		},
		"error": func(args ...interface{}) string {
			return os.Args[0] + " --builtin Error " + join(args)
		},
		"panic": func(args ...interface{}) string {
			return os.Args[0] + " --builtin Error " + join(args) + "; exit 1"
		},
		"warn": func(args ...interface{}) string {
			return os.Args[0] + " --builtin Warn " + join(args)
		},
		"info": func(args ...interface{}) string {
			return os.Args[0] + " --builtin Info " + join(args)
		},
		"made": func(args ...interface{}) string {
			return os.Args[0] + " --builtin Made " + join(args)
		},
		"ask": func(args ...interface{}) string {
			return os.Args[0] + " --builtin Ask " + join(args)
		},
		"askYN": func(args ...interface{}) string {
			return os.Args[0] + " --builtin AskYN " + join(args)
		},
		"read": func() string {
			return os.Args[0] + " --builtin Read ."
		},
		"readSecure": func() string {
			return os.Args[0] + " --builtin ReadSecure ."
		},
		"AskList": func(args ...interface{}) string {
			return os.Args[0] + " --builtin AskList " + join(args)
		},
		"bell": func() string {
			return os.Args[0] + " --builtin Bell ."
		},
	}
}

// join joins args with a space.
// /!\ there is a trailing space at the end.
func join(args []interface{}) string {
	var res string
	for _, arg := range args {
		res += fmt.Sprint(EscapeArg(arg)) + " "
	}
	return res
}
