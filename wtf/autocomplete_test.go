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

	res := autocomplete(groups, "wtf ", []string{"wtf"}, 4)
	sort.Strings(res)
	Made("autocomplete", res)

	res = autocomplete(groups, "wtf te", []string{"wtf", "te"}, 6)
	sort.Strings(res)
	Made("autocomplete", res)

	res = autocomplete(groups, "wtf forat", []string{"wtf", "forat"}, 7)
	sort.Strings(res)
	Made("autocomplete", res)

	res = autocomplete(groups, "wtf forat omg", []string{"wtf", "forat", "omg"}, 7)
	sort.Strings(res)
	Made("autocomplete", res)

	res = autocomplete(groups, "wtf      forat      omg", []string{"wtf", "forat", "omg"}, 11)
	sort.Strings(res)
	Made("autocomplete", res)

	res = autocomplete(groups, "wtf      \"forat\"      omg", []string{"wtf", "\"forat\"", "omg"}, 13)
	sort.Strings(res)
	Made("autocomplete", res)

	res = autocomplete(groups, "wtf  omg", []string{"wtf", "", "omg"}, 4)
	sort.Strings(res)
	Made("autocomplete", res)

	// Linux/osx send an empty string at the end
	res = autocomplete(groups, "wtf ", []string{"wtf", ""}, 4)
	sort.Strings(res)
	Made("autocomplete", res)

	// Output: [+] autocomplete [build format test]
	// [+] autocomplete [test]
	// [+] autocomplete [form]
	// [+] autocomplete [form]
	// [+] autocomplete [format]
	// [+] autocomplete []
	// [+] autocomplete [build format test]
	// [+] autocomplete [build format test]
}
