package storage

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type Storage struct {
	bucket string
	client *s3.S3
}

func New(host, bucket string) *Storage {
	cfg := &aws.Config{
		Endpoint:         aws.String(host),
		Region:           aws.String("us-east-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewEnvCredentials(),
	}
	sess := session.Must(session.NewSession(cfg))

	return &Storage{
		bucket: bucket,
		client: s3.New(sess),
	}
}

func (s *Storage) Store(filename string, buf *bytes.Buffer) (url string, err error) {
	uploader := s3manager.NewUploaderWithClient(s.client)
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
		Body:   buf,
	})
	if err != nil {
		return "", err
	}

	return res.Location, nil
}

func (s *Storage) Get(filename string) (io.ReadCloser, error) {
	res, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func (s *Storage) Exist(filename string) bool {
	_, err := s.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return false
	}

	return true
}
