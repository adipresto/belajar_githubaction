package file

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/arieffadhlan/go-fitbyte/internal/config"
	minioUploader "github.com/arieffadhlan/go-fitbyte/internal/pkg/minio"
	"github.com/minio/minio-go/v7"
)

type useCase struct {
	config config.Config
	minio  *minio.Client
}

func NewUseCase(config config.Config) UseCase {
	minioConfig := &minioUploader.MinioConfig{
		AccessKeyID:     config.Minio.AccessKeyID,
		SecretAccessKey: config.Minio.SecretAccessKey,
		UseSSL:          config.Minio.UseSSL,
		Endpoint:        config.Minio.Endpoint,
	}

	minioClient, _ := minioUploader.NewUploader(minioConfig)

	return &useCase{
		config: config,
		minio:  minioClient,
	}
}

func (uc *useCase) UploadFile(
	ctx context.Context,
	file *multipart.FileHeader,
	src multipart.File,
	fileName string,
) (string, error) {

	// Switch via env
	if uc.config.LOCAL_FILE == "true" {
		return uc.uploadLocal(ctx, file, src, fileName)
	}

	// default: MinIO (existing behavior)
	return uc.uploadMinio(ctx, file, src, fileName)
}

func (uc *useCase) GetFile(
	ctx context.Context,
	fileName string,
) (*os.File, error) {

	fullPath := filepath.Join("./storage/uploads", fileName)

	fmt.Println(fullPath)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (uc *useCase) uploadLocal(
	ctx context.Context,
	file *multipart.FileHeader,
	src multipart.File,
	fileName string,
) (string, error) {

	basePath := "./storage/uploads"

	if err := os.MkdirAll(basePath, 0755); err != nil {
		return "", err
	}

	dstPath := filepath.Join(basePath, fileName)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	if v := ctx.Value("request_url"); v != nil {
		if s, ok := v.(string); ok {
			return s + fileName, nil
		}
	}

	// return identifier, NOT absolute path
	return fileName, nil
}

func (uc *useCase) uploadMinio(
	ctx context.Context,
	file *multipart.FileHeader,
	src multipart.File,
	fileName string,
) (string, error) {

	bucketName := uc.config.Minio.BucketName

	_, err := uc.minio.PutObject(
		ctx,
		bucketName,
		fileName,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", err
	}

	presignedURL, err := uc.minio.PresignedGetObject(
		ctx,
		bucketName,
		fileName,
		time.Hour*24*7,
		nil,
	)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
