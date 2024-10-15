package mvc

import (
	"github.com/gin-gonic/gin"
)

var GlobalAllRoutes []WGRoute

type WGRoute struct {
	Path     string
	Method   string
	Handlers []gin.HandlerFunc
}

func AppendGlobalRoutes(routes ...[]WGRoute) {
	for _, route := range routes {
		GlobalAllRoutes = append(GlobalAllRoutes, route...)
	}
}
