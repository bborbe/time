# Changelog

All notable changes to this project will be documented in this file.

Please choose versions by [Semantic Versioning](http://semver.org/).

* MAJOR version when you make incompatible API changes,
* MINOR version when you add functionality in a backwards-compatible manner, and
* PATCH version when you make backwards-compatible bug fixes.

## v1.21.0

- update go and deps

## v1.20.0
- Add GitHub Actions workflows for CI, code review, and automated testing
- Add golangci-lint configuration for enhanced code quality checks
- Add security scanning tools: gosec, osv-scanner, and trivy
- Improve Makefile with additional security and quality checks
- Integrate golines for automatic line length formatting (max 100 characters)
- Update goimports-reviser to v3 for better import organization
- Update Go module dependencies for testing and tooling
- Improve code formatting across all range validation methods

## v1.19.2
- Add Max() method to DateTimeRanges, UnixTimeRanges, and DateRanges for finding maximum encompassing range
- Add Min() method to DateTimeRanges, UnixTimeRanges, and DateRanges for finding minimum overlapping range

## v1.19.1
- Add AddDate() method to all types (DateTime, Date, UnixTime) to align with Go standard library naming conventions
- Deprecate AddTime() methods (use AddDate() instead)

## v1.19.0
- Add time range constructor functions for all range types (DayTimeRange, WeekTimeRange, MonthTimeRange, QuarterTimeRange, YearTimeRange)
- Add range constructor functions for DateRange (DayDateRange, WeekDateRange, MonthDateRange, QuarterDateRange, YearDateRange)
- Add range constructor functions for DateTimeRange (DayDateTimeRange, WeekDateTimeRange, MonthDateTimeRange, QuarterDateTimeRange, YearDateTimeRange)
- Add range constructor functions for UnixTimeRange (DayUnixTimeRange, WeekUnixTimeRange, MonthUnixTimeRange, QuarterUnixTimeRange, YearUnixTimeRange)
- Add TimeRange() conversion methods to DateRange, DateTimeRange, and UnixTimeRange types
- Add public period boundary calculation functions (BeginningOfDay, EndOfDay, BeginningOfWeek, EndOfWeek, etc.)
- Add comprehensive test coverage for all new range constructor functionality with consistency validation
- Support timezone preservation in all range calculations

## v1.18.1
- Fix WaiterUntil to handle equal times correctly (no longer waits when until time equals current time)
- Remove 10-second buffer from WaiterUntil for more precise timing behavior
- Add SpecTimeout to waiter tests to prevent hanging and improve test reliability
- Add support for int64 and *int64 in ParseDuration for direct nanosecond value handling
- Add support for uppercase and mixed case duration parsing (e.g., "1H30M", "1h30M")
- Add .PHONY directive to Makefile test target
- Update dependencies to latest versions

## v1.18.0
- Add IsZero method to DateTime, UnixTime, and Date types
- Add DateTimeRange, DateRange, and UnixTimeRange types with validation
- Add comprehensive tests for all new IsZero methods
- Add proper godoc comments following Go documentation best practices
- Implement validation framework patterns for range type validation

## v1.17.1

- add test.ParseDuration utility function

## v1.17.0

- add counterfeiter mocks for all interfaces
- add Layout types and parsing functionality
- go mod update

## v1.16.3

- add test.ParseDate

## v1.16.2

- add tests
- go mod update

## v1.16.1

- add example to readme
- go mod update

## v1.16.0

- add test package with utils for testing
- go mod update

## v1.15.2

- add UnixTime.Before and UnixTime.after
- go mod update

## v1.15.1

- add HasDuration and HasTime interfaces

## v1.15.0

- add LoadLocation with cache
- add ParseLocation

## v1.14.2

- add UTC() and Weekday() to Date, DateTime and UnixTime 

## v1.14.1
 
- add Weekdays.Weekdays()

## v1.14.0

- add Weekday, Weekdays, ParseWeekday, ParseWeekdays

## v1.13.1

- fix TimePtr on nil Date, DateTime or UnixTime
- add AddTime to Date, DateTime or UnixTime

## v1.13.0

- remove vendor
- go mod update

## v1.12.1

- add CurrentDateTimeGetterFunc and CurrentTimeGetterFunc
- go mod update

## v1.12.0

- refactor WaiterUntil
- add WaiterDuration

## v1.11.6

- allow NOW-1d
- go mod update

## v1.11.5

- add UnixTime.DateTime

## v1.11.4

- add UnixTime.Truncate

## v1.11.3

- add DateTime.Truncate

## v1.11.2

- improve ParseDuration

## v1.11.1

- DateTime.Add use Duration

## v1.11.0

- add CurrentDateTime

## v1.10.0

- clean Duration.String() (1h15m0s => 1h15m) 

## v1.9.1

- add Time.Sub 
- add Duration.Abs

## v1.9.0

- Duration.String output now days and weeks like: 10w5d23h59m30s
- go mod update

## v1.8.2

- add list types

## v1.8.1

- allow parse timeOfDay with seconds

## v1.8.0

- allow unmarshal NOW
- go mod update

## v1.7.4

- add parse time without seconds
- go mod update

## v1.7.3

- Add MarshalBinary and UnmarshalBinary
- go mod update

## v1.7.2

- add Date.Add() and UnixTime.Add()

## v1.7.1

- add Year(), Month(), Day(), Hour(), Minute(), Second() and Nanosecond()

## v1.7.0

- add ParseUnixTime
- go mod update

## v1.6.2

- add Before, After and Equal to TimeOfDay

## v1.6.1

- add Duration.String()

## v1.6.0

- add ParseDateTimeDefault
- add ParseDurationDefault
- add ParseTimeDefault
- add ParseTimeOfDayDefault

## v1.5.2

- remove error from DateTime
- add Time

## v1.5.1

- add DateTime to TimeOfDay

## v1.5.0

- add Duration
- marshal unmarshal Duration from duration string or number
- go mod update

## v1.4.2

- test marshal Date and DateTime

## v1.4.1

- fix parse empty Date and DateTime

## v1.4.0

- add DateTime
- add UnixTime
- test Date
- go mod update

## v1.3.0

- Add compare time
- go mod update

## v1.2.0

- Allow ParseTimeOfDay with Timezone. Example '13:37 Europe/Berlin'
- go mod update

## v1.1.1

- fix ParseDuration

## v1.1.0

- Add ParseDuration with support for d=day and w=week

## v1.0.0

- Initial Version
