package main

import "blunt.sh/wtfcmd/vault"

// GetVaultCommands lists vault commands.
func GetVaultCommands() []*Group {
	return []*Group{
		{
			Name:    "--vault",
			Aliases: []string{},
			Commands: []*Command{
				{
					Name:    "user-add",
					Aliases: []string{"u-add"},
					Config: &Config{
						Desc: "Add an user to the allowed list, and re-encrypt",
						Args: []*ArgOrFlag{
							{
								Name:     []string{"public_key"},
								Desc:     "User's public key. Get it with 'wtf --vault public-key' comamnd",
								Required: true,
								Default:  nil,
								Test:     "",
								IsArray:  false,
							},
							{
								Name:     []string{"name"},
								Desc:     "Username to be associated with the public key",
								Required: true,
								Default:  nil,
								Test:     "",
								IsArray:  false,
							},
						},
						Flags:            []*ArgOrFlag{},
						internalFunction: vault.UserAddCmd,
					},
				},
				{
					Name:    "user-remove",
					Aliases: []string{"u-rm"},
					Config: &Config{
						Desc: "Remove an user from the allowed list, and re-encrypt",
						Args: []*ArgOrFlag{
							{
								Name:     []string{"public_key"},
								Desc:     "User's public key, fingerprint, or name. Get it with 'wtf --vault user-list' comamnd",
								Required: true,
								Default:  nil,
								Test:     "",
								IsArray:  false,
							},
						},
						Flags:            []*ArgOrFlag{},
						internalFunction: vault.UserRemoveCmd,
					},
				},
				{
					Name:    "user-list",
					Aliases: []string{"u-ls"},
					Config: &Config{
						Desc:             "List allower users",
						Args:             []*ArgOrFlag{},
						Flags:            []*ArgOrFlag{},
						internalFunction: vault.UserListCmd,
					},
				},
				{
					Name:    "public-key",
					Aliases: []string{"key"},
					Config: &Config{
						Desc: "Print your public key to the console",
						Args: []*ArgOrFlag{},
						Flags: []*ArgOrFlag{
							{
								Name:     []string{"raw"},
								Desc:     "Only return the public key without formatting the output",
								Required: false,
								Default:  false,
								Test:     "$bool",
								IsArray:  false,
							},
						},
						internalFunction: vault.PublicKeyCmd,
					},
				},
				{
					Name:    "set",
					Aliases: []string{},
					Config: &Config{
						Desc: "Add a protected value inside the vault",
						Args: []*ArgOrFlag{
							{
								Name:     []string{"key"},
								Desc:     "Value's name",
								Required: true,
								Default:  nil,
								Test:     "",
								IsArray:  false,
							},
							{
								Name:     []string{"value"},
								Desc:     "Value to store, or empty to delete",
								Required: false,
								Default:  "",
								Test:     "",
								IsArray:  false,
							},
						},
						Flags: []*ArgOrFlag{
							{
								Name:     []string{"ask"},
								Desc:     "Ask for the value instead of take the second parameter",
								Required: false,
								Default:  "",
								Test:     "",
								IsArray:  false,
							},
						},
						internalFunction: vault.SetCmd,
					},
				},
				{
					Name:    "get",
					Aliases: []string{},
					Config: &Config{
						Desc: "Get a protected value from the vault",
						Args: []*ArgOrFlag{
							{
								Name:     []string{"key"},
								Desc:     "Value's name",
								Required: true,
								Default:  nil,
								Test:     "",
								IsArray:  false,
							},
						},
						Flags: []*ArgOrFlag{
							{
								Name:     []string{"required"},
								Desc:     "Crash if the value is not set in the vault",
								Required: false,
								Default:  false,
								Test:     "$bool",
								IsArray:  false,
							},
						},
						internalFunction: vault.GetCmd,
					},
				},
			},
		},
	}
}
