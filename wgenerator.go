package wgenerator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cabbagen/wgenerator/v2/caches"
	"github.com/cabbagen/wgenerator/v2/conf"
	"github.com/cabbagen/wgenerator/v2/databases"
	"github.com/gin-gonic/gin"
)

type WGEngine struct {
	*gin.Engine
}

func getApplicationConfs() conf.ApplicationConfig {
	settings, error := conf.ScanfBuildinYamlConfig()

	if error != nil {
		panic("配置文件不存在，请先添加项目配置文件")
	}
	return settings
}

func initApplicationPresets(engine *gin.Engine) {
	settings := getApplicationConfs()

	if settings.Server.Static != "" {
		engine.Static("public", settings.Server.Static)
	}

	if settings.Server.TemplateDir != "" {
		engine.LoadHTMLGlob(fmt.Sprintf("%s/**/*.html", settings.Server.TemplateDir))
	}

	if settings.Database.DBName != "" && settings.Database.Username != "" && settings.Database.Address != "" && settings.Database.Password != "" {
		databases.ConnectMysql(settings.Database.Username, settings.Database.Password, settings.Database.Address, settings.Database.DBName)
	}

	if settings.Cacher.Address != "" {
		caches.InitRedisCacherInstance(settings.Cacher.DB, settings.Cacher.Address, settings.Cacher.Password)
	}

	if settings.Server.TemplateDir != "" && settings.Server.IsOpenSupportSpa == 1 {
		handleRenderSPAHTMLFunc := func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", nil)
		}
		engine.Handle("GET", "/", handleRenderSPAHTMLFunc)
		engine.NoRoute(handleRenderSPAHTMLFunc)
	}

	// 添加日志中间件
	if mode := gin.Mode(); mode == "debug" || mode == "" {
		engine.Use(gin.Logger())
	}
}

func WGDefault(withFuncs ...gin.OptionFunc) WGEngine {
	return WGEngine{gin.New(append(withFuncs, initApplicationPresets)...)}
}

func (wge WGEngine) WGRun(shutdownHooks ...func()) {
	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%s", getApplicationConfs().Server.Port),
		Handler: wge.Engine,
	}

	errCh := make(chan error, 1)

	go func() {
		errCh <- httpSrv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		signal.Stop(quit)
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("http server error: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Printf("http server shutdown: %v", err)
	}

	for _, hook := range shutdownHooks {
		if hook != nil {
			hook()
		}
	}
}
