package s3helper

import (
	"mime"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadFile(t *testing.T) {
	// Skip this test if AWS credentials are not available
	if _, err := os.Stat("/Users/romi/.aws/credentials"); os.IsNotExist(err) {
		t.Skip("AWS credentials not found, skipping upload test")
	}

	uploader := S3Helper{
		ProfileName: "default",
		BucketName:  "your-bucket-name",
		EndpointURL: "https://s3.amazonaws.com",
		Region:      "us-west-2",
	}

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	// Write some content to the temp file
	if _, err := tempFile.WriteString("This is a test file."); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	// Close the temp file before uploading
	if err := tempFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	// Attempt to upload the file
	err = uploader.UploadFile(tempFile.Name(), "test/testfile.txt")
	if err != nil {
		t.Fatalf("failed to upload file: %v", err)
	}
}

func TestS3PathCleaning(t *testing.T) {
	// Test various path cleaning scenarios
	testCases := []struct {
		input    string
		expected string
	}{
		{"/path/to/file", "path/to/file"},
		{"path/to/file", "path/to/file"},
		{"//path//to//file", "path/to/file"},
		{"path/to/file/", "path/to/file"},
		{"/", ""},
		{"", ""},
	}

	for _, tc := range testCases {
		cleaned := cleanS3Path(tc.input)
		if cleaned != tc.expected {
			t.Errorf("cleanS3Path(%q) = %q, expected %q", tc.input, cleaned, tc.expected)
		}
	}
}

// Helper function to test S3 path cleaning (extracted from UploadFile logic)
func cleanS3Path(path string) string {
	cleaned := filepath.Clean(path)
	if cleaned == "." {
		return ""
	}
	return strings.TrimPrefix(cleaned, "/")
}

func TestContentTypeDetection(t *testing.T) {
	testCases := []struct {
		filename string
		expected string
	}{
		{"file.txt", "text/plain; charset=utf-8"},
		{"image.jpg", "image/jpeg"},
		{"image.jpeg", "image/jpeg"},
		{"image.png", "image/png"},
		{"document.pdf", "application/pdf"},
		{"data.json", "application/json"},
		{"unknown.xyz", "application/octet-stream"}, // This might vary by system
		{"file", "application/octet-stream"},
	}

	for _, tc := range testCases {
		contentType := detectContentType(tc.filename)
		// For unknown.xyz, accept either chemical/x-xyz or application/octet-stream
		if tc.filename == "unknown.xyz" {
			if contentType != "chemical/x-xyz" && contentType != "application/octet-stream" {
				t.Errorf("detectContentType(%q) = %q, expected either 'chemical/x-xyz' or 'application/octet-stream'", tc.filename, contentType)
			}
		} else if contentType != tc.expected {
			t.Errorf("detectContentType(%q) = %q, expected %q", tc.filename, contentType, tc.expected)
		}
	}
}

// Helper function to test content type detection
func detectContentType(filename string) string {
	contentType := mime.TypeByExtension(filepath.Ext(filename))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return contentType
}

func TestUploadFileWithDifferentContentTypes(t *testing.T) {
	// Skip this test if AWS credentials are not available
	if _, err := os.Stat("/Users/romi/.aws/credentials"); os.IsNotExist(err) {
		t.Skip("AWS credentials not found, skipping upload test")
	}

	uploader := S3Helper{
		ProfileName: "default",
		BucketName:  "test-bucket",
		EndpointURL: "https://s3.amazonaws.com",
		Region:      "us-west-2",
	}

	// Test with different file types
	testFiles := []struct {
		content  string
		filename string
		s3Path   string
	}{
		{"Hello World", "test.txt", "text/test.txt"},
		{"{}", "data.json", "json/data.json"},
		{"<html></html>", "page.html", "html/page.html"},
	}

	for _, tf := range testFiles {
		tempFile, err := os.CreateTemp("", tf.filename)
		if err != nil {
			t.Fatalf("failed to create temp file %s: %v", tf.filename, err)
		}
		defer os.Remove(tempFile.Name())

		if _, err := tempFile.WriteString(tf.content); err != nil {
			t.Fatalf("failed to write to temp file %s: %v", tf.filename, err)
		}
		tempFile.Close()

		err = uploader.UploadFile(tempFile.Name(), tf.s3Path)
		if err != nil {
			t.Fatalf("failed to upload file %s: %v", tf.filename, err)
		}
	}
}

func TestDownloadFile(t *testing.T) {
	// Skip this test if AWS credentials are not available
	if _, err := os.Stat("/Users/romi/.aws/credentials"); os.IsNotExist(err) {
		t.Skip("AWS credentials not found, skipping download test")
	}

	downloader := S3Helper{
		ProfileName: "default",
		BucketName:  "test-bucket",
		EndpointURL: "https://s3.amazonaws.com",
		Region:      "us-west-2",
	}

	// Create a temporary file for download destination
	tempDir, err := os.MkdirTemp("", "download-test")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	localPath := filepath.Join(tempDir, "downloaded-file.txt")

	// Attempt to download a file from S3
	err = downloader.DownloadFile("text/test.txt", localPath)
	if err != nil {
		t.Fatalf("failed to download file: %v", err)
	}

	// Verify the file was downloaded
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		t.Fatalf("downloaded file does not exist at %s", localPath)
	}

	// Read and verify the content
	content, err := os.ReadFile(localPath)
	if err != nil {
		t.Fatalf("failed to read downloaded file: %v", err)
	}

	// Check if the content is not empty (assuming the test file exists in S3)
	if len(content) == 0 {
		t.Errorf("downloaded file is empty")
	}
}
