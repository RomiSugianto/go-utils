# Go Utils

A collection of reusable utility packages for Go applications. This repository is designed to help you quickly add common functionality—such as logging, configuration, and more—to your Go projects.

## Installation

To use a utility package from this module, install it with:

```bash
go get github.com/romisugianto/go-utils/utils/{utils}
```

Replace `{utils}` with any other utility package you want to use as your library grows.

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

### Housekeeper

A simple and effective housekeeper for Go applications.

#### Usage

```go
package main

import (
    "github.com/romisugianto/go-utils/utils/housekeeper"
)

func main() {
    hk, err := housekeeper.NewHousekeeper("myApp")
    if err != nil {
        panic(err)
    }
    defer hk.Close()

		dir := "./test"

    if err := hk.HousekeepFilesByAge(dir, 1); err != nil {
        hk.logger.Error("HousekeepFilesByAge failed: %v", err)
    }

		if err := hk.HousekeepFilesByCount(dir, 5); err != nil {
        hk.logger.Error("HousekeepFilesByCount failed: %v", err)
    }
}
```

#### Housekeeper Methods

- **NewHousekeeper(appName string)**: Creates a new housekeeper instance. If no application name is provided, it defaults to "script".
- **HousekeepFilesByAge(dir string, maxAgeDays int)**: Manages the housekeeping of files in a directory based on their age.
- **HousekeepFilesByCount(dir string, maxFiles int)**: Manages the housekeeping of files in a directory based on a maximum count.

### Splitter

A simple and effective splitter for Go applications.

#### Usage

```go
package main

import (
    "github.com/romisugianto/go-utils/utils/splitter"
)

func main() {
    sp, err := splitter.NewSplitter("myApp")
    if err != nil {
        panic(err)
    }
    defer sp.Close()

    testFile := "./test.txt"
    outputDir := "./output"
    processedDir := "./processed"

    if err := sp.SplitFileByLines(testFile, 2, outputDir, processedDir); err != nil {
        sp.logger.Error("SplitFileByLines failed: %v", err)
    }
}
```

#### Splitter Methods

- **NewSplitter(appName string)**: Creates a new splitter instance. If no application name is provided, it defaults to "script".
- **SplitFileByLines(filePath string, linesPerFile int, outputDir string, processedDir string)**: Splits a file into multiple files based on the number of lines specified.