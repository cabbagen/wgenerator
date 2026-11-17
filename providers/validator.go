package providers

import (
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"strconv"

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

func handleApplyDefaultTags(params any) {
	value := reflect.ValueOf(params)

	if value.Kind() != reflect.Ptr || value.IsNil() {
		return
	}

	value = value.Elem()

	if value.Kind() != reflect.Struct {
		return
	}

	structType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		if !field.CanSet() || !field.IsZero() {
			continue
		}

		defaultValue := structType.Field(i).Tag.Get("default")

		if defaultValue == "" {
			continue
		}

		switch field.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if parsed, err := strconv.ParseInt(defaultValue, 10, 64); err == nil {
				field.SetInt(parsed)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if parsed, err := strconv.ParseUint(defaultValue, 10, 64); err == nil {
				field.SetUint(parsed)
			}
		case reflect.Float32, reflect.Float64:
			if parsed, err := strconv.ParseFloat(defaultValue, 64); err == nil {
				field.SetFloat(parsed)
			}
		case reflect.Bool:
			if parsed, err := strconv.ParseBool(defaultValue); err == nil {
				field.SetBool(parsed)
			}
		case reflect.String:
			field.SetString(defaultValue)
		}
	}
}

// gin 解析校验传参
func HandleValidateRequestParamsWithGin[T any](context *gin.Context, params *T) error {
	// 解析参数
	if error := context.ShouldBind(params); error != nil {
		return error
	}

	handleApplyDefaultTags(params)

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
