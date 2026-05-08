---
status: prompted
tags:
    - dark-factory
    - spec
approved: "2026-05-08T16:37:45Z"
generating: "2026-05-08T16:37:57Z"
prompted: "2026-05-08T16:40:25Z"
branch: dark-factory/add-date-or-datetime-type
---

## Summary

- Add a polymorphic date/datetime type to `github.com/bborbe/time` that accepts either `YYYY-MM-DD` or RFC3339 on read.
- Round-trip rule: midnight-UTC values serialize as `YYYY-MM-DD`; everything else as RFC3339.
- Lives alongside existing `Date` and `DateTime` — does not replace them.
- Source of truth is the existing 37-LOC implementation in `vault-cli/pkg/domain/date_or_datetime.go` (port file + tests).
- Enables vault-cli, agent task-controller, OpenClaw, dark-factory, and trading services to share one primitive instead of redefining it.

## Problem

Multiple bborbe projects need to accept user-supplied date fields where the value can be either a calendar date (`2026-01-15`) or a precise instant (`2026-01-15T14:30:00Z`), and serialize back in whichever form preserves the original intent. Today this lives only in vault-cli as `pkg/domain/DateOrDateTime`. As more projects (agent, OpenClaw, dark-factory, trading services) grow the same need, each will either copy the file or invent its own incompatible variant. The `bborbe/time` library is the natural home — it already owns `Date`, `DateTime`, and the parser they all depend on.

## Goal

`github.com/bborbe/time` exposes `DateOrDateTime` as a third date primitive that consumer projects can import directly instead of vendoring their own copy.

## Non-goals

- Migrating vault-cli or any other consumer to the new type (separate work in those repos).
- Changing or deprecating existing `Date` or `DateTime`.
- Adding new parse formats beyond what `ParseTime` already supports.
- JSON-specific marshalers (text marshaler is sufficient — `encoding/json` falls back to it for string values).

## Desired Behavior

1. The library exposes a new exported type representing a value that is either a date or a datetime.
2. Unmarshaling text accepts any input that `ParseTime` accepts (notably `YYYY-MM-DD` and RFC3339); empty input yields the zero value without error.
3. Marshaling text returns `YYYY-MM-DD` when the value is exactly midnight UTC (hour, minute, second, nanosecond all zero in UTC); otherwise returns RFC3339.
4. The zero value behaves as zero (`IsZero()` returns true; round-trips through marshal/unmarshal cleanly).
5. The type implements the **full common API surface** shared by `Date`, `DateTime`, `UnixTime` so consumers treat it as a peer:
   - **Date components**: `Year`, `Month`, `Day`, `Weekday`
   - **Time components**: `Hour`, `Minute`, `Second`, `Nanosecond` (return 0 when value is date-only midnight-UTC)
   - **Identity**: `String`, `Validate(ctx)`, `Ptr`, `IsZero`, `UTC`, `Clone`, `ClonePtr`
   - **Marshaling**: `MarshalJSON`/`UnmarshalJSON`, `MarshalText`/`UnmarshalText`, `MarshalBinary`
   - **Conversion**: `Time`, `TimePtr`, `Format(layout)`, `Unix`, `UnixMicro`
   - **Comparison**: `Compare`/`ComparePtr`, `Before(HasTime)`, `After(HasTime)`, `Equal`, `EqualPtr`
   - **Arithmetic**: `Add(HasDuration)`, `Sub(HasTime)`, `AddDate(y,m,d)`, `AddTime(y,m,d)`, `Truncate(HasDuration)`
6. File and type naming are consistent with the existing `Date` and `DateTime` conventions in the repository.
7. Tests cover: date-only input, RFC3339 input, midnight-UTC round-trip as date, non-midnight round-trip as RFC3339, empty input, zero-value marshal, parse errors, AND every method in the common API surface (positive + zero-value cases).

## Constraints

- Existing `Date` and `DateTime` public APIs must not change.
- `make precommit` must pass (format + lint + test + security).
- No new external dependencies — implementation uses only stdlib plus existing `bborbe/time`, `bborbe/errors`, `bborbe/parse`, `bborbe/validation` already in `go.mod`.
- License header and package name match other files in this repo.
- Tests use Ginkgo/Gomega per `docs/dod.md`.
- The round-trip rule (midnight-UTC → date-only, else RFC3339) is part of the public contract — once shipped, changing it is a breaking change for consumers.

## Failure Modes

| Trigger | Expected behavior | Recovery |
|---------|-------------------|----------|
| Input string is neither `YYYY-MM-DD` nor RFC3339 nor a recognized `ParseTime` form | Unmarshal returns wrapped error | Caller surfaces parse error to user |
| Input is empty string / empty bytes | Unmarshal succeeds, leaves value at zero | None needed |
| Marshaling the zero value | Returns `nil, nil` matching `Date.MarshalText` repo convention | None needed |
| Value is midnight in non-UTC zone (e.g. `2026-01-15T00:00:00+02:00`) | Treated as non-midnight UTC, serializes as RFC3339 | Documented behavior — caller normalizes if they want date form |
| Non-zero sub-second precision | Serializes as RFC3339 with nanoseconds | None needed |

## Security / Abuse Cases

Not applicable — this is a parsing/formatting primitive. Input is text; `ParseTime` is the existing trust boundary and already used widely in the library.

## Resolved Design Decisions

1. **Type name**: `DateOrDateTime` (zero-friction migration for vault-cli; descriptive).
2. **File name**: `time_date-or-date-time.go` matching the `time_<concept>.go` repo pattern.
3. **Constructors**: match existing `Date` / `DateTime` constructor surface for symmetry.
4. **Sibling conversions**: include `AsDate() Date` and `AsDateTime() DateTime` (symmetric with `UnixTime.DateTime()`).
5. **`IsDateOnly() bool` discriminator**: included.
6. **Arithmetic form-preservation**: arithmetic results promote to RFC3339 if non-midnight-UTC. The serialization rule handles this naturally — no special-case logic.
7. **Zero-value marshal output**: returns `nil, nil` matching `Date.MarshalText` convention.
8. **`ParseTime` boundary**: type lives in `bborbe/time` package and calls `ParseTime` directly. The `libtime` alias from the vault-cli port is removed.

## Acceptance Criteria

- [ ] New type is exported from package `time` (i.e. `github.com/bborbe/time`).
- [ ] Type implements `encoding.TextMarshaler` and `encoding.TextUnmarshaler` (verified with compile-time `var _` assertions, matching the pattern used for `Date`).
- [ ] Marshaling a midnight-UTC value produces `YYYY-MM-DD`.
- [ ] Marshaling a non-midnight-UTC value produces RFC3339.
- [ ] Unmarshaling `YYYY-MM-DD` and RFC3339 both succeed and round-trip.
- [ ] Unmarshaling empty input succeeds and yields the zero value.
- [ ] Date components: `Year`, `Month`, `Day`, `Weekday`.
- [ ] Time components: `Hour`, `Minute`, `Second`, `Nanosecond` (return 0 for date-only midnight-UTC).
- [ ] Identity: `String`, `Validate(ctx)`, `Ptr`, `IsZero`, `UTC`, `Clone`, `ClonePtr`.
- [ ] Marshaling: `MarshalJSON`/`UnmarshalJSON`, `MarshalText`/`UnmarshalText`, `MarshalBinary`.
- [ ] Conversion: `Time`, `TimePtr`, `Format(layout)`, `Unix`, `UnixMicro`.
- [ ] Comparison: `Compare`/`ComparePtr`, `Before(HasTime)`, `After(HasTime)`, `Equal`, `EqualPtr`.
- [ ] Arithmetic: `Add(HasDuration)`, `Sub(HasTime)`, `AddDate(y,m,d)`, `AddTime(y,m,d)`, `Truncate(HasDuration)`.
- [ ] Sibling conversions: `AsDate() Date`, `AsDateTime() DateTime`, `IsDateOnly() bool`.
- [ ] File name follows the `time_<concept>.go` pattern.
- [ ] Tests cover all behaviors in "Desired Behavior" item 7.
- [ ] `make precommit` passes.
- [ ] Existing `Date` and `DateTime` tests still pass unchanged.
- [ ] No new entries in `go.mod` direct dependencies.

No scenario test required — this is a pure-Go primitive verifiable with unit tests.

## Verification

```
make precommit
```

Expected: format, lint, test, security checks all pass. No changes to existing test output for `Date` / `DateTime`.

## Do-Nothing Option

If we do nothing, vault-cli keeps its private copy and each new consumer either re-implements the type, copies the file, or invents an incompatible variant. Cost grows linearly with the number of consumers (agent, OpenClaw, dark-factory, trading services are all queued up). The library is the right home; deferring only multiplies the migration burden later.
