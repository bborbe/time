<objective>
Add encoding.TextMarshaler and encoding.TextUnmarshaler interfaces to the Date type. This enables YAML (gopkg.in/yaml.v3) and other encoding formats to serialize Date as "2026-03-04" (date-only) instead of a full timestamp.
</objective>

<context>
Read time_date.go — the Date type and existing MarshalJSON/UnmarshalJSON methods.
The Date type already has:
- `MarshalJSON` → outputs `"2026-03-04"` using `stdtime.DateOnly`
- `UnmarshalJSON` → parses via ParseTime which supports DateOnly, RFC3339, NOW expressions
- `String()` → returns `d.Format(stdtime.DateOnly)`
- `MarshalBinary` / `DateFromBinary` for binary encoding

Read time_date_test.go — understand existing test patterns (Ginkgo/Gomega, BeforeEach/JustBeforeEach style).
</context>

<phase1>
## Phase 1: Regression Tests (BEFORE any production code changes)

Add regression tests to time_date_test.go that lock down current JSON marshaling behavior. These must pass BEFORE Phase 2.

1. **JSON MarshalJSON regression** — verify exact output for:
   - Non-zero Date → `"2023-06-19"` (quoted string)
   - Zero Date → `null`
   - `*Date` nil pointer → `null`

2. **JSON UnmarshalJSON regression** — verify round-trip for:
   - `"2023-06-19"` → correct Date, marshal back → identical string
   - `null` → zero Date
   - `""` → zero Date
   - RFC3339 `"2023-06-19T07:56:34Z"` → date-only (2023-06-19)

3. **JSON struct regression with omitempty** — verify exact JSON output for struct with:
   - `Date` field (zero) → `null`
   - `*Date` field (nil) → `null`
   - `Date` field with `omitempty` (zero) → omitted from output
   - `*Date` field with `omitempty` (nil) → omitted from output
   - `Date` field (non-zero) → `"2023-06-19"`
   - `*Date` field (non-zero) → `"2023-06-19"`

4. **JSON round-trip regression** — marshal a struct with Date fields to JSON, unmarshal back, verify equality.

After adding these tests, run `make test` and confirm ALL tests pass (existing + new regression tests). Do NOT proceed to Phase 2 until this succeeds.
</phase1>

<phase2>
## Phase 2: Add TextMarshaler/TextUnmarshaler (AFTER Phase 1 passes)

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

3. Add compile-time interface checks in `time_date.go`:
   ```go
   var _ encoding.TextMarshaler = Date{}
   var _ encoding.TextUnmarshaler = (*Date)(nil)
   ```

4. Add new tests in `time_date_test.go`:
   - **MarshalText**: non-zero date → `[]byte("2023-06-19")`, zero date → `nil, nil`
   - **UnmarshalText**: `"2023-06-19"` → correct Date, empty string → zero Date, RFC3339 string → correct Date (date portion only)
   - **YAML round-trip**: struct with `Date` and `*Date` fields + yaml tags, marshal to YAML, verify output contains `date: 2023-06-19`, unmarshal back, verify equality
   - **YAML omitempty**: zero Date with `omitempty` yaml tag → field omitted, nil `*Date` with `omitempty` → field omitted

5. Run `make test` — ALL tests must pass (Phase 1 regression + Phase 2 new + all existing).

6. **Critical verification**: Re-run the Phase 1 JSON regression tests mentally — confirm they still pass unchanged. The addition of TextMarshaler MUST NOT alter JSON marshaling behavior. encoding/json prefers MarshalJSON over MarshalText when both exist, so this should be safe, but the regression tests prove it.
</phase2>

<constraints>
- Do NOT modify existing MarshalJSON/UnmarshalJSON methods
- Do NOT change the Date type definition
- Do NOT change any existing tests — only ADD new test contexts
- Follow existing code style and patterns in time_date.go and time_date_test.go
- Phase 1 must pass before starting Phase 2
</constraints>

<verification>
Run `make test` after Phase 1 (regression tests only).
Run `make test` after Phase 2 (all tests).
Both must pass with zero failures.
</verification>

<success_criteria>
1. All existing tests pass unchanged.
2. JSON regression tests prove marshaling behavior is identical before and after TextMarshaler addition.
3. YAML round-trip works: `Date` fields serialize as `2023-06-19` (not full timestamp).
4. encoding.TextMarshaler and encoding.TextUnmarshaler compile-time checks pass.
</success_criteria>
