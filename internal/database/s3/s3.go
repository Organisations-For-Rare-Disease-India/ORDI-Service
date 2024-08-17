package s3

import (
	"ORDI/internal/database"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type s3Service struct {
	uploader *s3manager.Uploader
}

var s3instance *s3Service

var (
	s3Region = os.Getenv("REGION")
	s3bucket = os.Getenv("S3BUCKET")
)

func NewS3ServiceConnection() database.Service {
	// Reuse Connection
	if s3instance != nil {
		return s3instance
	}

	// Initialize the S3 session once
	s3Sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s3Region),
	})
	if err != nil {
		log.Fatalf("Error creating AWS session: %v", err)
	}

	s3instance = &s3Service{
		uploader: s3manager.NewUploader(s3Sess),
	}

	return s3instance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *s3Service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	health := make(map[string]string)

	// Check if S3 bucket is available and we have access to it or not
	_, err := s.uploader.S3.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s3bucket),
	})

	if err != nil {
		health["status"] = "down"
		health["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("S3 down: %v", err)) // Log the error and terminate the program
		return health
	}

	// Database is up, add more statistics
	health["status"] = "up"
	health["message"] = "S3 bucket is present"

	return health
}

func (s *s3Service) Close() error {
	// nothing to close
	return nil
}
