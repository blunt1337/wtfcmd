package main

// Group of commands.
type Group struct {
	Name     string
	Aliases  []string
	Commands []*Command
}

// Command structure.
type Command struct {
	Name    string
	Aliases []string
	Config  *Config
}

// BuildHierarchy will create Groups and Commands depending of the order of the Configs.
// If same group name => merge commands.
// If same group alias (but not same group name) => remove alias.
// If same command name => remove parent one.
// If same command alias (but not same command name) => remove alias.
func BuildHierarchy(cfgs []*Config) []*Group {
	// Group alias to group name
	aliasToName := map[string]string{}
	names := map[string]bool{}

	for _, cfg := range cfgs {
		if cfg.Group != nil {
			if lg := len(cfg.Group); lg > 0 {
				name := cfg.Group[0]
				names[name] = true

				// Aliases
				if lg > 1 {
					for i := 1; i < lg; i++ {
						alias := cfg.Group[i]
						name2, ok := aliasToName[alias]

						if !ok {
							// First time for this alias
							aliasToName[alias] = name
						} else if name != name2 {
							// Alias already in use, but not same group name, remove alias
							aliasToName[alias] = ""
						}
					}
				}
			}
		}
	}

	// Build those groups
	groups := []*Group{}
	for name := range names {
		group := new(Group)
		group.Name = name
		groups = append(groups, group)

		// Build aliases
		aliases := []string{}
		for alias, name2 := range aliasToName {
			if name == name2 {
				aliases = append(aliases, alias)
			}
		}
		group.Aliases = aliases

		// Build commands
		group.Commands = buildCommandHierarchy(cfgs, name)
	}

	// Root group
	group := new(Group)
	group.Aliases = []string{}
	group.Commands = buildCommandHierarchy(cfgs, "")
	if len(group.Commands) > 0 {
		groups = append(groups, group)
	}

	return groups
}

// buildCommandHierarchy will create Commands depending of the order of the Configs.
// If same command name => remove parent one.
// If same command alias (but not same command name) => remove alias.
func buildCommandHierarchy(cfgs []*Config, groupname string) []*Command {
	// Command alias to command name
	aliasToName := map[string]string{}
	names := map[string]*Config{}

	for _, cfg := range cfgs {
		// Filter by group name
		isEmptyGroup := groupname == "" && (cfg.Group == nil || len(cfg.Group) == 0)
		isSameGroup := groupname != "" && cfg.Group != nil && len(cfg.Group) > 0 && cfg.Group[0] == groupname

		if isEmptyGroup || isSameGroup {
			name := cfg.Name[0]
			if _, ok := names[name]; !ok {
				names[name] = cfg
			}

			// Aliases
			if lg := len(cfg.Name); lg > 1 {
				for i := 1; i < lg; i++ {
					alias := cfg.Name[i]
					name2, ok := aliasToName[alias]

					if !ok {
						// First time for this alias
						aliasToName[alias] = name
					} else if name != name2 {
						// Alias already in use, but not same group name, remove alias
						aliasToName[alias] = ""
					}
				}
			}
		}
	}

	// Build those commands
	commands := []*Command{}
	for name, cfg := range names {
		command := new(Command)
		command.Name = name
		command.Config = cfg
		commands = append(commands, command)

		// Aliases
		aliases := []string{}
		for alias, name2 := range aliasToName {
			if name == name2 {
				aliases = append(aliases, alias)
			}
		}
		command.Aliases = aliases
	}

	return commands
}
