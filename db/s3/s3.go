package s3

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// StorageByS3 embeds s3.Bucket struct
type StorageByS3 struct {
	*s3.S3
	Bucket *string
}

// NewStorageByS3 allocates and returns new S3Storage
func NewStorageByS3(c *aws.Config, bucket string) *StorageByS3 {
	return &StorageByS3{
		S3:     s3.New(session.New(c)),
		Bucket: aws.String(bucket),
	}
}

// Put put or update object into AWS S3 storage.
func (s *StorageByS3) Put(path string, data []byte, ct string) (err error) {
	_, err = s.PutObject(&s3.PutObjectInput{
		Bucket:      s.Bucket,
		Key:         aws.String(path),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(ct),
		Body:        bytes.NewReader(data),
	})
	return
}

// Del delete object from AWS S3 storage.
func (s *StorageByS3) Del(path string) (err error) {
	_, err = s.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s.Bucket,
		Key:    aws.String(path),
	})
	return
}
