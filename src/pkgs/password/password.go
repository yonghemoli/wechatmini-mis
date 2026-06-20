package password

import (
	"golang.org/x/crypto/bcrypt"
)

// Hash 对明文密码进行 bcrypt 加密
func Hash(raw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	return string(bytes), err
}

// Verify 验证密码是否匹配
func Verify(hashed, raw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(raw))
	return err == nil
}
