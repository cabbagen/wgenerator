package wgenerator

import (
	"fmt"
	"net/http"
	"wgenerator/caches"
	"wgenerator/conf"
	"wgenerator/databases"

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

func WGDefault(withFuncs ...gin.OptionFunc) WGEngine {
	engine, settings := gin.Default(withFuncs...), getApplicationConfs()

	// 静态目录
	if settings["server"]["StaticDir"] != "" {
		engine.Static("public", settings["server"]["StaticDir"].(string))
	}

	// 模板文件
	if settings["server"]["TemplateDir"] != "" {
		engine.LoadHTMLGlob(fmt.Sprintf("%s/**/*.html", settings["server"]["TemplateDir"].(string)))
	}

	// 数据库支持
	if settings["database"]["dbname"] != "" && settings["database"]["username"] != "" && settings["database"]["password"] != "" {
		databases.ConnectMysql(settings["database"]["username"].(string), settings["database"]["password"].(string), settings["database"]["dbname"].(string))
	}

	// redis 支持
	if settings["cacher"]["address"] != "" {
		caches.InitRedisCacherInstance(settings["cacher"]["db"].(int), settings["cacher"]["address"].(string), settings["cacher"]["password"].(string))
	}

	// SPA 支持
	if settings["server"]["TemplateDir"] != "" && settings["server"]["IsOpenSupportSpa"].(int) == 1 {
		handleRenderSPAHTMLFunc := func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", nil)
		}
		engine.Handle("GET", "/", handleRenderSPAHTMLFunc)
		engine.NoRoute(handleRenderSPAHTMLFunc)
	}

	return WGEngine{engine}
}

func (wge WGEngine) Bootstrap() {
	wge.Run(getApplicationConfs()["server"]["port"].(string))
}
