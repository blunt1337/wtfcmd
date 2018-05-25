# Command's definition
---

Configuration files are stored in `.wtfcmd.json`, `.wtfcmd.yaml` or `.wtfcmd.yml` files. They are loaded from each directory and parent directories from the current working dir.

For example, if you run `wtf` in `/Users/blunt/projects`, it will try to load from `/Users/blunt/projects/.wtfcmd.ext` and `/Users/blunt/.wtfcmd.ext`, and `/Users/.wtfcmd.ext` and `/.wtfcmd.ext`.

A configuration can override a command from a parent configuration if they share the same group and name.

To test your commands, you can put as first argument `--debug` to print the command before running it.  
For example `wtf --debug commandname arg1 arg2`.

---

# File structure
The .wtfcmd file must contain an array of [command](#commands) objects defined like bellow.

## Commands
A command is an object with the following properties.

#### group
group's name *(optional)*

To run the command it will be `wtf groupname name`. Without a group's name it is just `wtf name`.
Group names can be either a string, or an array ["fullname", "alias"].
Only utf8 alphanumeric characters and `:._-` are allowed.

#### name
command's name *(required)*

To run the command it will be `wtf name`, or `wtf groupname name` if a group is mentioned.
Command names can be either a string, or an array ["fullname", "a"], where 'a' is a one character alias.
Only utf8 alphanumeric characters and `:._-` are allowed.

#### desc
description or help message *(optional, but recommended)*

Description can be either a string, or an array for multiple lines.

#### cmd
command template *(required)*

[More info for the template section here](/template).
String arrays will be concatenated with "\n".
If the command is compatible with bash/powershell, you can just write a string or string[] as a command,
but for specific terminals, you can write an object like:
``` json
{
	"bash": [
		"echo this command will run in bash",
		"echo multilines works too"
	],
	"powershell": "echo this command will run with powershell.exe"
}
```

#### cwd
current working dir of the command *(optional)*

Like commands, you can put one for "bash" and one for "powershell".

- If it starts with a dot, it will run in the config's directory + cwd.
	For example, in /my_projects/wtfcmd.json, there is cwd = ./public; then running `wtf` in /my_projects/awesome/, will put the current working dir to /my_projects/public

- If it starts with '/' or 'x:', an absolute path.
	For example, in /my_projects/wtfcmd.json, there is cwd = /public; then running `wtf` in /my_projects/awesome/, will put the current working dir to /public

- If none of above, the directory where wtf was called + cwd.
	For example, in /my_projects/wtfcmd.json, there is cwd = public; then running `wtf` in /my_projects/awesome/, will put the current working dir to /my_projects/awesome/public

The default is the current directory where `wtf` is executed.
For more complex cases, in the command you can use `cd {{configdir}}`.

#### args
arguments *(optional)* 

Arguments must contain an array of [argument](#arguments) objects defined like below.

#### flags
flags *(optional)*

Flags must contain an array of [flag](#flags-2) objects defined like below.

---

### Arguments
Arguments are filled from the command line when they don't start with `-`.
For example, `wtf cmdname arg1 arg2`.
If an argument starts with a `-`, you can stop flags parsing with `--`.
For example, `wtf cmdname --flag -- -arg1`.

#### name
argument name *(required)*

It should be unique of course, and contain only utf8 alphanumeric characters and `:._-`.

#### desc
description or help message *(optional, but recommended)*

Description can be either a string, or an array for multiple lines.

#### required
required to run the command *(optional, default false)*

Set to true if the argument is required and should print an error if missing. Only first arguments can be required.

#### default
default value if the argument is missing *(optional, default nil)* 

The value if the argument is not required and not filled. Can be anything except an object.
If is_array is set and the default value is not an array, the default value will be inserted in an empty array.

#### is_array
True if the value is an array *(optional, default false)*

Only the last argument can be an array.
For example `wtf cmd a b c d` can be 1 argument with `[a, b, c, d]`.

#### test
test the value before running the command *(optional)*

The value is a string, and it can be:
- a regex like `^[a-z0-9]+$`, for case insensitive: `^(?i)[a-z0-9]+$`,
- `$int` to check for integers
- `$uint` to check for positive integers
- `$bool` to check for booleans (It accepts 1, t, TRUE, true, True, 0, f, F, FALSE, false, False)
- `$float` or "$number" to check for floating numbers
- `$file` to check for an existing files
- `$dir` to check for an existing directory
- `$dir/file` to check for an existing file or directory
- `$json` to check and parse the argument as a json object

---

### Flags
Flags are filled from the command line when they start with `-`.
If the flag starts with a double `-`, it uses it's full name.
If the flag starts with a single `-`, it uses it's short name.
For example, a flag named `["super", "s"]` can be filled with `wtf cmdname --super value` or `wtf cmdname -s value`.

#### name
flag name and alias *(required)*

Can be either a string, or an array ["fullname", "f"], where 'f' is the alias. Only UTF8 alphanumeric characters and `:._-` are allowed.
Can be used with `--dir=value`, or `--dir value`, or `-d=value`, or `-d value` (if not a boolean for the last one).

#### desc
description or help message *(optional, but recommanded)*

Description can be either a string, or an array for multiple lines.

#### default
default value if the argument is missing *(optional, default nil)*

The default value if the flag not mentioned. Can be anything except an object.
If is_array is set and the default value is not an array, the default value will be inserted in an empty array.

#### is_array
True if the value is an array *(optional, default false)*

Only the last argument can be an array.
For example `wtf cmd --flag value 1 --flag value2` make an array flag=[value1, value2].

#### test
test the value before running the command *(optional)*

The value is a string, and it can be:
- a regex like `^[a-z0-9]+$`, for case insensitive: `^(?i)[a-z0-9]+$`,
- `$int` to check for integers
- `$uint` to check for positive integers
- `$bool` to check for booleans (It accepts 1, t, TRUE, true, True, 0, f, F, FALSE, false, False)
- `$float` or "$number" to check for floating numbers
- `$file` to check for an existing files
- `$dir` to check for an existing directory
- `$dir/file` to check for an existing file or directory
- `$json` to check and parse the argument as a json object

To suggest a feature, report a bug, or general discussion: http://github.com/blunt1337/wtfcmd/issues/