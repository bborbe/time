---
status: created
created: "2026-03-07T21:50:00Z"
---

<summary>
- Adds Before, After, and Truncate methods to the Date type
- Adds Equal and EqualPtr methods to the Date type
- Adds Compare and ComparePtr methods to the UnixTime type
- Follows the same patterns already used by DateTime which has all these methods
- Enables direct comparisons without calling .Time() first
- Adds tests for all new methods on both types
</summary>

<objective>
Align the `Date` and `UnixTime` APIs with `DateTime` by adding missing comparison and utility methods. Currently `Date` lacks `Before`, `After`, `Truncate`, `Equal`, and `EqualPtr`, while `UnixTime` lacks `Compare` and `ComparePtr`. Adding these eliminates the need to call `.Time()` for common operations.
</objective>

<context>
Read CLAUDE.md for project conventions.
Read `time_date.go` — the `Date` type (line 81), `Compare` (line 179), `Add` (line 196), `Sub` (line 200).
Read `time_date-time.go` — `DateTime.Before` (line 236), `DateTime.After` (line 240), `DateTime.Equal` (line 121), `DateTime.EqualPtr` (line 125), `DateTime.Compare` (line 252), `DateTime.ComparePtr` (line 256), `DateTime.Truncate` (line 269) for the patterns to follow.
Read `time_unix-time.go` — `UnixTime.Equal` (line 134), `UnixTime.EqualPtr` (line 138), `UnixTime.Before` (line 219), `UnixTime.After` (line 223), `UnixTime.Sub` (line 231), `UnixTime.Truncate` (line 243) for existing methods. Note: `UnixTime` already has `Equal`, `EqualPtr`, `Before`, `After`, `Sub`, `Truncate` — only `Compare` and `ComparePtr` are missing.
Read `time_date_test.go` for existing Date test patterns.
Read `time_unix-time_test.go` for existing UnixTime test patterns.
Read `~/.claude-yolo/docs/go-testing.md` for Ginkgo patterns.
</context>

<requirements>
1. In `time_date.go`, add five methods after the existing `ComparePtr` method (line 194):

   ```go
   func (d Date) Before(other HasTime) bool {
       return d.Time().Before(other.Time())
   }

   func (d Date) After(other HasTime) bool {
       return d.Time().After(other.Time())
   }

   func (d Date) Equal(other Date) bool {
       return d.Time().Equal(other.Time())
   }

   func (d *Date) EqualPtr(other *Date) bool {
       if d == nil && other == nil {
           return true
       }
       if d != nil && other != nil {
           return d.Equal(*other)
       }
       return false
   }

   func (d Date) Truncate(duration HasDuration) Date {
       return Date(d.Time().Truncate(duration.Duration()))
   }
   ```

   - `Before`/`After` use `HasTime` interface parameter — matches `DateTime` pattern at lines 236 and 240.
   - `Equal`/`EqualPtr` use concrete `Date` parameter — matches `DateTime.Equal`/`EqualPtr` at lines 121 and 125.
   - `Truncate` uses `HasDuration` interface parameter — matches `DateTime.Truncate` at line 269.

2. In `time_unix-time.go`, add two methods after the existing `Truncate` method (line 245):

   ```go
   func (u UnixTime) Compare(other UnixTime) int {
       return Compare(u.Time(), other.Time())
   }

   func (u *UnixTime) ComparePtr(other *UnixTime) int {
       if u == nil && other == nil {
           return 0
       }
       if u == nil {
           return -1
       }
       if other == nil {
           return 1
       }
       return u.Compare(*other)
   }
   ```

   - Uses `Compare` function — matches `DateTime.Compare`/`ComparePtr` at lines 252 and 256.
   - Uses `Date.Compare`/`ComparePtr` at lines 179 and 183 as additional reference.

3. In `time_date_test.go`, add tests for the five new Date methods. Follow existing Ginkgo `Describe`/`It` patterns:
   - `Before`: earlier returns true, later returns false, same returns false
   - `After`: later returns true, earlier returns false, same returns false
   - `Equal`: same date returns true, different date returns false
   - `EqualPtr`: both nil returns true, both non-nil same returns true, one nil returns false
   - `Truncate`: truncating to 24h returns same date

4. In `time_unix-time_test.go`, add tests for the two new UnixTime methods:
   - `Compare`: earlier returns -1, later returns 1, same returns 0
   - `ComparePtr`: both nil returns 0, nil receiver returns -1, nil argument returns 1, both non-nil delegates to Compare
</requirements>

<constraints>
- Do NOT change any existing methods or signatures
- Do NOT change any files other than `time_date.go`, `time_unix-time.go`, `time_date_test.go`, `time_unix-time_test.go`
- Follow the exact interface/concrete parameter patterns from `DateTime`
- Existing tests must still pass
- Do NOT commit — dark-factory handles git
</constraints>

<verification>
Run `make precommit` — must pass.
</verification>
