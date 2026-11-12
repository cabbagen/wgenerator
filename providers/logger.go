package providers

import (
	"log"
	"strings"
	"sync"

	"github.com/cabbagen/wgenerator/v2/conf"
	"github.com/cabbagen/wgenerator/v2/utils"
	"github.com/gookit/rotatefile"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

var (
	once         sync.Once
	singleLogger ILogger
)

type ILogger interface {
	FlushAll()
	Slog() *slog.Logger
	Logger(level slog.Level, args ...any)
	WithFieldsLogger(infos map[string]interface{}, level slog.Level, args ...any)
}

type logger struct {
	// 日志实例
	logger *slog.Logger

	// 输出类型
	output []string

	// 日志级别
	levels []slog.Level

	// file 日志文件路径
	file string

	// console 模板
	template string
}

func (l *logger) FlushAll() {
	if err := l.logger.FlushAll(); err != nil {
		log.Printf("Error flushing logs: %v", err)
	}
}

func (l *logger) Slog() *slog.Logger {
	return l.logger
}

func (l *logger) Logger(level slog.Level, args ...any) {
	l.logger.Log(level, args...)
}

func (l *logger) WithFieldsLogger(infos map[string]interface{}, level slog.Level, args ...any) {
	l.logger.WithFields(infos).Log(level, args...)
}

func appendConsoleLoggerHandler(slogger *slog.Logger, levels []slog.Level, settings conf.LoggerConfig) {
	consoleHandler := handler.ConsoleWithLevels(levels)

	if settings.Template != "" {
		consoleHandler.Formatter().(*slog.TextFormatter).SetTemplate(settings.Template)
	}
	slogger.AddHandler(consoleHandler)
}

func appendFileLoggerHandler(slogger *slog.Logger, levels []slog.Level, settings conf.LoggerConfig) {
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

	fileHandler.SetFormatter(slog.NewJSONFormatter())

	slogger.AddHandler(fileHandler)
}

func NewLoggerFactory() ILogger {
	once.Do(func() {
		if singleLogger != nil {
			return
		}

		settings, _ := conf.ScanfBuildinYamlConfig()

		slogger, levels := slog.New(), utils.MapBySlice(settings.Logger.Levels, func(level string, _ int) slog.Level {
			return slog.LevelByName(strings.ToUpper(level))
		})

		for _, output := range settings.Logger.Output {
			if output == "file" {
				appendFileLoggerHandler(slogger, levels, settings.Logger)
			}
			if output == "console" {
				appendConsoleLoggerHandler(slogger, levels, settings.Logger)
			}
		}

		singleLogger = &logger{
			logger:   slogger,
			output:   settings.Logger.Output,
			levels:   levels,
			file:     settings.Logger.File,
			template: settings.Logger.Template,
		}
	})
	return singleLogger
}
