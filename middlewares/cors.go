package middlewares

import (
	"net/http"

	"github.com/cabbagen/wgenerator/context"
)

var defaultCorsOptions map[string]string = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Headers": "content-type, token, mode, credentials, uid, x-requested-with",
}

func HandleCorsMiddleware(ctx *context.WGContext) {
	for key, value := range defaultCorsOptions {
		ctx.Header(key, value)
	}

	if ctx.Request.Method == "OPTIONS" {
		ctx.String(http.StatusOK, "true")
		ctx.AbortWithStatus(http.StatusOK)
		return
	}
	ctx.Next()
}
