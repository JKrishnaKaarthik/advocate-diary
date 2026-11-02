package storage

import (
	"bytes"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	client  *s3.Client
	presign *s3.PresignClient
	bucket  string
}

func New(endpoint, region, bucket, accessKey, secretKey string, pathStyle bool) *S3 {
	client := s3.New(s3.Options{
		Region:       region,
		Credentials:  credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		BaseEndpoint: endpointPtr(endpoint),
		UsePathStyle: pathStyle,
	})
	return &S3{client: client, presign: s3.NewPresignClient(client), bucket: bucket}
}

func endpointPtr(endpoint string) *string {
	if endpoint == "" {
		return nil
	}
	return aws.String(endpoint)
}

func (s *S3) Put(ctx context.Context, key, contentType string, body []byte) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(contentType),
	})
	return err
}

func (s *S3) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (s *S3) PresignGet(ctx context.Context, key string) (string, error) {
	out, err := s.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(5*time.Minute))
	if err != nil {
		return "", err
	}
	return out.URL, nil
}
