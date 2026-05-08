---
status: committing
summary: Added Clone and ClonePtr methods to Date and UnixTime types, mirroring the DateTime pattern, with Ginkgo/Gomega tests covering zero, non-zero, and nil receiver cases.
container: time-008-add-clone-to-date-and-unixtime
dark-factory-version: v0.156.1-1-g04f3863-dirty
created: "2026-05-08T16:37:45Z"
queued: "2026-05-08T16:37:45Z"
started: "2026-05-08T16:37:58Z"
---
<summary>
- Adds Clone and ClonePtr methods to the Date type
- Adds Clone and ClonePtr methods to the UnixTime type
- Mirrors the existing DateTime.Clone / DateTime.ClonePtr implementation exactly
- ClonePtr returns nil when called on a nil receiver
- Adds Ginkgo/Gomega tests covering zero value, non-zero value, and nil receiver
- No changes to DateTime or any other type
</summary>

<objective>
Align the `Date` and `UnixTime` APIs with `DateTime` by adding the missing `Clone()` and `ClonePtr()` methods. The pattern already exists on `DateTime` and should be copied verbatim with the receiver type swapped.
</objective>

<context>
Read CLAUDE.md for project conventions.
Read `docs/dod.md` for testing conventions (Ginkgo/Gomega).
Read `time_date-time.go` — `DateTime.Clone` and `DateTime.ClonePtr` (around line 217-226) for the exact pattern to mirror, including doc-comment style (or lack thereof — note these methods currently have no doc comments) and nil-receiver handling.
Read `time_date.go` — locate `Date.Ptr` (around line 110) and `Date.ComparePtr` to see existing receiver-style and method placement.
Read `time_unix-time.go` — locate `UnixTime.Ptr` (around line 159) and surrounding methods to see existing receiver-style and method placement.
Read `time_date_test.go` and `time_unix-time_test.go` for existing Ginkgo `Describe`/`Context`/`It` test patterns.
</context>

<requirements>
1. In `time_date.go`, add the following two methods (place them adjacent to other utility methods such as `MarshalBinary` / `Compare`, matching the relative position of `Clone`/`ClonePtr` in `time_date-time.go`):

   ```go
   func (d Date) Clone() Date {
       return d
   }

   func (d *Date) ClonePtr() *Date {
       if d == nil {
           return nil
       }
       return d.Clone().Ptr()
   }
   ```

   - Match the doc-comment style of `DateTime.Clone` / `DateTime.ClonePtr` exactly (if those have no doc comments, do not add any here either).
   - Use the existing `Ptr()` method on `Date` (defined around line 110).

2. In `time_unix-time.go`, add the equivalent methods:

   ```go
   func (u UnixTime) Clone() UnixTime {
       return u
   }

   func (u *UnixTime) ClonePtr() *UnixTime {
       if u == nil {
           return nil
       }
       return u.Clone().Ptr()
   }
   ```

   - Match the doc-comment style of `DateTime.Clone` / `DateTime.ClonePtr` exactly.
   - Use the existing `Ptr()` method on `UnixTime` (defined around line 159).

3. In `time_date_test.go`, add Ginkgo tests for the two new `Date` methods. Follow the existing `Describe`/`Context`/`It` patterns in the file:
   - `Clone` on a zero value returns a zero `Date`
   - `Clone` on a non-zero value returns an equal `Date` (use `Equal` or compare via `Time()`)
   - `ClonePtr` on a nil `*Date` returns nil
   - `ClonePtr` on a non-nil `*Date` returns a non-nil pointer to an equal `Date`
   - `ClonePtr` returns a different pointer than the receiver (i.e. is a copy, not the same address)

4. In `time_unix-time_test.go`, add Ginkgo tests for the two new `UnixTime` methods covering the same four cases:
   - `Clone` on a zero value
   - `Clone` on a non-zero value
   - `ClonePtr` on a nil `*UnixTime` returns nil
   - `ClonePtr` on a non-nil `*UnixTime` returns a non-nil pointer to an equal `UnixTime` with a different address
</requirements>

<constraints>
- Do NOT change any existing methods, signatures, or files other than `time_date.go`, `time_unix-time.go`, `time_date_test.go`, `time_unix-time_test.go`.
- Do NOT add `Clone`/`ClonePtr` to any other type (`DateRange`, `DateTimeRange`, `Duration`, `DateOrDateTime`, etc.) — out of scope.
- Doc-comment style must match `DateTime.Clone` / `DateTime.ClonePtr` exactly.
- Tests must use Ginkgo/Gomega per `docs/dod.md`.
- Existing tests must still pass.
- Do NOT commit — dark-factory handles git.
</constraints>

<verification>
Run `make precommit` — must pass.
</verification>
