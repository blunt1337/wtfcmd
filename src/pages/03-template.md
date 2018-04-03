# Commands template usage

---

Commands are in [golang template](https://golang.org/pkg/text/template/) code. Here are some simple examples to get you started:
|     |     |
| --- | --- |
| Print an argument | `echo {{ .varname }}`
| Print a string argument escaped | `echo {{ esc .varname }}`
| Remove spaces between template code | `"{{ 23 -}} < {{- 45 }}"` = 23<45
| Condition with a boolean argument | `{{ if .varname }}echo "TRUE"{{end}}`
| Condition with a string argument | `{{ if -eq .varname "yes" }}echo "TRUE"{{ end }}`
| Multiple conditions | `{{ if and (-eq .varname1 "yes") (-eq .varname2 "no")}}echo "TRUE"{{ end }}`
| Loops | `{{ range $index, $element := .array_variable }}{{ $element }}{{ else }}Empty array{{ end }}`
| Use functions | `{{ replace "foobar" "foo" "bar" -1 }}` = barbar
| Use utilities functions | `{{ info "My name is" .name }}`

All functions available are listed bellow.

---

## String functions

All [golang strings](https://golang.org/pkg/strings/) function are available (but with a lower case first character):
- compare, contains, containsAny, containsRune, count, equalFold, fields, fieldsFunc, hasPrefix, hasSuffix, index, indexAny, indexByte, indexFunc, indexRune, join, lastIndex, lastIndexAny, lastIndexByte, lastIndexFunc, map, repeat, replace, split, splitAfter, splitAfterN, splitN, title, toLower, toLowerSpecial, toTitle, toTitleSpecial, toUpper, toUpperSpecial, trim, trimFunc, trimLeft, trimLeftFunc, trimPrefix, trimRight, trimRightFunc, trimSpace, trimSuffix,

## Extra string related functions

| escape (param *) * |
|---|
| Escape the argument to pass it as argument for bash or powershell. If the argument is a string, else it will return it as it is. |

| esc (param *) * |
|---|
| Alias of escape. |

| unescape (param *) * |
|---|
| Unescape an argument from bash or powershell. If the argument is a string, else it will return it as it is. |

| raw (param *) * |
|---|
| Alias of unescape. |

| json (object *, pretty bool = false) string |
|---|
| Convert first argument to JSON. Return false on error. |

| jsonParse (json string) object |
|---|
| Convert the json string argument to golang interface{}. Return false on error. |

# Utilities functions

| configdir () string |
|---|
| Return the directory of the configuration file of the command running. |

| error (args... *) string |
|---|
| Return a command that will print any arguments on stderr prefixed by a red [x]. |

| panic (args... *) string |
|---|
| Return a command that will print any arguments on stderr prefixed by a red [x]. Then exit with status 1. |

| warn (args... *) string |
|---|
| Return a command that will print any arguments on stderr prefixed by an orange [-]. |

| info (args... *) string |
|---|
| Return a command that will print any arguments on stdout prefixed by a blue [>]. |

| made (args... *) string |
|---|
| Return a command that will print any arguments on stdout prefixed by a green [+].

| ask (args... *) string |
|---|
| Return a command that will print any arguments on stdout prefixed by a purple [?], then it will print the answer.<br>For example: `{{ ask "How old are you?" }}; echo "you are $({{ read }}) years old"`<br>or for secure questions: `{{ ask "Password?" }}; echo "your new password: $({{ readSecure }})"`. |

| read () string |
|---|
| Return a command that will read from stdin, for bash and powershell.<br>For example: `{{ ask "How old are you?" }}; echo "you are $({{ read }}) years old"`. |

| readSecure () string |
|---|
| Return a command that will read a hidden text from stdin, for bash and powershell.<br>For example: `{{ ask "Password?" }}; echo "your new password: $({{ readSecure }})"`. |

| askYN (args... *) string |
|---|
| Return a command that will print any arguments on stdout prefixed by a purple [?]. Then wait for a yes/no response.<br>For example: `{{ askYN "Install some stuff" }}) && echo "installing..." || echo "installation skiped"`.

| bell () string |
|---|
| Return a command that will print the ASCII bell character. |