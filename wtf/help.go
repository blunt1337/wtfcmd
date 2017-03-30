package main

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"os"
	"sort"
	"strings"
	"text/template"
)

// Bold character code.
var bold = ""

// Reset all character code.
var reset = ""

// Red character code.
var red = ""

// Fill bold/reset/red if the terminal can handle colors.
func init() {
	if os.Getenv("TERM") != "dumb" && isatty.IsTerminal(os.Stdout.Fd()) && GetTerminal() != TermCmd {
		bold = "\033[01m"
		reset = "\033[00m"
		red = "\033[91m"
	}
}

// GroupWrapper are parameters passed to the help template.
type GroupWrapper struct {
	Bold   string
	Reset  string
	Groups map[string]*Group
}

// CommandWrapper are parameters passed to the help2 template.
type CommandWrapper struct {
	Bold    string
	Reset   string
	Command *Command
	Group   *Group
}

// SortByName is a sort interface by .Name.
type SortByName []*Command

func (me SortByName) Len() int {
	return len(me)
}
func (me SortByName) Swap(i int, j int) {
	me[i], me[j] = me[j], me[i]
}
func (me SortByName) Less(i int, j int) bool {
	return me[i].Name < me[j].Name
}

// ShowHelp shows a help page with a list of all commands.
func ShowHelp(groups []*Group) {
	tplTxt := "{{$B := .Bold}}{{$R := .Reset}}" +
		"\n" +
		"{{$B}}AVAILABLE COMMANDS{{$R}}\n" +
		"\n" +
		"{{range .Groups}}" +
		/*	*/ "{{if .Name}}" +
		/*		*/ "{{$B}}{{.Name}}{{$R}}" +
		/*		*/ "{{if .Aliases}}" +
		/*			*/ " ({{join .Aliases \", \"}})" +
		/*		*/ "{{end}}" +
		/*		*/ "\n" +
		/*	*/ "{{end}}" +
		/*	*/ "{{range .Commands}}" +
		/*		*/ " - " +
		/*		*/ "{{if .Config.Group}}{{index .Config.Group 0}} {{end}}" +
		/*		*/ "{{$B}}{{.Name}}{{$R}}" +
		/*		*/ "{{if .Aliases}}" +
		/*			*/ " ({{join .Aliases \", \"}})" +
		/*		*/ "{{end}}" +
		/*		*/ "{{availability .Config.Cmd}}" +
		/*		*/ "\n" +
		/*		*/ "{{if .Config.Desc}}" +
		/*			*/ "	{{replace .Config.Desc \"\\n\" \"\\n	\" -1}}\n" +
		/*		*/ "{{end}}" +
		/*		*/ "\n" +
		/*	*/ "{{end}}" +
		/*	*/ "\n" +
		"{{end}}" +
		"\n"

	// Group by names
	groupByName := map[string]*Group{}
	for _, group := range groups {
		groupByName[group.Name] = group

		// Sort commands
		sort.Sort(SortByName(group.Commands))
	}

	// Execute template with group_by_names
	tmpl, err := template.New("help").Funcs(template.FuncMap{
		"join":         strings.Join,
		"replace":      strings.Replace,
		"availability": CmdAvailability,
	}).Parse(tplTxt)
	if err != nil {
		Panic(err)
	}
	err = tmpl.Execute(os.Stdout, GroupWrapper{
		Groups: groupByName,
		Bold:   bold,
		Reset:  reset,
	})
	if err != nil {
		Panic(err)
	}
}

// getUsageTemplate returns the template string for command usage.
func getUsageTemplate() string {
	return "{{if .Group.Name}}{{.Group.Name}} {{end}}{{.Command.Name}}" +
		"{{if .Command.Config.Flags}} [flags]{{end}}" +
		"{{range .Command.Config.Args}}" +
		/*	*/ " {{ if .Required}}<{{index .Name 0}}>" +
		/*	*/ "{{else}}[{{index .Name 0}}={{if .Default}}{{.Default}}{{else}}\"\"{{end}}]" +
		/*	*/ "{{end}}" +
		"{{end}}"
}

// ShowHelpCommand shows detailed help for a command.
func ShowHelpCommand(group *Group, command *Command) {
	tplTxt := "{{$B := .Bold}}{{$R := .Reset}}" +
		"\n" +
		"{{$B}}SYNOPSIS{{$R}}\n" +
		"	" + getUsageTemplate() + "{{availability .Command.Config.Cmd}}" +
		"\n\n" +
		"{{if or .Group.Aliases .Command.Aliases}}" +
		/*	*/ "{{$B}}ALIASES{{$R}}\n" +
		/*	*/ "	{{if .Group.Name}}{{.Group.Name}}{{if .Group.Aliases}}|{{join .Group.Aliases \"|\"}}{{end}} {{end}}" +
		/*	*/ "{{.Command.Name}}{{if .Command.Aliases}}|{{join .Command.Aliases \"|\"}}{{end}}" +
		/*	*/ "\n\n" +
		"{{end}}" +
		"{{if .Command.Config.Desc}}" +
		/*	*/ "{{$B}}DESCRIPTION{{$R}}\n" +
		/*	*/ "	{{replace .Command.Config.Desc \"\\n\" \"\\n	\" -1}}" +
		/*	*/ "\n\n" +
		"{{end}}" +
		"{{if .Command.Config.Args}}" +
		/*	*/ "{{$B}}ARGUMENTS{{$R}}\n" +
		/*	*/ "{{range .Command.Config.Args}}" +
		/*		*/ "	{{$B}}{{index .Name 0}}{{$R}} ({{if .Test}}{{.Test}}{{else}}string{{end}}, {{if .Required}}required{{else}}default {{if .Default}}{{.Default}}{{else}}\"\"{{end}}{{end}})\n\n" +
		/*		*/ "{{if .Desc}}" +
		/*			*/ "	{{replace .Desc \"\\n\" \"\\n	\" -1}}\n" +
		/*		*/ "{{end}}\n" +
		/*	*/ "{{end}}" +
		"{{end}}" +
		"{{if .Command.Config.Flags}}" +
		/*	*/ "{{$B}}FLAGS{{$R}}\n" +
		/*	*/ "{{range .Command.Config.Flags}}" +
		/*		*/ "	{{$B}}--{{index .Name 0}}{{$R}}" +
		/*		*/ "{{$l := len .Name}}{{if gt $l 1}}{{range $index, $Name := .Name}}{{if ne $index 0}}, -{{$Name}}{{end}}{{end}}{{end}}" +
		/*		*/ " ({{if .Test}}{{.Test}}{{else}}string{{end}}{{if .Default}}, default {{.Default}}{{end}})\n\n" +
		/*		*/ "{{if .Desc}}" +
		/*			*/ "	{{replace .Desc \"\\n\" \"\\n	\" -1}}\n" +
		/*		*/ "{{end}}\n" +
		/*	*/ "{{end}}" +
		"{{end}}" +
		"\n\n"

	// Execute template with command
	tmpl, err := template.New("help2").Funcs(template.FuncMap{
		"join":         strings.Join,
		"replace":      strings.Replace,
		"availability": CmdAvailability,
	}).Parse(tplTxt)
	if err != nil {
		Panic(err)
	}
	err = tmpl.Execute(os.Stdout, CommandWrapper{
		Command: command,
		Group:   group,
		Bold:    bold,
		Reset:   reset,
	})
	if err != nil {
		Panic(err)
	}
}

// ShowCommandError prints an error for a command.
func ShowCommandError(msg string, group *Group, command *Command) {
	name := os.Args[0]

	// Error
	fmt.Printf("%s: error: %s.\n\n", name, msg)

	// Usage
	if group != nil && command != nil {
		tmpl, err := template.New("usage").Parse("Usage: " + getUsageTemplate() + "\n\n")
		if err != nil {
			Panic(err)
		}
		err = tmpl.Execute(os.Stdout, CommandWrapper{
			Command: command,
			Group:   group,
			Bold:    bold,
			Reset:   reset,
		})
		if err != nil {
			Panic(err)
		}
	}

	// Help
	fmt.Printf("To see help text, you can run:\n")
	fmt.Printf("%s help\n", name)
	fmt.Printf("%s help <command>\n", name)
	fmt.Printf("%s help <command> <subcommand>\n\n", name)

	os.Exit(1)
}
