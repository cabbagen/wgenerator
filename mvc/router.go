package mvc

import (
	"github.com/cabbagen/wgenerator/context"
)

var GlobalAllRoutes []WGRoute

type WGHandleFunc func(*context.WGContext)

type WGRoute struct {
	Method  string
	Path    string
	Handles []WGHandleFunc
}

func AppendGlobalRoutes(routes ...[]WGRoute) {
	for _, route := range routes {
		GlobalAllRoutes = append(GlobalAllRoutes, route...)
	}
}
