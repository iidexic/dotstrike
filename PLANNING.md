# dotstrike — Planning & Action Items

Living punch list. Update in the same change as any planning decision, scope change, or completion. Mark completed items with `[x]` and an ISO date. Record dropped/deferred items with a one-line reason.

Status legend: `[ ]` open · `[~]` in progress · `[x]` done · `[-]` dropped/deferred

Last reviewed: 2026-09-03

---

## Phase 0 — Cleanup (done, for reference)

- [x] 2026-09-03 delete merged branches (`barebones`, `more-barer-bones`, `i-branch-1`)
- [x] 2026-09-03 strip dead non-test symbols/files (see prior commits)
- [x] 2026-09-03 audit + drop stale test helpers
- [x] 2026-09-03 remove `pelletier/go-toml/v2` dep (only used in scratch comparison)
- [x] 2026-09-03 audit `.extra_code/` — verdict: no salvageable code

Pending user decision:
- [ ] flush `.extra_code/` directory (audit complete, awaiting go-ahead)
- [ ] decide fate of `tu.py` python helper (user pulled out; missing `pyproject.toml`, half-broken)

---

## Phase 1 — Nil-safety and small bugs (low blast radius)

Do these first. Each is a bounded fix, most 1–2 files.

- [ ] **`dscore/globalModify.go:419` — `prefs.setOpt` assigns to nil map.** If `p.Bools == nil`, initialize before write. This is the underlying issue behind `cmd/cmd-root_test.go:230` "Config Change; Prefs.setOpt() assigns to nil map".
- [ ] **`cmd/xcheck(debug).go:192` — `dirs`/`paths`/`sysdirs` prints `*pops.HomePath` etc.** Nil-check `HomePath`, `ConfigPath`, `CachePath` before dereference. Call `pops.PopulateSysDirs()` on the fly if unset, or print `<unset>`.
- [ ] **`dscore/spec.go:400` — `Spec.DeleteByPtr` nil ptr.** Iterating `components ...*PathComponent` and reading `c.Path` without nil-checking `c`. Skip nil entries and/or return error if all nil.
- [ ] **`cmd/cmd-config.go:232` → `dscore/globalModify.go:254` — `SetOptionString` nil ptr in `cfgApplyGlobalTargetCautious`.** Trace: probably `tempData` is nil (not initialized) or `gm.globalData` is nil. Add nil guard + return descriptive error, and ensure init path is exercised before this call.
- [ ] **`dscore/spec.go:184` — `Spec.IsPathChild` fails to match.** Currently compares `src.Path` against multiple transforms of `path` but never compares transforms of `src.Path` against `path`. Very likely the stored `src.Path` was TildeExpanded/MakeAbs'd at write time, so a raw incoming `path` never equals it. Fix: normalize both sides (probably: compare `pops.MakeAbs(src.Path) == pops.MakeAbs(path)` and treat `src.Abspath` as authoritative when set).

## Phase 2 — Config lookup ambiguity

- [ ] **`config/config_test.go:12`-`:13` — `LookupOption` maps ambiguous substrings to wrong option.** `"nohiddenrepo"` matches both `IgnoreHidden` and `IgnoreRepo`; `"useglobaltgtdir"` matches both `UseGlobalTarget` and `GlobalTargetPath`. Decide policy: (a) require exact/prefix match, (b) score by longest-match, (c) reject ambiguous input with an error listing candidates. Recommend (c) — silent misroute is the worst outcome for a config command.

## Phase 3 — Panic → error (higher blast radius)

Signatures change for several of these; do one at a time with `go build ./...` between.

- [ ] **Delete `pathops/pathops.go:49` — `ce()` helper.** Dead (zero call sites); it panics only if you pass a message, which is inverted logic anyway. Safe delete.
- [ ] **`pathops/pathops.go:286` — `MakeAbs` panics.** ~10 call sites across `dscore/`, `cmd/`, `magefiles/`, and internal `pathops/`. Options: (a) rename to `MakeAbsE(string) (string, error)` and add call-site fixes, (b) keep name and just log + return `""` on failure. Recommend (a); many callers already have error paths. Note the existing `MakeAbsIfPathlike` already returns `(string, error)`, use it as reference.
- [ ] **`pathops/pathops.go:570` — `CalledFrom` panics.** Only 3 call sites, all in `cmd/xcheck(debug).go`. Return `(string, error)`; xcheck can print `<err: ...>`.
- [ ] **`dscore/globals.go:78` — `globalsFilepath` panics if unset.** One caller (`dscore/globalToml.go:72`). Return `(string, error)`; propagate.
- [ ] **`dscore/dsconfig.go:190` — `decodeRawData` panics on toml decode failure.** One test caller (`globals_test.go:54`). Return error; already have `decodeAsConfig` as the pattern.
- [ ] **`cmd/cmd-root.go:136` — `configLoadInit` panics on `LoadGlobals` error.** Called via `cobra.OnInitialize`. Convert to soft-fail: print to stderr, set a global flag, let subsequent commands decide whether to bail.
- [ ] **`cmd/user_confirmation.go:32` — `askConfirmf` panics on stdin read error.** Return `(bool, error)`; callers (`checkConfirm`, `checkConfirmF`) also update.

## Phase 4 — Housekeeping

- [ ] **Copyright header cleanup.** `main.go` and several `cmd/*.go` still say `Copyright © 2025 NAME HERE <EMAIL ADDRESS>`; `cmd/cmd-source.go` uses `Copyright © 2025 derek :)`. Pick one format and apply uniformly. **Awaiting user preference** — see Open Questions.
- [ ] **`cmd/cmd-root.go:82` — bare `cmd.Printf("DEBUG")` with no newline.** Change to `Println("DEBUG")` or drop entirely (`DumpGlobals` output that follows is already labeled).
- [ ] **`cmd/cmd-source.go:36` — same bare `Printf("DEBUG")` pattern.** Same fix.

## Deferred / low-priority tracked TODOs

Left in-source, not scheduled:

- `dscore/dsconfig.go:52` — `ConfigOption = config.OptionKey` type alias slated for full replacement with direct `config.OptionKey`.
- `dscore/dsconfig.go:246` — `standardizeAlias` cleanup pass.
- `dscore/globalModify.go:234` — global prefs missing keys should populate as `false`.
- `pathops/pathops.go:172` — dedupe multiple system-dir helpers (`SystemDirectories` vs `PopulateSysDirs` vs `GetSysDirs`).

## Open Questions

- **Copyright boilerplate**: which format? Suggested options — `// Copyright © 2025 Derek <derekqvandam@gmail.com>` (formal), `// Copyright © 2025 Derek` (name-only), or delete header entirely (Go convention rarely requires it).
- **`.extra_code/` flush**: `git rm -r .extra_code` now, or leave since already committed as backup?
- **`tu.py`**: bring back into `scripts/` with `pyproject.toml`, or drop entirely?
