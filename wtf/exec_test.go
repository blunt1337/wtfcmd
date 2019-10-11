package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestEscapeInt(t *testing.T) {
	echo(1337, t)
}
func TestEscapeFloat(t *testing.T) {
	echo(1234.567890, t)
}
func TestEscapeString1(t *testing.T) {
	echo("simple test first", t)
}
func TestEscapeString2(t *testing.T) {
	echo("Complexe string\n$varname %varname /><$%$5!@#$%^&*()_+=4````36`573\n'@456@7345 dsfsd??//:;\"\\\"\\'''dsfsdf", t)
}
func TestEscape(t *testing.T) {
	value := "Complexe string\n$varname %varname /><$%$5!@#$%^&*()_+=4````36`573456@7345 dsfsd??//:;\"\\\"\\'''dsfsdf"
	res := UnescapeArg(EscapeArg(value))
	if value != res {
		t.Errorf("failed to unescape escaped string: %s\n%s", value, res)
	}
}

// Run the echo command.
func echo(value interface{}, t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		// Run the command
		params := map[string]interface{}{
			"test": value,
		}
		cmdTpl := &TermDependant{
			Bash:       "function _test {\necho \"$1\"\n}\n_test {{esc .test}}",
			Powershell: "function _test { echo $args[0] } _test {{esc .test}}",
		}
		cwd := &TermDependant{
			Bash:       "",
			Powershell: "",
		}
		config := &Config{"/", []string{}, []string{}, cmdTpl, "", []*ArgOrFlag{}, []*ArgOrFlag{}, cwd}
		command := &Command{"command", []string{"cmd"}, config}
		group := &Group{"group", []string{"g"}, []*Command{command}}

		ExecCmd(group, command, params, false)
		return
	}

	// Run the test as sub process to catch the exit
	cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		t.Errorf("exit code %v\n%s", err, output.String())
	}

	// Check output
	if fmt.Sprintf("%v", value) != strings.TrimRight(output.String(), "\n\r") {
		t.Errorf("wrong output\nwanted: %v\ngot: %s", value, output.String())
	}
}
