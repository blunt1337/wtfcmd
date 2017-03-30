package main

import (
	"os"
)

// Disable colors
func init() {
	noColor = true
}

func ExampleError() {
	os.Stderr = os.Stdout

	Error("super", "cool", "error", "message", 1337)
	// Output: [x] super cool error message 1337
}

func ExampleWarn() {
	os.Stderr = os.Stdout

	Warn("super", "cool", "warning", "message", 1337)
	// Output: [-] super cool warning message 1337
}

func ExampleInfo() {
	Info("super", "information", 1337, "!")
	// Output: [>] super information 1337 !
}

func ExampleMade() {
	Made("you parsed", 1337, "files !")
	// Output: [+] you parsed 1337 files !
}

func ExampleAsk() {
	// Fake input
	r, w, _ := os.Pipe()
	os.Stdin = r

	w.WriteString("Olivier\n")
	name := Ask("What's your name?")
	Info("You told me your name was", name)

	// Output: [?] What's your name?
	// [>] You told me your name was Olivier
}

func ExampleAskYN() {
	// Fake input
	r, w, _ := os.Pipe()
	os.Stdin = r
	os.Stderr = os.Stdout

	// Yes
	w.WriteString("yes\n")
	if AskYN("Are lizards small?") {
		Info("I agree!")
	}

	// No
	w.WriteString("no\n")
	if !AskYN("Can you stretch a rock?") {
		Info("+1")
	}

	// Ask again if not yes/y/no/n
	w.WriteString("maybe\n")
	w.WriteString("y\n")
	if AskYN("Can you grow a mustache on your foot?") {
		Info("x)")
	}

	// Default answer
	w.WriteString("\n")
	if AskYN("Is it sunny today?", true) {
		Info(":)")
	}

	// Not the default answer
	w.WriteString("n\n")
	if !AskYN("Do you speak english?", true) {
		Info("Why am i even talking to you!")
	}

	// Output: [?] Are lizards small? [yes|no]
	// [>] I agree!
	// [?] Can you stretch a rock? [yes|no]
	// [>] +1
	// [?] Can you grow a mustache on your foot? [yes|no]
	// [-] Please answer with "yes" or "no"
	// [>] x)
	// [?] Is it sunny today? [yes (default)|no]
	// [>] :)
	// [?] Do you speak english? [yes (default)|no]
	// [>] Why am i even talking to you!
}

func ExampleAskList() {
	// Fake input
	r, w, _ := os.Pipe()
	os.Stdin = r
	os.Stderr = os.Stdout

	choices := []string{
		"Europe",
		"America",
		"China :o",
		"Out of space",
		"Other",
	}

	// Simple choices
	w.WriteString("3\n")
	index := AskList("Where do you come from?", choices)
	Info("Oh so you're from", choices[index])

	// Some fails
	w.WriteString("\n")
	w.WriteString("12\n")
	w.WriteString("2\n")
	index = AskList("Where do you come from?", choices)
	Info("Oh so you're from", choices[index])

	// Default value
	w.WriteString("\n")
	index = AskList("Where do you come from?", choices, 1)
	Info("Oh so you're from", choices[index])

	// Output: [?] Where do you come from?
	//     [0] Europe
	//     [1] America
	//     [2] China :o
	//     [3] Out of space
	//     [4] Other
	// [>] Oh so you're from Out of space
	// [?] Where do you come from?
	//     [0] Europe
	//     [1] America
	//     [2] China :o
	//     [3] Out of space
	//     [4] Other
	// [-] Please answer with a number between 0 and 5
	// [-] Please answer with a number between 0 and 5
	// [>] Oh so you're from China :o
	// [?] Where do you come from?
	//     [0] Europe
	//     [1] America (default)
	//     [2] China :o
	//     [3] Out of space
	//     [4] Other
	// [>] Oh so you're from America
}

/*
Go test doesn't like our \r

func ExampleLoading() {
	Loading("Loading...")
	fmt.Println() // Prevent loadbar to override last one
	Loading("Loading (", 1, "/", 5, " items)", 1 / 5.0)
	fmt.Println() // Prevent loadbar to override last one
	Loading("Loading (", 2, "/", 5, " items)", 2 / 5.0)
	fmt.Println() // Prevent loadbar to override last one
	Loading("Loading (", 3, "/", 5, " items)", 3 / 5.0)
	fmt.Println() // Prevent loadbar to override last one
	Loading("Loading (", 4, "/", 5, " items)", 4 / 5.0)
	fmt.Println() // Prevent loadbar to override last one
	Loading("Loading (", 5, "/", 5, " items)", 5 / 5.0)

	// Output: [-] Loading...                                                                  [  0%]
	// [\] Loading (1/5 items)                                                         [ 20%]
	// [|] Loading (2/5 items)                                                         [ 40%]
	// [/] Loading (3/5 items)                                                         [ 60%]
	// [-] Loading (4/5 items)                                                         [ 80%]
	// [+] Loading (5/5 items)                                                         [100%]
}*/
