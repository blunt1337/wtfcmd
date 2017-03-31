# What's The Fu*king command!

- You want to run a command and you don't remember it?
- You have co-workers asking you for commands every time?
- You have a ton of shell scripts and don't remember what each one is for?

Build a single json/yaml file with all the commands available in your project, and anyone working with you, know in seconds how to run everything.

# Features

- Run commands with 'wtf my_command_name' defined in a json/yaml file.
- Auto-generated help pages! so you know what they are for and how to use them.
- Autocomplete all of your commands and flags!
- Possibility to have custom commands for bash and powershell.exe.

# Installation

If you have `go` installed, just run `go install blunt.sh/wtfcmd/wtf`
  Also [binaries](https://github.com/blunt1337/wtfcmd/releases) are available. Download and put it somewhere in your $PATH.
  If `wtf` is already in use, you can rename it whatever you like :)

To install the autocomplete, run `wtf --autocomplete install`

# Simple example

Now, let's say we use docker. To run our docker machine, we use the command:
`docker run -it --rm -p 8080:80 -v .:/app --name myproject myimage`

Let's see how we would use it with `wtf` instead:
`wtf docker start`

And how you'd need to configure it in the .cmds.json file:
```json
[
	{
		"group": ["docker", "dkr"],
        "name": ["start", "s"],
        "desc": [
			"Start an http server on the port 8080 by default.",
			"Files from the current directory are mapped to /app"
		],
        "cmd": "docker run -it --rm -p {{.port}}:80 -v .:/app --name myproject myimage",
        "flags": [
            {
                "name": ["port", "p"],
                "desc": "Port number",
                "test": "$uint",
                "default": "8080"
            }
        ]
	}
]
```

## Step by step

Base structure with your command:
```json
[
	{
		"cmd": "docker run -it --rm -p 8080:80 -v .:/app --name myproject myimage",
	}
]
```

The name of my command will be `wtf docker start`, but want an alias `wtf dkr s` too.
So i add to my object:
```json
"group": ["docker", "dkr"],
"name": ["start", "s"],
```

I want my team to know about this command, so i add:
```json
"desc": [
	"Start an http server on the port 8080 by default.",
	"Files from the current directory are mapped to /app"
],
```

I want the port number to be a parameter:
```json
"flags": [
	{
		"name": ["port", "p"],
		"desc": "Port number",
		"test": "$uint",
		"default": "8080"
	}
]
```
and change the 8080 of my command to {{.port}}:
`"cmd": "docker run -it --rm -p {{.port}}:80 -v .:/app --name myproject myimage",`

Check the [full configuration documentation](https://github.com/blunt1337/wtfcmd/CMDS.md) for more.

# TODOs

- [ ] More tests
- [ ] UI to build commands file
- [ ] Raw option, to forward all option/arguments not defined into a single string variable
- [ ] Global configurations in a folder, defined by an env variable
- [ ] Array arguments (last argument can be arg* or arg+, for an array)
- [ ] Array flags (multiple --flag value --flag value)
- [ ] Maybe, print the command documentation after a double tab?

To suggest a feature, report a bug, or general discussion:
http://github.com/blunt1337/wtfcmd/issues/