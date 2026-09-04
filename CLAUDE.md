# CLAUDE.md

Guidance for Claude Code when working in this repository.

## Project

`dotstrike` (`iidexic.dotstrike`) is a hand-written Go 1.25 CLI for copying source directories to target directories according to reusable named specs. Config and user data live in a single TOML file (`dotstrikeData.toml`) under the OS user config dir.

## Layout

- `main.go` — entry, calls `cmd.Execute()`
- `cmd/` — Cobra commands (`cmd-root.go`, `cmd-spec.go`, `cmd-source.go`, `cmd-target.go`, `cmd-config.go`, `cmd-run.go`, `utilCmd-*.go`, `xcheck(debug).go`)
- `dscore/` — user data, spec/prefs types, TOML decode/encode, job manager, execution
- `pathops/` — filesystem + path helpers (`pops` alias), copy primitives
- `config/` — option key registry, lookup, type flags (`Tbool`/`Tstring`)
- `match/` — pattern/glob matching
- `uout/` — indented output builder used across packages
- `magefiles/` — Mage build tasks
- `.extra_code/` — pre-release backup, unused (see PLANNING.md for flush status)

## Build / Run / Test

Windows PowerShell environment. From repo root:

```powershell
go build ./...
go vet ./...
go test ./...
mage -l              # list mage targets
```

The compiled binary is named `ds` (see magefile). Test data may reference `../_xtra/[samplefiles]` which no longer exists — some tests are expected to fail until refactored.

## Conventions

- Import alias: `pops "iidexic.dotstrike/pathops"`
- Errors are usually returned; a few legacy call sites still `panic` (tracked in PLANNING.md).
- Spec + component aliases are standardized via `dscore.standardizeAlias` (lower-case, trim escape chars).
- User data mutation goes through `tempData` (`dscore.TempData()`) and is flushed to disk in `dscore.EndEncode` on Cobra finalize. Modifying methods must call `tempData.Modify()` (or `gm.Modify()`) before touching the underlying fields.
- Copyright header format currently inconsistent across `cmd/*.go` — see PLANNING.md housekeeping section.

## Working notes

- **Keep `PLANNING.md` current.** Whenever we scope new work, decide on an approach, or complete an action item, update `PLANNING.md` in the same change: add new tasks, tick off completed ones with the date, and record any deferred/dropped items with a one-line reason. Treat it as the single source of truth for "what's next" on this project.
- Prefer editing existing files over introducing new ones. This project is small; premature abstraction has bitten it already (see `.extra_code`).
- When removing dead code, verify with a grep across the whole tree — the `dscore` package re-exports config aliases, so a symbol may look unused within its package but still be referenced from `cmd` or tests.
- Do not touch anything under `.extra_code/` — it is a committed backup slated for flush; see PLANNING.md.
