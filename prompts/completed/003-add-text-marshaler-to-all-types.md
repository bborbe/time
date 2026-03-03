---
status: completed
summary: Added encoding.TextMarshaler/TextUnmarshaler to DateTime, UnixTime, Duration, TimeOfDay with comprehensive JSON and YAML regression tests
container: time-003-add-text-marshaler-to-all-types
dark-factory-version: v0.14.1
created: "2026-03-03T20:12:12Z"
queued: "2026-03-03T20:12:12Z"
started: "2026-03-03T20:12:12Z"
completed: "2026-03-03T20:22:57Z"
---
<objective>
Add encoding.TextMarshaler and encoding.TextUnmarshaler interfaces to DateTime, UnixTime, Duration, and TimeOfDay types. Date already has these (use as reference pattern). Also add comprehensive 4-field struct regression and round-trip tests for ALL types (Date, DateTime, UnixTime, Duration, TimeOfDay) covering JSON and YAML marshal/unmarshal.
</objective>

<context>
Read these files first:
- time_date.go — Date already has MarshalText/UnmarshalText (reference pattern)
- time_date_test.go — Date already has some tests (reference pattern)
- time_date-time.go — DateTime type, has MarshalJSON/UnmarshalJSON, uses RFC3339Nano format
- time_unix-time.go — UnixTime type, has MarshalJSON/UnmarshalJSON, marshals as integer (Unix seconds)
- time_duration.go — Duration type, has MarshalJSON/UnmarshalJSON, marshals as Go duration string
- time_time-of-day.go — TimeOfDay type, has MarshalJSON/UnmarshalJSON, marshals as "15:04" string

Key formats per type:
- Date: "2023-06-19" (DateOnly)
- DateTime: "2023-06-19T07:56:34Z" (RFC3339Nano)
- UnixTime: JSON=integer (1687161394), Text=RFC3339Nano string ("2023-06-19T07:56:34Z")
- Duration: Go duration string ("1h30m")
- TimeOfDay: "15:04" layout string

Note on UnixTime: JSON marshals as integer, but TextMarshaler must produce a string. Use RFC3339Nano format for text representation (same as String() method). UnmarshalText should accept RFC3339Nano strings.
</context>

<phase1>
## Phase 1: Regression Tests (BEFORE any production code changes)

For EACH of the 5 types (Date, DateTime, UnixTime, Duration, TimeOfDay), add a test context called "struct marshal regression" that tests a struct with 4 fields:

```go
type testStruct struct {
    Field         T     `json:"field" yaml:"field"`
    FieldPtr      *T    `json:"fieldPtr" yaml:"fieldPtr"`
    FieldOmit     T     `json:"fieldOmit,omitempty" yaml:"fieldOmit,omitempty"`
    FieldPtrOmit  *T    `json:"fieldPtrOmit,omitempty" yaml:"fieldPtrOmit,omitempty"`
}
```

Test cases per type (define inline structs, don't create named types):

### A. JSON Marshal — non-zero values set for Field and FieldPtr only
- Verify exact JSON output string
- Zero/nil fields produce `null` or are omitted with omitempty

### B. JSON Unmarshal — round-trip
- Marshal the struct to JSON, unmarshal back, verify equality of non-zero fields
- Verify zero fields remain zero after round-trip

### C. JSON Marshal — all fields set (non-zero)
- All 4 fields have values
- Verify all appear in output

### D. JSON Marshal — all fields zero/nil
- Verify omitempty fields are omitted, non-omitempty produce null/zero

After adding ALL regression tests for all 5 types, run `make test`. ALL must pass before Phase 2.
</phase1>

<phase2>
## Phase 2: Add TextMarshaler/TextUnmarshaler to DateTime, UnixTime, Duration, TimeOfDay

Use the existing Date implementation as the reference pattern.

### DateTime (time_date-time.go)

```go
func (d DateTime) MarshalText() ([]byte, error) {
    t := d.Time()
    if t.IsZero() {
        return nil, nil
    }
    return []byte(t.Format(stdtime.RFC3339Nano)), nil
}

func (d *DateTime) UnmarshalText(b []byte) error {
    str := string(b)
    if len(str) == 0 {
        *d = DateTime(stdtime.Time{})
        return nil
    }
    t, err := ParseTime(context.Background(), str)
    if err != nil {
        return errors.Wrapf(context.Background(), err, "parse time failed")
    }
    *d = DateTime(*t)
    return nil
}
```

Add compile-time checks:
```go
var _ encoding.TextMarshaler = DateTime{}
var _ encoding.TextUnmarshaler = (*DateTime)(nil)
```

### UnixTime (time_unix-time.go)

```go
func (u UnixTime) MarshalText() ([]byte, error) {
    t := u.Time()
    if t.IsZero() {
        return nil, nil
    }
    return []byte(t.Format(stdtime.RFC3339Nano)), nil
}

func (u *UnixTime) UnmarshalText(b []byte) error {
    str := string(b)
    if len(str) == 0 {
        *u = UnixTime(stdtime.Time{})
        return nil
    }
    t, err := ParseTime(context.Background(), str)
    if err != nil {
        return errors.Wrapf(context.Background(), err, "parse time failed")
    }
    *u = UnixTime(*t)
    return nil
}
```

Add compile-time checks:
```go
var _ encoding.TextMarshaler = UnixTime{}
var _ encoding.TextUnmarshaler = (*UnixTime)(nil)
```

### Duration (time_duration.go)

```go
func (d Duration) MarshalText() ([]byte, error) {
    if d.Duration() == 0 {
        return nil, nil
    }
    return []byte(d.Duration().String()), nil
}

func (d *Duration) UnmarshalText(b []byte) error {
    str := string(b)
    if len(str) == 0 {
        *d = Duration(0)
        return nil
    }
    ctx := context.Background()
    duration, err := ParseDuration(ctx, str)
    if err != nil {
        return errors.Wrapf(ctx, err, "parse duration failed")
    }
    *d = *duration
    return nil
}
```

Add compile-time checks:
```go
var _ encoding.TextMarshaler = Duration(0)
var _ encoding.TextUnmarshaler = (*Duration)(nil)
```

### TimeOfDay (time_time-of-day.go)

```go
func (t TimeOfDay) MarshalText() ([]byte, error) {
    return []byte(t.String()), nil
}

func (t *TimeOfDay) UnmarshalText(b []byte) error {
    str := string(b)
    if len(str) == 0 {
        *t = TimeOfDay{}
        return nil
    }
    parsed, err := ParseTimeOfDay(context.Background(), str)
    if err != nil {
        return errors.Wrapf(context.Background(), err, "parse time of day failed")
    }
    *t = *parsed
    return nil
}
```

Add compile-time checks:
```go
var _ encoding.TextMarshaler = TimeOfDay{}
var _ encoding.TextUnmarshaler = (*TimeOfDay)(nil)
```

### Add "encoding" import

Each file needs `"encoding"` added to the import block.

### Phase 2 Tests

For EACH of the 5 types, add YAML tests using the same 4-field struct pattern:

### E. YAML Marshal — non-zero values for Field and FieldPtr only
- Verify YAML output contains correct format
- Zero/nil fields with omitempty are omitted

### F. YAML Unmarshal — round-trip
- Marshal to YAML, unmarshal back, verify equality

### G. YAML Marshal — all fields set
- All 4 fields have values, all appear in output

### H. YAML Marshal — all fields zero/nil
- omitempty fields omitted

Also add unit tests for MarshalText/UnmarshalText directly:
- Non-zero value → correct bytes
- Zero value → nil
- UnmarshalText with valid string → correct value
- UnmarshalText with empty string → zero value

Run `make test` — ALL tests must pass (Phase 1 JSON regression + Phase 2 YAML + all existing).
</phase2>

<constraints>
- Do NOT modify existing MarshalJSON/UnmarshalJSON methods
- Do NOT change any type definitions
- Do NOT change any existing tests — only ADD new test contexts
- Follow existing code style in each file
- Add `"encoding"` import where needed
- Phase 1 must pass before starting Phase 2
- Use the same test value consistently per type:
  - Date: time.Unix(1687161394, 0) → "2023-06-19"
  - DateTime: time.Unix(1687161394, 0) → "2023-06-19T07:56:34Z"
  - UnixTime: time.Unix(1687161394, 0) → "2023-06-19T07:56:34Z" (text) / 1687161394 (json)
  - Duration: 1h30m → "1h30m0s"
  - TimeOfDay: {Hour: 13, Minute: 37} → "13:37"
</constraints>

<verification>
Run `make test` after Phase 1 (regression tests only).
Run `make test` after Phase 2 (all tests).
Both must pass with zero failures.
</verification>

<success_criteria>
1. All existing tests pass unchanged.
2. JSON regression tests prove marshaling behavior is identical for all 5 types.
3. YAML round-trip works for all 5 types with correct format.
4. encoding.TextMarshaler and encoding.TextUnmarshaler compile-time checks for DateTime, UnixTime, Duration, TimeOfDay.
5. All 5 types serialize correctly in YAML (not full timestamp or unexpected format).
</success_criteria>
