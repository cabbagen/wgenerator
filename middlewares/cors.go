package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var defaultCorsOptions map[string]string = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Headers": "content-type, token, mode, credentials, uid, x-requested-with",
}

func HandleCorsMiddleware(ctx *gin.Context) {
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
