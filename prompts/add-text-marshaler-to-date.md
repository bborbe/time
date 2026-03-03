<objective>
Add encoding.TextMarshaler and encoding.TextUnmarshaler interfaces to the Date type. This enables YAML (gopkg.in/yaml.v3) and other encoding formats to serialize Date as "2026-03-04" (date-only) instead of a full timestamp.
</objective>

<context>
Read time_date.go — the Date type and existing MarshalJSON/UnmarshalJSON methods.
The Date type already has:
- `MarshalJSON` → outputs `"2026-03-04"` using `stdtime.DateOnly`
- `UnmarshalJSON` → parses via ParseTime which supports DateOnly, RFC3339, NOW expressions
- `String()` → returns `d.Format(stdtime.DateOnly)`
</context>

<requirements>
1. Add `MarshalText` method to `Date` in `time_date.go`:
   ```go
   func (d Date) MarshalText() ([]byte, error) {
       t := d.Time()
       if t.IsZero() {
           return nil, nil
       }
       return []byte(t.Format(stdtime.DateOnly)), nil
   }
   ```

2. Add `UnmarshalText` method to `*Date` in `time_date.go`:
   ```go
   func (d *Date) UnmarshalText(b []byte) error {
       str := string(b)
       if len(str) == 0 {
           *d = Date(stdtime.Time{})
           return nil
       }
       t, err := ParseTime(context.Background(), str)
       if err != nil {
           return errors.Wrapf(context.Background(), err, "parse time failed")
       }
       *d = ToDate(*t)
       return nil
   }
   ```

3. Add tests in `time_date_test.go` (or the appropriate existing test file):
   - MarshalText: non-zero date → `[]byte("2026-03-04")`, zero date → `nil, nil`
   - UnmarshalText: `"2026-03-04"` → correct Date, empty string → zero Date, RFC3339 string → correct Date (date portion only)
   - YAML round-trip test: define a struct with `Date` field + yaml tag, marshal to YAML, verify output is `date: 2026-03-04`, unmarshal back, verify equality

4. Add compile-time interface checks:
   ```go
   var _ encoding.TextMarshaler = Date{}
   var _ encoding.TextUnmarshaler = (*Date)(nil)
   ```
</requirements>

<constraints>
- Do NOT modify existing MarshalJSON/UnmarshalJSON methods
- Do NOT change the Date type definition
- Do NOT change any existing tests
- Follow existing code style and patterns in time_date.go
</constraints>

<verification>
Run: `make test`
Confirm: all tests pass including the new text marshaler tests.
</verification>

<success_criteria>
After this change, any code using gopkg.in/yaml.v3 with a `Date` field will serialize as `2026-03-04` instead of a full RFC3339 timestamp.
</success_criteria>
