package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zuyatna/shop-retail-employee-service/internal/config"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
	endpoint   string
	useSSL     bool
}

func NewMinioStorage(cfg *config.Config) (*MinioStorage, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, err
	}

	storage := &MinioStorage{
		client:     client,
		bucketName: cfg.MinioBucket,
		endpoint:   cfg.MinioEndpoint,
		useSSL:     cfg.MinioUseSSL,
	}

	if err := storage.initBucket(context.Background()); err != nil {
		return nil, err
	}

	return storage, nil
}

func (m *MinioStorage) initBucket(ctx context.Context) error {
	exists, err := m.client.BucketExists(ctx, m.bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = m.client.MakeBucket(ctx, m.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	policy := fmt.Sprintf(`{
	"Version":"2012-10-17",
	"Statement":[
		{
			"Effect":"Allow",
			"Principal":{"AWS":["*"]},
			"Action":["s3:GetObject"],
			"Resource":["arn:aws:s3:::%s/*"]
		}
	]}`, m.bucketName)

	err = m.client.SetBucketPolicy(ctx, m.bucketName, policy)
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}

	return nil
}

func (m *MinioStorage) UploadFile(ctx context.Context, fileName string, contentType string, content io.Reader, size int64) (string, error) {
	file, err := m.client.PutObject(ctx, m.bucketName, fileName, content, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	protocol := "http"
	if m.useSSL {
		protocol = "https"
	}

	url := fmt.Sprintf("%s://%s/%s/%s", protocol, m.endpoint, m.bucketName, file.Key)
	return url, nil
}
