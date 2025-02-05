package oss

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
	Endpoint        string
}

var config Config

func init() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		fmt.Println("未能加载 .env 文件，使用系统环境变量")
	}

	if err := LoadConfig(); err != nil {
		panic(fmt.Sprintf("加载OSS配置失败: %v", err))
	}
}

func LoadConfig() error {
	config = Config{
		AccessKeyId:     os.Getenv("OSS_ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("OSS_ACCESS_KEY_SECRET"),
		BucketName:      os.Getenv("OSS_BUCKET_NAME"),
		Endpoint:        os.Getenv("OSS_ENDPOINT"),
	}

	if config.AccessKeyId == "" || config.AccessKeySecret == "" || config.BucketName == "" || config.Endpoint == "" {
		return fmt.Errorf("OSS配置未正确设置，请检查环境变量")
	}
	return nil
}

func GetConfig() Config {
	return config
}
