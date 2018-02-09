# Full documentation of the configuration file

- The configuration file is loaded from each directory and parent directories from the current working dir.

  e.g. If you run `wtf` in `/Users/blunt/projects`, it will load from `/Users/blunt/projects` and `/Users/blunt`, and `/Users` and `/`.

- Loaded files are `.wtfcmd.json`, `.wtfcmd.yaml` or `.wtfcmd.yml`.
- A configuration can override a command from a parent configuration if they share the same group and name.

- To test your commands, you can put as first argument `--debug` to print the command before running it.  
E.g. `wtf --debug commandname arg1 arg2`

### Format
The file must contain an array of [command](#command) objects.
TODO: summary

### command
A command is an object with the following properties

##### group
> **Optional** group name
> - to run the command it will be `wtf group name`. Without a group name it is just `wtf name`.
> - Group names can be either a string, or an array ["fullname", "alias"].
> - Only utf8 alphanumeric characters and `:._-` are allowed.

##### name
> **Required** command name
> - to run the command it will be `wtf group name`.
> - Command names can be either a string, or an array ["fullname", "a"], where 'a' is the one character alias.
> - Only utf8 alphanumeric characters and `:._-` are allowed.

##### desc
> **Optional, but recommanded** description or help message  
> Description can be either a string, or an array for multiple lines.

##### cmd
> **Required** command template  
> commands can go from a simple strings, to a complexe template.
> - If the command is compatible with bash/powershell, you can just write a string or string[] as a command.  
>   But for specific terminals, you can write an object:  
>   ```js
    {  
    	"bash": [  
    		"echo this command will run in bash",  
    		"echo multilines works too"  
    	],  
    	"powershell": "echo this command will run with powershell.exe"  
    }  
    ```
> - The command format is **go template** https://golang.org/pkg/text/template/.
>   - go template can do so many things, like variables, loops, functions, etc.  
>     E.g. to print an argument, just write {{.argumentName}}, or to print a string safely (escaped): {{esc .argumentName}}
>   - We made all strings functions available too https://golang.org/pkg/strings/.  
>     E.g. {{replace .myFlagName "old" "new" -1}}
>   - We also added some message functions:  
>     {{made "some success message"}}  
>     {{error "some error message"}}  
>     {{warn "a warning here"}}  
>     {{askYN "is it true?"}} && echo "You said yes" || echo "You said no"  
>     For all functions, more information at TODO:doc

##### cwd
> **Optional**, the current working dir of the command. Like commands, you can put one for "bash" and one for "powershell".
> - If it starts with a dot, it will run in the config's directory + cwd.  
>   E.g. In /my_projects/wtfcmd.json, there is cwd = ./public; then running `wtf` in /my_projects/awesome/, will put the current working dir to /my_projects/public
> - If it starts with '/' or 'x:', an absolute path.  
>   E.g. In /my_projects/wtfcmd.json, there is cwd = /public; then running `wtf` in /my_projects/awesome/, will put the current working dir to /public
> - If none of above, the directory where wtf was called + cwd.  
>   E.g. In /my_projects/wtfcmd.json, there is cwd = public; then running `wtf` in /my_projects/awesome/, will put the current working dir to /my_projects/awesome/public  
>
> The default is the current directory where `wtf` is executed.  
> For more complex cases, you can use in the command a `cd {{configdir}}`.

##### args
> **Optional** arguments  
> Arguments must contain an array of [Argument](#Argument) objects.

##### flags
> **Optional** flags  
> Flags must contain an array of [Flag](#Flag) objects.

### Argument
An argument is an object with the following properties

##### name
> **Required** argument name
> - Should be unique of course.
> - Only utf8 alphanumeric characters and `:._-` are allowed.

##### desc
> **Optional, but recommanded** description or help message  
> Description can be either a string, or an array for multiple lines.

##### required
> **Optional, default false**
> - Set to true if the argument is required and should print an error if missing.
> - Only first arguments can be required

##### default
> **Optional, default nil** default value
> - The value if the argument is not required and not filled.
> - Can be anything except an object or array.

##### is_array
> **Optional, default false** the value is an array  
> Only the last argument can be an array. E.g. "wtf cmd a b c d" can be 1 argument with [a, b, c, d].

##### test
> **Optional** test code/regex  
> The test will show an error message if the argument doesn't pass it.  
> The value is a string, it can be:
> - a regex like `^[a-z0-9]+$`, for case insensitive: `^(?i)[a-z0-9]+$`,
> - `$int` to check for integers
> - `$uint` to check for positive integers
> - `$bool` to check for booleans (It accepts 1, t, TRUE, true, True, 0, f, F, FALSE, false, False)
> - `$float` or "$number" to check for floating numbers
> - `$file` to check for an existing files
> - `$dir` to check for an existing directory
> - `$dir/file` to check for an existing file or directory

### Flags
A flag is an object with the following properties

##### name
> **Required** flag name
> - Can be either a string, or an array ["fullname", "f"], where 'f' is the alias.
> - Can be used with `--dir=value`, or `--dir value`, or `-d=value`, or `-d value` (if not a boolean for the last one).
> - Only UTF8 alphanumeric characters and `:._-` are allowed.

##### desc
> **Optional, but recommanded** description or help message  
> Description can be either a string, or an array for multiple lines.

##### default
> **Optional, default nil** default value
> - The default value if the flag not mentioned.
> - Can be anything except an object or array.

##### is_array
> **Optional, default false** the value is an array  
> So multiple --flag value1 --flag value2 make an array flag=[value1, value2].

##### test
> **Optional** test code/regex  
> The test will show an error message if the flag doesn't pass it.  
> The value is a string, it can be:
> - a regex like `^[a-z0-9]+$`, for case insensitive: `^(?i)[a-z0-9]+$`,
> - `$int` to check for integers
> - `$uint` to check for positive integers
> - `$bool` to check for booleans (It accepts 1, t, TRUE, true, True, 0, f, F, FALSE, false, False)
> - `$float` or "$number" to check for floating numbers
> - `$file` to check for an existing files
> - `$dir` to check for an existing directory
> - `$dir/file` to check for an existing file or directory
