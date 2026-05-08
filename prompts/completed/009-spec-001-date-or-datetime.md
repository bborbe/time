---
status: completed
spec: [001-add-date-or-datetime-type]
summary: Created DateOrDateTime type with full API surface (marshaling, comparison, arithmetic, discriminators) and comprehensive Ginkgo/Gomega tests covering all 40 required scenarios.
container: time-009-spec-001-date-or-datetime
dark-factory-version: v0.156.1-1-g04f3863-dirty
created: "2026-05-08T17:00:00Z"
queued: "2026-05-08T16:46:32Z"
started: "2026-05-08T16:46:33Z"
completed: "2026-05-08T16:50:27Z"
branch: dark-factory/add-date-or-datetime-type
---

<summary>
- A new exported type `DateOrDateTime` is added to `github.com/bborbe/time`, sitting alongside the existing `Date`, `DateTime`, and `UnixTime` types.
- The type accepts either `YYYY-MM-DD` or RFC3339 input on unmarshal, delegating directly to the existing `ParseTime` function.
- Marshaling follows a round-trip rule: values that are exactly midnight UTC serialize as `YYYY-MM-DD`; all others serialize as RFC3339Nano.
- The zero value marshals as `nil, nil` (matching `Date.MarshalText` convention) and `IsZero()` returns true for it.
- The full common API surface is implemented: date components, time components, identity, marshaling, conversion, comparison, and arithmetic — mirroring `Date` and `DateTime` method for method.
- Two discriminators are added: `IsDateOnly() bool` and sibling conversions `AsDate() Date` and `AsDateTime() DateTime`.
- Compile-time interface assertions confirm `encoding.TextMarshaler` and `encoding.TextUnmarshaler` are satisfied.
- Unit tests cover: date-only input, RFC3339 input, midnight-UTC round-trip as date, non-midnight round-trip as RFC3339, empty input, zero-value marshal, parse errors, and every method in the common API surface.
- No existing files are modified; `make precommit` passes with all existing tests unchanged.
</summary>

<objective>
Implement `DateOrDateTime` as a first-class date primitive in `github.com/bborbe/time`, enabling any project to import one canonical type instead of vendoring private copies. The type accepts both date-only and full datetime inputs and serializes back in whichever form preserves original intent.
</objective>

<context>
Read CLAUDE.md for project conventions and build commands.
Read `time_date.go` — study every method: constructors (`ParseDate`, `DatePtr`, `ToDate`, `NewDate`), type definition, compile-time assertions, all methods including `MarshalText`, `UnmarshalText`, `MarshalJSON`, `UnmarshalJSON`, `MarshalBinary`, `Compare`, `ComparePtr`, `Before`, `After`, `Equal`, `EqualPtr`, `Truncate`, `Add`, `Sub`, `AddDate`, `AddTime`, `UTC`, `IsZero`, `Clone`, `ClonePtr`, `Time`, `TimePtr`, `Format`, `Unix`, `UnixMicro`, `Weekday`, `Year`, `Month`, `Day`, `String`, `Validate`, `Ptr`.
Read `time_date-time.go` — study the same methods on `DateTime`, in particular: `Hour`, `Minute`, `Second`, `Nanosecond`, the `String()` implementation (RFC3339Nano), and `MarshalJSON`/`MarshalText` pattern.
Read `time_unix-time.go` lines 275-278 — the `DateTime() DateTime` sibling conversion for the pattern to apply for `AsDate()` and `AsDateTime()`.
Read `time_date_test.go` — study the full Ginkgo/Gomega test structure used for `Date` (contexts, BeforeEach, JustBeforeEach, nested Describe/Context/It blocks, import alias `libtime "github.com/bborbe/time"`).
Read `time_date-time_test.go` — for additional test patterns covering time components.
Read `time_suite_test.go` — for the test suite bootstrap (do NOT add a new suite file; the existing one covers all `_test.go` files in the package).
Read `go-patterns.md` in `~/.claude/plugins/marketplaces/coding/docs/` for interface → constructor → struct conventions and error wrapping rules.
Read `go-testing-guide.md` in `~/.claude/plugins/marketplaces/coding/docs/` for Ginkgo/Gomega patterns and coverage requirements.
Read `test-pyramid-triggers.md` in `~/.claude/plugins/marketplaces/coding/docs/` for which test types to write for each code change.
Read `docs/dod.md` for Definition of Done criteria.
</context>

<requirements>
## File 1: `time_date-or-date-time.go`

Create this file at the repo root (alongside `time_date.go`). Start with the BSD license header exactly as it appears in `time_date.go`.

Package: `package time`

### Imports
```go
import (
    "context"
    "encoding"
    "encoding/json"
    "strings"
    stdtime "time"

    "github.com/bborbe/errors"
    "github.com/bborbe/parse"
    "github.com/bborbe/validation"
)
```

### Plural slice type

```go
type DateOrDateTimes []DateOrDateTime

func (d DateOrDateTimes) Interfaces() []interface{} {
    result := make([]interface{}, len(d))
    for i, ss := range d {
        result[i] = ss
    }
    return result
}

func (d DateOrDateTimes) Strings() []string {
    result := make([]string, len(d))
    for i, ss := range d {
        result[i] = ss.String()
    }
    return result
}
```

### Constructor functions

```go
func DateOrDateTimeFromBinary(ctx context.Context, value []byte) (*DateOrDateTime, error) {
    var t stdtime.Time
    if err := t.UnmarshalBinary(value); err != nil {
        return nil, errors.Wrapf(ctx, err, "unmarshalBinary failed")
    }
    return DateOrDateTime(t).Ptr(), nil
}

func ParseDateOrDateTimeDefault(ctx context.Context, value interface{}, defaultValue DateOrDateTime) DateOrDateTime {
    result, err := ParseDateOrDateTime(ctx, value)
    if err != nil {
        return defaultValue
    }
    return *result
}

func ParseDateOrDateTime(ctx context.Context, value interface{}) (*DateOrDateTime, error) {
    str, err := parse.ParseString(ctx, value)
    if err != nil {
        return nil, errors.Wrapf(ctx, err, "parse value failed")
    }
    t, err := ParseTime(ctx, str)
    if err != nil {
        return nil, errors.Wrapf(ctx, err, "parse time failed")
    }
    return DateOrDateTimePtr(t), nil
}

func DateOrDateTimePtr(value *stdtime.Time) *DateOrDateTime {
    if value == nil {
        return nil
    }
    return DateOrDateTime(*value).Ptr()
}

// NewDateOrDateTime creates a DateOrDateTime representing the date and time specified by the given parameters.
// It wraps the standard library's time.Date function with the same parameter signature.
func NewDateOrDateTime(
    year int,
    month stdtime.Month,
    day, hour, min, sec, nsec int,
    loc *stdtime.Location,
) DateOrDateTime {
    return DateOrDateTime(stdtime.Date(year, month, day, hour, min, sec, nsec, loc))
}
```

### Type definition and compile-time assertions

```go
type DateOrDateTime stdtime.Time

var _ encoding.TextMarshaler = DateOrDateTime{}

var _ encoding.TextUnmarshaler = (*DateOrDateTime)(nil)
```

### Unexported helper

Add this unexported helper function (not a method) immediately after the compile-time assertions:

```go
// isMidnightUTC reports whether t is exactly midnight UTC (all time components zero in UTC).
// This is the key discriminator for the round-trip serialization rule.
func isMidnightUTC(t stdtime.Time) bool {
    u := t.UTC()
    return u.Hour() == 0 && u.Minute() == 0 && u.Second() == 0 && u.Nanosecond() == 0
}
```

### Date component methods

```go
func (d DateOrDateTime) Year() int {
    return d.Time().Year()
}

func (d DateOrDateTime) Month() stdtime.Month {
    return d.Time().Month()
}

func (d DateOrDateTime) Day() int {
    return d.Time().Day()
}

func (d DateOrDateTime) Weekday() Weekday {
    return Weekday(d.Time().Weekday())
}
```

### Time component methods

```go
func (d DateOrDateTime) Hour() int {
    return d.Time().Hour()
}

func (d DateOrDateTime) Minute() int {
    return d.Time().Minute()
}

func (d DateOrDateTime) Second() int {
    return d.Time().Second()
}

func (d DateOrDateTime) Nanosecond() int {
    return d.Time().Nanosecond()
}
```

### Identity methods

```go
func (d DateOrDateTime) String() string {
    t := d.Time()
    if t.IsZero() {
        return ""
    }
    if isMidnightUTC(t) {
        return t.UTC().Format(stdtime.DateOnly)
    }
    return t.Format(stdtime.RFC3339Nano)
}

func (d DateOrDateTime) Validate(ctx context.Context) error {
    if d.Time().IsZero() {
        return errors.Wrapf(ctx, validation.Error, "time is zero")
    }
    return nil
}

func (d DateOrDateTime) Ptr() *DateOrDateTime {
    return &d
}

// IsZero reports whether d represents the zero time instant.
func (d DateOrDateTime) IsZero() bool {
    return d.Time().IsZero()
}

func (d DateOrDateTime) UTC() DateOrDateTime {
    return DateOrDateTime(d.Time().UTC())
}

func (d DateOrDateTime) Clone() DateOrDateTime {
    return d
}

func (d *DateOrDateTime) ClonePtr() *DateOrDateTime {
    if d == nil {
        return nil
    }
    return d.Clone().Ptr()
}
```

### Marshaling methods

```go
func (d *DateOrDateTime) UnmarshalJSON(b []byte) error {
    str := strings.Trim(string(b), `"`)
    switch str {
    case "", "null":
        *d = DateOrDateTime(stdtime.Time{})
        return nil
    default:
        t, err := ParseTime(context.Background(), str)
        if err != nil {
            return errors.Wrapf(context.Background(), err, "parse time failed")
        }
        *d = DateOrDateTime(*t)
        return nil
    }
}

func (d DateOrDateTime) MarshalJSON() ([]byte, error) {
    t := d.Time()
    if t.IsZero() {
        return json.Marshal(nil)
    }
    if isMidnightUTC(t) {
        return json.Marshal(t.UTC().Format(stdtime.DateOnly))
    }
    return json.Marshal(t.Format(stdtime.RFC3339Nano))
}

func (d DateOrDateTime) MarshalText() ([]byte, error) {
    t := d.Time()
    if t.IsZero() {
        return nil, nil
    }
    if isMidnightUTC(t) {
        return []byte(t.UTC().Format(stdtime.DateOnly)), nil
    }
    return []byte(t.Format(stdtime.RFC3339Nano)), nil
}

func (d *DateOrDateTime) UnmarshalText(b []byte) error {
    str := string(b)
    if len(str) == 0 {
        *d = DateOrDateTime(stdtime.Time{})
        return nil
    }
    t, err := ParseTime(context.Background(), str)
    if err != nil {
        return errors.Wrapf(context.Background(), err, "parse time failed")
    }
    *d = DateOrDateTime(*t)
    return nil
}

func (d DateOrDateTime) MarshalBinary() ([]byte, error) {
    return d.Time().MarshalBinary()
}
```

### Conversion methods

```go
func (d DateOrDateTime) Time() stdtime.Time {
    return stdtime.Time(d)
}

func (d *DateOrDateTime) TimePtr() *stdtime.Time {
    if d == nil {
        return nil
    }
    t := stdtime.Time(*d)
    return &t
}

func (d DateOrDateTime) Format(layout string) string {
    return d.Time().Format(layout)
}

func (d DateOrDateTime) Unix() int64 {
    return d.Time().Unix()
}

func (d DateOrDateTime) UnixMicro() int64 {
    return d.Time().UnixMicro()
}
```

### Comparison methods

```go
func (d DateOrDateTime) Compare(other DateOrDateTime) int {
    return Compare(d.Time(), other.Time())
}

func (d *DateOrDateTime) ComparePtr(other *DateOrDateTime) int {
    if d == nil && other == nil {
        return 0
    }
    if d == nil {
        return -1
    }
    if other == nil {
        return 1
    }
    return d.Compare(*other)
}

func (d DateOrDateTime) Before(other HasTime) bool {
    return d.Time().Before(other.Time())
}

func (d DateOrDateTime) After(other HasTime) bool {
    return d.Time().After(other.Time())
}

func (d DateOrDateTime) Equal(other DateOrDateTime) bool {
    return d.Time().Equal(other.Time())
}

func (d *DateOrDateTime) EqualPtr(other *DateOrDateTime) bool {
    if d == nil && other == nil {
        return true
    }
    if d != nil && other != nil {
        return d.Equal(*other)
    }
    return false
}
```

### Arithmetic methods

```go
func (d DateOrDateTime) Add(duration HasDuration) DateOrDateTime {
    return DateOrDateTime(d.Time().Add(duration.Duration()))
}

func (d DateOrDateTime) Sub(other HasTime) Duration {
    return Duration(d.Time().Sub(other.Time()))
}

func (d DateOrDateTime) AddDate(years int, months int, days int) DateOrDateTime {
    return DateOrDateTime(d.Time().AddDate(years, months, days))
}

// Deprecated: Use AddDate instead.
// AddTime adds the given years, months, and days to the DateOrDateTime but will be removed in future versions.
func (d DateOrDateTime) AddTime(years int, months int, days int) DateOrDateTime {
    return d.AddDate(years, months, days)
}

func (d DateOrDateTime) Truncate(duration HasDuration) DateOrDateTime {
    return DateOrDateTime(d.Time().Truncate(duration.Duration()))
}
```

### Discriminator and sibling conversions

```go
// IsDateOnly reports whether d represents a date-only value (midnight UTC).
func (d DateOrDateTime) IsDateOnly() bool {
    return !d.IsZero() && isMidnightUTC(d.Time())
}

// AsDate returns the date component of d as a Date value.
func (d DateOrDateTime) AsDate() Date {
    return ToDate(d.Time())
}

// AsDateTime returns d as a DateTime value.
func (d DateOrDateTime) AsDateTime() DateTime {
    return DateTime(d.Time())
}
```

---

## File 2: `time_date-or-date-time_test.go`

Create this file at the repo root. Use `package time_test` (external test package matching all other `_test.go` files). Do NOT add a new suite bootstrap — `time_suite_test.go` already covers the whole package.

Start with the BSD license header exactly as it appears in `time_date_test.go`.

Imports pattern (mirror `time_date_test.go`):
```go
import (
    "bytes"
    "context"
    "encoding/json"
    "time"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    libtime "github.com/bborbe/time"
)
```

### Required test coverage

Write a top-level `Describe("DateOrDateTime", ...)` block. Inside it, cover all of the following — use the exact same nested `Context` / `BeforeEach` / `JustBeforeEach` / `It` pattern as `time_date_test.go`:

1. **MarshalBinary & DateOrDateTimeFromBinary** — round-trip: marshal a value, unmarshal it back, confirm the Unix timestamp matches.

2. **MarshalJSON — midnight UTC** — produces `"2026-01-15"` (date-only format).

3. **MarshalJSON — non-midnight UTC** — produces an RFC3339Nano string.

4. **MarshalJSON — zero value** — produces `null`.

5. **UnmarshalJSON — date-only input** — `"2026-01-15"` unmarshals and re-marshals as `"2026-01-15"`.

6. **UnmarshalJSON — RFC3339 input** — `"2026-01-15T14:30:00Z"` unmarshals and re-marshals as RFC3339Nano.

7. **UnmarshalJSON — empty string** — yields zero value, no error.

8. **UnmarshalJSON — "null"** — yields zero value, no error.

9. **UnmarshalJSON — invalid input** — returns non-nil error.

10. **MarshalText — midnight UTC** — produces `2026-01-15` (no quotes, it's []byte).

11. **MarshalText — non-midnight** — produces RFC3339Nano bytes.

11a. **MarshalText — midnight in non-UTC zone** — `2026-01-15T00:00:00+02:00` (which is NOT midnight in UTC) serializes as RFC3339Nano, not date-only. Verifies the spec's Failure Mode contract: midnight-UTC check is strictly UTC.

12. **MarshalText — zero value** — produces `nil, nil`.

13. **UnmarshalText — date-only input** — succeeds, value round-trips as date-only.

14. **UnmarshalText — RFC3339 input** — succeeds, value round-trips as RFC3339Nano.

15. **UnmarshalText — empty bytes** — yields zero value, no error.

16. **IsDateOnly** — true for midnight-UTC, false for non-midnight, false for zero value.

17. **IsZero** — true for zero value, false for non-zero.

18. **AsDate** — midnight-UTC value converts to the correct `Date`; verify with `.String()`.

19. **AsDateTime** — any value converts to `DateTime`; verify with `.Time().Equal(...)`.

20. **String** — midnight UTC returns `YYYY-MM-DD`; non-midnight returns RFC3339Nano; zero returns `""`.

21. **Validate** — zero value returns non-nil error; non-zero returns nil.

22. **UTC** — value in non-UTC zone converts to UTC.

23. **Clone / ClonePtr** — cloned value equals original; ClonePtr of nil pointer returns nil.

24. **Year / Month / Day / Weekday** — verify with a known date (e.g. `2026-01-15` → year 2026, month January, day 15, weekday Thursday).

25. **Hour / Minute / Second / Nanosecond** — for midnight-UTC all return 0; for a datetime value returns correct components.

26. **Time / TimePtr** — `Time()` returns correct `time.Time`; `TimePtr()` on nil pointer returns nil.

27. **Format** — `.Format(time.DateOnly)` on a midnight value returns `"2026-01-15"`.

28. **Unix / UnixMicro** — verify with a known epoch-based value.

29. **Compare / ComparePtr** — earlier < 0, equal = 0, later > 0; ComparePtr handles nil cases.

30. **Before / After / Equal / EqualPtr** — standard comparison cases; EqualPtr handles both-nil case.

31. **Add** — add 1 hour to a midnight value, confirm it is no longer midnight.

32. **Sub** — difference between two values equals expected duration.

33. **AddDate** — add 1 year/1 month/1 day, confirm each component increments.

34. **Truncate** — truncate a non-midnight value to 24h, confirm result is midnight UTC for that day.

35. **JSON struct round-trip** — define an inline struct with a `DateOrDateTime` field tagged `json:"date"`, marshal and unmarshal it, confirm the value survives.

36. **ParseDateOrDateTime — valid date string** — returns non-nil pointer, no error.

37. **ParseDateOrDateTime — valid RFC3339 string** — returns non-nil pointer, no error.

38. **ParseDateOrDateTime — invalid string** — returns nil, non-nil error.

39. **ParseDateOrDateTimeDefault — invalid input** — returns the supplied default.

40. **DateOrDateTimePtr — nil input** — returns nil.
</requirements>

<constraints>
- Do NOT modify any existing files — only create `time_date-or-date-time.go` and `time_date-or-date-time_test.go`.
- Existing `Date` and `DateTime` public APIs must not change.
- No new entries in `go.mod` direct dependencies — use only stdlib plus the packages already imported in `time_date.go` (`github.com/bborbe/errors`, `github.com/bborbe/parse`, `github.com/bborbe/validation`).
- License header and package name must match other files in this repo.
- Tests use Ginkgo/Gomega per `docs/dod.md` — do NOT add a new suite file.
- The round-trip rule (midnight-UTC → date-only, else RFC3339Nano) is part of the public contract and must be implemented exactly.
- Zero-value marshal returns `nil, nil` from `MarshalText` — matching `Date.MarshalText` convention.
- Do NOT commit — dark-factory handles git.
</constraints>

<verification>
Run `make precommit` — must pass (format + lint + test + security).

Additionally verify:
```bash
# Confirm new file exists with correct package
grep -n "^package time$" time_date-or-date-time.go

# Confirm compile-time assertions are present
grep -n "var _ encoding" time_date-or-date-time.go

# Confirm no new direct deps added to go.mod
grep "^require" go.mod

# Confirm existing tests still pass
go test -count=1 -run "Date[^O]" ./...
```
</verification>
