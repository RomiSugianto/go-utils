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
    "github.com/romisugianto/go-utils/utils/logger"
    "github.com/romisugianto/go-utils/utils/housekeeper"
)

func main() {
    logger.NewLogger("myApp")
    hk, err := housekeeper.NewHousekeeper(logger)
    if err != nil {
        panic(err)
    }

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

- **NewHousekeeper(appName string)**: Creates a new housekeeper instance.
- **HousekeepFilesByAge(dir string, maxAgeDays int), recursive ...bool)**: Manages the housekeeping of files in a directory based on their age, recrusive or not.
- **HousekeepFilesByCount(dir string, maxFiles int)**: Manages the housekeeping of files in a directory based on a maximum count.

### Splitter

A simple and effective splitter for Go applications.

#### Usage

```go
package main

import (
    "github.com/romisugianto/go-utils/utils/logger"
    "github.com/romisugianto/go-utils/utils/splitter"
)

func main() {
    logger.NewLogger("myApp")
    sp, err := splitter.NewSplitter(logger)
    if err != nil {
        panic(err)
    }

    testFile := "./test.txt"
    outputDir := "./output"
    processedDir := "./processed"

    if err := sp.SplitFileByLines(testFile, 2, outputDir, processedDir); err != nil {
        sp.logger.Error("SplitFileByLines failed: %v", err)
    }
}
```

#### Splitter Methods

- **NewSplitter(appName string)**: Creates a new splitter instance.
- **SplitFileByLines(filePath string, linesPerFile int, outputDir string, processedDir string)**: Splits a file into multiple files based on the number of lines specified.

### S3Helper

A simple and effective AWS S3 utility for Go applications. Provides operations for uploading, downloading, listing, and deleting files from Amazon S3.

#### Prerequisites

Before using S3Helper, ensure you have AWS credentials configured. You can set them up using:

```bash
# Install AWS CLI
brew install awscli

# Configure AWS credentials
aws configure
```

#### Usage

```go
package main

import (
    "github.com/romisugianto/go-utils/utils/s3helper"
)

func main() {
    // Initialize S3Helper with your AWS configuration
    s3 := s3helper.S3Helper{
        ProfileName: "default",      // AWS profile name
        BucketName:  "your-bucket",   // S3 bucket name
        EndpointURL: "https://s3.amazonaws.com", // S3 endpoint
        Region:      "us-west-2",    // AWS region
    }

    // Upload a file to S3
    err := s3.UploadFile("local-file.txt", "s3/path/file.txt")
    if err != nil {
        panic(err)
    }

    // List files in S3
    files, err := s3.ListFiles("s3/path/")
    if err != nil {
        panic(err)
    }

    // Download a file from S3
    err = s3.DownloadFile("s3/path/file.txt", "downloaded-file.txt")
    if err != nil {
        panic(err)
    }

    // Delete a file from S3
    err = s3.DeleteFile("s3/path/file.txt")
    if err != nil {
        panic(err)
    }
}
```

#### S3Helper Methods

- **UploadFile(filePath string, s3Path string) error**: Uploads a local file to the specified S3 path.
- **DownloadFile(s3Path string, localPath string) error**: Downloads a file from S3 to the local filesystem.
- **ListFiles(prefix string) ([]string, error)**: Lists all files in the specified S3 path prefix.
- **DeleteFile(s3Path string) error**: Deletes a file from S3.

#### Configuration Fields

- **ProfileName**: AWS profile name (defaults to "default")
- **BucketName**: S3 bucket name (required)
- **EndpointURL**: S3 endpoint URL (defaults to AWS standard endpoints)
- **Region**: AWS region (required)
