# Time

[![Go Reference](https://pkg.go.dev/badge/github.com/bborbe/time.svg)](https://pkg.go.dev/github.com/bborbe/time)
[![Go Report Card](https://goreportcard.com/badge/github.com/bborbe/time)](https://goreportcard.com/report/github.com/bborbe/time)

A comprehensive Go library providing enhanced time and date utilities with type safety, dependency injection support, and extended functionality beyond the standard `time` package.

## Features

- **🎯 Type-Safe Time Operations** - Strongly typed wrappers around Go's time types
- **💉 Dependency Injection Ready** - Interfaces for testable time operations
- **📅 Rich Time Types** - DateTime, Date, TimeOfDay, Duration, UnixTime, and Timezone types
- **⏱️ Extended Duration Support** - Parse human-readable durations like "1w2d3h4m5s"
- **🔄 JSON Marshaling** - Built-in JSON support for all time types
- **✅ Validation** - Input validation with meaningful error messages
- **🧪 Testing Utilities** - Helper functions for controlled time in tests
- **🌍 Timezone Handling** - Enhanced timezone operations with caching

## Installation

```bash
go get github.com/bborbe/time
```

## Quick Start

### Basic Usage

```go
import libtime "github.com/bborbe/time"

// Parse various time formats
dateTime, _ := libtime.ParseDateTime(ctx, "2023-12-25T15:30:00Z")
date, _ := libtime.ParseDate(ctx, "2023-12-25")
duration, _ := libtime.ParseDuration(ctx, "1w2d3h") // 1 week, 2 days, 3 hours

// Create time types
now := libtime.DateTime(time.Now())
unixTime := libtime.UnixTime(1703520600)
timeOfDay, _ := libtime.ParseTimeOfDay(ctx, "15:30:45")
```

### Dependency Injection Pattern

```go
type Service struct {
    currentDateTime libtime.CurrentDateTime
}

// Production code
func NewService() *Service {
    return &Service{
        currentDateTime: libtime.NewCurrentDateTime(),
    }
}

func (s *Service) ProcessOrder() {
    now := s.currentDateTime.Now()
    // Use now for timestamps...
}

// Test code
func TestService(t *testing.T) {
    service := &Service{
        currentDateTime: libtime.NewCurrentDateTime(),
    }
    
    // Control time in tests
    fixedTime := libtimetest.ParseDateTime("2023-12-25T00:00:00Z")
    service.currentDateTime.SetNow(fixedTime)
    
    // Test with predictable time
    service.ProcessOrder()
}
```

## Core Types

### DateTime
Enhanced `time.Time` wrapper with validation and JSON support:

```go
dt := libtime.DateTime(time.Now())
json, _ := dt.MarshalJSON()           // RFC3339Nano format
formatted := dt.Format("2006-01-02") // Standard Go formatting
ptr := dt.Ptr()                      // Get pointer for optional fields
```

### Duration
Extended duration with weeks and days support:

```go
// Parse human-readable durations
duration, _ := libtime.ParseDuration(ctx, "2w3d4h30m")  // 2 weeks, 3 days, 4.5 hours
duration, _ := libtime.ParseDuration(ctx, "1.5h")       // 1.5 hours
duration, _ := libtime.ParseDuration(ctx, "30s")        // 30 seconds

// Use constants
totalTime := libtime.Week + 2*libtime.Day + 3*libtime.Hour
```

### Date
Date-only type without time component:

```go
date, _ := libtime.ParseDate(ctx, "2023-12-25")
tomorrow := date.AddDate(0, 0, 1)
isWeekend := date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
```

### TimeOfDay
Time component without date:

```go
timeOfDay, _ := libtime.ParseTimeOfDay(ctx, "15:30:45")
hour := timeOfDay.Hour()
minute := timeOfDay.Minute()
second := timeOfDay.Second()
```

### UnixTime
Unix timestamp handling:

```go
unixTime := libtime.UnixTime(1703520600)
dateTime := unixTime.DateTime()
formatted := unixTime.String() // RFC3339Nano format
```

## Testing Support

The library provides extensive testing utilities in the `/test` package:

```go
import libtimetest "github.com/bborbe/time/test"

// Parse without error handling (panics on error - use only in tests)
dt := libtimetest.ParseDateTime("2023-12-25T15:30:00Z")
date := libtimetest.ParseDate("2023-12-25")
duration := libtimetest.ParseDuration("1h30m")

// Control time in dependency injection
currentDateTime := libtime.NewCurrentDateTime()
currentDateTime.SetNow(libtimetest.ParseDateTime("2023-12-25T00:00:00Z"))
```

## Advanced Features

### Validation
All types implement validation interfaces:

```go
dateTime, err := libtime.ParseDateTime(ctx, "invalid")
if err != nil {
    // Handle validation error with context
}
```

### Interfaces for Polymorphism

```go
// HasTime interface - for types containing time information
var hasTime libtime.HasTime = libtime.DateTime(time.Now())
timeValue := hasTime.Time()

// HasDuration interface - for types representing durations  
var hasDuration libtime.HasDuration = libtime.Duration(time.Hour)
durationValue := hasDuration.Duration()
```

### Timezone Operations

```go
tz, _ := libtime.ParseTimezone(ctx, "America/New_York")
location := tz.Location()
dateTimeInTZ := dateTime.In(location)
```

## Development

### Running Tests
```bash
make test                    # Run all tests with coverage
go test -cover -race ./...   # Manual test execution
```

### Code Quality
```bash
make precommit              # Complete workflow (format, test, check, etc.)
make format                 # Format code
make check                  # Static analysis
```

### Mock Generation
```bash
make generate               # Generate mocks using counterfeiter
```

## Dependencies

- **Runtime**: `github.com/bborbe/collection`, `github.com/bborbe/errors`, `github.com/bborbe/parse`, `github.com/bborbe/validation`
- **Testing**: Ginkgo v2, Gomega, Counterfeiter
- **Development**: goimports-reviser, addlicense, govulncheck

## License

BSD-style license. See [LICENSE](LICENSE) file for details.
