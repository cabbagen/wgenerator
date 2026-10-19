package providers

import (
	"fmt"
	"log"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// 全局自定义校验
var GlobalCustomValidator map[string]validator.Func = map[string]validator.Func{
	"mobile":   handleMobileRegister,
	"password": HandlePasswordRegister,
}

// 定义一个全局翻译器T
var (
	Trans    ut.Translator
	Validate *validator.Validate
)

// InitTrans 初始化翻译器
func InitTrans(locale string) error {

	zhT, enT, ok := zh.New(), en.New(), false

	uni := ut.New(enT, zhT, enT)

	Validate = validator.New()

	// 绑定自定义校验
	for validatorName, validatorFunc := range GlobalCustomValidator {
		Validate.RegisterValidation(validatorName, validatorFunc)
	}

	Trans, ok = uni.GetTranslator(locale)

	if !ok {
		return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
	}

	if locale == "zh" {
		return zhTranslations.RegisterDefaultTranslations(Validate, Trans)
	}

	return enTranslations.RegisterDefaultTranslations(Validate, Trans)
}

func init() {
	if error := InitTrans("zh"); error != nil {
		log.Fatalf("init trans failed, err:%s\n", error.Error())
		return
	}
}
