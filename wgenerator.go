package wgenerator

import (
	"fmt"
	"net/http"

	"github.com/cabbagen/wgenerator/caches"
	"github.com/cabbagen/wgenerator/conf"
	"github.com/cabbagen/wgenerator/databases"

	"github.com/gin-gonic/gin"
)

type WGEngine struct {
	*gin.Engine
}

func getApplicationConfs() map[string]map[string]interface{} {
	settings, error := conf.ScanfBuildinYamlConfig()

	if error != nil {
		panic("配置文件不存在，请先添加项目配置文件")
	}
	return settings
}

func initApplicationPresets(engine *gin.Engine) {
	settings := getApplicationConfs()

	if settings["server"]["static"] != nil {
		engine.Static("public", settings["server"]["static"].(string))
	}

	if settings["server"]["templateDir"] != nil {
		engine.LoadHTMLGlob(fmt.Sprintf("%s/**/*.html", settings["server"]["templateDir"].(string)))
	}

	if settings["database"]["dbname"] != nil && settings["database"]["username"] != nil && settings["database"]["address"] != nil && settings["database"]["password"] != nil {
		databases.ConnectMysql(settings["database"]["username"].(string), settings["database"]["password"].(string), settings["database"]["address"].(string), settings["database"]["dbname"].(string))
	}

	if settings["cacher"]["address"] != nil {
		caches.InitRedisCacherInstance(settings["cacher"]["db"].(int), settings["cacher"]["address"].(string), settings["cacher"]["password"].(string))
	}

	if settings["server"]["templateDir"] != nil && settings["server"]["isOpenSupportSpa"] != nil && settings["server"]["isOpenSupportSpa"].(int) == 1 {
		handleRenderSPAHTMLFunc := func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", nil)
		}
		engine.Handle("GET", "/", handleRenderSPAHTMLFunc)
		engine.NoRoute(handleRenderSPAHTMLFunc)
	}
}

func WGDefault(withFuncs ...gin.OptionFunc) WGEngine {
	return WGEngine{gin.Default(append(withFuncs, initApplicationPresets)...)}
}

func (wge WGEngine) WGRun() {
	wge.Run(fmt.Sprintf(":%s", getApplicationConfs()["server"]["port"].(string)))
}
