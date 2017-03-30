package main

import (
	"blunt.sh/wtfcmd/assets"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

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
	if len(args) == 1 && args[0] == "setup" {
		setupAutocomplete()
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
	res = []string{"super", "souper"}
	// Print
	fmt.Print(strings.Join(res, resultSeparator))
	os.Exit(0)
}

// autocomplete returns an array of words to aucomplete the search.
func autocomplete(groups []*Group, cmdline string, words []string, cursorPosition int) []string {
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

	// Bash
	switch GetTerminal() {
	case TermBash:
		cmd = assets.Get("autocomplete.sh")
	case TermPowershell:
		cmd = assets.Get("autocomplete.ps1")
	default:
		Panic("No autocomplete for this terminal")
	}

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

	// Replace variables
	cmd = strings.Replace(cmd, "CMDPATH", cmdpath, -1)
	cmd = strings.Replace(cmd, "CMDNAME", cmdname, -1)
	cmd = strings.Replace(cmd, "SEPARATOR", resultSeparator, -1)

	fmt.Print(cmd)
	os.Exit(0)
}
