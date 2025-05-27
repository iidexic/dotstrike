# Dotstrike Notes


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
