package main

import (
	"os"
	"path/filepath"
	"strings"
)

// findConfigFiles finds all .wtfcmd.[ext] from current directory and all parents.
func findConfigFiles() []string {
	res := make([]string, 0)

	// Current working dir
	dir, err := os.Getwd()
	if err != nil {
		return res
	}

	maxDepth := 10
	for {
		// Check for all extensions
		exts := []string{"json", "jsonc", "json5", "yaml", "yml"}
		for _, ext := range exts {
			path := dir + "/.wtfcmd." + ext
			if _, err := os.Stat(path); err == nil {
				res = append(res, path)
			}
		}

		// Parent dir
		maxDepth--
		if maxDepth == 0 {
			break
		}
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

		cfgs, err = ParseConfigs(reader, cfgs, file, file[strings.LastIndex(file, ".")+1:])
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
