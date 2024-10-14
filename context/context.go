package context

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type WGContext struct {
	gin.Context
}

func (wgc *WGContext) QueryByInt(key string) (int, error) {
	return strconv.Atoi(wgc.Query(key))
}

func (wgc *WGContext) DefaultQueryByInt(key string, defaultIntValue int) (int, error) {
	if wgc.Query(key) == "" {
		return defaultIntValue, nil
	}
	return wgc.QueryByInt(wgc.Query(key))
}
