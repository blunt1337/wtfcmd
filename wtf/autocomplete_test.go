package main

import (
	"os"
)

// Sample config
func sampleConfig() []*Group {
	// Open the file
	reader, err := os.Open("../.cmds.yaml")
	if err != nil {
		Warn("can't open .cmds.yaml : ", err)
	}

	// Parse the config
	var cfgs []*Config
	cfgs, err = ParseConfigs(reader, cfgs, ".cmds.yaml", "yaml")
	if err != nil {
		Panic("Error parsing ", ".cmds.yaml : ", err)
	}

	return BuildHierarchy(cfgs)
}

func Example1() {
	groups := sampleConfig()

	Made("autocomplete", autocomplete(groups, "act ", []string{"act"}, 4))
	Made("autocomplete", autocomplete(groups, "act te", []string{"act", "te"}, 6))
	Made("autocomplete", autocomplete(groups, "act forat", []string{"act", "forat"}, 7))
	Made("autocomplete", autocomplete(groups, "act forat omg", []string{"act", "forat", "omg"}, 7))
	Made("autocomplete", autocomplete(groups, "act      forat      omg", []string{"act", "forat", "omg"}, 11))
	Made("autocomplete", autocomplete(groups, "act      \"forat\"      omg", []string{"act", "\"forat\"", "omg"}, 13)) //TODO
	Made("autocomplete", autocomplete(groups, "act  omg", []string{"act", "", "omg"}, 4))

	// Output: [+] autocomplete [build test format]
	// [+] autocomplete [test]
	// [+] autocomplete [form]
	// [+] autocomplete [form]
	// [+] autocomplete [format]
	// [+] autocomplete []
	// [+] autocomplete [build test format]
}
