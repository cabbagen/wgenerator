package middlewares

import (
	"fmt"
	"net/http"

	"github.com/cabbagen/wgenerator/v2/definitions"
	"github.com/gin-gonic/gin"
)

func HandlePanicRecoverMiddleware(ctx *gin.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			ctx.Abort()
			panicText := fmt.Sprintf("系统错误 - [url]: %s, [method]: %s, [error]: %s", ctx.Request.URL.String(), ctx.Request.Method, rec)
			ctx.JSON(http.StatusOK, definitions.NewResponse(definitions.ResponseCodeMap["fail"], nil, panicText))
		}
	}()
	ctx.Next()
}
