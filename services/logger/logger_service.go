package logger_service

import (
	"fmt"
	env_models "go_boilerplate_project/models/env"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger is a wrapper around zap.SugaredLogger for ease of use
var Logger *zap.SugaredLogger

func InitLogger(cfg *env_models.LoggerEnv) (*zap.SugaredLogger, error) {
	logLevel := parseLogLevel(cfg.Level)

	consoleEncoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    coloredLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	fileEncoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(consoleEncoderCfg),
		zapcore.AddSync(ConsoleWriter{}),
		logLevel,
	)

	fileWriter := &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(fileEncoderCfg),
		zapcore.AddSync(fileWriter),
		logLevel,
	)

	core := zapcore.NewTee(consoleCore, fileCore)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	Logger = zapLogger.Sugar()
	return Logger, nil
}

func GetLogger() *zap.SugaredLogger {
	return Logger
}

func parseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// coloredLevelEncoder provides color output for different log levels
func coloredLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	var c string
	switch level {
	case zapcore.DebugLevel:
		c = "\033[36mDEBUG\033[0m" // Cyan
	case zapcore.InfoLevel:
		c = "\033[32mINFO\033[0m" // Green
	case zapcore.WarnLevel:
		c = "\033[33mWARN\033[0m" // Yellow
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		c = "\033[31mERROR\033[0m" // Red
	default:
		c = level.String()
	}
	enc.AppendString(c)
}

// ConsoleWriter writes logs to standard output, could be expanded for file-based logging
type ConsoleWriter struct{}

func (cw ConsoleWriter) Write(p []byte) (n int, err error) {
	return fmt.Print(string(p))
}

// GinLoggingMiddleware logs incoming HTTP requests in a pretty, colored style
func GinLoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		var statusColor, methodColor, resetColor string
		statusColor, methodColor, resetColor = colorForStatus(statusCode), colorForMethod(method), "\033[0m"

		if raw != "" {
			path = path + "?" + raw
		}

		Logger.Infof("|%s %3d %s| %13v |%s %-7s %s| %15s | %s | %s",
			statusColor, statusCode, resetColor,
			latency,
			methodColor, method, resetColor,
			clientIP,
			path,
			errorMessage,
		)
	}
}

// Helper functions for coloring output in Gin middleware
func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[32m" // Green
	case code >= 300 && code < 400:
		return "\033[36m" // Cyan
	case code >= 400 && code < 500:
		return "\033[33m" // Yellow
	default:
		return "\033[31m" // Red
	}
}
func colorForMethod(method string) string {
	switch method {
	case "GET":
		return "\033[32m" // Green
	case "POST":
		return "\033[34m" // Blue
	case "PUT":
		return "\033[33m" // Yellow
	case "DELETE":
		return "\033[31m" // Red
	case "PATCH":
		return "\033[36m" // Cyan
	case "OPTIONS":
		return "\033[35m" // Magenta
	default:
		return "\033[0m"
	}
}
