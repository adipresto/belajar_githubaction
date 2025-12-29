package config

import "os"

type aws struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	BucketName      string
	AccountID       string
}

func loadAwsConfig() aws {
	return aws{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Region:          os.Getenv("AWS_REGION"),
		AccountID:       os.Getenv("AWS_ACCOUNT_ID"),
		BucketName:      os.Getenv("S3_BUCKET_NAME"),
	}
}
