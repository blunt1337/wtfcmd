package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// noColor means that the terminal doesn't handle color codes
var noColor = false

// loadingStep is a static state of the loadbar loader
var loadingStep = 0

func init() {
	noColor = os.Getenv("TERM") == "dumb" || GetTerminal() == TermCmd
}

// Error prints an error.
func Error(msg ...interface{}) {
	if noColor {
		fmt.Fprint(os.Stderr, "[x] ")
	} else {
		fmt.Fprint(os.Stderr, "\033[91;01m[x]\033[00m ")
	}
	fmt.Fprintln(os.Stderr, msg...)
}

// Panic prints an error then exit with status 1.
func Panic(msg ...interface{}) {
	Error(msg...)
	os.Exit(1)
}

// Warn prints a warning.
func Warn(msg ...interface{}) {
	if noColor {
		fmt.Fprint(os.Stderr, "[-] ")
	} else {
		fmt.Fprint(os.Stderr, "\033[38;5;208;01m[-]\033[00m ")
	}
	fmt.Fprintln(os.Stderr, msg...)
}

// Info prints a message.
func Info(msg ...interface{}) {
	if noColor {
		fmt.Print("[>] ")
	} else {
		fmt.Print("\033[38;5;85m[>]\033[00m ")
	}
	fmt.Println(msg...)
}

// Made prints a message for an action made.
func Made(msg ...interface{}) {
	if noColor {
		fmt.Print("[+] ")
	} else {
		fmt.Print("\033[38;5;70;01m[+]\033[00m ")
	}
	fmt.Println(msg...)
}

// Ask prints a question and return it's answer.
func Ask(question ...interface{}) string {
	if noColor {
		fmt.Print("[?] ")
	} else {
		fmt.Print("\033[38;5;99;01m[?]\033[00m ")
	}
	fmt.Println(question...)

	// Read response
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return text
}

// ReadSecure prints a question and return it's answer, hidden when typing.
func ReadSecure() string {
	// Command to read from the term
	cmdWrapper, cmd := GetLangAndCommandTemplate(&TermDependant{
		Bash: "read -s password; echo \"$password\"",
		Powershell: "$password = Read-Host -AsSecureString; " +
			"$password = [System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($password); " +
			"$password = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto($password); " +
			"echo $password",
	})

	// Read a secure response
	process := exec.Command(cmdWrapper[0], cmdWrapper[1], cmd)

	// Pipes
	var output bytes.Buffer
	process.Stdout = &output
	process.Stdin = os.Stdin

	// Run
	if err := process.Run(); err != nil {
		Panic("Failed to read secure password")
	}
	return output.String()
}

// AskYN prints a yes/no question and return it's answer.
// If the last parameter is either "yes"/"y"/true or "no"/"n"/false, then it is treated as the default answer.
// So parameter are "question0 question1 ... optional_default_value.
func AskYN(questionAndFefault ...interface{}) bool {
	var hasDefault bool
	var dflt bool

	// Get the last param if it's the default value
	last := len(questionAndFefault) - 1
	switch questionAndFefault[last] {
	case "yes", "true", true:
		hasDefault = true
		dflt = true
		questionAndFefault = questionAndFefault[:last]
	case "no", "false", false:
		hasDefault = true
		dflt = false
		questionAndFefault = questionAndFefault[:last]
	}

	// Show the question
	if noColor {
		fmt.Print("[?] ")
	} else {
		fmt.Print("\033[38;5;99;01m[?]\033[00m ")
	}
	fmt.Print(questionAndFefault...)

	// Print the default
	if hasDefault {
		if dflt {
			if noColor {
				fmt.Print(" [yes (default)|no]\n")
			} else {
				fmt.Print(" [yes \033[38;5;99m(default)\033[00m|no]\n")
			}
		} else {
			if noColor {
				fmt.Print(" [yes|no (default)]\n")
			} else {
				fmt.Print(" [yes|no \033[38;5;99m(default)\033[00m]\n")
			}
		}
	} else {
		fmt.Print(" [yes|no]\n")
	}

	// Read response
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		text = strings.ToLower(text)

		switch text {
		case "":
			if hasDefault {
				return dflt
			}
		case "yes", "y":
			return true
		case "no", "n":
			return false
		}
		Warn("Please answer with \"yes\" or \"no\"")
	}
}

// AskList prints a question and a list of answers, and return it's index.
// If the last or before last parameter must be a string[], it is treated as the answer list.
// If the last parameter is an integer, then it is treated as the default answer index.
// So parameter are "question0 question1 ... answers_array optional_default_index.
func AskList(questionAndAnswersAndDefault ...interface{}) int {
	var hasDefault bool
	var dflt int
	var answers []string

	// Get the last param if it's the default value
	last := len(questionAndAnswersAndDefault) - 1
	if i, ok := questionAndAnswersAndDefault[last].(int); ok {
		hasDefault = true
		dflt = i
		questionAndAnswersAndDefault = questionAndAnswersAndDefault[:last]
		last--
	}

	// List of answers
	if arr, ok := questionAndAnswersAndDefault[last].([]string); ok && len(arr) != 0 {
		answers = arr
		questionAndAnswersAndDefault = questionAndAnswersAndDefault[:last]
	} else {
		Warn("No choices available")
		return -1
	}

	// Show the question
	if noColor {
		fmt.Print("[?] ")
	} else {
		fmt.Print("\033[38;5;99;01m[?]\033[00m ")
	}
	fmt.Println(questionAndAnswersAndDefault...)

	// Check for default
	max := len(answers)
	hasDefault = hasDefault && dflt >= 0 && dflt < max

	// Print answers
	for i, res := range answers {
		fmt.Print("    [", i, "] ", res)
		if hasDefault && i == dflt {
			if noColor {
				fmt.Print(" (default)")
			} else {
				fmt.Print(" \033[38;5;99m(default)\033[00m")
			}
		}
		fmt.Println()
	}

	// Read response
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		text = strings.ToLower(text)

		if text == "" && hasDefault == true {
			return dflt
		}

		i, err := strconv.Atoi(text)
		if err == nil && i >= 0 && i < max {
			return i
		}
		Warn("Please answer with a number between 0 and", max)
	}
}

// Bell ring a bell in the terminal/
func Bell() {
	if !noColor {
		fmt.Print("\007")
		time.Sleep(100)
		fmt.Print("\007")
	}
}

// Jsonp prints every argument as pretty json.
func Jsonp(objs ...interface{}) {
	for _, obj := range objs {
		bytes, err := json.MarshalIndent(obj, "", "	")
		if err != nil {
			Panic(err)
		}
		Info(string(bytes[:]))
	}
}

// Loading prints a loading bar with a small message.
// Last parameter must be: 0 to 1.
func Loading(msgAndPercent ...interface{}) {
	var percent float32

	// Get percent parameter
	last := len(msgAndPercent) - 1
	switch value := msgAndPercent[last].(type) {
	case float64:
		percent = float32(value)
		msgAndPercent = msgAndPercent[:last]
	case float32:
		percent = value
		msgAndPercent = msgAndPercent[:last]
	case int:
		percent = float32(value)
		msgAndPercent = msgAndPercent[:last]
	case string:
		if f, err := strconv.ParseFloat(value, 32); err == nil && f >= 0 && f <= 1 {
			percent = float32(f)
			msgAndPercent = msgAndPercent[:last]
		}
	}

	// Fix the message to 75 chars
	msg := fmt.Sprint(msgAndPercent...)
	msg = regexp.MustCompile("[\r\n]+").ReplaceAllString(msg, " ")

	if len(msg) > 75 {
		msg = msg[0:74] + "…"
	} else {
		for len(msg) < 75 {
			msg += " "
		}
	}

	// Finished
	if percent >= 1 {
		if noColor {
			fmt.Print("\r[+] ")
		} else {
			fmt.Print("\r\033[38;5;70;01m[+]\033[00m ")
		}
		fmt.Print(msg + " [100%]\n")
		loadingStep = 0
	} else {
		if noColor {
			fmt.Print("\r[")
		} else {
			fmt.Print("\r\033[01m[")
		}

		switch loadingStep {
		case 0:
			fmt.Print("-")
		case 1:
			fmt.Print("\\")
		case 2:
			fmt.Print("|")
		case 3:
			fmt.Print("/")
			loadingStep = -1
		}
		loadingStep++

		// Percent
		if noColor {
			fmt.Printf("] %s [%3d%%]", msg, int(percent*100))
		} else {
			fmt.Printf("]\033[00m %s [%3d%%]", msg, int(percent*100))
		}
	}
}
