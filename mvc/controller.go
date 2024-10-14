package mvc

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"wgenerator/context"
	"wgenerator/definitions"
	"wgenerator/providers"
)

type BaseController struct {
}

// 响应失败
func (tbc BaseController) OfFailResponse(ctx *context.WGContext, message string) {
	ctx.JSON(http.StatusOK, definitions.NewResponse(definitions.ResponseCodeMap["fail"], nil, message))
}

// 响应成功
func (tbc BaseController) OfSuccessResponse(ctx *context.WGContext, data interface{}) {
	ctx.JSON(http.StatusOK, definitions.NewResponse(definitions.ResponseCodeMap["succees"], data, ""))
}

// 响应二进制数据
func (tbc BaseController) OfSuccessBytesResponse(ctx *context.WGContext, filename string, thunk []byte) {
	if filename == "" {
		ctx.Header("Content-Transfer-Encoding", "binary")
	} else {
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
	}
	ctx.Data(http.StatusOK, "application/octet-stream", thunk)
}

// 校验结构体
func (tbc BaseController) HandleValidateRequestParams(ctx *context.WGContext, target *interface{}) error {
	return providers.HandleValidateRequestParamsWithGin(ctx, target)
}

// 文件上传
func (tbc BaseController) HandleFileMultipartForm(ctx *context.WGContext) (*multipart.Form, error) {
	form, error := ctx.MultipartForm()

	if error != nil {
		return nil, fmt.Errorf("接口上传参数读取失败: %s", error.Error())
	}
	return form, nil
}
