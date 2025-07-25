# Go Logger Module

This repository contains a reusable logging module for Go applications. The `logger` package provides a simple and effective way to log messages at different levels (Info, Error, Warning) and manage log files.

## Installation

To use the logger module in your Go project, you can import it using the following command:

```bash
go get github.com/romisugianto/go-logger/logger
```

## Usage

Here’s a quick example of how to use the logger module in your application:

```go
package main

import (
	"github.com/romisugianto/go-logger/logger"
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

### Logger Methods

- **NewLogger(appName string)**: Creates a new logger instance. If no application name is provided, it defaults to "script".
- **Info(format string, args ...any)**: Logs an informational message.
- **Warning(format string, args ...any)**: Logs a warning message.
- **Error(format string, args ...any)**: Logs an error message.
- **Close() error**: Closes the logger's file handle.
- **GetLogFilePath() string**: Returns the path to the current log file.
- **DisplayCredits(banner string, appName string, appVersion string)**: Displays the application credits/banner.

## Log File Location

Log files are created in the `logs` directory within the project root. The log file is named using the application name and the current date.

## License

This project is licensed under the MIT License. See the LICENSE file for details.