package main

import (
	"fmt"
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

// Orange character code.
var orange = ""

// Green character code.
var green = ""

// Cyan character code.
var cyan = ""

// Fill color codes if stdout is an ANSI-capable terminal.
func init() {
	if !noColor {
		bold = "\033[01m"
		reset = "\033[00m"
		red = "\033[91m"
		orange = "\033[33m"
		green = "\033[32m"
		cyan = "\033[36m"
	}
}

// GroupWrapper are parameters passed to the help template.
type GroupWrapper struct {
	Bold   string
	Reset  string
	Orange string
	Green  string
	Cyan   string
	Groups map[string]*Group
}

// CommandWrapper are parameters passed to the help2 template.
type CommandWrapper struct {
	Bold    string
	Reset   string
	Orange  string
	Green   string
	Cyan    string
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
func ShowHelp(groups []*Group, std *os.File) {
	tplTxt := "{{$B := .Bold}}{{$R := .Reset}}{{$O := .Orange}}{{$G := .Green}}{{$C := .Cyan}}" +
		"\n" +
		"{{$O}}{{$B}}AVAILABLE COMMANDS{{$R}}\n" +
		"\n" +
		"{{range .Groups}}" +
		/*	*/ "{{if .Name}}" +
		/*		*/ "{{$O}}{{$B}}{{.Name}}{{$R}}" +
		/*		*/ "{{if .Aliases}}" +
		/*			*/ " ({{join .Aliases \", \"}})" +
		/*		*/ "{{end}}" +
		/*		*/ "\n" +
		/*	*/ "{{end}}" +
		/*	*/ "{{range .Commands}}" +
		/*		*/ " - " +
		/*		*/ "{{if .Config.Group}}{{index .Config.Group 0}} {{end}}" +
		/*		*/ "{{$G}}{{$B}}{{.Name}}{{$R}}" +
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
	err = tmpl.Execute(std, GroupWrapper{
		Groups: groupByName,
		Bold:   bold,
		Reset:  reset,
		Orange: orange,
		Green:  green,
		Cyan:   cyan,
	})
	if err != nil {
		Panic(err)
	}
}

// getDefaultValueTemplate returns the template string for the default value.
func getDefaultValueTemplate() string {
	return "{{if .IsArray}}[" +
		/*	*/ "{{range $index, $default := .Default}}" +
		/*		*/ "{{if $index}}, {{end}}" +
		/*		*/ "{{printf \"%#v\" $default}}" +
		/*	*/ "{{end}}" +
		"]{{else}}" +
		/*	*/ "{{printf \"%#v\" .Default}}" +
		"{{end}}"
}

// getUsageTemplate returns the template string for command usage.
func getUsageTemplate() string {
	return "{{if .Group.Name}}{{.Group.Name}} {{end}}{{$G}}{{$B}}{{.Command.Name}}{{$R}}" +
		"{{if .Command.Config.Flags}} [flags]{{end}}" +
		"{{range .Command.Config.Args}}" +
		/*	*/ " {{if .Required}}<{{$G}}{{index .Name 0}}{{$R}}>{{if .IsArray}}...{{end}}" +
		/*	*/ "{{else}}[{{$G}}{{index .Name 0}}{{$R}}=" + getDefaultValueTemplate() + "]{{if .IsArray}}...{{end}}" +
		/*	*/ "{{end}}" +
		"{{end}}"
}

// ShowHelpCommand shows detailed help for a command.
func ShowHelpCommand(group *Group, command *Command) {
	tplTxt := "{{$B := .Bold}}{{$R := .Reset}}{{$O := .Orange}}{{$G := .Green}}{{$C := .Cyan}}" +
		"\n" +
		"{{$O}}{{$B}}SYNOPSIS{{$R}}\n" +
		"	" + getUsageTemplate() + "{{availability .Command.Config.Cmd}}" +
		"\n\n" +
		"{{if or .Group.Aliases .Command.Aliases}}" +
		/*	*/ "{{$O}}{{$B}}ALIASES{{$R}}\n" +
		/*	*/ "	{{if .Group.Name}}{{.Group.Name}}{{if .Group.Aliases}}|{{join .Group.Aliases \"|\"}}{{end}} {{end}}" +
		/*	*/ "{{.Command.Name}}{{if .Command.Aliases}}|{{join .Command.Aliases \"|\"}}{{end}}" +
		/*	*/ "\n\n" +
		"{{end}}" +
		"{{if .Command.Config.Desc}}" +
		/*	*/ "{{$O}}{{$B}}DESCRIPTION{{$R}}\n" +
		/*	*/ "	{{replace .Command.Config.Desc \"\\n\" \"\\n	\" -1}}" +
		/*	*/ "\n\n" +
		"{{end}}" +
		"{{if .Command.Config.Cwd}}" +
		/*	*/ "{{$O}}{{$B}}WORKING DIRECTORY{{$R}}\n" +
		/*	*/ "	{{cwd .Command.Config}}" +
		/*	*/ "\n\n" +
		"{{end}}" +
		"{{if .Command.Config.StopOnError}}" +
		/*	*/ "{{$O}}{{$B}}STOP ON ERROR:{{$R}} yes\n\n" +
		"{{end}}" +
		"{{if .Command.Config.Args}}" +
		/*	*/ "{{$O}}{{$B}}ARGUMENTS{{$R}}\n" +
		/*	*/ "{{range .Command.Config.Args}}" +
		/*		*/ "	{{$G}}{{$B}}{{index .Name 0}}{{$R}} ({{$C}}" +
		/*		*/ "{{if .IsArray}}array of {{end}}" +
		/*		*/ "{{if .Test}}" +
		/*			*/ "{{if hasPrefix .Test \"$\"}}" +
		/*				*/ "{{trimLeft .Test \"$\"}}" +
		/*			*/ "{{else}}" +
		/*				*/ "/{{replace .Test \"/\" \"\\\\/\" -1}}/" +
		/*			*/ "{{end}}" +
		/*		*/ "{{else}}string{{end}}" +
		/*		*/ "{{$R}}, {{if .Required}}required{{else}}default " + getDefaultValueTemplate() + "{{end}}" +
		/*		*/ ")\n\n" +
		/*		*/ "{{if .Desc}}" +
		/*			*/ "	{{replace .Desc \"\\n\" \"\\n	\" -1}}\n" +
		/*		*/ "{{end}}\n" +
		/*	*/ "{{end}}" +
		"{{end}}" +
		"{{if .Command.Config.Flags}}" +
		/*	*/ "{{$O}}{{$B}}FLAGS{{$R}}\n" +
		/*	*/ "{{range .Command.Config.Flags}}" +
		/*		*/ "	{{$G}}{{$B}}--{{index .Name 0}}{{$R}}" +
		/*		*/ "{{$l := len .Name}}{{if gt $l 1}}{{range $index, $Name := .Name}}{{if ne $index 0}}, {{$G}}-{{$Name}}{{$R}}{{end}}{{end}}{{end}} ({{$C}}" +
		/*		*/ "{{if .IsArray}}array of {{end}}" +
		/*		*/ "{{if .Test}}" +
		/*			*/ "{{if hasPrefix .Test \"$\"}}" +
		/*				*/ "{{trimLeft .Test \"$\"}}" +
		/*			*/ "{{else}}" +
		/*				*/ "/{{replace .Test \"/\" \"\\\\/\" -1}}/" +
		/*			*/ "{{end}}" +
		/*		*/ "{{else}}string{{end}}" +
		/*		*/ "{{$R}}, default " + getDefaultValueTemplate() +
		/*		*/ ")\n\n" +
		/*		*/ "{{if .Desc}}" +
		/*			*/ "	{{replace .Desc \"\\n\" \"\\n	\" -1}}\n" +
		/*		*/ "{{end}}\n" +
		/*	*/ "{{end}}" +
		"{{end}}" +
		"{{if .Command.Config.Envs}}" +
		/*	*/ "{{$O}}{{$B}}ENVIRONMENTS{{$R}}\n" +
		/*	*/ "{{range .Command.Config.Envs}}" +
		/*		*/ "	- {{$G}}{{replace . \"=\" (printf \"%s: \" $R) 1}}\n" +
		/*	*/ "{{end}}" +
		"{{end}}" +
		"\n\n"

	// Execute template with command
	tmpl, err := template.New("help2").Funcs(template.FuncMap{
		"join":         strings.Join,
		"replace":      strings.Replace,
		"availability": CmdAvailability,
		"hasPrefix":    strings.HasPrefix,
		"trimLeft":     strings.TrimLeft,
		"cwd": func(cfg *Config) string {
			switch GetTerminal() {
			case TermBash:
				return cfg.Cwd.Bash + " => " + ResolveCwd(cfg)
			case TermCmd, TermPowershell:
				return cfg.Cwd.Powershell + " => " + strings.Replace(ResolveCwd(cfg), "\\", "/", -1)
			}
			return ""
		},
	}).Parse(tplTxt)
	if err != nil {
		Panic(err)
	}
	err = tmpl.Execute(os.Stdout, CommandWrapper{
		Command: command,
		Group:   group,
		Bold:    bold,
		Reset:   reset,
		Orange:  orange,
		Green:   green,
		Cyan:    cyan,
	})
	if err != nil {
		Panic(err)
	}
}

// ShowCommandError prints an error for a command.
func ShowCommandError(msg string, group *Group, command *Command, groups []*Group) {
	name := os.Args[0]

	// Error
	fmt.Fprintf(os.Stderr, "%s: error: %s.\n\n", name, msg)

	// Usage
	if group != nil {
		if command != nil {
			// Usage of group command
			tplTxt := "{{$B := .Bold}}{{$R := .Reset}}{{$O := .Orange}}{{$G := .Green}}{{$C := .Cyan}}" +
				"Usage: " + getUsageTemplate() + "\n" +
				"{{if .Command.Config.Desc}}" +
				/*	*/ "Description: {{replace .Command.Config.Desc \"\\n\" \"\\n	\" -1}}\n" +
				"{{end}}" +
				"\n"

			tmpl, err := template.New("usage").Funcs(template.FuncMap{
				"replace": strings.Replace,
			}).Parse(tplTxt)
			if err != nil {
				Panic(err)
			}
			err = tmpl.Execute(os.Stderr, CommandWrapper{
				Command: command,
				Group:   group,
				Bold:    bold,
				Reset:   reset,
				Orange:  orange,
				Green:   green,
				Cyan:    cyan,
			})
			if err != nil {
				Panic(err)
			}
		} else {
			// Show group doc
			ShowHelp([]*Group{group}, os.Stderr)
		}
	} else {
		// Show all
		ShowHelp(groups, os.Stderr)
	}
	os.Exit(1)
}
