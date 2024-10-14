package providers

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 手机号码验证
func handleMobileRegister(vf validator.FieldLevel) bool {
	value := vf.Field().String()
	regexp := regexp.MustCompile(`^1[345789]{1}\d{9}$`)

	return regexp.MatchString(value)
}
