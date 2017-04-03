package main

import (
	"os"
	"sort"
)

// Sample config
func sampleConfig() []*Group {
	// Open the file
	reader, err := os.Open("../.wtfcmd.yaml")
	if err != nil {
		Warn("can't open .wtfcmd.yaml : ", err)
	}

	// Parse the config
	var cfgs []*Config
	cfgs, err = ParseConfigs(reader, cfgs, ".wtfcmd.yaml", "yaml")
	if err != nil {
		Panic("Error parsing ", ".wtfcmd.yaml : ", err)
	}

	return BuildHierarchy(cfgs)
}

func ExampleAutocomplete1() {
	groups := sampleConfig()

	res := autocomplete(groups, "act ", []string{"act"}, 4)
	sort.Strings(res)
	Made("autocomplete", res)

	res = autocomplete(groups, "act te", []string{"act", "te"}, 6)
	sort.Strings(res)
	Made("autocomplete", res)

	res = autocomplete(groups, "act forat", []string{"act", "forat"}, 7)
	sort.Strings(res)
	Made("autocomplete", res)

	res = autocomplete(groups, "act forat omg", []string{"act", "forat", "omg"}, 7)
	sort.Strings(res)
	Made("autocomplete", res)

	res = autocomplete(groups, "act      forat      omg", []string{"act", "forat", "omg"}, 11)
	sort.Strings(res)
	Made("autocomplete", res)

	res = autocomplete(groups, "act      \"forat\"      omg", []string{"act", "\"forat\"", "omg"}, 13)
	sort.Strings(res)
	Made("autocomplete", res)

	res = autocomplete(groups, "act  omg", []string{"act", "", "omg"}, 4)
	sort.Strings(res)
	Made("autocomplete", res)

	// Output: [+] autocomplete [build format test]
	// [+] autocomplete [test]
	// [+] autocomplete [form]
	// [+] autocomplete [form]
	// [+] autocomplete [format]
	// [+] autocomplete []
	// [+] autocomplete [build format test]
}
