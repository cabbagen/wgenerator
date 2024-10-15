package middlewares

import (
	"fmt"
	"net/http"

	"github.com/cabbagen/wgenerator/definitions"
	"github.com/gin-gonic/gin"
)

func HandlePanicRecover(ctx *gin.Context) {
	defer func() {
		if error := recover(); error != nil {
			var panicText = fmt.Sprintf("系统错误 - [url]: %s, [method]: %s, [error]: %s", ctx.Request.URL.String(), ctx.Request.Method, error)
			ctx.JSON(http.StatusOK, definitions.NewResponse(definitions.ResponseCodeMap["fail"], nil, panicText))
		}
	}()
	ctx.Next()
}
