package mvc

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/cabbagen/wgenerator/definitions"
	"github.com/cabbagen/wgenerator/providers"
	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

// 响应失败
func (tbc BaseController) OfFailResponse(ctx *gin.Context, message string) {
	ctx.JSON(http.StatusOK, definitions.NewResponse(definitions.ResponseCodeMap["fail"], nil, message))
}

// 响应成功
func (tbc BaseController) OfSuccessResponse(ctx *gin.Context, data interface{}) {
	fmt.Printf("========> %d\n\n\n", definitions.ResponseCodeMap["succees"])
	ctx.JSON(http.StatusOK, definitions.NewResponse(definitions.ResponseCodeMap["succees"], data, ""))
}

// 响应二进制数据
func (tbc BaseController) OfSuccessBytesResponse(ctx *gin.Context, filename string, thunk []byte) {
	if filename != "" {
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
	}

	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Data(http.StatusOK, "application/octet-stream", thunk)
}

// 校验结构体
func (tbc BaseController) HandleValidateRequestParams(ctx *gin.Context, target *interface{}) error {
	return providers.HandleValidateRequestParamsWithGin(ctx, target)
}

// 文件上传
func (tbc BaseController) HandleFileMultipartForm(ctx *gin.Context) (*multipart.Form, error) {
	form, error := ctx.MultipartForm()

	if error != nil {
		return nil, fmt.Errorf("接口上传参数读取失败: %s", error.Error())
	}
	return form, nil
}
