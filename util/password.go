package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// 使用bctypt散列函数加密用户密码
func HashPassword(password string) (string, error) {
	hashword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashword), nil
}

// 检查散列函数
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
