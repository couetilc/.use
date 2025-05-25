# `.use`

A utility to make small bash scripts in my projects convenient and concise.
`.use` parses usage comments and validates scripts arguments declaratively.

## TODO

- need to parse docopt strings
- assertions will map 1-to-1 to `test` "conditional evaluation utility", in natural language instead.
- I think I need a default interpreter for the script e.g. defaults to "bash" or can do `#!/usr/bin/env .use python`

## Installation

`curl http://thisrepo/.use > bin/.use`

or

`curl http://thisrepo/release/v0.Whatever > bin/.use`

## Usage

Invoking `.use` directly should scan all scripts in its directory, printing
each one's usage. It will act as a overall view of the CLI available to a
particular project.
- note: shebang lines are limited to 255 characters since Linux 5.1 (see `man
  execve`), so only need to parse 255 bytes from start of file to determine if
  uses `.use`

Invoking a particular script will take advantageo of `.use` features
transparently, as long as it specifies `.use` in the shebang line.
- the script can now print usage with `-h`, `-help`, and `--help` without it
  being explicitly checked for.
- all arguments will be validated according to docopt grammar.
- Is there a way to summarize it in 80 characters?

## File Structure

File: `.use`

```sh
#!/usr/bin/env bash
wrapper="$0"
target="$1"
shift

# now $@ contains script arguments, $wrapper contains the path to the .use
# script, and $target contains the path to the invoked script e.g. .example 

# TODO: parse usage from $target
# TODO: validate argument format
# TODO: run assertions against validated arguments
# TODO: check for -h, -help, --help and print usage
# TODO: directly invoking .use should print usage for all scripts with shebang `.use`.
```

File: `.example`

```sh
#!/usr/bin/env .use
# Some example helper script
# Usage: .example <some-arg> [<some-optional-arg>]
#   .example -option 
# Options:
#   <some-arg>          Is some argument
#   <some-optional-arg> Is some optional argument
# Assertions:
#   <some-arg> is a file
#   <some-optional-arg> is a non-zero number
```

## Contributing

Install
- go
- direnv
