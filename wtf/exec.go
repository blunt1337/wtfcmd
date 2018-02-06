package main

import (
	"bytes"
	"fmt"
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
func ExecCmd(group *Group, command *Command, params map[string]interface{}, debug bool) {
	config := command.Config

	// Get the right command
	cmdWrapper, cmdTpl := GetLangAndCommandTemplate(config.Cmd)

	// Generate the command from the template
	tmpl, err := template.New("cmd").Funcs(getTplFuncs(command.Config)).Parse(cmdTpl)
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
	process := exec.Command(cmdWrapper[0], cmdWrapper[1], cmd)

	// Pipes
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr
	process.Stdin = os.Stdin
	process.Dir = ResolveCwd(command.Config)

	// Start
	if err := process.Start(); err != nil {
		Panic(err.Error())
	}

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
}

// ResolveCwd finds where the current working dir will be
func ResolveCwd(cfg *Config) string {
	var to_resolve string
	if cfg.Cwd == nil {
		to_resolve = ""
	} else {
		switch term {
		case TermBash:
			to_resolve = cfg.Cwd.Bash
		case TermCmd, TermPowershell:
			to_resolve = cfg.Cwd.Powershell
		}
	}
	lg := len(to_resolve)

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

	// Starting with dot: config dir + to_resolve
	if to_resolve[0] == '.' {
		return path.Join(filepath.Dir(cfg.File), to_resolve)
	}

	// Starting with / or x:/ absolute path
	switch term {
	case TermBash:
		if to_resolve[0] == '/' {
			return to_resolve
		}
	case TermCmd, TermPowershell:
		if lg >= 3 && to_resolve[1] == ':' && (to_resolve[1] == '/' || to_resolve[1] == '\\') {
			return to_resolve
		}
	}

	// Default: cwd + to_resolve
	return path.Join(cwd, to_resolve)
}

// getTplFuncs creates a big funcMap with all strings functions and more for the template.
func getTplFuncs(config *Config) template.FuncMap {
	return template.FuncMap{
		// All strings
		"contains":    strings.Contains,
		"containsAny": strings.ContainsAny,
		"count":       strings.Count,
		"equalFold":   strings.EqualFold,
		"fields":      strings.Fields,
		"fieldsFunc":  strings.FieldsFunc,
		"hasPrefix":   strings.HasPrefix,
		"hasSuffix":   strings.HasSuffix,
		"index":       strings.Index,
		"indexAny":    strings.IndexAny,
		"indexFunc":   strings.IndexFunc,
		"indexRune":   strings.IndexRune,
		"join":        strings.Join,
		"lastIndex":   strings.LastIndex,
		"map":         strings.Map,
		"newReplacer": strings.NewReplacer,
		"repeat":      strings.Repeat,
		"replace":     strings.Replace,
		"split":       strings.Split,
		"splitAfter":  strings.SplitAfter,
		"splitAfterN": strings.SplitAfterN,
		"splitN":      strings.SplitN,
		"title":       strings.Title,
		"toLower":     strings.ToLower,
		"toTitle":     strings.ToTitle,
		"toUpper":     strings.ToUpper,
		"trim":        strings.Trim,
		"trimPrefix":  strings.TrimPrefix,
		"trimSpace":   strings.TrimSpace,
		"trimSuffix":  strings.TrimSuffix,
		"esc":         EscapeArg,
		"escape":      EscapeArg,
		"raw":         UnescapeArg,
		"unescape":    UnescapeArg,
		"configdir": func() string {
			return filepath.Dir(config.File)
		},
		"error": func(args ...interface{}) string {
			return os.Args[0] + " --builtin Error " + join(args)
		},
		"panic": func(args ...interface{}) string {
			return os.Args[0] + " --builtin Error " + join(args) + "&& exit 1"
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
			res := os.Args[0] + " --builtin Ask " + join(args)
			switch term {
			case TermBash:
				res += "\nread response "
			/*case TermCmd:
			res += "\nset /p response=\"\" "*/
			case TermPowershell, TermCmd:
				res += "\n$response = Read-Host -Prompt '' "
			}
			return res
		},
		"askYN": func(args ...interface{}) string {
			return os.Args[0] + " --builtin AskYN " + join(args)
		},
		/*"AskList": func (args... interface{}) string {
		    return os.Args[0] + " --builtin AskList " + join(args)
		},*/
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
		res += fmt.Sprint(arg) + " "
	}
	return res
}
