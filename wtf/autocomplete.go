package main

import (
	"blunt.sh/wtfcmd/assets"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// version is the version of autocomplete scripts
var version = "1"

// resultSeparator is the separator to split results
var resultSeparator = "\n"

// AutocompleteResult are results for autocompletes.
type AutocompleteResult struct {
	score   int
	start   string
	end     string
	Results []string
}

// checkAndAdd will add the word to the result if matching.
// If a result match both start and end, results with only start matching will be removed.
func (me *AutocompleteResult) checkAndAdd(word string) {
	// Made("tested: " + word + ", with: " + me.start + "|" + me.end)
	if strings.HasPrefix(word, me.start) {
		// Match both start + end
		if strings.HasSuffix(word, me.end) {
			res := word[0 : len(word)-len(me.end)]

			if me.score < 2 {
				me.score = 2
				me.Results = []string{res}
			} else {
				me.Results = append(me.Results, res)
			}
			return
		}

		// Match start only
		if me.score < 2 {
			me.Results = append(me.Results, word)
		}
	}
}

// AutocompleteCall will parse arguments into cursor position, cmdline, words...
func AutocompleteCall(groups []*Group, args []string) {
	if len(args) == 1 {
		switch args[0] {
		case "setup":
			setupAutocomplete()
		case "install":
			installAutocomplete()
		case "uninstall":
			uninstallAutocomplete()
		}
		os.Exit(0)
	}

	if len(args) < 3 {
		Panic("missing parameters")
	}

	// Parse cursor position
	cursorPosition, err := strconv.Atoi(args[0])
	if err != nil {
		Panic("missing cursor position", err)
	}

	res := autocomplete(groups, args[1], args[2:], cursorPosition)

	// Print
	fmt.Print(strings.Join(res, resultSeparator))
	os.Exit(0)
}

// autocomplete returns an array of words to aucomplete the search.
func autocomplete(groups []*Group, cmdline string, words []string, cursorPosition int) []string {
	// Remove empty words
	var realwords []string
	for _, word := range words {
		if word != "" {
			realwords = append(realwords, word)
		}
	}
	words = realwords

	// Find the word index and starting index of our word
	wordPosition := 0
	wordIndex := -1
	lastEnd := 0
	for i, word := range words {
		start := lastEnd + strings.Index(cmdline[lastEnd:], word)
		end := start + len(word)

		if start <= cursorPosition && cursorPosition <= end {
			wordIndex = i
			wordPosition = cursorPosition - start
			break
		}

		lastEnd = end
	}

	// No word
	if wordIndex == -1 {
		wordIndex = len(words)
		words = append(words, "")
	}

	// Start and end of the word
	word := words[wordIndex]
	start := word[0:wordPosition]
	end := word[wordPosition:]

	// Autocomplete the command, wtf no.
	if wordIndex == 0 {
		return []string{}
	}

	// Autocomplete command/group
	if wordIndex == 1 {
		return autocompleteCommandOrGroup(groups, start, end)
	}

	// Find what the first argument is
	group, command := findGroupAndCommand(groups, words[1:wordIndex])
	if group == nil {
		return []string{}
	}

	// Find the command name
	if command == nil {
		return autocompleteCommand(group, start, end)
	}

	// Autocomplete flags
	if strings.HasPrefix(start, "--") {
		return autocompleteFlags(command.Config.Flags, start, end)
	}

	return []string{}
}

// autocompleteCommandOrGroup autocompletes commands and groups.
func autocompleteCommandOrGroup(groups []*Group, start string, end string) []string {
	results := AutocompleteResult{start: start, end: end}

	// Find a group
	for _, group := range groups {
		if group.Name == "" {
			// Find a command
			for _, command := range group.Commands {
				results.checkAndAdd(command.Name)
			}
		} else {
			results.checkAndAdd(group.Name)
		}
	}

	return results.Results
}

// autocompleteCommand autocompletes commands.
func autocompleteCommand(group *Group, start string, end string) []string {
	results := AutocompleteResult{start: start, end: end}

	// Find a command
	for _, command := range group.Commands {
		results.checkAndAdd(command.Name)
	}

	return results.Results
}

// autocompleteFlags autocompletes flags.
func autocompleteFlags(flags []*ArgOrFlag, start string, end string) []string {
	results := AutocompleteResult{start: start, end: end}

	// Find a flag
	for _, flag := range flags {
		results.checkAndAdd("--" + flag.Name[0])
	}

	return results.Results
}

// setupAutocomplete print the command to install the autocomplete.
func setupAutocomplete() {
	var cmd string
	var separator string

	// Bash
	switch GetTerminal() {
	case TermBash:
		cmd = assets.Get("autocomplete.sh")
		separator = resultSeparator
	case TermPowershell:
		cmd = assets.Get("autocomplete.ps1")
		cmd = strings.Replace(cmd, "\r", "", -1)
		cmd = strings.Replace(cmd, "\n", "\r\n", -1)
		separator = strings.Replace(resultSeparator, "\n", "\\n", -1)
	default:
		Panic("No autocomplete for this terminal")
	}

	// Command path/name
	cmdpath, cmdname := getCmdNameAndPath()

	// Replace variables
	cmd = strings.Replace(cmd, "CMDPATH", cmdpath, -1)
	cmd = strings.Replace(cmd, "CMDNAME", cmdname, -1)
	cmd = strings.Replace(cmd, "SEPARATOR", separator, -1)

	fmt.Print(cmd)
	os.Exit(0)
}

// installAutocomplete installs the command `wtf --autocomplete setup` at bash/powershell startup.
func installAutocomplete() {
	// Command name
	_, cmdname := getCmdNameAndPath()

	switch GetTerminal() {
	case TermCmd:
		Panic("There is no autocomplete inside cmd.exe")
	case TermPowershell:
		// Find $PROFILE
		bytes, err := exec.Command("powershell.exe", "-command", "echo $PROFILE").Output()
		if err != nil {
			Panic("Your $PROFILE variable was not found :o", err)
		}
		profile := strings.TrimSpace(string(bytes))

		// Prompt
		if !AskYN("It will add a script into $PROFILE and enable execution of powershell scripts.\n    Continue?", true) {
			Panic("Then we can't install the autocomplete.")
		}

		// Enable execution of scripts
		exec.Command("powershell.exe", "-command", "Set-ExecutionPolicy -Scope CurrentUser -Force RemoteSigned").Run()

		// Content to add to $profile
		script := strings.Replace(""+
			"# CMDNAME autocomplete\r\n"+
			"$CMDNAME_autcomplete_path = \"$env:temp\\CMDNAME_autocomplete_"+version+".ps1\"\r\n"+
			"if (!(Test-Path $CMDNAME_autcomplete_path)) {\r\n"+
			"    CMDNAME --autocomplete setup > $CMDNAME_autcomplete_path\r\n"+
			"}\r\n"+
			". $CMDNAME_autcomplete_path\r\n"+
			"# CMDNAME autocomplete end\r\n",
			"CMDNAME", cmdname, -1)

		// Open $profile
		data, err := ioutil.ReadFile(profile)
		if err != nil {
			if os.IsNotExist(err) {
				// Create sub folders
				os.MkdirAll(filepath.Dir(profile), 0777)

				// Create the file
				err = ioutil.WriteFile(profile, []byte(script), 0777)
				if err != nil {
					Panic("Cannot write the $PROFILE script", err)
				}
			} else {
				Panic(err)
			}
		} else {
			content := string(data)

			// Remove old code
			content = regexp.MustCompile(strings.Replace("(?s)(\\n|\\r)*"+ // All space before
				"# CMDNAME autocomplete\\r?\\n"+ // Start comment
				".*?"+ // Everything between
				"# CMDNAME autocomplete end", // End comment
				"CMDNAME", cmdname, -1)).ReplaceAllString(content, "")

			// Add our script
			content += "\n" + script

			// Rewrite the file
			err = ioutil.WriteFile(profile, []byte(content), 0777)
			if err != nil {
				Panic("Cannot rewrite the $PROFILE script", err)
			}
		}
	case TermBash:
		// Find ~/.bash_profile
		user, err := user.Current()
		if err != nil {
			Panic(err)
		}
		profile := user.HomeDir + "/.bash_profile"

		// Prompt
		if !AskYN("It will add a line into your .bash_profile.\n    Continue?", true) {
			Panic("Then we can't install the autocomplete.")
		}

		// Content to add to .bash_profile
		script := strings.Replace(""+
			"# CMDNAME autocomplete\n"+
			"CMDNAME_autcomplete_path=\"$TMPDIR/CMDNAME_autocomplete_"+version+".sh\"\n"+
			"if [ ! -f \"$CMDNAME_autcomplete_path\" ]; then\n"+
			"    CMDNAME --autocomplete setup > \"$CMDNAME_autcomplete_path\"\n"+
			"fi\n"+
			"source \"$CMDNAME_autcomplete_path\"\n"+
			"# CMDNAME autocomplete end\n",
			"CMDNAME", cmdname, -1)

		// Open ~/.bash_profile
		data, err := ioutil.ReadFile(profile)
		if err != nil {
			if os.IsNotExist(err) {
				// Create the file
				err = ioutil.WriteFile(profile, []byte(script), 0644)
				if err != nil {
					Panic("Cannot write the .bash_profile script", err)
				}
			} else {
				Panic(err)
			}
		} else {
			content := string(data)

			// Remove old code
			content = regexp.MustCompile(strings.Replace("(?s)(\\n|\\r)*"+ // All space before
				"# CMDNAME autocomplete\\r?\\n"+ // Start comment
				".*?"+ // Everything between
				"# CMDNAME autocomplete end", // End comment
				"CMDNAME", cmdname, -1)).ReplaceAllString(content, "")

			// Add our script
			content += "\n" + script

			// Rewrite the file
			err = ioutil.WriteFile(profile, []byte(content), 0644)
			if err != nil {
				Panic("Cannot rewrite the .bash_profile script", err)
			}
		}
	}
	Made(cmdname + " autocomplete installed :)")
}

// uninstallAutocomplete removes the command `wtf --autocomplete setup` from startup.
func uninstallAutocomplete() {
	// Command name
	_, cmdname := getCmdNameAndPath()

	switch GetTerminal() {
	case TermCmd:
		Panic("There is no autocomplete inside cmd.exe")
	case TermPowershell:
		// Find $PROFILE
		bytes, err := exec.Command("powershell.exe", "-command", "echo $PROFILE").Output()
		if err != nil {
			Panic("Your $PROFILE variable was not found :o", err)
		}
		profile := strings.TrimSpace(string(bytes))

		// Open $profile
		data, err := ioutil.ReadFile(profile)
		if err != nil {
			if os.IsNotExist(err) {
				Made("No autocomplete installed")
				return
			} else {
				Panic(err)
			}
		} else {
			content := string(data)

			// Remove old code
			content = regexp.MustCompile(strings.Replace("(?s)(\\n|\\r)*"+ // All space before
				"# CMDNAME autocomplete\\r?\\n"+ // Start comment
				".*?"+ // Everything between
				"# CMDNAME autocomplete end", // End comment
				"CMDNAME", cmdname, -1)).ReplaceAllString(content, "")

			// Rewrite the file
			err = ioutil.WriteFile(profile, []byte(content), 0777)
			if err != nil {
				Panic("Cannot rewrite the $PROFILE script", err)
			}
		}
	case TermBash:
		// Find ~/.bash_profile
		user, err := user.Current()
		if err != nil {
			Panic(err)
		}
		profile := user.HomeDir + "/.bash_profile"

		// Open ~/.bash_profile
		data, err := ioutil.ReadFile(profile)
		if err != nil {
			if os.IsNotExist(err) {
				Made("No autocomplete installed")
				return
			} else {
				Panic(err)
			}
		} else {
			content := string(data)

			// Remove old code
			content = regexp.MustCompile(strings.Replace("(?s)(\\n|\\r)*"+ // All space before
				"# CMDNAME autocomplete\\r?\\n"+ // Start comment
				".*?"+ // Everything between
				"# CMDNAME autocomplete end", // End comment
				"CMDNAME", cmdname, -1)).ReplaceAllString(content, "")

			// Rewrite the file
			err = ioutil.WriteFile(profile, []byte(content), 0644)
			if err != nil {
				Panic("Cannot rewrite the .bash_profile script", err)
			}
		}
	}
	Made(cmdname + " autocomplete uninstalled :)")
}

// getCmdNameAndPath returns the path and the name of the command. In case they renamed it.
func getCmdNameAndPath() (string, string) {
	// Command path/name
	cmdname := os.Args[0]
	cmdpath, err := os.Executable()
	if err != nil {
		cmdpath = cmdname
	}

	// Windows path
	if runtime.GOOS == "windows" {
		cmdname = cmdname[strings.LastIndexByte(cmdname, '\\')+1:]
		index := strings.LastIndex(cmdname, ".exe")
		if index > 0 {
			cmdname = cmdname[0:index]
		}
	}

	return cmdpath, cmdname
}
