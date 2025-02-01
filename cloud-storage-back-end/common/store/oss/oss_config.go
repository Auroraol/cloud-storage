package oss

import (
	"os"
)

type OSSConfig struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
}

var Config = OSSConfig{
	Endpoint:        os.Getenv("OSS_ENDPOINT"),
	AccessKeyId:     os.Getenv("OSS_ACCESS_KEY_ID"),
	AccessKeySecret: os.Getenv("OSS_ACCESS_KEY_SECRET"),
	BucketName:      os.Getenv("OSS_BUCKET_NAME"),
}
