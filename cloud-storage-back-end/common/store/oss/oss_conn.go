package oss

import (
	"fmt"
	"io"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossCli *oss.Client

// Client : 创建oss client对象
func Client() *oss.Client {
	if ossCli != nil {
		return ossCli
	}
	var err error
	ossCli, err = oss.New(OSSEndpoint, OSSAccesskeyID, OSSAccessKeySecret)
	if err != nil {
		fmt.Printf("创建OSS客户端失败: %v\n", err)
		return nil
	}
	return ossCli
}

// Bucket : 获取bucket存储空间
func Bucket() *oss.Bucket {
	cli := Client()
	if cli != nil {
		bucket, err := cli.Bucket(OSSBucket)
		if err != nil {
			fmt.Printf("获取Bucket失败: %v\n", err)
			return nil
		}
		return bucket
	}
	fmt.Println("OSS客户端未初始化")
	return nil
}

// DownloadURL(使用签名URL): 临时授权下载url
func DownloadURL(objName string) string {
	bucket := Bucket()
	if bucket == nil {
		fmt.Println("获取Bucket失败，无法生成下载URL")
		return ""
	}
	signedURL, err := bucket.SignURL(objName, oss.HTTPGet, 3600)
	if err != nil {
		fmt.Printf("生成下载URL失败: %v\n", err)
		return ""
	}
	return signedURL
}

// BuildLifecycleRule : 针对指定bucket设置生命周期规则
func BuildLifecycleRule(bucketName string) {
	ruleTest1 := oss.BuildLifecycleRuleByDays("rule1", "test/", true, 30)
	rules := []oss.LifecycleRule{ruleTest1}

	if err := Client().SetBucketLifecycle(bucketName, rules); err != nil {
		fmt.Printf("设置生命周期规则失败: %v\n", err)
	}
}

// GenFileMeta : 构造文件元信息
func GenFileMeta(metas map[string]string) []oss.Option {
	options := make([]oss.Option, 0, len(metas))
	for k, v := range metas {
		options = append(options, oss.Meta(k, v))
	}
	return options
}

// 文件上传
func Upload(inputStream io.Reader, path string, contentType string) (string, error) {
	metadata := GenFileMeta(map[string]string{
		"Content-Type": contentType,
	})

	err := Bucket().PutObject(path, inputStream, metadata...)
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %v", err)
	}

	return fmt.Sprintf("https://%s.%s/%s", OSSBucket, OSSEndpoint, path), nil
}

// 删除文件
//
// @param fullPath 文件完整路径
// @return 是否删除成功
func Delete(fullPath string) (bool, error) {
	if strings.TrimSpace(fullPath) == "" {
		return false, nil
	}

	err := Bucket().DeleteObject(getFileNameFromFullPath(fullPath))
	if err != nil {
		fmt.Printf("删除文件失败: %v\n", err)
		return false, err
	}
	return true, nil
}

// 获取文件名，如 1134664743.png
//
// @param fullPath 完整路径。如 http://qiniu.poile.cn/1134664743.png
// @return 文件名
func getFileNameFromFullPath(fullPath string) string {
	if fullPath != "" {
		return fullPath[strings.LastIndex(fullPath, "/")+1:]
	}
	return ""
}
