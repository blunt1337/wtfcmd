package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"regexp"
	"strings"
)

// Config is the first level of the json.
type Config struct {
	File        string
	Group       []string
	Name        []string
	Cmd         *TermDependant
	Desc        string
	Args        []*ArgOrFlag
	Flags       []*ArgOrFlag
	Cwd         *TermDependant
	StopOnError bool
}

// ArgOrFlag holds an arguments or flags from the config.
type ArgOrFlag struct {
	Name     []string
	Desc     string
	Required bool
	Default  interface{}
	Test     string
	IsArray  bool
}

// TermDependant holds a string for bash, cmd, powershell from the config.
type TermDependant struct {
	Bash       string
	Powershell string
}

var nameRegex = regexp.MustCompile("^[\\p{L}0-9][\\p{L}0-9:._-]*$")
var aliasRegex = regexp.MustCompile("^[\\p{L}0-9]$")

// ParseConfigs parses the JSON configuration to structs.
// Checks some fields too.
// Returns an array of Configs or the parsing/checking error.
func ParseConfigs(input io.Reader, cfgs []*Config, file string, format string) ([]*Config, error) {
	var data interface{}

	switch format {
	case "json":
		// Parse the json
		dec := json.NewDecoder(input)
		if err := dec.Decode(&data); err != nil {
			return nil, err
		}
	case "yaml", "yml":
		// Read the file
		buf := new(bytes.Buffer)
		buf.ReadFrom(input)

		if err := yaml.Unmarshal(buf.Bytes(), &data); err != nil {
			return nil, err
		}

		// Fix maps
		data = fixYamlMaps(data)
	}

	// Assertion of array
	if array, ok := data.([]interface{}); ok {
		// Convert to structs
		for i, v := range array {
			// Parse config
			cfg, err := parseConfig(v)
			if err != nil {
				return nil, fmt.Errorf("[%d]%s", i, err.Error())
			}

			cfg.File = file
			cfgs = append(cfgs, cfg)
		}
		return cfgs, nil
	}
	return nil, errors.New("the configuration must be an array of objects")
}

// parseConfig parses a config data.
func parseConfig(data interface{}) (*Config, error) {
	// Assertion of object
	if hash, ok := data.(map[string]interface{}); ok {
		res := new(Config)

		// Foreach key => value
		for k, v := range hash {
			switch k {
			case "group", "name", "desc":
				value, err := parseStringArray(v)
				if err != nil {
					return nil, fmt.Errorf(".%s%s", k, err.Error())
				}

				switch k {
				case "group":
					res.Group = filterArray(value)
				case "name":
					res.Name = filterArray(value)
				case "desc":
					res.Desc = strings.Join(value, "\n")
				}
			case "cmd":
				value, err := parseTermDependant(v, "\n")
				if err != nil {
					return nil, fmt.Errorf(".%s%s", k, err.Error())
				}
				res.Cmd = value
			case "cwd":
				value, err := parseTermDependant(v, string(os.PathSeparator))
				if err != nil {
					return nil, fmt.Errorf(".%s%s", k, err.Error())
				}
				res.Cwd = value
			case "args", "flags":
				value, err := parseArgOrFlagArray(v, k == "args")
				if err != nil {
					return nil, fmt.Errorf(".%s%s", k, err.Error())
				}

				switch k {
				case "args":
					res.Args = value
				case "flags":
					res.Flags = value
				}
			case "stopOnError":
				res.StopOnError, ok = v.(bool)
				if !ok {
					return nil, errors.New(".stopOnError : must be a boolean")
				}
			default:
				return nil, fmt.Errorf(".%s : unknown property", k)
			}
		}

		// Name required
		if len(res.Name) == 0 {
			return nil, errors.New(".name : is required")
		}

		// Alphanum name
		for _, name := range res.Name {
			if !nameRegex.MatchString(name) {
				return nil, errors.New(".name : contain invalid character")
			}
		}

		// Alphanum group
		if res.Group != nil {
			for _, name := range res.Group {
				if !nameRegex.MatchString(name) {
					return nil, errors.New(".group : contain invalid character")
				}
			}
		}

		// Cmd required
		if res.Cmd == nil || (res.Cmd.Bash == "" && res.Cmd.Powershell == "" /*&& res.Cmd.Cmd == ""*/) {
			return nil, errors.New(".cmd : is required")
		}

		// Not 2 args/flags with the same name
		if name, ok := areNamesUnique(res.Args, res.Flags); !ok {
			return nil, fmt.Errorf(".args/flags : the name %s is used more than once", name)
		}

		// Not an arg required = false before a required = true
		if name, ok := checkArgOrder(res.Args); !ok {
			return nil, fmt.Errorf(".args : the argument %s cannot be required after an optionnal one", name)
		}

		// Only the last argument can be an array
		if name, ok := checkArgIsArray(res.Args); !ok {
			return nil, fmt.Errorf(".args : the argument %s cannot be an array, only last argument can", name)
		}

		return res, nil
	}
	return nil, errors.New(" : the configuration must be an object")
}

// parseStringArray checks and returns a string array from a string or string array.
func parseStringArray(data interface{}) ([]string, error) {
	switch value := data.(type) {
	case string:
		return []string{value}, nil
	case []interface{}:
		res := []string{}
		for _, v := range value {
			switch subvalue := v.(type) {
			case string:
				res = append(res, subvalue)
			default:
				return nil, errors.New(" : must be a string or an array of strings")
			}
		}
		return res, nil
	default:
		return nil, errors.New(" : must be a string or an array of strings")
	}
}

// parseTermDependant checks and returns an os definition from the json interface.
func parseTermDependant(data interface{}, jointer string) (*TermDependant, error) {
	res := new(TermDependant)

	// Map
	if obj, ok := data.(map[string]interface{}); ok {
		for k, subdata := range obj {
			switch k {
			case "bash" /*"cmd",*/, "powershell":
				// String/array
				value, err := parseStringArray(subdata)
				if err != nil {
					return nil, fmt.Errorf("[%s]%s", k, err.Error())
				}

				switch k {
				case "bash":
					res.Bash = strings.TrimSpace(strings.Join(value, jointer))
				case "powershell":
					res.Powershell = strings.TrimSpace(strings.Join(value, jointer))
				}
				// case "cmd":
				// 	res.Cmd = strings.TrimSpace(strings.Join(value, "\n"))
			default:
				return nil, fmt.Errorf("[%s] : unrecognized terminal type", k)
			}
		}
	} else {
		// String/array
		value, err := parseStringArray(data)
		if err != nil {
			return nil, err
		}

		tmp := strings.TrimSpace(strings.Join(value, "\n"))
		res.Bash = tmp
		res.Powershell = tmp
		// res.Cmd = tmp
	}

	return res, nil
}

// parseArgOrFlagArray checks and parses array of args/flags from the json interface.
func parseArgOrFlagArray(data interface{}, isArg bool) ([]*ArgOrFlag, error) {
	switch value := data.(type) {
	case []interface{}:
		res := []*ArgOrFlag{}
		for i, v := range value {
			// Parse arg or flag
			subvalue, err := parseArgOrFlag(v, isArg)
			if err != nil {
				return nil, fmt.Errorf("[%d]%s", i, err.Error())
			}
			res = append(res, subvalue)
		}

		return res, nil
	default:
		return nil, errors.New(" : must be a array of objects")
	}
}

// parseArgOrFlag checks and parses arg/flag from the json struct.
func parseArgOrFlag(jsonInterface interface{}, isArg bool) (*ArgOrFlag, error) {
	// Assertion of object
	if data, ok := jsonInterface.(map[string]interface{}); ok {
		res := new(ArgOrFlag)

		// Foreach key => value
		for k, v := range data {
			switch k {
			case "name", "desc":
				value, err := parseStringArray(v)
				if err != nil {
					return nil, fmt.Errorf(".%s%s", k, err.Error())
				}

				switch k {
				case "name":
					res.Name = filterArray(value)
				case "desc":
					res.Desc = strings.Join(value, "\n")
				}
			case "required", "is_array", "array":
				if subvalue, ok := v.(bool); ok {
					switch k {
					case "required":
						res.Required = subvalue
					case "is_array", "array":
						res.IsArray = subvalue
					}
				} else {
					return nil, fmt.Errorf(".%s : must be a boolean", k)
				}
			case "default":
				// Ignore now, do it when we have all other info
			case "test":
				if subvalue, ok := v.(string); ok {
					res.Test = subvalue
				} else {
					return nil, fmt.Errorf(".%s : must be a string regex", k)
				}
			default:
				return nil, fmt.Errorf(".%s : unknown property", k)
			}
		}

		// Name required
		l := len(res.Name)
		if l == 0 {
			return nil, errors.New(".name : is required")
		}

		// Alphanum name
		if !nameRegex.MatchString(res.Name[0]) {
			return nil, errors.New(".name : contain invalid character")
		}
		// Alphanum alias
		for i := 1; i < l; i++ {
			if !aliasRegex.MatchString(res.Name[i]) {
				return nil, errors.New(".name : aliases are 1 character alphanumeric")
			}
		}

		// Args have no aliases
		if isArg && len(res.Name) > 1 {
			return nil, errors.New(".name : argument cannot have aliases")
		}

		// No required on flags
		if !isArg && res.Required {
			return nil, errors.New(".required : flags cannot be required")
		}

		// Prevent is_array and json together
		if res.IsArray && res.Test == "$json" {
			return nil, errors.New(".test : cannot be both $json and is_array at the same time")
		}

		// Default
		if rawDefault, ok := data["default"]; ok {
			// Not a required = true with a "default"
			if res.Required {
				return nil, errors.New(".default : cannot have a default value if required")
			}

			// Read value
			switch value := rawDefault.(type) {
			case []interface{}:
				// Array only valid for is_array
				if !res.IsArray {
					return nil, errors.New(".default : must be a number, a string or a boolean")
				}

				for index, v2 := range value {
					switch v2.(type) {
					case int, int32, int64, float64, float32, string, bool:
						res.Default = value
					default:
						return nil, fmt.Errorf(".default[%d] : must be a number, a string or a boolean", index)
					}
				}
			case int, int32, int64, float64, float32, bool, string:
				if str, isStr := value.(string); isStr && res.Test == "$json" {
					var obj interface{}
					if err := json.Unmarshal([]byte(str), &obj); err != nil {
						return nil, fmt.Errorf(".default : cannot decode json: %s", err.Error())
					}
					res.Default = obj
				} else if res.IsArray {
					res.Default = []interface{}{value}
				} else {
					res.Default = value
				}
			default:
				return nil, errors.New(".default : must be a number, a string or a boolean")
			}
		} else if !res.Required {
			// Default 'Default'
			if res.IsArray {
				res.Default = []interface{}{}
			} else {
				switch res.Test {
				case "$bool", "$json":
					res.Default = false
				case "$int", "$uint", "$float", "$number":
					res.Default = 0
				default:
					res.Default = ""
				}
			}
		}

		return res, nil
	}
	return nil, errors.New(" : must be an object")
}

// filterArray removes duplicates and empty strings from the array.
func filterArray(array []string) []string {
	result := []string{}

	for i, v := range array {
		// Non empty string
		if len(v) == 0 {
			continue
		}

		// Scan slice for a previous element of the same value
		exists := false
		for j := 0; j < i; j++ {
			if array[j] == v {
				exists = true
				break
			}
		}

		// If no previous element exists, append this one
		if !exists {
			result = append(result, v)
		}
	}
	return result
}

// areNamesUnique returns true if all names from args/flags are unique.
func areNamesUnique(args []*ArgOrFlag, flags []*ArgOrFlag) (string, bool) {
	encountered := map[string]bool{}

	for _, arg := range args {
		for _, name := range arg.Name {
			if _, ok := encountered[name]; ok {
				return name, false
			}
			encountered[name] = true
		}
	}
	for _, flag := range flags {
		for _, name := range flag.Name {
			if _, ok := encountered[name]; ok {
				return name, false
			}
			encountered[name] = true
		}
	}

	return "", true
}

// checkArgOrder returns true if no "required = false" are before a "required = true".
func checkArgOrder(args []*ArgOrFlag) (string, bool) {
	lastRequired := true
	for _, arg := range args {
		if arg.Required && !lastRequired {
			return arg.Name[0], false
		}
		lastRequired = arg.Required
	}
	return "", true
}

// checkArgIsArray returns true if no argument has "is_array = true", except last one
func checkArgIsArray(args []*ArgOrFlag) (string, bool) {
	lastIndex := len(args) - 1
	for index, arg := range args {
		if index != lastIndex && arg.IsArray {
			return arg.Name[0], false
		}
	}
	return "", true
}

// fixYamlMaps convert map[interface{}] to map[string].
func fixYamlMaps(data interface{}) interface{} {
	if imap, ok := data.(map[interface{}]interface{}); ok {
		res := map[string]interface{}{}
		for key, val := range imap {
			res[fmt.Sprint(key)] = fixYamlMaps(val)
		}
		return res
	}
	if array, ok := data.([]interface{}); ok {
		for i, val := range array {
			array[i] = fixYamlMaps(val)
		}
		return array
	}
	return data
}
