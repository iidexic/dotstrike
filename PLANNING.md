# dotstrike — Planning & Action Items

Living punch list. Update in the same change as any planning decision, scope change, or completion. Mark completed items with `[x]` and an ISO date. Record dropped/deferred items with a one-line reason.

Status legend: `[ ]` open · `[~]` in progress · `[x]` done · `[-]` dropped/deferred

Last reviewed: 2026-09-03 (Phase 1 complete)

---

## Phase 0 — Cleanup (done, for reference)

- [x] 2026-09-03 delete merged branches (`barebones`, `more-barer-bones`, `i-branch-1`)
- [x] 2026-09-03 strip dead non-test symbols/files (see prior commits)
- [x] 2026-09-03 audit + drop stale test helpers
- [x] 2026-09-03 remove `pelletier/go-toml/v2` dep (only used in scratch comparison)
- [x] 2026-09-03 audit `.extra_code/` — verdict: no salvageable code

- [x] 2026-09-03 flush `.extra_code/` (done by user)
- [x] 2026-09-03 remove `tu.py` (done by user)

---

## Phase 1 — Nil-safety and small bugs (low blast radius) — complete

- [x] 2026-09-03 **`dscore/globalModify.go` — `prefs.setOpt` assigns to nil map.** Init `p.Bools` if nil before write. Also cleared mirror `BUG:` comment in `cmd/cmd-root_test.go`.
- [x] 2026-09-03 **`cmd/xcheck(debug).go` — `dirs`/`paths`/`sysdirs` nil deref.** Now calls `pops.PopulateSysDirs()` first and prints `<unset>` via new `ptrOrUnset` helper if any path pointer is still nil.
- [x] 2026-09-03 **`dscore/spec.go` — `Spec.DeleteByPtr` nil ptr.** Nil-checks receiver + components; skips nils; errors on all-nil; iterates in reverse to survive `slices.Delete` index shift.
- [x] 2026-09-03 **`cmd/cmd-config.go` → `dscore/globalModify.go` — `SetOptionString` nil deref.** Receiver-side guard on `gm`/`gm.globalData`; parallel guard added to `SetOptionBool` (plus nil-map init); caller `cfgApplyGlobalTargetCautious` now bails early with a user-facing message if `dscore.TempData()` is nil.
- [x] 2026-09-03 **`dscore/spec.go` — `Spec.IsPathChild` match logic.** New shared helper `componentMatchesPath` normalizes both sides using `pc.Abspath` (authoritative when set) or `pops.MakeAbs(pc.Path)` compared against `pops.MakeAbs(incoming)`. Refactored `IsPathSource`, `IsPathTarget`, and `GetIfChild` onto the same helper — `GetIfChild` also fixed to return `&S.Sources[i]` (real slice element) instead of the pre-Go-1.22-era `&src` bug that returned a pointer to a per-iteration copy.

### Phase 1 collateral fixes

- [x] 2026-09-03 `pathops/pathops.go` — `TildeCheck` and `TildeExpand` panicked on empty string (`ospath[0]` on zero-length). Both now short-circuit on empty. Surfaced by the new `componentMatchesPath` running against zero-value `PathComponent{}` entries in `TestAddComponent`.

### Phase 1 verification

- `go build ./...` — clean
- `go vet ./...` — clean
- `go test ./...` — the panic in `TestAddComponent` from `TildeCheck` on empty is gone. Many pre-existing test failures remain (reference the now-deleted `_xtra/[samplefiles]` fixtures, or use `Spec` zero-value shortcuts that mask other pre-existing bugs). None of the remaining failures are Phase 1 regressions. **Follow-up:** dedicated test-suite rehab pass — track separately below.

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

- [x] 2026-09-03 **Copyright header cleanup.** Unified all `cmd/*.go` + `main.go` to `Copyright © 2025 Derek`.
- [ ] **`cmd/cmd-root.go:82` — bare `cmd.Printf("DEBUG")` with no newline.** Change to `Println("DEBUG")` or drop entirely (`DumpGlobals` output that follows is already labeled).
- [ ] **`cmd/cmd-source.go:36` — same bare `Printf("DEBUG")` pattern.** Same fix.

## Phase 5 — Test suite rehab (added 2026-09-03)

Pre-existing failures uncovered while verifying Phase 1. Not regressions, but blocking a clean CI baseline.

- [ ] Delete or reroute tests that require `_xtra/[samplefiles]` fixtures (deleted by user): `TestEncodeHardAssign`, `TestEncodeToBuffer`, `TestForceEncodeDefaults`, and any others under `dscore/` that hard-code that path.
- [ ] `TestAddComponent` (`dscore/spec_test.go:9`) — both `if`/`else` branches call `t.Errorf`; test is un-passable. Rewrite assertions to reflect real intent.
- [ ] `TestRunMultiSource` — panics via `pathops.(*JobGroup).ConfigToJobs` (`copyjobGroup.go:64`), probably nil `spec.group`. Investigate under Phase 2/3 depending on scope.
- [ ] `TestRunFSdirs` — panics inside `pathops.testing_job` (`moveData_test.go:23`). Missing fixture or nil setup.

## Deferred / low-priority tracked TODOs

Left in-source, not scheduled:

- `dscore/dsconfig.go:52` — `ConfigOption = config.OptionKey` type alias slated for full replacement with direct `config.OptionKey`.
- `dscore/dsconfig.go:246` — `standardizeAlias` cleanup pass.
- `dscore/globalModify.go:234` — global prefs missing keys should populate as `false`.
- `pathops/pathops.go:172` — dedupe multiple system-dir helpers (`SystemDirectories` vs `PopulateSysDirs` vs `GetSysDirs`).

## Open Questions

- **`LookupOption` policy (Phase 2)**: reject ambiguous input with "did you mean" (recommended, safest), longest-match, or prefix-only? Blocks Phase 2 start.
