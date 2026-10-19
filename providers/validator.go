package providers

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"

	"github.com/cabbagen/wgenerator/v2/definitions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const DefaultRequestJSONKey = "RequestJSON"

// 手机号码验证
func handleMobileRegister(vf validator.FieldLevel) bool {
	value := vf.Field().String()
	regexp := regexp.MustCompile(`^1[345789]{1}\d{9}$`)

	return regexp.MatchString(value)
}

// 密码校验 密码规则：大于 6 位、小于 18 位，可包含小写字母、大写字母、数字、特殊符号
func HandlePasswordRegister(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[\x21-\x7e]{7,17}$`).MatchString(fl.Field().String())
}

// 校验传参
func HandleValidateRequestParams[T any](params T) error {
	if error := Validate.Struct(params); error != nil {
		if _, ok := error.(*validator.InvalidValidationError); ok {
			return error
		}

		for _, error := range error.(validator.ValidationErrors) {
			field, ok := reflect.TypeOf(params).FieldByName(error.Field())

			if ok && field.Tag.Get("error") != "" {
				return errors.New(field.Tag.Get("error"))
			} else {
				return errors.New(error.Translate(Trans))
			}
		}
	}
	return nil
}

// gin 解析校验传参
func HandleValidateRequestParamsWithGin[T any](context *gin.Context, params *T) error {
	// 解析参数
	if error := context.ShouldBind(params); error != nil {
		return error
	}

	// 参数校验
	if error := HandleValidateRequestParams(*params); error != nil {
		return error
	}
	return nil
}

// gin 校验传参中间件
func HandleValidateRequestParamsWithGinMiddleware[T any]() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params T

		if error := HandleValidateRequestParamsWithGin(ctx, &params); error != nil {
			ctx.JSON(http.StatusOK, definitions.NewResponse(definitions.ResponseCodeMap["fail"], nil, error.Error()))
			ctx.Abort()
		} else {
			ctx.Set(DefaultRequestJSONKey, params)
			ctx.Next()
		}
	}
}
