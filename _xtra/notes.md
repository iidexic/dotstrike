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
