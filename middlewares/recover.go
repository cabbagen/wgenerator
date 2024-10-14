package middlewares

import (
	"fmt"
	"net/http"
	"wgenerator/context"
	"wgenerator/definitions"
)

func HandlePanicRecover(c *context.WGContext) {
	defer func() {
		if error := recover(); error != nil {
			var panicText = fmt.Sprintf("系统错误 - [url]: %s, [method]: %s, [error]: %s", c.Request.URL.String(), c.Request.Method, error)
			c.JSON(http.StatusOK, definitions.NewResponse(definitions.ResponseCodeMap["fail"], nil, panicText))
		}
	}()
	c.Next()
}
