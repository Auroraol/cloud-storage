package sms

import (
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

// Config 阿里云短信配置
type SmsConfig struct {
	AccessKeyID     string // 阿里云 AccessKey ID
	AccessKeySecret string // 阿里云 AccessKey Secret
	SignName        string // 短信签名
	TemplateCode    string // 短信模板ID
	RegionID        string // 地域ID，默认为 cn-hangzhou
	ExpiredAt       int64  // 过期时间(s)
}

// Client 短信客户端
type SmsClient struct {
	config SmsConfig
	client *dysmsapi.Client
	cache  *collection.Cache
}

// NewClient 创建短信客户端
func NewClient(config SmsConfig) (SmsClient, error) {
	if config.AccessKeyID == "" || config.AccessKeySecret == "" {
		return SmsClient{}, errors.New("AccessKeyID and AccessKeySecret cannot be empty")
	}
	if config.SignName == "" {
		return SmsClient{}, errors.New("SignName cannot be empty")
	}
	if config.TemplateCode == "" {
		return SmsClient{}, errors.New("TemplateCode cannot be empty")
	}

	// 设置默认区域
	if config.RegionID == "" {
		config.RegionID = "cn-hangzhou"
	}

	// 创建阿里云短信客户端
	client, err := dysmsapi.NewClientWithAccessKey(
		config.RegionID,
		config.AccessKeyID,
		config.AccessKeySecret,
	)
	if err != nil {
		return SmsClient{}, err
	}
	smsCache, err := collection.NewCache(time.Duration(config.ExpiredAt))
	if err != nil {
		zap.S().Errorf("创建SMS缓存失败: %s", err.Error())
		panic(err)
	}
	return SmsClient{
		config: config,
		client: client,
		cache:  smsCache,
	}, nil
}

// SendVerificationCode 发送验证码
// phone: 手机号码
func (c *SmsClient) SendVerificationCode(phone string) (string, error) {
	// 生成6位随机验证码
	verifyCode := generateRandomCode(6)

	// 构建短信发送请求
	request := dysmsapi.CreateSendMessageWithTemplateRequest()
	request.To = phone
	request.From = c.config.SignName             //短信签名
	request.TemplateCode = c.config.TemplateCode //短信模板编码
	request.TemplateParam = `{"code":"` + verifyCode + `"}`

	// 发送短信
	response, err := c.client.SendMessageWithTemplate(request)
	if err != nil {
		return "", err
	}

	// 检查发送结果
	if response.ResponseCode != "OK" {
		return "", fmt.Errorf("send sms failed: %s", response.ResponseDescription)
	}

	// 存储验证码
	c.saveVerificationCode(phone, verifyCode)

	return verifyCode, nil
}

// VerifyCode 验证短信验证码
func (c *SmsClient) VerifyCode(phone, inputVerifyCode string) bool {
	cachedCode, is := c.cache.Get(fmt.Sprintf("verifyCode_%s", phone))
	if !is {
		logx.Error("获取缓存验证码失败")
		return false
	}
	if cachedCode == inputVerifyCode {
		return true
	} else {
		return false
	}
}

// RemoveVerificationCode 移除验证码（验证成功后调用）
func (c *SmsClient) RemoveVerificationCode(phone string) {
	c.cache.Del(fmt.Sprintf("verifyCode_%s", phone))
}

// 保存验证码到内存(可以使用内存缓存(go-zero提供的)或者redis)
func (c *SmsClient) saveVerificationCode(phone, verifyCode string) {
	c.cache.Set(fmt.Sprintf("verifyCode_%s", phone), verifyCode)
}

// 生成指定长度的随机数字验证码
func generateRandomCode(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	code := ""
	for i := 0; i < length; i++ {
		code += fmt.Sprintf("%d", rand.Intn(10))
	}
	return code
}
