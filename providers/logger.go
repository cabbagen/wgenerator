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

func newConsoleLogger(slogger *slog.Logger, levels []slog.Level, settings conf.LoggerConfig) {
	consoleHandler := handler.ConsoleWithLevels(levels)

	if settings.Template != "" {
		consoleHandler.Formatter().(*slog.TextFormatter).SetTemplate(settings.Template)
	}
	slogger.AddHandler(consoleHandler)
}

func newFileLogger(slogger *slog.Logger, levels []slog.Level, settings conf.LoggerConfig) {
	confFns := []handler.ConfigFn{
		handler.WithLogLevels(levels),
		handler.WithMaxSize(50 * 1024 * 1024),
		handler.WithRotateMode(rotatefile.ModeCreate),
		handler.WithCompress(true),
	}

	fileHandler, error := handler.NewRotateFileHandler(settings.File, rotatefile.EveryDay, confFns...)

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

	slogger, levels := slog.New(), utils.MapBySlice(settings.Logger.Levels, func(level string, _ int) slog.Level {
		return slog.LevelByName(strings.ToUpper(level))
	})

	// 控制台日志
	if settings.Logger.Type == "console" {
		newConsoleLogger(slogger, levels, settings.Logger)
	}

	// 文件日志
	if settings.Logger.Type == "file" {
		newFileLogger(slogger, levels, settings.Logger)
	}

	singleLogger = &logger{
		logger: slogger,
	}

	return singleLogger
}
