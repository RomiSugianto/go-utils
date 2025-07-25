# Go Utils

A collection of reusable utility packages for Go applications. This repository is designed to help you quickly add common functionality—such as logging, configuration, and more—to your Go projects.

## Installation

To use a utility package from this module, install it with:

```bash
go get github.com/romisugianto/go-utils/utils/logger
```

Replace `logger` with any other utility package you want to use as your library grows.

## Packages

### Logger

A simple and effective logger for Go applications.

#### Usage

```go
package main

import (
    "github.com/romisugianto/go-utils/utils/logger"
)

func main() {
    log, err := logger.NewLogger("myApp")
    if err != nil {
        panic(err)
    }
    defer log.Close()

    log.Info("Application started")
    log.Warning("This is a warning message")
    log.Error("An error occurred")
}
```

#### Logger Methods

- **NewLogger(appName string)**: Creates a new logger instance. If no application name is provided, it defaults to "script".
- **Info(format string, args ...any)**: Logs an informational message.
- **Warning(format string, args ...any)**: Logs a warning message.
- **Error(format string, args ...any)**: Logs an error message.
- **Close() error**: Closes the logger's file handle.