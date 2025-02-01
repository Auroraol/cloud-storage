package oss

import "os"

const (
	// OSSBucket : oss bucket名
	OSSBucket = "cloud-storage-1"
	// OSSEndpoint : oss endpoint
	OSSEndpoint = "oss-cn-beijing.aliyuncs.com"
)

// GetOSSAccessKeyID 从环境变量获取AccessKeyID
func GetOSSAccessKeyID() string {
	return os.Getenv("ALIYUN_OSS_ACCESS_KEY_ID")
}

// GetOSSAccessKeySecret 从环境变量获取AccessKeySecret
func GetOSSAccessKeySecret() string {
	return os.Getenv("ALIYUN_OSS_ACCESS_KEY_SECRET")
} 