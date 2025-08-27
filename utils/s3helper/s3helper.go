package s3helper

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Helper holds the configuration for S3 operations.
type S3Helper struct {
	ProfileName	string
	BucketName	string
	EndpointURL	string
	Region			string
}

func (u *S3Helper) UploadFile(filePath, s3Path string) error {
	// Create a new AWS session with the specified profile and region
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(u.Region),
			Endpoint:    aws.String(u.EndpointURL),
			Credentials: credentials.NewSharedCredentials("", u.ProfileName),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	// Create an S3 service client
	s3Client := s3.New(sess)

	// Open the file for reading
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %v", filePath, err)
	}
	defer file.Close()

	// Get the file info to obtain the size
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for %q: %v", filePath, err)
	}

	// Clean the S3 path (remove leading/trailing slashes)
	s3Path = strings.TrimPrefix(filepath.Clean(s3Path), "/")

	// Determine the content type based on the file extension
	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType == "" {
		contentType = "application/octet-stream" // Default content type
	}

	// Upload the file to S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(u.BucketName),
		Key:           aws.String(s3Path),
		Body:          file,
		ContentLength: aws.Int64(fileInfo.Size()),
		ContentType:   aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %v", err)
	}

	log.Printf("Successfully uploaded %q to s3://%s/%s", filePath, u.BucketName, s3Path)
	return nil
}

// ListFiles lists all files in the specified S3 path prefix
func (u *S3Helper) ListFiles(prefix string) ([]string, error) {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(u.Region),
			Endpoint:    aws.String(u.EndpointURL),
			Credentials: credentials.NewSharedCredentials("", u.ProfileName),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	s3Client := s3.New(sess)

	var files []string
	err = s3Client.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket: aws.String(u.BucketName),
		Prefix: aws.String(prefix),
	}, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			files = append(files, *obj.Key)
		}
		return !lastPage
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list files: %v", err)
	}

	return files, nil
}

// DeleteFile deletes a file from S3
func (u *S3Helper) DeleteFile(s3Path string) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(u.Region),
			Endpoint:    aws.String(u.EndpointURL),
			Credentials: credentials.NewSharedCredentials("", u.ProfileName),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	s3Client := s3.New(sess)

	_, err = s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(u.BucketName),
		Key:    aws.String(s3Path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file %q: %v", s3Path, err)
	}

	log.Printf("Successfully deleted s3://%s/%s", u.BucketName, s3Path)
	return nil
}

// DownloadFile downloads a file from S3 to the local filesystem
func (u *S3Helper) DownloadFile(s3Path, localPath string) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String(u.Region),
			Endpoint:    aws.String(u.EndpointURL),
			Credentials: credentials.NewSharedCredentials("", u.ProfileName),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	s3Client := s3.New(sess)

	// Get the object from S3
	result, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(u.BucketName),
		Key:    aws.String(s3Path),
	})
	if err != nil {
		return fmt.Errorf("failed to get object %q from S3: %v", s3Path, err)
	}
	defer result.Body.Close()

	// Create the directory for the local file if it doesn't exist
	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %q: %v", dir, err)
	}

	// Create the local file
	file, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file %q: %v", localPath, err)
	}
	defer file.Close()

	// Copy the S3 object content to the local file
	_, err = file.ReadFrom(result.Body)
	if err != nil {
		return fmt.Errorf("failed to write to local file %q: %v", localPath, err)
	}

	log.Printf("Successfully downloaded s3://%s/%s to %s", u.BucketName, s3Path, localPath)
	return nil
}
