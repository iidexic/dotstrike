# Todo:
## structure
- [ ] clarify difference + relationship between dscore package files
	- [ ] specifically, globals + dsconfig have a lot of overlap

# Dotstrike Structure

## Current Structure
### Package: main
Serves exclusively to kick off the cmd package (cobra-based cli)
### Package: cmd
CLI functionality & direct user interface code
using Cobra; made up of individual command .go files
#### Files:
- root.go
- cfg.go
- config.go (remove)
- init.go
- show(debug).go
#### Commands
##### cfg.go
#####  config.go (remove)
#####  init.go
#####  show(debug).go
### Package: dscore
Structure of components required to perform actual functionality of building filesets and copying them between locations
### Package: pathops
Performs actual filesystem operations required; this includes:
* creating/reading/writing user data files (dotstrike config and user structure data)
* create/copy files; to be used by dscore for core functionality

# General Dotstrike Notes
## Cobra

## Planning
### Names for shit
#### groups
I am not the biggest fan of group as the name (specifically dg group) as a subcommand

**alternatives**:
---
- set
- lot
- kit
- bag
- pile
- pack
- pot? wad? box? lump?

So far, probably leaning toward kit, bag, pack?

## Ancient Confucian Wisdom

### Filesys Ops

1. Don't check for a dir's existence, just try to make it. 
    - `os.MkdirAll` works even if the path exists
2. Don't check for a file's existence. just try to make it.
    - if the file exists, it will error, and we catch with `os.IsExist(err)`
    - we attempt to create with `os.OpenFile`, using the `os.O_CREATE` flag

### Running other shit
* os/exec package is how we would run other executable thingz

# Extras:
## Random go code
### Get environment vars, split into var name and value, print var names:
```go

		fmt.Println("env:")
		env := os.Environ()
		evars := []string{}
		evals := []string{}
		for _, v := range env {
			vs := strings.Split(v, "=")
			if len(vs) == 1 {
				print("noeq -> ", vs[0])
			} else if len(vs) > 2 {
				vs0 := vs[0]
				vs1 := strings.Join(vs[1:], "|")
				vs = []string{vs0, vs1}
			}
			evars = append(evars, vs[0])
			evals = append(evals, vs[1])
		}
		for _, v := range evars {
			print(v, "|")
```

# toml format + notes

toml structure:

storagePath:string
prefs: struct
|- keepRepo:bool
|- keepHidden:bool
|- storedataSourcedir:bool (do not remember purpose)
|- globalTarget:bool

cfgs: []struct
|- alias:string
|- sources:[]string
|- targets:[]string

## Types
[toml table] --> go struct or map
[toml table array] --> go []struct or []map

## KEYS:
[toml key] -> go map key or go struct field name
### Map a non-matching field name to toml key:
```Go
server struct {
	//omitempty avoids encoding blank val if struct field empty
	IP     string       `toml:"ip,omitempty"` 
	Config serverConfig `toml:"config"`
}
```
---------
* use backtick as doublequote of tomlkey seems necessary (though technically can make a tag with any quotation mark)
- most of our fields may need omitempty as we will be doing many partial
	edits to toml file; do not want to overwrite toml data with blank values
 ───────────────────────────────────────
*/
