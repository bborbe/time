---
status: idea
kind: bug
---

## Summary

- `ParseDuration` lowercases its input before unit lookup, so `M` (intended as month) silently becomes `m` (minute). Callers get a valid `time.Time` that is months off, with no error or warning.
- The library has no month or year unit at all. `Y` falls through to error today; `M` is the silently-wrong case.
- The bug is data-corruption shaped, not a crash. Consumers see plausible but wrong durations until they spot-check the math.
- Confirmed in production: a downstream weekly cron passing `NOW-6M` got a 6-minute window for weeks, producing empty backtests that looked superficially valid.
- A consumer-side workaround (using week-unit syntax) is already deployed and unblocks the immediate consumer. This spec is about the latent bug itself, not the consumer fix.

## Problem

`ParseDuration` is the canonical entry point for human-written duration strings across the bborbe ecosystem (`ParseTime`, `NOW-NX` arithmetic, config files, cron specs). The implementation calls `strings.ToLower` on the entire input before regex matching, which collapses any case-sensitive unit distinction the regex could otherwise enforce. Today there is exactly one case-sensitive distinction that matters in practice — `M` (month) vs. `m` (minute) — and it is the one most likely to bite: month is a natural human unit for lookback windows, minute is the only single-letter unit currently mapped, and the collision produces no error.

The blast radius is "any caller writing `NOW-NM` or `NOW-NY`." Today that set is small (one trading cron, mitigated via workaround), but the library is the canonical home for these primitives — every new consumer is one `NOW-6M` away from silently broken behavior. Documentation in consumer repos still shows `NOW-1M  # 1 month ago` as an example because the bug is invisible from outside.

## Goal

`ParseDuration` either correctly parses month and year units, or rejects them as unknown. Silent collapse to a different unit is not a possible outcome. A human reading `NOW-6M` and getting a result that is anything other than "six months" or "error" never happens again.

## Non-goals

- Calendar-aware month arithmetic. A constant 30-day month and 365-day year are acceptable, consistent with how `Day` and `Week` already work in `UnitMap`.
- Reworking `ParseTime` semantics beyond what falls out of fixing `ParseDuration`.
- Updating consumer pins or downstream documentation — that is per-consumer work after a release.
- Adding new unit letters beyond `M` and `Y` (no quarter, no fortnight, no decade).
- Changing existing lowercase unit behavior. `m` must still be Minute, `h` Hour, `d` Day, `w` Week, etc.

## Desired Behavior

1. `ParseDuration("6M")` returns exactly 180 days (4320h), not 6 minutes.
2. `ParseDuration("1Y")` returns exactly 365 days (8760h).
3. `ParseDuration("6m")` still returns 6 minutes — lowercase semantics unchanged.
4. `ParseDuration("5Q")` (or any unknown unit) returns a non-nil error. No silent zero, no silent fallthrough.
5. `ParseTime("NOW-6M")` returns a timestamp roughly six months before `Now()` (180 days, within sub-second tolerance for a fixed clock).
6. Mixed expressions like `1Y2M3w4d5h6m7s` parse correctly and additively, with `M` always meaning month and `m` always meaning minute regardless of position.
7. All previously-passing parse tests continue to pass without modification — the only test that must flip is any existing assertion that `1M` equals 1 minute (that assertion was encoding the bug).

## Constraints

- **Backward compatibility for lowercase units**: `m`, `s`, `h`, `d`, `w`, `ns`, `us`, `ms` semantics frozen. No surprises for existing callers using lowercase.
- **No new dependencies**: solution must use stdlib only (current dependency set is unchanged).
- **Error type compatibility**: unknown-unit errors must flow through the existing `github.com/bborbe/errors` wrapping pattern (`errors.Errorf(ctx, ...)`), not raw `fmt.Errorf`.
- **Public API surface unchanged**: no new exported functions, no signature changes on `ParseDuration` / `ParseTime`. The fix lives entirely in the implementation.
- **`UnitMap` stays exported**: it is part of the public surface today; additions are additive only.
- **Test framework**: new test cases use the existing `DescribeTable` / `Entry` pattern in `time_duration_test.go` — no `t.Run` table loops, no parallel-only sub-tests.
- **License headers and `make precommit` cleanliness**: the existing `addlicense` + lint + vulncheck pipeline must stay green.

## Failure Modes

| Trigger | Expected behavior | Recovery |
|---------|-------------------|----------|
| Caller writes `NOW-6M` (intended: 6 months) | Returns 6 months offset. No silent collapse to 6 minutes. | N/A — the bug is the recovery. |
| Caller writes `NOW-5Q` (unknown unit) | `ParseTime` returns a wrapped error naming the unit. Caller is forced to handle it. | Caller fixes the input or catches the error. |
| Caller writes mixed-case `6M30d` | Parses as 6 months + 30 days = 210 days. `M` does not collapse. | N/A. |
| Caller writes `6m` (lowercase, intended: 6 minutes) | Returns 6 minutes — unchanged from today. | N/A — this is the contract. |
| Caller writes `6mM` or `6MM` (typo-shaped input) | Returns an error (regex mismatch or unknown unit), never a silent partial parse. | Caller fixes input. |
| Caller writes an empty string | Existing behavior preserved (whatever it returns today — zero or error). | N/A — out of scope. |

## Do-Nothing Option

The immediate consumer (trading weekly backtest cron) is already mitigated via week-unit syntax (`NOW-26W`). The bug becomes purely latent — no live caller is wrong.

**Cost of doing nothing:**
- Every future consumer that writes `NOW-NM` or `NOW-NY` will silently corrupt. No error, no log line, no test failure unless someone specifically asserts the resulting span.
- Documentation in downstream repos that shows `NOW-1M  # 1 month ago` (as one trading-repo doc already does) will be silently misleading forever.
- The bug is harder to discover the second time it bites because the canonical example (the trading cron) has been worked around — there is no live failure to surface it.
- The library's API surface continues to claim it supports human-readable durations (`README` mentions "weeks and days: 1w2d3h4m5s") while having a sharp edge where `M` looks like it should mean month.

**Cost of fixing:**
- Two-file change in this repo (`time_duration.go`, `time_duration_test.go`), plus CHANGELOG.
- One minor-version bump (additive feat + behavioral fix).
- No consumer pin bumps required for current callers — week-unit workaround stays valid regardless of the fix.

The asymmetry favors fixing. The latent footgun is cheap to eliminate now; the next discovery may not have a convenient workaround.

## Reproduction

**Library version:** any tag through `v1.27.0` (current latest, as of 2026-05-10).

**Minimum reproduction:**

```go
package main

import (
	"fmt"
	libtime "github.com/bborbe/time"
)

func main() {
	d, err := libtime.ParseDuration("6M")
	fmt.Printf("input=%q duration=%s err=%v\n", "6M", d, err)
	// Observed: input="6M" duration=6m0s err=<nil>
	// Expected: input="6M" duration=4320h0m0s err=<nil>   (180 days)
	//   OR:     input="6M" duration=0s err=unknown unit 'M'
}
```

**Production reproduction (the way this was discovered):**

- A Kubernetes CronJob in the trading repo (`core/backtest/cron/overlays/prod/weekly-cronjob.yaml`) passes `"from": "NOW-6M"` to a backtest service.
- The backtest service calls `libtime.ParseTime("NOW-6M")`, intending six months of historical data.
- Observed: `from` and `until` end up six minutes apart instead of six months. Resulting backtest is empty (zero trades, null profit factor).
- Confirmed via `mcp__trading-prod__get-backtest cc98a622-25a7-49f8-9d81-a61d3ac3d9e4` (May 2 weekly DE40 V25 run): `from=2026-05-02T01:09:04Z`, `until=2026-05-02T01:15:04Z`, span = 6 minutes.
- After switching the cron to `"from": "NOW-26W"` (week-unit syntax) and re-triggering, the same backtest produced a 182-day span — proving the library, not the cron container, was the source of corruption.

**Source-level evidence:**

- `time_duration.go:32-41` — `UnitMap` contains `ns us ms s m h d w` only. No `M` (month), no `Y` (year).
- `time_duration.go:43-45` — `durationRegexp` has no capture group for uppercase `M` or `Y`.
- `time_duration.go:109` — `str = strings.ToLower(str)` is called on the entire input before regex matching, so any uppercase `M` is folded to `m` before lookup, and `UnitMap["m"] = Minute` wins.

## Expected vs Actual

| Input | Expected (per a reasonable reading of "human-readable duration") | Actual |
|-------|------------------------------------------------------------------|--------|
| `"6M"` | 180 days (6 × 30d), or a clear "unknown unit" error | **6 minutes** — silently wrong |
| `"1Y"` | 365 days, or "unknown unit" error | **Error** — regex doesn't match `Y` (this is the *less bad* of the two; at least it's loud) |
| `"6m"` | 6 minutes | 6 minutes ✓ |
| `"NOW-6M"` (via `ParseTime`) | Roughly six months before `Now()` | **Six minutes before `Now()`** — silently wrong |
| `"5Q"` (clearly invalid) | Error naming the unit | Currently returns an error via the regex no-match path (acceptable, but the path is incidental, not designed) |

The asymmetry between `M` (silently wrong) and `Y` (loudly wrong) is the central defect. A fix that adds both as proper case-sensitive units gives both the same shape: known unit → correct result, unknown unit → error.

## Workaround

**For consumers that need month-scale lookback today**, use week-unit syntax instead:

| Intent | Don't write | Write instead |
|--------|-------------|---------------|
| 1 month back | `NOW-1M` | `NOW-4W` (28d, ~7% short) |
| 3 months back | `NOW-3M` | `NOW-13W` (91d, ~1% over) |
| 6 months back | `NOW-6M` | `NOW-26W` (182d, <1d off true average) |
| 1 year back | `NOW-1Y` | `NOW-52W` (364d, ~1d short) |

Week-unit (`W`) is in `UnitMap` as 7 × Day, has no case-collision with another unit, and is unambiguous through `strings.ToLower`. This is the workaround already deployed for the trading weekly cron (PR #119 in `bborbe/trading`, merged 2026-05-10).

**The workaround is not a substitute for the fix** — it requires every future consumer to remember the gotcha. The fix removes the gotcha.

## Acceptance Criteria

- [ ] `ParseDuration("6M")` returns 4320h0m0s exactly, error nil.
- [ ] `ParseDuration("1Y")` returns 8760h0m0s exactly, error nil.
- [ ] `ParseDuration("6m")` returns 6m0s exactly, error nil (unchanged).
- [ ] `ParseDuration("5Q")` returns a non-nil error whose message names the offending unit.
- [ ] `ParseDuration("1Y2M3w4d5h6m7s")` returns a duration equal to the sum of each component (365d + 60d + 21d + 4d + 5h + 6m + 7s).
- [ ] `ParseTime("NOW-6M")` returns a timestamp where `Now().Sub(result)` is within 1 second of 180 days.
- [ ] All existing tests in `time_duration_test.go` and `time_parse-time_test.go` still pass, with the single exception of any assertion that previously expected `1M` to equal 1 minute (that assertion is updated to expect 30 days).
- [ ] `make precommit` passes (format, vet, errcheck, vulncheck, addlicense, test, race).
- [ ] CHANGELOG.md has an `Unreleased` entry combining `feat:` (new units) and `fix:` (case-sensitive parsing), triggering a minor version bump on release.

## Verification

```bash
cd ~/Documents/workspaces/time
make precommit
go test -run 'TestDuration|TestParseTime' -race ./...
```

The `make precommit` target is the canonical pre-merge check (per `CLAUDE.md`). The targeted `go test` invocation is a faster inner-loop iteration during development.

After release, the natural verification is that a downstream consumer can revert their workaround (e.g. switch a cron back from `NOW-26W` to `NOW-6M`) and observe identical span behavior. This is consumer-side work and out of scope for this spec.

## Open Questions

- **Month/year duration semantics**: 30 days and 365 days are the obvious choices (consistent with the existing `Day` and `Week` constants in `UnitMap`). Calendar-aware arithmetic would require a reference time, breaking the pure-`Duration` signature. Recommend: ship the constant version; if calendar-aware is ever needed, it belongs on `DateTime`/`Date`, not `Duration`.
- **Implementation shape — drop `ToLower` vs. split into two maps**: either approach satisfies the contract. Dropping `ToLower` is simpler but means typo-shaped inputs (`6H` instead of `6h`) start returning errors where they used to work. Splitting the map (lowercase-friendly subset + case-sensitive `M`/`Y`) preserves typo tolerance for the lowercase set. The prompt that lands the fix should make this decision explicit; both are acceptable from the spec's perspective.
- **Should `Year` mean 365 or 365.25 days?** Probably 365 for consistency with what most cron users mean. 365.25 introduces a fractional-day boundary that's surprising in a `Duration` context. Document the choice in the CHANGELOG.

## Related

- `time_duration.go:32-41`, `:43-45`, `:109` — the three bug sites.
- `time_duration_test.go` — the existing `DescribeTable` to extend.
- `bborbe/trading` PR #119 — the consumer-side workaround that mitigates the immediate impact.
- `bborbe/trading` PR #121 — README documentation update reflecting the chosen workaround.
- Obsidian: `[[Switch Trading Backtest Cron to Week-Unit Syntax]]` — the task page recording the discovery, decision, and verification.
