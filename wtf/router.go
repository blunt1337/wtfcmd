package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Route shows a help message, or find the command by the os.Args to execute, and executes it.
func Route(groups []*Group) {
	args := os.Args

	// No args
	if len(args) == 1 {
		ShowCommandError("no command given", nil, nil, groups)
	}

	// Special first flags
	showHelp := false
	debug := false
	switch os.Args[1] {
	case "help", "--help":
		// Show help
		args = args[2:]
		showHelp = true

		// Global help
		if len(args) == 0 {
			ShowHelp(groups, os.Stdout)
			return
		}
	case "--debug":
		args = args[2:]
		debug = true
	case "--autocomplete":
		AutocompleteCall(groups, args[2:])
	case "--builtin":
		ExecBuiltin(args[2:])
	default:
		args = args[1:]
	}

	group, command := findGroupAndCommand(groups, args)
	if group == nil {
		if showHelp {
			fmt.Fprintf(os.Stderr, "%s: error: command not found.\n\n", os.Args[0])
			ShowHelp(groups, os.Stdout)
			return
		}
		ShowCommandError("command not found", nil, nil, groups)
	} else {
		if command == nil {
			if showHelp {
				ShowHelp([]*Group{group}, os.Stdout)
				return
			}
			ShowCommandError("command not found", group, nil, nil)
		}

		// Execute or show help
		if showHelp {
			ShowHelpCommand(group, command)
		} else if group.Name == "" {
			parseAndExecuteCommand(group, command, args[1:], debug)
		} else {
			parseAndExecuteCommand(group, command, args[2:], debug)
		}
	}
}

// findGroupAndCommand finds the group and command matching the arguments.
func findGroupAndCommand(groups []*Group, args []string) (*Group, *Command) {
	// Find a root command
	for _, group := range groups {
		if group.Name == "" {
			command := findCommand(group, args[0])
			if command != nil {
				return group, command
			}
		}
	}

	// Find the group
	first := args[0]
	for _, group := range groups {
		if group.Name == first || inArray(first, group.Aliases) {
			if len(args) < 2 {
				return group, nil
			}

			command := findCommand(group, args[1])
			if command != nil {
				return group, command
			}
			return group, nil
		}
	}
	return nil, nil
}

// findCommand finds for the command by name.
func findCommand(group *Group, name string) *Command {
	for _, command := range group.Commands {
		if command.Name == name || inArray(name, command.Aliases) {
			return command
		}
	}
	return nil
}

// parseAndExecuteCommand parses the arguments left, then execute the real command with ExecCmd.
func parseAndExecuteCommand(group *Group, command *Command, args []string, debug bool) {
	// Parse the parameters from flags/args
	params := parseParams(group, command, args)

	// Execute the command
	ExecCmd(group, command, params, debug)
}

// parseParams parses/checks a command's arguments and build a parameter map.
func parseParams(group *Group, command *Command, args []string) map[string]interface{} {
	res := map[string]interface{}{}

	argIndex := 0
	canBeFlag := true
	l := len(args)
	for i := 0; i < l; i++ {
		arg := args[i]

		// -- = end of flags
		if canBeFlag && arg == "--" {
			canBeFlag = false
			continue
		}

		// Weird single -
		if canBeFlag && arg == "-" {
			continue
		}

		// Flags
		if canBeFlag && arg[0] == '-' {
			// Split by name=value, or name value
			index := strings.Index(arg, "=")

			var name string
			var value string
			var hasValue string
			var incI bool
			var isArray bool
			if index > 0 {
				name = arg[1:index]
				value = arg[index+1:]
				hasValue = "="
			} else {
				name = arg[1:]

				if i+1 < l {
					value = args[i+1]
					hasValue = "after"
				} else {
					hasValue = "no"
				}
			}

			// Long name
			if name[0] == '-' {
				var resValue interface{}
				name, resValue, incI, isArray = parseFlag(group, command, name[1:], value, hasValue)
				addParamValue(res, name, resValue, isArray)
			} else {
				// Aliases
				var resName string
				var resValue interface{}

				lg := len(name)
				for j := 0; j < lg; j++ {
					resName, resValue, incI, isArray = parseFlag(group, command, string(name[j]), value, hasValue)
					addParamValue(res, resName, resValue, isArray)
				}
			}

			// When the next argument is used as value
			if incI {
				i++
			}
			continue
		}

		// Arg
		if len(command.Config.Args) > argIndex {
			var value interface{}
			argCfg := command.Config.Args[argIndex]
			value = arg
			name := argCfg.Name[0]

			// Check the value
			if len(argCfg.Test) > 0 {
				var err error
				value, err = checkValue(arg, argCfg.Test)
				if err != nil {
					ShowCommandError(fmt.Sprintf("argument %s: %s", name, err.Error()), group, command, nil)
				}
			}

			addParamValue(res, name, value, argCfg.IsArray)
			if !argCfg.IsArray {
				argIndex++
			}
			continue
		}
		ShowCommandError("too many arguments", group, command, nil)
	}

	// Defaults & required
	l = len(command.Config.Args)
	for ; argIndex < l; argIndex++ {
		arg := command.Config.Args[argIndex]

		// Ignore if last arg is a filled array
		if arg.IsArray {
			if _, ok := res[arg.Name[0]]; ok {
				break
			}
		}

		if arg.Required {
			msg := "missing required argument: " + arg.Name[0]
			if len(arg.Desc) > 0 {
				msg += ".\n" + arg.Desc
			}
			ShowCommandError(msg, group, command, nil)
		} else {
			addParamValue(res, arg.Name[0], arg.Default, false)
		}
	}

	// Default flags
	for _, flag := range command.Config.Flags {
		if _, ok := res[flag.Name[0]]; !ok {
			addParamValue(res, flag.Name[0], flag.Default, false)
		}
	}

	return res
}

// parseFlag parses and checks a flag.
// Returns the flag name, flag value, true if used the next arguement, and true if the flag is an array
func parseFlag(group *Group, command *Command, name string, value string, hasValue string) (string, interface{}, bool, bool) {
	// Find the flag
	for _, flag := range command.Config.Flags {
		if inArray(name, flag.Name) {
			// Flag found
			nextValueUsed := false

			// Ignore the value after if not used with =
			if flag.Test == "$bool" {
				if hasValue != "=" {
					value = "true"
				}
			} else {
				if hasValue == "after" {
					nextValueUsed = true
				} else if hasValue == "no" {
					ShowCommandError(fmt.Sprintf("flag %s requires a value", name), group, command, nil)
				}
			}

			// Check the value
			if len(flag.Test) > 0 {
				var err error
				var resValue interface{}
				resValue, err = checkValue(value, flag.Test)
				if err != nil {
					ShowCommandError(fmt.Sprintf("flag %s %s", name, err.Error()), group, command, nil)
				}

				return flag.Name[0], resValue, nextValueUsed, flag.IsArray
			}
			return flag.Name[0], value, nextValueUsed, flag.IsArray
		}
	}
	ShowCommandError(fmt.Sprintf("flag %s not found", name), group, command, nil)
	return "", "", false, false
}

// checkValue run the test string to see if value match.
// It transform the value from string to the real type depending on the test.
func checkValue(value string, test string) (interface{}, error) {
	switch test {
	case "$bool":
		// Booleans
		b, err := strconv.ParseBool(value)
		if err != nil {
			return "", errors.New("should be a boolean")
		}
		if b {
			return true, nil
		}
		return false, nil
	case "$int":
		// Integers
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i, nil
		}
		return "", errors.New("should be an integer")
	case "$uint":
		// Unsigned integers
		if ui, err := strconv.ParseUint(value, 10, 64); err == nil {
			return ui, nil
		}
		return "", errors.New("should be a positive integer")
	case "$float", "$number":
		// Numbers
		if f, err := strconv.ParseFloat(value, 64); err == nil {
			return f, nil
		}
		return "", errors.New("should be a number")
	case "$file":
		// File exists
		if stat, err := os.Stat(value); os.IsNotExist(err) {
			return "", errors.New("file not found")
		} else if stat.IsDir() {
			return "", errors.New("must be a file")
		}
		return value, nil
	case "$dir":
		// Dir exists
		if stat, err := os.Stat(value); os.IsNotExist(err) {
			return "", errors.New("directory not found")
		} else if !stat.IsDir() {
			return "", errors.New("must be a directory")
		}
		return value, nil
	case "$dir/file":
		// File or dir exists
		if _, err := os.Stat(value); os.IsNotExist(err) {
			return "", errors.New("file or directory not found")
		}
		return value, nil
	case "$json":
		// Parse the json
		var data interface{}
		if err := json.Unmarshal([]byte(value), &data); err != nil {
			return "", fmt.Errorf("cannot decode json: %s", err.Error())
		}
		return data, nil
	default:
		// Regex
		regex := regexp.MustCompile(test)
		if regex.MatchString(value) {
			return value, nil
		}
		return "", errors.New("must match " + test)
	}
}

// in_array returns true if needle is in the stack.
func inArray(needle string, stack []string) bool {
	if stack == nil {
		return false
	}
	for _, value := range stack {
		if value == needle {
			return true
		}
	}
	return false
}

// addParamValue adds a value to the param list.
func addParamValue(res map[string]interface{}, name string, value interface{}, isArray bool) {
	if isArray {
		if _, ok := res[name]; !ok {
			res[name] = []interface{}{value}
		} else {
			if array, ok := res[name].([]interface{}); ok {
				res[name] = append(array, value)
			}
		}
	} else {
		res[name] = value
	}
}
