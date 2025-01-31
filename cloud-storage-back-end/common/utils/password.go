package utils

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
)

// 使用 bcrypt是一种单向哈希算法，只用于加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// md5加密
func Md5ByString(str string) string {
	m := md5.New()
	_, err := io.WriteString(m, str)
	if err != nil {
		panic(err)
	}
	arr := m.Sum(nil)
	return fmt.Sprintf("%x", arr)
}

func Md5ByBytes(b []byte) string {
	return fmt.Sprintf("%x", md5.Sum(b))
}

// ComparePassword 比较密码是否匹配
func ComparePassword(md5Password, password string) bool {
	return Md5ByString(password) == md5Password
}
