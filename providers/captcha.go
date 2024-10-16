package providers

import "github.com/mojocn/base64Captcha"

type DigitCaptcha struct {
	B64s      string `json:"b64s"`
	CaptchaId string `json:"captchaId"`
}

var store base64Captcha.Store = base64Captcha.DefaultMemStore

var driver base64Captcha.Driver = base64Captcha.DefaultDriverDigit

func HandleGenerateDigitCaptcha() (DigitCaptcha, error) {
	captchaId, b64s, _, error := base64Captcha.NewCaptcha(driver, store).Generate()

	if error != nil {
		return DigitCaptcha{}, error
	}
	return DigitCaptcha{b64s, captchaId}, nil
}

func HandleVerifyDigitCaptcha(captchaId, answer string) bool {
	return store.Verify(captchaId, answer, true)
}
