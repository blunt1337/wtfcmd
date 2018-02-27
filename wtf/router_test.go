package main

import (
	"fmt"
	"sort"
	"strings"
)

func ExampleRouter1() {
	// Sample config
	group := &Group{}
	command := &Command{
		Config: &Config{
			Args: []*ArgOrFlag{
				{
					Name:     []string{"a_req"},
					Required: true,
					IsArray:  false,
				},
				{
					Name:     []string{"b_not_req"},
					Required: false,
					IsArray:  false,
					Default:  "",
				},
			},
			Flags: []*ArgOrFlag{
				{
					Name:     []string{"cflag", "c"},
					Required: false,
					Default:  1337,
					Test:     "$int",
					IsArray:  false,
				},
				{
					Name:     []string{"dflag", "d"},
					Required: false,
					Default:  true,
					Test:     "$bool",
					IsArray:  false,
				},
				{
					Name:     []string{"dflag2", "D"},
					Required: false,
					Default:  false,
					Test:     "$bool",
					IsArray:  false,
				},
				{
					Name:     []string{"eflag", "e"},
					Required: false,
					Default:  []string{"e"},
					IsArray:  true,
				},
			},
		},
	}

	echoParams(group, command, "param1 --cflag 123 --dflag --eflag element")
	echoParams(group, command, "param1 param2 -c 123 -d -e element0 -e element1")
	echoParams(group, command, "param1 --cflag=123 --dflag=1 --eflag=element")
	echoParams(group, command, "param1 param2 -c=123 -d=0 -e=element0 -e=element1")
	echoParams(group, command, "param1 -d -D")
	echoParams(group, command, "param1 -dD")
	echoParams(group, command, "param1 -dD=0")
	echoParams(group, command, "param1 -cd=1")

	// Output: [+] parse params a_req=param1, b_not_req=, cflag=123, dflag=true, dflag2=false, eflag=[element]
	// [+] parse params a_req=param1, b_not_req=param2, cflag=123, dflag=true, dflag2=false, eflag=[element0 element1]
	// [+] parse params a_req=param1, b_not_req=, cflag=123, dflag=true, dflag2=false, eflag=[element]
	// [+] parse params a_req=param1, b_not_req=param2, cflag=123, dflag=false, dflag2=false, eflag=[element0 element1]
	// [+] parse params a_req=param1, b_not_req=, cflag=1337, dflag=true, dflag2=true, eflag=[e]
	// [+] parse params a_req=param1, b_not_req=, cflag=1337, dflag=true, dflag2=true, eflag=[e]
	// [+] parse params a_req=param1, b_not_req=, cflag=1337, dflag=false, dflag2=false, eflag=[e]
	// [+] parse params a_req=param1, b_not_req=, cflag=1, dflag=true, dflag2=false, eflag=[e]
}

func echoParams(group *Group, command *Command, args string) {
	res := parseParams(group, command, strings.Split(args, " "))

	// Build sorted keys
	var keys []string
	for key := range res {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Print
	str := ""
	for _, k := range keys {
		str += ", " + k + "=" + fmt.Sprint(res[k])
	}
	Made("parse params", str[2:])
}
