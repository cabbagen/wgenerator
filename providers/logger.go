package providers

import (
	"log"
	"strings"

	"github.com/cabbagen/wgenerator/v2/conf"
	"github.com/cabbagen/wgenerator/v2/utils"
	"github.com/gookit/rotatefile"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

var singleLogger ILogger = nil

type ILogger interface {
	FlushAll()
	Logger(level slog.Level, args ...any)
	WithFieldsLogger(infos map[string]interface{}, level slog.Level, args ...any)
}

type logger struct {
	logger *slog.Logger
}

func (l *logger) FlushAll() {
	if err := l.logger.FlushAll(); err != nil {
		log.Printf("Error flushing logs: %v", err)
	}
}

func (l *logger) Logger(level slog.Level, args ...any) {
	l.logger.Log(level, args...)
}

func (l *logger) WithFieldsLogger(infos map[string]interface{}, level slog.Level, args ...any) {
	l.logger.WithFields(infos).Log(level, args...)
}

func newConsoleLogger(slogger *slog.Logger, levels []slog.Level, settings map[string]interface{}) {
	consoleHandler := handler.ConsoleWithLevels(levels)

	if settings["template"] != nil {
		consoleHandler.Formatter().(*slog.TextFormatter).SetTemplate(settings["template"].(string))
	}
	slogger.AddHandler(consoleHandler)
}

func newFileLogger(slogger *slog.Logger, levels []slog.Level, settings map[string]interface{}) {
	confFns := []handler.ConfigFn{
		handler.WithLogLevels(levels),
		handler.WithMaxSize(50 * 1024 * 1024),
		handler.WithRotateMode(rotatefile.ModeCreate),
		handler.WithCompress(true),
	}

	fileHandler, error := handler.NewRotateFileHandler(settings["file"].(string), rotatefile.EveryDay, confFns...)

	if error != nil {
		panic(error)
	}

	// NewRotateFileHandler 不会读取 WithUseJSON，需要单独设置 formatter。
	fileHandler.SetFormatter(slog.NewJSONFormatter())

	slogger.AddHandler(fileHandler)
}

func NewLoggerFactory() ILogger {
	if singleLogger != nil {
		return singleLogger
	}

	settings, _ := conf.ScanfBuildinYamlConfig()

	slogger, levels := slog.New(), utils.MapBySlice(settings["logger"]["levels"].([]interface{}), func(level interface{}, _ int) slog.Level {
		return slog.LevelByName(strings.ToUpper(level.(string)))
	})

	// 控制台日志
	if settings["logger"]["type"] == "console" {
		newConsoleLogger(slogger, levels, settings["logger"])
	}

	// 文件日志
	if settings["logger"]["type"] == "file" {
		newFileLogger(slogger, levels, settings["logger"])
	}

	singleLogger = &logger{
		logger: slogger,
	}

	return singleLogger
}
