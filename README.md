# Time

Time and date utils.

## Examples

### NewCurrentDateTime

```go
import (
    libtime "github.com/bborbe/time"
    libtimetest "github.com/bborbe/time/test"
)
dt := libtime.NewCurrentDateTime()
dt.SetNow(libtimetest.ParseDateTime("2006-01-02 15:04:05"))
fmt.Println(dt.Format("2006-01-02 15:04:05"))
```

### NewCurrentTime

```go
import (
    libtime "github.com/bborbe/time"
    libtimetest "github.com/bborbe/time/test"
)
dt := libtime.NewCurrentTime()
dt.SetNow(libtimetest.ParseTime("2006-01-02 15:04:05"))
fmt.Println(dt.Format("2006-01-02 15:04:05"))
```
