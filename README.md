
Dotstrike is used to define and run repeatable file copy jobs to/from one or more paths.

It currently works for my use cases but is a bit rough around the edges, so use at your own risk.

Built in Go with [Cobra](https://github.com/spf13/cobra) and [Toml](https://github.com/BurntSushi/toml)

###  Quickstart

1. Create a `spec` and select it.
	Requres a *unique* alias (in this case "myspec")
		`> ds spec myspec`
	Specs define the details of a copy job.
	When a new spec is made, it is automatically selected.

2. Next, add source (`src`) and target (`tgt`) paths to the spec. Add a source (*copyjob input path*) with `src` command::
		`> ds src c:/my_files/`
	Add a target (*copyjob output path*) with `tgt` command':
		`> ds tgt 'd:/backups/personal files/'`
3. Run the spec at any time with the `run` command:
	Note that when a spec is run, **all** spec source paths are copied to **all** spec target paths.
		`> ds run myspec`


Most commands will apply to the currently selected spec by default.
The selected spec will be shown or denoted when running the `sel`, `list`, or `spec` commands with no arguments.
To select a different spec, use the `sel` command:
	`> ds sel myspec`
This will select 'myspec' if it exists. If not, it will search for the first spec with 'myspec' in its alias.

To see a list of commands and their usage, run `ds help`.
For detailed help on a specific command, run `ds help [command]` or `ds [command] -h`.

