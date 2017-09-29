package main

import (
	"os"
	"path/filepath"
)

// findConfigFiles finds all .cmd.json from current directory and all parents.
func findConfigFiles() []string {
	res := make([]string, 0)

	// Current working dir
	dir, err := os.Getwd()
	if err != nil {
		return res
	}

	for {
		// Check if a .wtfcmd.json exists
		path := dir + "/.wtfcmd.json"
		if _, err := os.Stat(path); err == nil {
			res = append(res, path)
		}

		// Check if a .wtfcmd.yaml exists
		path = dir + "/.wtfcmd.yaml"
		if _, err := os.Stat(path); err == nil {
			res = append(res, path)
		}

		// Parent dir
		olddir := dir
		dir = filepath.Dir(dir)
		if olddir == dir {
			break
		}
	}

	return res
}

// main is the starting point.
func main() {
	files := findConfigFiles()

	// Parse all configs
	var cfgs []*Config
	for _, file := range files {
		// Open the file
		reader, err := os.Open(file)
		if err != nil {
			Warn("can't open", file, ":", err)
		}

		cfgs, err = ParseConfigs(reader, cfgs, file, file[len(file)-4:])
		if err != nil {
			reader.Close()
			Panic("Error parsing", file, ":", err)
		}
		reader.Close()
	}

	// Remove overrided functions
	groups := BuildHierarchy(cfgs)

	// Execute the router
	Route(groups)
}
